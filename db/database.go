package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	RedisClient *redis.Client
	PgPool      *pgxpool.Pool
	ctx         = context.TODO()
)

// InitDB khởi tạo kết nối đến Redis và PostgreSQL
func InitDB() error {
	// Khởi tạo kết nối Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Địa chỉ Redis mặc định
		Password: "",               // Không có mật khẩu
		DB:       0,                // DB mặc định
	})

	// Kiểm tra kết nối Redis
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("không thể kết nối đến Redis: %v", err)
	}

	// Khởi tạo kết nối PostgreSQL
	connStr := "postgres://{username}:{password}@localhost:{port}/{databaseName}?sslmode=disable"
	PgPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("không thể kết nối đến PostgreSQL: %v", err)
	}

	fmt.Println("Kết nối cơ sở dữ liệu đã được khởi tạo thành công")
	return nil
}

// CloseDB đóng kết nối cơ sở dữ liệu
func CloseDB() {
	if RedisClient != nil {
		RedisClient.Close()
	}
	if PgPool != nil {
		PgPool.Close()
	}
}
