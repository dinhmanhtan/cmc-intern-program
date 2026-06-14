package handler

import (
	"encoding/json"
	"mini-asm/internal/service"
	"net/http"
	"time"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	startTime time.Time
	service   *service.AssetService
}

// NewHealthHandler creates a new health check handler
// Accepts an AssetService to retrieve storage information
func NewHealthHandler(svc *service.AssetService) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		service:   svc,
	}
}

// storageInfo represents storage details in the health response
type storageInfo struct {
	Type       string `json:"type"`
	AssetCount int    `json:"asset_count"`
}

// healthResponse represents the health check response
type healthResponse struct {
	Status        string      `json:"status"`
	Storage       storageInfo `json:"storage"`
	UptimeSeconds float64     `json:"uptime_seconds"`
	Timestamp     string      `json:"timestamp"`
}

// Check handles GET /health
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	storageType := "in-memory"
	assetCount := 0

	if h.service != nil {
		stype, count, err := h.service.GetStorageInfo()
		if err == nil {
			storageType = stype
			assetCount = count
		}
	}

	response := healthResponse{
		Status: "ok",
		Storage: storageInfo{
			Type:       storageType,
			AssetCount: assetCount,
		},
		UptimeSeconds: time.Since(h.startTime).Seconds(),
		Timestamp:     time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

/*
🎓 NOTES:

Refactored từ Buổi 1:
- Buổi 1: Health check logic trong main.go
- Buổi 2: Extracted to separate handler
- Homework Day 1: Added storage info (type + asset_count)

Benefits:
- Consistent with other handlers
- Can add more health checks (database, etc.) in Buổi 3
- Reusable and testable
*/
