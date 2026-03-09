package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// ConnectWithRetry attempts to connect to the database with exponential backoff
//
// 🔑 CONCEPT: Exponential Backoff
// - Lần 1 fail → đợi 1s  (1 << 0 = 1)
// - Lần 2 fail → đợi 2s  (1 << 1 = 2)
// - Lần 3 fail → đợi 4s  (1 << 2 = 4)
// - Lần 4 fail → đợi 8s  (1 << 3 = 8)
// - Lần 5 fail → đợi 16s (1 << 4 = 16)
//
// Tại sao? Nếu DB vừa restart, cần thời gian để sẵn sàng.
// Đợi lâu hơn sau mỗi lần thất bại để tránh "storm" kết nối.
func ConnectWithRetry(connStr string, maxRetries int) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("🔄 Database connection attempt %d/%d...", attempt, maxRetries)

		// Mở connection
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("⚠️  Failed to open database: %v", err)
		} else {
			// Kiểm tra kết nối thực sự hoạt động
			err = db.Ping()
			if err == nil {
				log.Println("✅ Database connected successfully!")
				return db, nil
			}
			// Đóng connection nếu ping fail
			db.Close()
			log.Printf("⚠️  Connection failed: %v", err)
		}

		// Nếu chưa phải lần cuối → đợi rồi thử lại
		if attempt < maxRetries {
			// Exponential backoff: 1s, 2s, 4s, 8s, 16s
			waitTime := time.Duration(1<<uint(attempt-1)) * time.Second
			log.Printf("   Retrying in %v...", waitTime)
			time.Sleep(waitTime)
		}
	}

	return nil, fmt.Errorf("❌ failed to connect to database after %d attempts: %w", maxRetries, err)
}
