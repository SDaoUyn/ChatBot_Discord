package db

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

const (
	redisSubscriberKey = "discord:subscribers"
	redisAdminKey      = "discord:admins"
	redisTTL           = 24 * time.Hour
)

// Subscriber đại diện cho người dùng đã đăng ký nhận thông báo
type Subscriber struct {
	ID            int
	DiscordUserID int64
	SubscribedAt  time.Time
	UpdatedAt     time.Time
	IsActive      bool
	IsAdmin       bool
}

// AddSubscriber thêm người đăng ký mới vào PostgreSQL và Redis
func AddSubscriber(discordUserID int64) error {
	// Chuyển đổi discord_user_id thành string để lưu trong Redis
	userIDStr := strconv.FormatInt(discordUserID, 10)

	// Thêm vào PostgreSQL
	_, err := PgPool.Exec(context.Background(), `
		INSERT INTO subscribers (discord_user_id)
		VALUES ($1)
		ON CONFLICT (discord_user_id)
		DO UPDATE SET is_active = TRUE
	`, discordUserID)
	if err != nil {
		return fmt.Errorf("lỗi khi thêm người đăng ký vào cơ sở dữ liệu: %v", err)
	}

	// Thêm vào Redis cache
	err = RedisClient.SAdd(ctx, redisSubscriberKey, userIDStr).Err()
	if err != nil {
		return fmt.Errorf("lỗi khi thêm người đăng ký vào Redis: %v", err)
	}

	// Đặt TTL cho Redis key
	RedisClient.Expire(ctx, redisSubscriberKey, redisTTL)

	return nil
}

// RemoveSubscriber hủy đăng ký người dùng từ PostgreSQL và Redis
func RemoveSubscriber(discordUserID int64) error {
	// Chuyển đổi discord_user_id thành string để sử dụng trong Redis
	userIDStr := strconv.FormatInt(discordUserID, 10)

	// Đặt is_active = FALSE thay vì xóa bản ghi
	_, err := PgPool.Exec(context.Background(), `
		UPDATE subscribers
		SET is_active = FALSE
		WHERE discord_user_id = $1
	`, discordUserID)
	if err != nil {
		return fmt.Errorf("lỗi khi cập nhật trạng thái người đăng ký trong cơ sở dữ liệu: %v", err)
	}

	// Xóa khỏi Redis cache
	err = RedisClient.SRem(ctx, redisSubscriberKey, userIDStr).Err()
	if err != nil {
		return fmt.Errorf("lỗi khi xóa người đăng ký khỏi Redis: %v", err)
	}

	return nil
}

// IsSubscribed kiểm tra xem người dùng có đang đăng ký và hoạt động không
func IsSubscribed(discordUserID int64) (bool, error) {
	// Chuyển đổi discord_user_id thành string để sử dụng trong Redis
	userIDStr := strconv.FormatInt(discordUserID, 10)

	// Kiểm tra Redis trước
	exists, err := RedisClient.SIsMember(ctx, redisSubscriberKey, userIDStr).Result()
	if err == nil && exists {
		return true, nil
	}

	// Truy vấn PostgreSQL nếu không tìm thấy trong Redis
	var isActive bool
	err = PgPool.QueryRow(context.Background(), `
		SELECT is_active FROM subscribers
		WHERE discord_user_id = $1
	`, discordUserID).Scan(&isActive)

	if err != nil {
		// Nếu không tìm thấy bản ghi, người dùng chưa đăng ký
		return false, nil
	}

	// Nếu tìm thấy trong database và is_active=true nhưng không có trong Redis, cập nhật Redis
	if isActive && !exists {
		RedisClient.SAdd(ctx, redisSubscriberKey, userIDStr)
		RedisClient.Expire(ctx, redisSubscriberKey, redisTTL)
	}

	return isActive, nil
}

