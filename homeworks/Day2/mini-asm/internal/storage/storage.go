package storage

import "mini-asm/internal/model"

// Storage defines the interface for data access operations
// This is the contract that any storage implementation must fulfill
//
// Why interface?
// - Allows multiple implementations (memory, postgres, mongodb, etc.)
// - Makes testing easy (can mock)
// - Follows Dependency Inversion Principle
type Storage interface {
	// Create adds a new asset to storage
	Create(asset *model.Asset) error

	// GetAll retrieves all assets
	// Returns slice of pointers for efficiency (no copying)
	GetAll() ([]*model.Asset, error)

	// GetByID retrieves a single asset by its ID
	// Returns ErrNotFound if asset doesn't exist
	GetByID(id string) (*model.Asset, error)

	// Update modifies an existing asset
	// Returns ErrNotFound if asset doesn't exist
	Update(id string, asset *model.Asset) error

	// Delete removes an asset from storage
	// Returns ErrNotFound if asset doesn't exist
	Delete(id string) error

	// Filter returns assets matching the given criteria
	// Empty string parameters are ignored (match all)
	Filter(assetType, status string) ([]*model.Asset, error)

	// Search finds assets by partial name match
	Search(query string) ([]*model.Asset, error)

	// === HOMEWORK: New methods ===

	// GetStats returns statistics about all assets (Bài 1)
	GetStats() (*model.AssetStats, error)

	// CountByFilter counts assets matching the given criteria (Bài 1)
	CountByFilter(assetType, status string) (int, error)

	// BatchCreate inserts multiple assets in a single transaction (Bài 2)
	// All or nothing: if one fails, all are rolled back
	BatchCreate(assets []*model.Asset) error

	// BatchDelete removes multiple assets by IDs (Bài 3)
	// Returns count of deleted and not found
	BatchDelete(ids []string) (deleted int, notFound int, err error)

	// ListWithPagination returns paginated assets with filters (Bài 6)
	// Returns: assets slice, total count, error
	ListWithPagination(page, limit int, assetType, status string) ([]*model.Asset, int, error)
}

/*
🎓 NOTES:

1. Interface Design:
   - Define behavior, not implementation
   - Methods should be atomic and clear
   - Return errors, don't panic

2. Why Pointers?
   - []*model.Asset vs []model.Asset
   - Pointers avoid copying large structs
   - Allows modification through reference
   - More memory efficient for large datasets

3. Error Handling:
   - Return error as last value
   - Use model.ErrNotFound for consistency
   - Caller decides how to handle

4. Method Signatures:
   - Create(asset *model.Asset) - pointer: will be modified (ID, timestamps)
   - GetByID(id string) - string: immutable lookup
   - Filter/Search - flexible parameters

5. Interface Benefits:

   Buổi 2: MemoryStorage implements this
   type MemoryStorage struct { ... }
   func (m *MemoryStorage) Create(asset *model.Asset) error { ... }

   Buổi 3: PostgresStorage implements the SAME interface
   type PostgresStorage struct { ... }
   func (p *PostgresStorage) Create(asset *model.Asset) error { ... }

   Service layer doesn't change!
   type AssetService struct {
       storage Storage  // Works with ANY implementation!
   }

6. Testing Benefits:
   type MockStorage struct { ... }
   func (m *MockStorage) Create(asset *model.Asset) error {
       return nil // or test-specific behavior
   }

📝 COMPARISON:

Without Interface (BAD):
    type AssetService struct {
        storage *MemoryStorage  // Coupled to specific implementation
    }
    // Can't swap to database without changing service!

With Interface (GOOD):
    type AssetService struct {
        storage Storage  // Any implementation works
    }
    // Easy to swap implementations!
*/
