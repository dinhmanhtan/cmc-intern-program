package service

import (
	"errors"

	"assets-api/internal/domain"
	"assets-api/internal/repository"
)

type AssetService struct {
	repo *repository.AssetRepository
}

func NewAssetService(r *repository.AssetRepository) *AssetService {
	return &AssetService{r}
}

func (s *AssetService) GetStats() (*domain.Stats, error) {
	return s.repo.GetStats()
}

func (s *AssetService) Count(t, status string) (int, error) {
	return s.repo.Count(t, status)
}

func (s *AssetService) BatchCreate(assets []domain.Asset) ([]string, error) {

	if len(assets) > 100 {
		return nil, errors.New("max 100 assets per request")
	}

	for _, a := range assets {

		if a.Type != "domain" && a.Type != "ip" && a.Type != "service" {
			return nil, errors.New("invalid asset type")
		}
	}

	return s.repo.BatchCreate(assets)
}

func (s *AssetService) BatchDelete(ids []string) (int, int, error) {
	return s.repo.BatchDelete(ids)
}

func (s *AssetService) List(page, limit int, t, status string) ([]domain.Asset, int, error) {
	return s.repo.List(page, limit, t, status)
}

func (s *AssetService) Search(q string) ([]domain.Asset, error) {
	return s.repo.Search(q)
}
