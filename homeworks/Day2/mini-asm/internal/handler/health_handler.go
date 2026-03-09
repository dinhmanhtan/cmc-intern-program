package handler

import (
	"database/sql"
	"net/http"
	"time"
)

// HealthHandler handles health check requests
// Updated for Bài 5: now includes database health information
type HealthHandler struct {
	startTime time.Time
	db        *sql.DB // Database connection for health check
}

// NewHealthHandler creates a new health check handler
// Bài 5: Now accepts *sql.DB to check database status
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		db:        db,
	}
}

// DatabaseHealth represents database health information
type DatabaseHealth struct {
	Status          string `json:"status"`
	OpenConnections int    `json:"open_connections"`
	InUse           int    `json:"in_use"`
	Idle            int    `json:"idle"`
	MaxOpen         int    `json:"max_open"`
}

// HealthResponse represents the full health check response
type HealthResponse struct {
	Status    string         `json:"status"`
	Database  DatabaseHealth `json:"database"`
	Timestamp time.Time      `json:"timestamp"`
}

// Check handles GET /health
// Bài 5: Returns 200 if DB is connected, 503 if DB is down
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Timestamp: time.Now(),
	}

	// Check database connection using db.Ping()
	dbStatus := "connected"
	httpStatus := http.StatusOK
	appStatus := "ok"

	if err := h.db.Ping(); err != nil {
		dbStatus = "disconnected"
		httpStatus = http.StatusServiceUnavailable // 503
		appStatus = "degraded"
	}

	// Get connection pool stats using db.Stats()
	stats := h.db.Stats()

	response.Status = appStatus
	response.Database = DatabaseHealth{
		Status:          dbStatus,
		OpenConnections: stats.OpenConnections,
		InUse:           stats.InUse,
		Idle:            stats.Idle,
		MaxOpen:         stats.MaxOpenConnections,
	}

	RespondJSON(w, httpStatus, response)
}
