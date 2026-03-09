package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"mini-asm/internal/database"
	"mini-asm/internal/handler"
	"mini-asm/internal/service"
	"mini-asm/internal/storage/postgres"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	log.Println("🚀 Starting Mini ASM Server (Homework Edition)...")

	// ============================================
	// CONFIGURATION - Load from environment
	// ============================================

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "mini_asm")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	log.Printf("📊 Connecting to database: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)

	// ============================================
	// DATABASE CONNECTION (Bài 4: Retry Logic)
	// ============================================

	// Bài 4: Sử dụng ConnectWithRetry thay vì sql.Open + db.Ping
	// Tự động retry 5 lần với exponential backoff
	db, err := database.ConnectWithRetry(connStr, 5)
	if err != nil {
		log.Fatal("❌ ", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * 60 * 1000)

	// ============================================
	// DEPENDENCY INJECTION - Wire up all layers
	// ============================================

	// 1. Storage Layer
	store := postgres.NewPostgresStorage(db)
	log.Println("✅ Storage initialized: PostgreSQL")

	// 2. Service Layer
	assetService := service.NewAssetService(store)
	log.Println("✅ Service initialized: AssetService")

	// 3. Handler Layer
	assetHandler := handler.NewAssetHandler(assetService)
	healthHandler := handler.NewHealthHandler(db) // Bài 5: pass db for health check
	log.Println("✅ Handlers initialized")

	// ============================================
	// ROUTING - Register HTTP endpoints
	// ============================================

	mux := http.NewServeMux()

	// Health check (Bài 5: includes DB status)
	mux.HandleFunc("GET /health", healthHandler.Check)

	// Original CRUD operations
	mux.HandleFunc("POST /assets", assetHandler.CreateAsset)
	mux.HandleFunc("GET /assets", assetHandler.ListAssets) // Bài 6: now supports ?page=&limit=
	mux.HandleFunc("GET /assets/{id}", assetHandler.GetAsset)
	mux.HandleFunc("PUT /assets/{id}", assetHandler.UpdateAsset)
	mux.HandleFunc("DELETE /assets/{id}", assetHandler.DeleteAsset)

	// === HOMEWORK: New routes ===
	mux.HandleFunc("GET /assets/stats", assetHandler.GetStats)             // Bài 1
	mux.HandleFunc("GET /assets/count", assetHandler.CountAssets)          // Bài 1
	mux.HandleFunc("POST /assets/batch", assetHandler.BatchCreateAssets)   // Bài 2
	mux.HandleFunc("DELETE /assets/batch", assetHandler.BatchDeleteAssets) // Bài 3
	mux.HandleFunc("GET /assets/search", assetHandler.SearchAssets)        // Bài 7

	log.Println("✅ Routes registered:")
	log.Println("   GET    /health")
	log.Println("   POST   /assets")
	log.Println("   GET    /assets           (supports ?page=&limit=&type=&status=)")
	log.Println("   GET    /assets/{id}")
	log.Println("   PUT    /assets/{id}")
	log.Println("   DELETE /assets/{id}")
	log.Println("   GET    /assets/stats     [Bài 1]")
	log.Println("   GET    /assets/count     [Bài 1]")
	log.Println("   POST   /assets/batch     [Bài 2]")
	log.Println("   DELETE /assets/batch     [Bài 3]")
	log.Println("   GET    /assets/search    [Bài 7]")

	// ============================================
	// START SERVER
	// ============================================

	port := getEnv("SERVER_PORT", "8080")
	addr := ":" + port

	log.Printf("🌐 Server listening on http://localhost%s\n", addr)
	log.Println("Press Ctrl+C to stop")
	log.Println()

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("❌ Server failed to start:", err)
	}
}

// getEnv retrieves an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
