package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func HealthHandler(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		status := "ok"
		dbStatus := "connected"

		err := db.Ping()

		if err != nil {
			status = "degraded"
			dbStatus = "disconnected"
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		stats := db.Stats()

		response := map[string]interface{}{
			"status": status,
			"database": map[string]interface{}{
				"status":           dbStatus,
				"open_connections": stats.OpenConnections,
				"in_use":           stats.InUse,
				"idle":             stats.Idle,
				"max_open":         stats.MaxOpenConnections,
			},
			"timestamp": time.Now(),
		}

		json.NewEncoder(w).Encode(response)
	}
}
