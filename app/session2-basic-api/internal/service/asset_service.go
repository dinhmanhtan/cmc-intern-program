package service

import (
	"mini-asm/internal/model"
	"mini-asm/internal/storage"
	"time"

	"github.com/google/uuid"
)

// AssetService handles business logic for asset operations
// It sits between handlers (HTTP layer) and storage (data layer)
type AssetService struct {
	storage storage.Storage // Dependency injection - any Storage implementation
}

// NewAssetService creates a new asset service
// Takes Storage interface - can be memory, database, or mock
func NewAssetService(storage storage.Storage) *AssetService {
	return &AssetService{
		storage: storage,
	}
}

// CreateAsset creates a new asset with validation
// Returns the created asset or an error
func (s *AssetService) CreateAsset(name, assetType string) (*model.Asset, error) {
	// Validation - business rules enforcement
	if name == "" {
		return nil, model.ErrEmptyName
	}

	if !model.IsValidType(assetType) {
		return nil, model.ErrInvalidType
	}

	// Business logic - create asset with defaults
	asset := &model.Asset{
		ID:        uuid.New().String(), // Auto-generate UUID
		Name:      name,
		Type:      assetType,
		Status:    model.StatusActive, // Default status
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Delegate to storage layer
	if err := s.storage.Create(asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// GetAllAssets retrieves all assets
func (s *AssetService) GetAllAssets() ([]*model.Asset, error) {
	return s.storage.GetAll()
}

// GetAssetByID retrieves a single asset by ID
func (s *AssetService) GetAssetByID(id string) (*model.Asset, error) {
	if id == "" {
		return nil, model.ErrInvalidInput
	}

	return s.storage.GetByID(id)
}

// UpdateAsset updates an existing asset
// Only updates provided fields (partial update)
func (s *AssetService) UpdateAsset(id string, name, assetType, status string) (*model.Asset, error) {
	// Validate ID
	if id == "" {
		return nil, model.ErrInvalidInput
	}

	// Get existing asset
	existing, err := s.storage.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Apply updates (only if provided)
	if name != "" {
		existing.Name = name
	}

	if assetType != "" {
		if !model.IsValidType(assetType) {
			return nil, model.ErrInvalidType
		}
		existing.Type = assetType
	}

	if status != "" {
		if !model.IsValidStatus(status) {
			return nil, model.ErrInvalidStatus
		}
		existing.Status = status
	}

	// Update timestamp
	existing.UpdatedAt = time.Now()

	// Save to storage
	if err := s.storage.Update(id, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteAsset removes an asset
func (s *AssetService) DeleteAsset(id string) error {
	if id == "" {
		return model.ErrInvalidInput
	}

	return s.storage.Delete(id)
}

// FilterAssets returns assets matching criteria
func (s *AssetService) FilterAssets(assetType, status string) ([]*model.Asset, error) {
	// Validate filters if provided
	if assetType != "" && !model.IsValidType(assetType) {
		return nil, model.ErrInvalidType
	}

	if status != "" && !model.IsValidStatus(status) {
		return nil, model.ErrInvalidStatus
	}

	return s.storage.Filter(assetType, status)
}

// SearchAssets finds assets by name
func (s *AssetService) SearchAssets(query string) ([]*model.Asset, error) {
	if query == "" {
		return nil, model.ErrInvalidInput
	}

	return s.storage.Search(query)
}

// GetStatistics returns aggregated asset statistics
func (s *AssetService) GetStatistics() (*model.Statistics, error) {
	return s.storage.GetStatistics()
}

// CountAssets returns the number of assets matching the given filters
func (s *AssetService) CountAssets(assetType, status string) (int, error) {
	if assetType != "" && !model.IsValidType(assetType) {
		return 0, model.ErrInvalidType
	}
	if status != "" && !model.IsValidStatus(status) {
		return 0, model.ErrInvalidStatus
	}
	return s.storage.Count(assetType, status)
}

// BatchCreateAssets creates multiple assets with all-or-nothing semantics
// Validates ALL items first, then inserts only if all are valid
func (s *AssetService) BatchCreateAssets(items []model.BatchCreateItem) ([]string, error) {
	// Limit check
	if len(items) > 100 {
		return nil, model.ErrBatchLimitExceeded
	}
	if len(items) == 0 {
		return nil, model.ErrInvalidInput
	}

	// Phase 1: Validate ALL assets before inserting anything
	now := time.Now()
	assets := make([]*model.Asset, 0, len(items))
	for _, item := range items {
		if item.Name == "" {
			return nil, model.ErrEmptyName
		}
		if !model.IsValidType(item.Type) {
			return nil, model.ErrInvalidType
		}
		assets = append(assets, &model.Asset{
			ID:        uuid.New().String(),
			Name:      item.Name,
			Type:      item.Type,
			Status:    model.StatusActive,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	// Phase 2: All valid → insert all
	return s.storage.BatchCreate(assets)
}

// BatchDeleteAssets removes multiple assets by their IDs
func (s *AssetService) BatchDeleteAssets(ids []string) (int, int, error) {
	if len(ids) == 0 {
		return 0, 0, model.ErrInvalidInput
	}
	return s.storage.BatchDelete(ids)
}

// ListAssetsPaginated returns assets with pagination and optional filters
func (s *AssetService) ListAssetsPaginated(assetType, status string, page, limit int) (*model.PaginatedResponse, error) {
	// Validate filters
	if assetType != "" && !model.IsValidType(assetType) {
		return nil, model.ErrInvalidType
	}
	if status != "" && !model.IsValidStatus(status) {
		return nil, model.ErrInvalidStatus
	}

	// Get total count for pagination metadata
	total, err := s.storage.Count(assetType, status)
	if err != nil {
		return nil, err
	}

	// Get filtered assets
	allFiltered, err := s.storage.Filter(assetType, status)
	if err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Apply pagination (slice bounds check)
	var pageData []*model.Asset
	if offset >= len(allFiltered) {
		pageData = []*model.Asset{}
	} else {
		end := offset + limit
		if end > len(allFiltered) {
			end = len(allFiltered)
		}
		pageData = allFiltered[offset:end]
	}

	totalPages := (total + limit - 1) / limit

	return &model.PaginatedResponse{
		Data: pageData,
		Pagination: model.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetStorageInfo returns information about the storage layer
func (s *AssetService) GetStorageInfo() (string, int, error) {
	count, err := s.storage.Count("", "")
	if err != nil {
		return "", 0, err
	}
	return "in-memory", count, nil
}

/*
🎓 NOTES:

1. Service Layer Responsibilities:
   ✅ Business logic
   ✅ Validation
   ✅ Orchestration (coordinate multiple operations)
   ✅ Default values
   ❌ HTTP concerns (status codes, JSON)
   ❌ Database details (SQL, queries)

2. Dependency Injection:
   func NewAssetService(storage storage.Storage) *AssetService

   Q: Tại sao không dùng global variable?
   A: Testability! Có thể inject mock storage trong tests

   Example:
   // Production
   service := NewAssetService(memory.NewMemoryStorage())

   // Testing
   service := NewAssetService(&MockStorage{})

3. Validation Strategy:
   - Validate BEFORE business logic
   - Return specific errors (ErrEmptyName, ErrInvalidType)
   - Handler layer maps to HTTP status codes

4. Business Logic Examples:
   - Auto-generate UUID
   - Set default status = active
   - Auto-set timestamps
   - Partial updates (only update provided fields)

5. Error Propagation:
   if err := s.storage.Create(asset); err != nil {
       return nil, err  // Let caller handle
   }

   Q: Tại sao không handle error ở đây?
   A: Service không biết context (HTTP? CLI? gRPC?)
      Handler layer sẽ decide status code

6. Comparison với "Fat Controller":

   ❌ BAD (All in handler):
   func CreateAssetHandler(w, r) {
       // Parse JSON
       // Validate
       // Generate UUID
       // Set defaults
       // Save to DB
       // Return response
   }
   → Hard to test, hard to reuse

   ✅ GOOD (Layered):
   Handler: Parse JSON, call service, return HTTP response
   Service: Validate, business logic
   Storage: Data persistence
   → Easy to test each layer, reusable logic

7. UUID Generation:
   uuid.New().String() → "550e8400-e29b-41d4-a716-446655440000"
   - Globally unique
   - No need for database auto-increment
   - Can generate offline
   - URL-safe

8. Timestamps:
   CreatedAt: time.Now()  // Set once
   UpdatedAt: time.Now()  // Update on every change

   Q: Tại sao không dùng int64 (Unix timestamp)?
   A: time.Time có timezone, human-readable, JSON support

📝 TESTING SERVICE LAYER:

func TestCreateAsset(t *testing.T) {
    // Arrange
    mockStorage := &MockStorage{}
    service := NewAssetService(mockStorage)

    // Act
    asset, err := service.CreateAsset("example.com", "domain")

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, asset.ID)
    assert.Equal(t, "active", asset.Status) // Default
}

❓ QUESTIONS TO ASK:

1. Tại sao CreateAsset return (*Asset, error) thay vì chỉ error?
   → Need to return created asset with ID to client

2. Tại sao UpdateAsset là partial update?
   → Flexible - client chỉ gửi fields muốn update

3. Service layer có nên biết về HTTP không?
   → KHÔNG! Separation of concerns

4. Làm sao test service layer mà không cần database?
   → Mock storage! (Buổi 5)
*/
