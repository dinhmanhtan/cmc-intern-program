package repository

import (
	"assets-api/internal/domain"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type AssetRepository struct {
	db *sql.DB
}

func NewAssetRepository(db *sql.DB) *AssetRepository {
	return &AssetRepository{db}
}

// Bài 1
func (r *AssetRepository) GetStats() (*domain.Stats, error) {

	stats := &domain.Stats{
		ByType:   map[string]int{},
		ByStatus: map[string]int{},
	}

	err := r.db.QueryRow("SELECT COUNT(*) FROM assets").Scan(&stats.Total)
	if err != nil {
		return nil, err
	}

	rows, _ := r.db.Query("SELECT type,COUNT(*) FROM assets GROUP BY type")

	for rows.Next() {
		var t string
		var c int
		rows.Scan(&t, &c)
		stats.ByType[t] = c
	}

	rows, _ = r.db.Query("SELECT status,COUNT(*) FROM assets GROUP BY status")

	for rows.Next() {
		var s string
		var c int
		rows.Scan(&s, &c)
		stats.ByStatus[s] = c
	}

	return stats, nil
}

// bài 1
func (r *AssetRepository) Count(t, status string) (int, error) {

	query := "SELECT COUNT(*) FROM assets WHERE 1=1"

	args := []interface{}{}
	i := 1

	if t != "" {
		query += fmt.Sprintf(" AND type=$%d", i)
		args = append(args, t)
		i++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status=$%d", i)
		args = append(args, status)
	}

	var count int

	err := r.db.QueryRow(query, args...).Scan(&count)

	return count, err
}

// bài2
func (r *AssetRepository) BatchCreate(assets []domain.Asset) ([]string, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	ids := []string{}

	for _, a := range assets {

		id := uuid.New().String()

		_, err := tx.Exec(
			"INSERT INTO assets(id,name,type,status) VALUES($1,$2,$3,$4)",
			id, a.Name, a.Type, "active",
		)

		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	err = tx.Commit()

	return ids, err
}

// bài 3
func (r *AssetRepository) BatchDelete(ids []string) (int, int, error) {

	validIDs := []string{}
	notFound := 0

	for _, id := range ids {

		_, err := uuid.Parse(id)

		if err != nil {
			notFound++
			continue
		}

		validIDs = append(validIDs, id)
	}

	deleted := 0

	for _, id := range validIDs {

		res, err := r.db.Exec("DELETE FROM assets WHERE id=$1", id)
		if err != nil {
			return 0, 0, err
		}

		rows, _ := res.RowsAffected()

		if rows == 0 {
			notFound++
		} else {
			deleted++
		}
	}

	return deleted, notFound, nil
}

// bonus
func (r *AssetRepository) List(page, limit int, t, status string) ([]domain.Asset, int, error) {

	offset := (page - 1) * limit

	query := "SELECT id,name,type,status,created_at FROM assets WHERE 1=1"

	args := []interface{}{}
	i := 1

	if t != "" {
		query += fmt.Sprintf(" AND type=$%d", i)
		args = append(args, t)
		i++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status=$%d", i)
		args = append(args, status)
		i++
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}

	var assets []domain.Asset

	for rows.Next() {

		var a domain.Asset

		rows.Scan(&a.ID, &a.Name, &a.Type, &a.Status, &a.CreatedAt)

		assets = append(assets, a)
	}

	total, _ := r.Count(t, status)

	return assets, total, nil
}

// bonus
func (r *AssetRepository) Search(q string) ([]domain.Asset, error) {

	rows, err := r.db.Query(
		"SELECT id,name,type,status,created_at FROM assets WHERE name ILIKE $1 LIMIT 100",
		"%"+q+"%",
	)

	if err != nil {
		return nil, err
	}

	var assets []domain.Asset

	for rows.Next() {

		var a domain.Asset

		rows.Scan(&a.ID, &a.Name, &a.Type, &a.Status, &a.CreatedAt)

		assets = append(assets, a)
	}

	return assets, nil
}
