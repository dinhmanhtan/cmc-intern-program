package handler

import (
	"assets-api/internal/domain"
	"assets-api/internal/service"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type AssetHandler struct {
	service *service.AssetService
}

func (h *AssetHandler) HealthHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	panic("unimplemented")
}

func NewAssetHandler(s *service.AssetService) *AssetHandler {
	return &AssetHandler{s}
}

// bài1
func (h *AssetHandler) Stats(w http.ResponseWriter, r *http.Request) {

	stats, _ := h.service.GetStats()

	json.NewEncoder(w).Encode(stats)
}

// bài 1
func (h *AssetHandler) Count(w http.ResponseWriter, r *http.Request) {

	t := r.URL.Query().Get("type")
	status := r.URL.Query().Get("status")

	count, err := h.service.Count(t, status)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"count": count,
		"filters": map[string]string{
			"type":   t,
			"status": status,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// bài2
func (h *AssetHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Assets []domain.Asset `json:"assets"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ids, err := h.service.BatchCreate(req.Assets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"created": len(ids),
		"ids":     ids,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// bài 3
func (h *AssetHandler) BatchDelete(w http.ResponseWriter, r *http.Request) {

	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		http.Error(w, "ids parameter required", http.StatusBadRequest)
		return
	}

	ids := strings.Split(idsParam, ",")

	deleted, notFound, err := h.service.BatchDelete(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]int{
		"deleted":   deleted,
		"not_found": notFound,
	}

	json.NewEncoder(w).Encode(response)
}

// bonus
func (h *AssetHandler) List(w http.ResponseWriter, r *http.Request) {

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 20
	}

	typeFilter := r.URL.Query().Get("type")
	statusFilter := r.URL.Query().Get("status")

	data, total, err := h.service.List(page, limit, typeFilter, statusFilter)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	response := map[string]interface{}{
		"data": data,
		"pagination": map[string]int{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// bonus
func (h *AssetHandler) Search(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")

	if query == "" {
		http.Error(w, "q parameter required", http.StatusBadRequest)
		return
	}

	results, err := h.service.Search(query)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(results)
}
