package repository

import (
	"context"

	"ect-borrow-be/internal/config"
	"ect-borrow-be/internal/dto"
)

type Repository interface {
	FetchInventory(ctx context.Context) ([]dto.InventoryRow, error)
}

// Factory
func NewSheetsRepository(cfg config.Config) (Repository, error) {
	return newSheetsRepository(cfg)
}
