package model

import (
	"errors"
	"time"
)

// Asset represents a public-facing resource (domain, IP, or service)
// This is our core domain entity - no dependencies on other layers
type Asset struct {
	ID        string    `json:"id"`         // UUID
	Name      string    `json:"name"`       // e.g., "example.com", "192.168.1.1"
	Type      string    `json:"type"`       // domain, ip, or service
	Status    string    `json:"status"`     // active or inactive
	CreatedAt time.Time `json:"created_at"` // Auto-set on creation
	UpdatedAt time.Time `json:"updated_at"` // Auto-updated
}

// Asset types - using constants for type safety
const (
	TypeDomain  = "domain"
	TypeIP      = "ip"
	TypeService = "service"
)

// Asset statuses
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

// IsValidType checks if the given type is valid
func IsValidType(t string) bool {
	return t == TypeDomain || t == TypeIP || t == TypeService
}

// IsValidStatus checks if the given status is valid
func IsValidStatus(s string) bool {
	return s == StatusActive || s == StatusInactive
}

// Statistics represents aggregated stats about assets
type Statistics struct {
	Total    int            `json:"total"`
	ByType   map[string]int `json:"by_type"`
	ByStatus map[string]int `json:"by_status"`
}

// CountResponse represents the response for count endpoint
type CountResponse struct {
	Count   int               `json:"count"`
	Filters map[string]string `json:"filters"`
}

// PaginatedResponse represents a paginated list response
type PaginatedResponse struct {
	Data       []*Asset   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Pagination holds pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// StorageInfo holds storage information for health check
type StorageInfo struct {
	Type       string `json:"type"`
	AssetCount int    `json:"asset_count"`
}

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status          string      `json:"status"`
	Storage         StorageInfo `json:"storage"`
	UptimeSeconds   float64     `json:"uptime_seconds"`
	Timestamp       string      `json:"timestamp"`
}

// BatchCreateRequest represents the request body for batch creating assets
type BatchCreateRequest struct {
	Assets []BatchCreateItem `json:"assets"`
}

// BatchCreateItem represents a single item in a batch create request
type BatchCreateItem struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// BatchCreateResponse represents the response for batch create
type BatchCreateResponse struct {
	Created int      `json:"created"`
	IDs     []string `json:"ids"`
}

// BatchDeleteResponse represents the response for batch delete
type BatchDeleteResponse struct {
	Deleted  int `json:"deleted"`
	NotFound int `json:"not_found"`
}

// ErrBatchLimitExceeded indicates too many items in batch request
var ErrBatchLimitExceeded = errors.New("batch limit exceeded: maximum 100 items per request")

/*
🎓 NOTES:

1. Pure Domain Entity:
   - No database tags (gorm, sql, etc.)
   - No HTTP concerns
   - Just business concepts
   - This is the "Entity Layer" in Clean Architecture

2. Struct Tags:
   - `json:"id"` - định nghĩa tên field trong JSON response
   - Quan trọng: nếu không có tag, Go sẽ export field name as-is
   - Example: ID → "ID" vs id → "id"

3. Constants vs Strings:
   - ✅ TypeDomain - compiler checked, typo-safe
   - ❌ "domain" - runtime error if typo
   - Best practice: use constants!

4. Helper Functions:
   - IsValidType(), IsValidStatus()
   - Used by service layer for validation
   - Keep validation logic reusable

5. Why time.Time?
   - Built-in JSON marshalling
   - Timezone aware
   - Easy comparison and manipulation


*/
