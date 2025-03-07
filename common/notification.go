package common

import (
	"ChatBotDiscord/db"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ProcessSubscribeCommand xử lý lệnh !subscribe
func ProcessSubscribeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Lấy thông tin người dùng
	userID, err := strconv.ParseInt(m.Author.ID, 10, 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lỗi khi xử lý ID người dùng: %v", err))
		return
	}

	// Kiểm tra xem đã đăng ký chưa
	isSubscribed, err := db.IsSubscribed(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lỗi khi kiểm tra trạng thái đăng ký: %v", err))
		return
	}

	if isSubscribed {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, bạn đã đăng ký nhận thông báo rồi!", m.Author.ID))
		return
	}

	// Thêm người đăng ký
	err = db.AddSubscriber(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lỗi khi đăng ký: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ <@%s>, bạn đã đăng ký nhận thông báo thành công!", m.Author.ID))
}

// ProcessUnsubscribeCommand xử lý lệnh !unsubscribe
func ProcessUnsubscribeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Lấy thông tin người dùng
	userID, err := strconv.ParseInt(m.Author.ID, 10, 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lỗi khi xử lý ID người dùng: %v", err))
		return
	}

	// Kiểm tra xem đã đăng ký chưa
	isSubscribed, err := db.IsSubscribed(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lỗi khi kiểm tra trạng thái đăng ký: %v", err))
		return
	}

	if !isSubscribed {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, bạn chưa đăng ký nhận thông báo!", m.Author.ID))
		return
	}

	// Hủy đăng ký
	err = db.RemoveSubscriber(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lỗi khi hủy đăng ký: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ <@%s>, bạn đã hủy đăng ký nhận thông báo thành công!", m.Author.ID))
}

// ProcessSetAdminCommand xử lý lệnh !setadmin
func ProcessSetAdminCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Kiểm tra xem người gửi có phải là admin không
	authorID, _ := strconv.ParseInt(m.Author.ID, 10, 64)
	isAdmin, err := db.IsAdmin(authorID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lỗi khi kiểm tra quyền admin: %v", err))
		return
	}

	// Nếu chưa có admin nào, người đầu tiên sử dụng lệnh này sẽ trở thành admin
	if !isAdmin {
		// Kiểm tra xem đã có admin nào chưa
		subscribers, err := db.GetAllActiveSubscribers()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lỗi khi kiểm tra danh sách admin: %v", err))
			return
		}

		hasAdmin := false
		for _, sub := range subscribers {
			if sub.IsAdmin {
				hasAdmin = true
				break
			}
		}

		if hasAdmin {
			s.ChannelMessageSend(m.ChannelID, "⛔ Bạn không có quyền quản trị để thực hiện lệnh này!")
			return
		}

		// Nếu chưa có admin, đặt người này làm admin
		err = db.SetAdmin(authorID, true)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lỗi khi đặt quyền admin: %v", err))
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ <@%s>, bạn đã trở thành admin đầu tiên của hệ thống!", m.Author.ID))
		return
	}

	// Phân tích lệnh
	parts := strings.SplitN(m.Content, " ", 3)
	if len(parts) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Sử dụng: !setadmin @user [true/false]")
		return
	}

	// Lấy ID người dùng từ mention
	mentions := m.Mentions
	if len(mentions) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Vui lòng đề cập đến người dùng bằng @username")
		return
	}

	targetUser := mentions[0]
	targetID, _ := strconv.ParseInt(targetUser.ID, 10, 64)

	// Xác định giá trị admin
	adminValue := strings.TrimSpace(parts[2])
	if adminValue != "true" && adminValue != "false" {
		s.ChannelMessageSend(m.ChannelID, "Giá trị phải là 'true' hoặc 'false'")
		return
	}

	setAdmin := adminValue == "true"

	// Kiểm tra xem người dùng đã đăng ký chưa
	isSubscribed, err := db.IsSubscribed(targetID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lỗi khi kiểm tra: %v", err))
		return
	}

	if !isSubscribed {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> chưa đăng ký nhận thông báo! Họ cần đăng ký trước bằng lệnh !subscribe", targetUser.ID))
		return
	}

	// Đặt quyền admin
	err = db.SetAdmin(targetID, setAdmin)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lỗi khi đặt quyền admin: %v", err))
		return
	}

	if setAdmin {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ <@%s> đã được đặt làm admin!", targetUser.ID))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ <@%s> đã bị hủy quyền admin!", targetUser.ID))
	}
}

