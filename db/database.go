package db

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	RedisClient *redis.Client
	PgPool      *pgxpool.Pool
	ctx         = context.TODO()
)

// getEnv lấy giá trị từ biến môi trường hoặc trả về giá trị mặc định
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// InitDB khởi tạo kết nối đến Redis và PostgreSQL
func InitDB() error {
	// Lấy thông tin kết nối Redis từ biến môi trường
	redisHost := getEnv("REDIS_HOST", "new_redis")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// Khởi tạo kết nối Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	// Kiểm tra kết nối Redis
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("không thể kết nối đến Redis (%s): %v", redisAddr, err)
	}

	// Lấy thông tin kết nối PostgreSQL từ biến môi trường
	pgHost := getEnv("POSTGRES_HOST", "my_postgres")
	pgPort := getEnv("POSTGRES_PORT", "5432")
	pgUser := getEnv("POSTGRES_USER", "admin")
	pgPassword := getEnv("POSTGRES_PASSWORD", "123456")
	pgDB := getEnv("POSTGRES_DB", "Subscriber")

	// Tạo connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgUser, pgPassword, pgHost, pgPort, pgDB)

	// Khởi tạo kết nối PostgreSQL
	PgPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("không thể kết nối đến PostgreSQL (%s): %v", connStr, err)
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
