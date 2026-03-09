package postgres

import (
	"database/sql"
	"fmt"

	"mini-asm/internal/config"
	"mini-asm/internal/model"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresStorage implements the Storage interface using PostgreSQL
// This is a concrete implementation that can be swapped with MemoryStorage
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage creates a new PostgreSQL storage instance
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func NewPostgresStorageFromConfig(config *config.PostgresConfig) (*PostgresStorage, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &PostgresStorage{db: db}, nil
}

// Create inserts a new asset into the database
func (p *PostgresStorage) Create(asset *model.Asset) error {
	query := `
		INSERT INTO assets (id, name, type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := p.db.Exec(
		query,
		asset.ID,
		asset.Name,
		asset.Type,
		asset.Status,
		asset.CreatedAt,
		asset.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create asset: %w", err)
	}

	return nil
}

// GetAll retrieves all assets from the database
func (p *PostgresStorage) GetAll() ([]*model.Asset, error) {
	query := `
		SELECT id, name, type, status, created_at, updated_at
		FROM assets
		ORDER BY created_at DESC
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	var assets []*model.Asset
	for rows.Next() {
		asset := &model.Asset{}
		err := rows.Scan(
			&asset.ID,
			&asset.Name,
			&asset.Type,
			&asset.Status,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}
		assets = append(assets, asset)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return assets, nil
}

// GetByID retrieves a single asset by its ID
func (p *PostgresStorage) GetByID(id string) (*model.Asset, error) {
	query := `
		SELECT id, name, type, status, created_at, updated_at
		FROM assets
		WHERE id = $1
	`

	asset := &model.Asset{}
	err := p.db.QueryRow(query, id).Scan(
		&asset.ID,
		&asset.Name,
		&asset.Type,
		&asset.Status,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	return asset, nil
}

// Update modifies an existing asset in the database
func (p *PostgresStorage) Update(id string, asset *model.Asset) error {
	query := `
		UPDATE assets
		SET name = $1, type = $2, status = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := p.db.Exec(
		query,
		asset.Name,
		asset.Type,
		asset.Status,
		asset.UpdatedAt,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update asset: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return model.ErrNotFound
	}

	return nil
}

// Delete removes an asset from the database
func (p *PostgresStorage) Delete(id string) error {
	query := `DELETE FROM assets WHERE id = $1`

	result, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return model.ErrNotFound
	}

	return nil
}

// Filter returns assets matching the given criteria
func (p *PostgresStorage) Filter(assetType, status string) ([]*model.Asset, error) {
	query := `
		SELECT id, name, type, status, created_at, updated_at
		FROM assets
		WHERE 1=1
	`

	var args []interface{}
	argCount := 1

	if assetType != "" {
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, assetType)
		argCount++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	query += " ORDER BY created_at DESC"

	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to filter assets: %w", err)
	}
	defer rows.Close()

	var assets []*model.Asset
	for rows.Next() {
		asset := &model.Asset{}
		err := rows.Scan(
			&asset.ID,
			&asset.Name,
			&asset.Type,
			&asset.Status,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}
		assets = append(assets, asset)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return assets, nil
}

// Search finds assets by partial name match
func (p *PostgresStorage) Search(query string) ([]*model.Asset, error) {
	sqlQuery := `
		SELECT id, name, type, status, created_at, updated_at
		FROM assets
		WHERE name ILIKE $1
		ORDER BY created_at DESC
	`

	// Add wildcards for partial matching
	searchPattern := "%" + query + "%"

	rows, err := p.db.Query(sqlQuery, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search assets: %w", err)
	}
	defer rows.Close()

	var assets []*model.Asset
	for rows.Next() {
		asset := &model.Asset{}
		err := rows.Scan(
			&asset.ID,
			&asset.Name,
			&asset.Type,
			&asset.Status,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}
		assets = append(assets, asset)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return assets, nil
}

// === HOMEWORK: New methods ===

// GetStats returns statistics about all assets (Bài 1)
// Uses COUNT + GROUP BY to get totals by type and status in one pass
func (p *PostgresStorage) GetStats() (*model.AssetStats, error) {
	stats := &model.AssetStats{
		ByType:   make(map[string]int),
		ByStatus: make(map[string]int),
	}

	// Count total
	err := p.db.QueryRow("SELECT COUNT(*) FROM assets").Scan(&stats.Total)
	if err != nil {
		return nil, fmt.Errorf("failed to count assets: %w", err)
	}

	// Count by type using GROUP BY
	rows, err := p.db.Query("SELECT type, COUNT(*) FROM assets GROUP BY type")
	if err != nil {
		return nil, fmt.Errorf("failed to count by type: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var assetType string
		var count int
		if err := rows.Scan(&assetType, &count); err != nil {
			return nil, fmt.Errorf("failed to scan type count: %w", err)
		}
		stats.ByType[assetType] = count
	}

	// Count by status using GROUP BY
	rows2, err := p.db.Query("SELECT status, COUNT(*) FROM assets GROUP BY status")
	if err != nil {
		return nil, fmt.Errorf("failed to count by status: %w", err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var status string
		var count int
		if err := rows2.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan status count: %w", err)
		}
		stats.ByStatus[status] = count
	}

	return stats, nil
}

// CountByFilter counts assets matching the given criteria (Bài 1)
// Builds dynamic WHERE clause based on which filters are provided
func (p *PostgresStorage) CountByFilter(assetType, status string) (int, error) {
	query := "SELECT COUNT(*) FROM assets WHERE 1=1"
	var args []interface{}
	argCount := 1

	if assetType != "" {
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, assetType)
		argCount++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	var count int
	err := p.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count assets: %w", err)
	}

	return count, nil
}

// BatchCreate inserts multiple assets in a single transaction (Bài 2)
// 🔑 KEY CONCEPT: Transaction = all-or-nothing
//   - db.Begin() starts a transaction
//   - tx.Exec() runs queries within the transaction
//   - tx.Commit() saves all changes
//   - tx.Rollback() undoes everything if any error occurs
func (p *PostgresStorage) BatchCreate(assets []*model.Asset) error {
	// Start transaction
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// defer Rollback will only execute if Commit hasn't been called
	defer tx.Rollback()

	query := `
		INSERT INTO assets (id, name, type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	// Insert each asset within the same transaction
	for _, asset := range assets {
		_, err := tx.Exec(query,
			asset.ID,
			asset.Name,
			asset.Type,
			asset.Status,
			asset.CreatedAt,
			asset.UpdatedAt,
		)
		if err != nil {
			// Any error → automatic rollback via defer
			return fmt.Errorf("failed to create asset '%s': %w", asset.Name, err)
		}
	}

	// If we reach here, all inserts succeeded → commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// BatchDelete removes multiple assets by IDs (Bài 3)
// Deletes all valid IDs, counts how many were not found
// Invalid UUID format IDs are counted as not_found (not error)
func (p *PostgresStorage) BatchDelete(ids []string) (deleted int, notFound int, err error) {
	deleted = 0
	notFound = 0

	for _, id := range ids {
		// Validate UUID format first to avoid PostgreSQL error
		if _, parseErr := parseUUID(id); parseErr != nil {
			notFound++ // Invalid UUID → treat as not found
			continue
		}

		result, err := p.db.Exec("DELETE FROM assets WHERE id = $1", id)
		if err != nil {
			return deleted, notFound, fmt.Errorf("failed to delete asset %s: %w", id, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return deleted, notFound, fmt.Errorf("failed to get affected rows: %w", err)
		}

		if rowsAffected > 0 {
			deleted++
		} else {
			notFound++
		}
	}

	return deleted, notFound, nil
}

// parseUUID validates if a string is a valid UUID format
func parseUUID(s string) (string, error) {
	// Simple UUID validation: must be 36 characters with hyphens in correct positions
	if len(s) != 36 {
		return "", fmt.Errorf("invalid UUID length")
	}
	for i, c := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return "", fmt.Errorf("invalid UUID format")
			}
		} else {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return "", fmt.Errorf("invalid UUID character")
			}
		}
	}
	return s, nil
}

// ListWithPagination returns paginated assets with optional filters (Bài 6)
// Uses LIMIT/OFFSET for pagination and a separate COUNT query for total
func (p *PostgresStorage) ListWithPagination(page, limit int, assetType, status string) ([]*model.Asset, int, error) {
	// Build WHERE clause (same pattern as Filter)
	whereClause := "WHERE 1=1"
	var args []interface{}
	argCount := 1

	if assetType != "" {
		whereClause += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, assetType)
		argCount++
	}

	if status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	// Step 1: COUNT total matching records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM assets %s", whereClause)
	var total int
	err := p.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count assets: %w", err)
	}

	// Step 2: SELECT with LIMIT and OFFSET
	offset := (page - 1) * limit
	dataQuery := fmt.Sprintf(
		"SELECT id, name, type, status, created_at, updated_at FROM assets %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		whereClause, argCount, argCount+1,
	)
	args = append(args, limit, offset)

	rows, err := p.db.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	var assets []*model.Asset
	for rows.Next() {
		asset := &model.Asset{}
		err := rows.Scan(
			&asset.ID,
			&asset.Name,
			&asset.Type,
			&asset.Status,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan asset: %w", err)
		}
		assets = append(assets, asset)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return assets, total, nil
}

/*
🎓 TEACHING NOTES:

1. Interface Implementation:
   - PostgresStorage implements Storage interface
   - Same methods as MemoryStorage
   - Main.go CHỈ CẦN THAY 1 DÒNG!

2. SQL Queries:
   - Parameterized queries ($1, $2) prevent SQL injection
   - NEVER concatenate user input into SQL!
   - ❌ BAD: query := "SELECT * FROM assets WHERE id = '" + id + "'"
   - ✅ GOOD: query := "SELECT * FROM assets WHERE id = $1"

3. Error Handling:
   - Check sql.ErrNoRows → return model.ErrNotFound
   - Wrap errors with fmt.Errorf("context: %w", err)
   - RowsAffected() to verify UPDATE/DELETE success

4. Scanning Rows:
   - rows.Scan() maps columns to struct fields
   - ORDER MATTERS! Must match SELECT order
   - Don't forget defer rows.Close()

5. Dynamic Query Building:
   - Filter() builds query dynamically based on parameters
   - Track parameter count for $1, $2, $3...
   - WHERE 1=1 trick for easier dynamic AND conditions

6. ILIKE vs LIKE:
   - LIKE: case-sensitive
   - ILIKE: case-insensitive (PostgreSQL specific)
   - % wildcards for partial matching

7. Connection Pool:
   - sql.DB maintains connection pool automatically
   - db.Exec(), db.Query(), db.QueryRow() reuse connections
   - Don't close db in storage methods!

8. Transaction Support (Buổi 4):
   - Current: each operation is auto-committed
   - Future: db.Begin() for multi-step operations

COMPARISON: Memory vs Postgres

MemoryStorage:
- data := make(map[string]*model.Asset)
- Fast: O(1) lookups

PostgresStorage:
- data in database on disk
- Slightly slower but persistent
- Can handle millions of records
- Support for advanced queries (JOIN, aggregation)

KEY POINT: Clean Architecture cho phép swap giữa 2 implementations này
mà KHÔNG THAY ĐỔI business logic!
*/