// SendNotificationToAll gửi thông báo đến tất cả người đăng ký
func SendNotificationToAll(s *discordgo.Session, m *discordgo.MessageCreate, message string) error {
	// Kiểm tra quyền admin
	userID, _ := strconv.ParseInt(m.Author.ID, 10, 64)
	isAdmin, err := db.IsAdmin(userID)
	if err != nil {
		return fmt.Errorf("lỗi khi kiểm tra quyền admin: %v", err)
	}

	if !isAdmin {
		s.ChannelMessageSend(m.ChannelID, "⛔ Bạn không có quyền gửi thông báo đến tất cả người dùng! Chỉ admin mới có thể thực hiện chức năng này.")
		return nil
	}

	// Lấy tất cả người đăng ký đang hoạt động
	subscribers, err := db.GetAllActiveSubscribers()
	if err != nil {
		return fmt.Errorf("lỗi khi lấy danh sách người đăng ký: %v", err)
	}

	// Không có người đăng ký, không cần gửi thông báo
	if len(subscribers) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Không có người dùng nào đăng ký nhận thông báo!")
		return nil
	}

	// Tạo tin nhắn với đề cập đến tất cả người đăng ký
	notificationMsg := "🔔 **Thông báo:** " + message + "\n"
	for _, sub := range subscribers {
		notificationMsg += fmt.Sprintf("<@%d> ", sub.DiscordUserID)
	}

	// Gửi tin nhắn thông báo
	_, err = s.ChannelMessageSend(m.ChannelID, notificationMsg)
	if err != nil {
		return fmt.Errorf("lỗi khi gửi thông báo: %v", err)
	}

	return nil
}

// SendDirectMessageToAll gửi tin nhắn trực tiếp đến tất cả người đăng ký
func SendDirectMessageToAll(s *discordgo.Session, m *discordgo.MessageCreate, message string) error {
	// Kiểm tra quyền admin
	userID, _ := strconv.ParseInt(m.Author.ID, 10, 64)
	isAdmin, err := db.IsAdmin(userID)
	if err != nil {
		return fmt.Errorf("lỗi khi kiểm tra quyền admin: %v", err)
	}

	if !isAdmin {
		s.ChannelMessageSend(m.ChannelID, "⛔ Bạn không có quyền gửi tin nhắn trực tiếp đến tất cả người dùng! Chỉ admin mới có thể thực hiện chức năng này.")
		return nil
	}

	// Lấy tất cả người đăng ký đang hoạt động
	subscribers, err := db.GetAllActiveSubscribers()
	if err != nil {
		return fmt.Errorf("lỗi khi lấy danh sách người đăng ký: %v", err)
	}

	// Không có người đăng ký, không cần gửi thông báo
	if len(subscribers) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Không có người dùng nào đăng ký nhận thông báo!")
		return nil
	}

	// Thông báo bắt đầu gửi
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("🔔 Đang gửi tin nhắn trực tiếp đến %d người dùng...", len(subscribers)))

	// Đếm số tin nhắn gửi thành công
	successCount := 0

	// Gửi tin nhắn trực tiếp đến từng người đăng ký
	for _, sub := range subscribers {
		// Tạo kênh DM với người dùng
		channel, err := s.UserChannelCreate(strconv.FormatInt(sub.DiscordUserID, 10))
		if err != nil {
			fmt.Printf("Lỗi khi tạo kênh DM đến %d: %v\n", sub.DiscordUserID, err)
			continue
		}

		// Gửi tin nhắn đến kênh DM
		_, err = s.ChannelMessageSend(channel.ID, fmt.Sprintf("🔔 **Thông báo từ bot:** %s", message))
		if err != nil {
			fmt.Printf("Lỗi khi gửi DM đến %d: %v\n", sub.DiscordUserID, err)
			continue
		}
		successCount++
	}

	// Thông báo kết quả
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Đã gửi tin nhắn trực tiếp thành công đến %d/%d người dùng.", successCount, len(subscribers)))

	return nil
}