// GetAllActiveSubscribers lấy tất cả người đăng ký đang hoạt động
func GetAllActiveSubscribers() ([]Subscriber, error) {
	// Truy vấn PostgreSQL để lấy danh sách đầy đủ
	rows, err := PgPool.Query(context.Background(), `
		SELECT id, discord_user_id, subscribed_at, updated_at, is_active, is_admin
		FROM subscribers
		WHERE is_active = TRUE
	`)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi lấy danh sách người đăng ký: %v", err)
	}
	defer rows.Close()

	var subscribers []Subscriber
	for rows.Next() {
		var sub Subscriber
		if err := rows.Scan(&sub.ID, &sub.DiscordUserID, &sub.SubscribedAt, &sub.UpdatedAt, &sub.IsActive, &sub.IsAdmin); err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu người đăng ký: %v", err)
		}
		subscribers = append(subscribers, sub)
	}

	// Cập nhật Redis cache
	if len(subscribers) > 0 {
		pipe := RedisClient.Pipeline()
		pipe.Del(ctx, redisSubscriberKey) // Xóa key cũ
		pipe.Del(ctx, redisAdminKey)      // Xóa key admin cũ

		for _, sub := range subscribers {
			userIDStr := strconv.FormatInt(sub.DiscordUserID, 10)
			pipe.SAdd(ctx, redisSubscriberKey, userIDStr)

			// Nếu là admin, thêm vào danh sách admin trong Redis
			if sub.IsAdmin {
				pipe.SAdd(ctx, redisAdminKey, userIDStr)
			}
		}
		pipe.Expire(ctx, redisSubscriberKey, redisTTL)
		pipe.Expire(ctx, redisAdminKey, redisTTL)
		_, err = pipe.Exec(ctx)
		if err != nil {
			// Ghi log lỗi nhưng không làm chức năng thất bại
			fmt.Printf("Cảnh báo: Không thể cập nhật bộ nhớ đệm Redis: %v\n", err)
		}
	}

	return subscribers, nil
}

// IsAdmin kiểm tra xem người dùng có phải là admin không
func IsAdmin(discordUserID int64) (bool, error) {
	// Chuyển đổi discord_user_id thành string để sử dụng trong Redis
	userIDStr := strconv.FormatInt(discordUserID, 10)

	// Kiểm tra Redis trước
	isAdmin, err := RedisClient.SIsMember(ctx, redisAdminKey, userIDStr).Result()
	if err == nil && isAdmin {
		return true, nil
	}

	// Truy vấn PostgreSQL nếu không tìm thấy trong Redis
	var admin bool
	err = PgPool.QueryRow(context.Background(), `
		SELECT is_admin FROM subscribers
		WHERE discord_user_id = $1 AND is_active = TRUE
	`, discordUserID).Scan(&admin)

	if err != nil {
		// Nếu không tìm thấy bản ghi, người dùng không phải admin
		return false, nil
	}

	// Nếu tìm thấy trong database và is_admin=true nhưng không có trong Redis, cập nhật Redis
	if admin && !isAdmin {
		RedisClient.SAdd(ctx, redisAdminKey, userIDStr)
		RedisClient.Expire(ctx, redisAdminKey, redisTTL)
	}

	return admin, nil
}

// SetAdmin đặt hoặc hủy quyền admin cho người dùng
func SetAdmin(discordUserID int64, isAdmin bool) error {
	// Chuyển đổi discord_user_id thành string để sử dụng trong Redis
	userIDStr := strconv.FormatInt(discordUserID, 10)

	// Cập nhật PostgreSQL
	_, err := PgPool.Exec(context.Background(), `
		UPDATE subscribers
		SET is_admin = $2, updated_at = NOW()
		WHERE discord_user_id = $1
		RETURNING id
	`, discordUserID, isAdmin)

	if err != nil {
		return fmt.Errorf("lỗi khi đặt quyền admin: %v", err)
	}

	// Cập nhật Redis cache
	if isAdmin {
		err = RedisClient.SAdd(ctx, redisAdminKey, userIDStr).Err()
	} else {
		err = RedisClient.SRem(ctx, redisAdminKey, userIDStr).Err()
	}

	if err != nil {
		return fmt.Errorf("lỗi khi cập nhật bộ nhớ đệm Redis: %v", err)
	}

	return nil
}

// MigrateAddAdminColumn thêm cột is_admin nếu chưa tồn tại
func MigrateAddAdminColumn() error {
	// Kiểm tra xem cột đã tồn tại chưa
	var columnExists bool
	err := PgPool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'subscribers' AND column_name = 'is_admin'
		)
	`).Scan(&columnExists)

	if err != nil {
		return fmt.Errorf("lỗi khi kiểm tra cột is_admin: %v", err)
	}

	// Nếu cột chưa tồn tại, thêm vào
	if !columnExists {
		_, err = PgPool.Exec(context.Background(), `
			ALTER TABLE subscribers 
			ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE
		`)
		if err != nil {
			return fmt.Errorf("lỗi khi thêm cột is_admin: %v", err)
		}
		fmt.Println("Đã thêm cột is_admin vào bảng subscribers")
	}

	return nil
}
