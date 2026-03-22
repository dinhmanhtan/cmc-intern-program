package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mini-asm/internal/service"
	"mini-asm/internal/storage"
)

// ExportHandler handles report export requests
//
// Bài 6 - Day 3 Bonus: Export Reports (CSV/JSON)
//
// Endpoints:
//
//	GET /assets/export?format=csv  → Download all assets as CSV file
//	GET /assets/export?format=json → Download all assets as JSON file
//	GET /assets/{id}/results/export → Export all scan results for an asset (JSON)
type ExportHandler struct {
	assetService *service.AssetService
	scanService  *service.ScanService
}

// NewExportHandler creates a new export handler
func NewExportHandler(assetService *service.AssetService, scanService *service.ScanService) *ExportHandler {
	return &ExportHandler{
		assetService: assetService,
		scanService:  scanService,
	}
}

// ExportAssets exports all assets as CSV or JSON (file download)
// GET /assets/export?format=csv|json
func (h *ExportHandler) ExportAssets(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json" // default
	}

	// Fetch all assets (large page size = export everything)
	result, err := h.assetService.ListAssets(storage.QueryParams{
		Page:     1,
		PageSize: 10000,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
		return
	}

	switch format {
	case "csv":
		// CSV export - good for Excel/Sheets
		filename := fmt.Sprintf("assets-export-%s.csv", time.Now().Format("2006-01-02"))
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// Header row
		writer.Write([]string{"id", "name", "type", "status", "created_at", "updated_at"})

		// Data rows
		for _, asset := range result.Data {
			writer.Write([]string{
				asset.ID,
				asset.Name,
				string(asset.Type),
				string(asset.Status),
				asset.CreatedAt.Format(time.RFC3339),
				asset.UpdatedAt.Format(time.RFC3339),
			})
		}

	case "json":
		// JSON export - full data with metadata
		filename := fmt.Sprintf("assets-export-%s.json", time.Now().Format("2006-01-02"))
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)

		type exportResponse struct {
			ExportedAt string                   `json:"exported_at"`
			Total      int64                    `json:"total"`
			Data       interface{}              `json:"data"`
		}

		resp := exportResponse{
			ExportedAt: time.Now().Format(time.RFC3339),
			Total:      result.Total,
			Data:       result.Data,
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(resp)

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"invalid format, use 'csv' or 'json'"}`)
	}
}

// ExportAssetResults exports all scan results for a specific asset as JSON
// GET /assets/{id}/results/export
func (h *ExportHandler) ExportAssetResults(w http.ResponseWriter, r *http.Request) {
	assetID := r.PathValue("id")
	if assetID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"asset ID required"}`)
		return
	}

	// Get combined scan results (DNS + WHOIS + Subdomains)
	results, err := h.scanService.GetAssetAllScanResults(assetID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
		return
	}

	// Get scan job history (non-fatal if fails)
	jobs, _ := h.scanService.ListScanJobs(assetID)

	// Build export payload
	type exportPayload struct {
		ExportedAt string      `json:"exported_at"`
		AssetID    string      `json:"asset_id"`
		ScanJobs   interface{} `json:"scan_jobs"`
		Results    interface{} `json:"results"`
	}

	payload := exportPayload{
		ExportedAt: time.Now().Format(time.RFC3339),
		AssetID:    assetID,
		ScanJobs:   jobs,
		Results:    results,
	}

	// Short ID for filename
	shortID := assetID
	if len(assetID) > 8 {
		shortID = assetID[:8]
	}

	filename := fmt.Sprintf("asset-%s-results-%s.json", shortID, time.Now().Format("2006-01-02"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(payload)
}
