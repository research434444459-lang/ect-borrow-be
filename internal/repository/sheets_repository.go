package repository

import (
	"context"
	"errors"
	"fmt"

	"ect-borrow-be/internal/config"
	"ect-borrow-be/internal/dto"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type sheetsRepository struct {
	svc     *sheets.Service
	sheetID string
	tab     string
}

func newSheetsRepository(cfg config.Config) (*sheetsRepository, error) {
	if cfg.SheetID == "" {
		return nil, errors.New("SHEET_ID is required")
	}
	ctx := context.Background()
	svc, err := sheets.NewService(ctx, option.WithScopes(sheets.SpreadsheetsReadonlyScope))
	if err != nil {
		return nil, fmt.Errorf("sheets.NewService: %w", err)
	}
	return &sheetsRepository{svc: svc, sheetID: cfg.SheetID, tab: cfg.SheetTab}, nil
}

// FetchInventory reads A2:Z to skip headers and map needed columns.
func (r *sheetsRepository) FetchInventory(ctx context.Context) ([]dto.InventoryRow, error) {
	rng := fmt.Sprintf("%s!A2:Z", r.tab)
	resp, err := r.svc.Spreadsheets.Values.Get(r.sheetID, rng).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("read sheet: %w", err)
	}

	rows := make([]dto.InventoryRow, 0, len(resp.Values))
	for _, v := range resp.Values {
		// Defensive indexing; empty cells may be missing.
		get := func(i int) string {
			if i < len(v) {
				if s, ok := v[i].(string); ok {
					return s
				}
			}
			return ""
		}
		row := dto.InventoryRow{
			Serial:   get(1), // B
			Group:    get(2), // C
			Device:   get(3), // D
			Status:   get(4), // E ⬅️ เพิ่ม
			ImageURL: get(7), // H
		}
		// Skip empty device name
		if row.Device == "" {
			continue
		}
		rows = append(rows, row)
	}
	return rows, nil
}
