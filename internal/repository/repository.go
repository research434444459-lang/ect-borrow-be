package repository

import (
	"context"

	"ect-borrow-be/internal/config"
	"ect-borrow-be/internal/dto"
)

type Repository interface {
	// เดิม
	FetchInventory(ctx context.Context) ([]dto.InventoryRow, error)
	FetchRequests(ctx context.Context) ([]dto.RequestRow, error)
	FetchRequestDetails(ctx context.Context) ([]dto.RequestDetailRow, error)
	FetchAdminTodayRows(ctx context.Context) ([]dto.AdminTodayRow, error)
	FetchAdminRequestRows(ctx context.Context) ([]dto.AdminRequestListRow, error)

	// ใหม่: สำหรับ GET/PATCH /api/admin/requests/{requestId}
	FetchAdminRequestDetailRows(ctx context.Context) ([]dto.AdminRequestDetailRow, error)
	UpdateAdminRequestRow(ctx context.Context, rowIndex int, payload dto.AdminRequestUpdatePayload) error
}

func NewSheetsRepository(cfg config.Config) (Repository, error) {
	return newSheetsRepository(cfg)
}
