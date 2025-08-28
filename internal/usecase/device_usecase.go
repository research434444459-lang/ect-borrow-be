package usecase

import (
	"context"
	"sort"
	"strings"

	"ect-borrow-be/internal/config"
	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
)

type DeviceUsecase interface {
	SearchDevices(ctx context.Context, category, q string) (dto.DevicesData, error)
}

type deviceUsecase struct {
	repo repository.Repository
	cfg  config.Config
}

func NewDeviceUsecase(repo repository.Repository, cfg config.Config) DeviceUsecase {
	return &deviceUsecase{repo: repo, cfg: cfg}
}

func (u *deviceUsecase) SearchDevices(ctx context.Context, category, q string) (dto.DevicesData, error) {
	norm := func(s string) string { return strings.TrimSpace(strings.ToLower(s)) }
	cat := norm(category)
	needle := norm(q)
	ready := norm(u.cfg.ReadyStatus)
	rows, err := u.repo.FetchInventory(ctx)
	if err != nil {
		return dto.DevicesData{}, err
	}

	//cat := strings.TrimSpace(strings.ToLower(category))
	//needle := strings.TrimSpace(strings.ToLower(q))

	// Filter by Column C (Group) using both `category` and `q` as per spec.
	filtered := make([]dto.InventoryRow, 0)
	for _, r := range rows {
		g := strings.ToLower(r.Group)
		if cat != "" && g != cat {
			continue
		}
		if needle != "" && g != needle {
			continue
		}

		filtered = append(filtered, r)
	}

	// Aggregate by device name (Column D)
	agg := map[string]struct {
		count  int
		serial string
		image  string
	}{}
	for _, r := range filtered {
		v, ok := agg[r.Device]
		if !ok {
			v = struct {
				count  int
				serial string
				image  string
			}{
				count: 0, serial: r.Serial, image: r.ImageURL,
			}
		}
		if norm(r.Status) == ready {
			v.count++
		}
		if v.image == "" && r.ImageURL != "" {
			v.image = r.ImageURL
		}
		agg[r.Device] = v
	}
	items := make([]dto.DeviceItem, 0, len(agg))
	for name, v := range agg {
		items = append(items, dto.DeviceItem{
			Name:      name,
			Serial:    v.serial,
			Available: v.count,
			ImageURL:  v.image,
		})
	}
	// Sort by name ASC for predictable output
	sort.Slice(items, func(i, j int) bool { return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name) })

	// Determine meta.title from Column C of first filtered row, else fallback to category
	title := strings.TrimSpace(category)
	if len(filtered) > 0 && strings.TrimSpace(filtered[0].Group) != "" {
		title = filtered[0].Group
	}

	data := dto.DevicesData{
		Meta:  dto.DevicesMeta{Title: title, Subtitle: u.cfg.Subtitle},
		Items: items,
	}
	return data, nil
}
