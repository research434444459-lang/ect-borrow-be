package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ect-borrow-be/internal/config"
	"ect-borrow-be/internal/dto"

	"google.golang.org/api/sheets/v4"
)

// ===== Interface =====

type ForecastRepository interface {
	Forecast(ctx context.Context, date time.Time, q string) (dto.ForecastData, error)
}

// ===== Impl =====

type forecastRepository struct {
	svc         *sheets.Service
	sheetID     string
	invTab      string // inventory tab name
	reqTab      string // responses tab name ("การตอบแบบฟอร์ม 1")
	timezone    string
	readyStatus string // e.g. "พร้อมใช้งาน"
}

func NewForecastRepository(svc *sheets.Service, cfg config.Config) ForecastRepository {
	return &forecastRepository{
		svc:         svc,
		sheetID:     cfg.SheetID,
		invTab:      cfg.SheetTabInventory,
		reqTab:      cfg.SheetTabRequests,
		timezone:    "Asia/Bangkok",
		readyStatus: cfg.ReadyStatus,
	}
}

func (r *forecastRepository) Forecast(ctx context.Context, date time.Time, q string) (dto.ForecastData, error) {
	target := dateTrunc(date)

	// 1) โหลด inventory แล้วรวม "พร้อมใช้งาน" ตามชื่อ (D) สถานะ (E)
	inv, err := r.loadInventory(ctx)
	if err != nil {
		return dto.ForecastData{}, fmt.Errorf("load inventory: %w", err)
	}
	type invAgg struct {
		ReadyCount int
		Serial     string
		ImageURL   string
	}
	agg := map[string]*invAgg{}
	for _, row := range inv {
		name := norm(row.Name)
		if name == "" {
			continue
		}
		if _, ok := agg[name]; !ok {
			agg[name] = &invAgg{Serial: row.Serial, ImageURL: row.ImageURL}
		}
		// นับเฉพาะสถานะพร้อมใช้งาน
		if equalsFold(row.Status, r.readyStatus) {
			agg[name].ReadyCount++
		}
		// เติม serial / image ถ้ายังว่าง
		if agg[name].Serial == "" && row.Serial != "" {
			agg[name].Serial = row.Serial
		}
		if agg[name].ImageURL == "" && row.ImageURL != "" {
			agg[name].ImageURL = row.ImageURL
		}
	}

	// 2) โหลดคำขอ แล้วนับรายการ (O..X) ที่ช่วงวัน [borrow, return) ทับกับ target
	reqs, err := r.loadRequests(ctx)
	if err != nil {
		return dto.ForecastData{}, fmt.Errorf("load requests: %w", err)
	}

	reserved := map[string]int{} // name -> count on target
	for _, rq := range reqs {
		// หากต้องการนับทุกสถานะ ให้เปลี่ยนเป็น: if !statusAffects(rq.Status) { continue }
		if !statusAffects(rq.Status) {
			continue
		}
		start, okS := parseDateAny(rq.DateBorrow)
		end, okE := parseDateAny(rq.DateReturn) // Z
		if !okS {
			continue
		}
		start = dateTrunc(start)
		// กรณีไม่มีวันคืน ให้ถือว่าใช้ 1 วัน
		if !okE {
			end = start.Add(24 * time.Hour)
		} else {
			end = dateTrunc(end)
			// กันข้อมูลเพี้ยน: ถ้า end <= start ให้ดันเป็น start+1d
			if !end.After(start) {
				end = start.Add(24 * time.Hour)
			}
		}

		// นับเฉพาะวันที่ target ทับกับช่วง [start, end)
		if target.Before(end) && (target.Equal(start) || target.After(start)) {
			for _, it := range rq.Items {
				n := norm(it)
				if n != "" {
					reserved[n]++
				}
			}
		}
	}

	// 3) สร้างผลลัพธ์
	out := dto.ForecastData{
		Date:  target.Format("2006-01-02"),
		Items: []dto.ForecastItem{},
	}
	for name, a := range agg {
		rem := a.ReadyCount - reserved[name]
		if rem < 0 {
			rem = 0
		}
		item := dto.ForecastItem{
			Name:      name,
			Serial:    a.Serial,
			ImageURL:  a.ImageURL,
			Remaining: rem,
		}
		if q != "" {
			qq := strings.ToLower(q)
			if !strings.Contains(strings.ToLower(item.Name), qq) &&
				!strings.Contains(strings.ToLower(item.Serial), qq) {
				continue
			}
		}
		out.Items = append(out.Items, item)
	}
	return out, nil
}

// ===== loaders =====

type inventoryRow struct {
	Serial   string
	Name     string // D
	Status   string // E
	ImageURL string // H
}

// อ่าน B:H แล้ว map B,D,E,H (อิงตามโครงชีตเดิม)
func (r *forecastRepository) loadInventory(ctx context.Context) ([]inventoryRow, error) {
	a1 := fmt.Sprintf("%s!B:H", r.invTab)
	resp, err := r.svc.Spreadsheets.Values.Get(r.sheetID, a1).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	out := make([]inventoryRow, 0, len(resp.Values))
	for _, v := range resp.Values {
		get := func(i int) string {
			if i < len(v) {
				return strings.TrimSpace(fmt.Sprint(v[i]))
			}
			return ""
		}
		row := inventoryRow{
			Serial:   get(0), // B
			Name:     get(2), // D
			Status:   get(3), // E
			ImageURL: get(6), // H
		}
		if row.Name == "" || equalsFold(row.Name, "ชื่ออุปกรณ์") {
			continue
		}
		out = append(out, row)
	}
	return out, nil
}

type requestRow struct {
	DateBorrow string   // D
	DateReturn string   // Z
	Status     string   // F (ถ้าไม่ต้องการฟิลเตอร์สถานะ ไปแก้ที่ statusAffects)
	Items      []string // O..X
}

func (r *forecastRepository) loadRequests(ctx context.Context) ([]requestRow, error) {
	a1 := fmt.Sprintf("%s!A:AE", r.reqTab)
	resp, err := r.svc.Spreadsheets.Values.Get(r.sheetID, a1).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	out := make([]requestRow, 0, len(resp.Values))
	for _, v := range resp.Values {
		get := func(i int) string {
			if i < len(v) {
				return strings.TrimSpace(fmt.Sprint(v[i]))
			}
			return ""
		}
		items := make([]string, 0, 10)
		for i := 14; i <= 23; i++ { // O..X
			if n := get(i); n != "" {
				items = append(items, n)
			}
		}
		row := requestRow{
			DateBorrow: get(3),  // D
			Status:     get(5),  // F
			DateReturn: get(28), // Z
			Items:      items,
		}
		if len(row.Items) == 0 {
			continue
		}
		out = append(out, row)
	}
	return out, nil
}

// ===== utils =====

func norm(s string) string {
	x := strings.TrimSpace(s)
	x = strings.ReplaceAll(x, "–", "-")
	x = strings.Join(strings.Fields(x), " ")
	return x
}

func equalsFold(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func parseDateAny(s string) (time.Time, bool) {
	ss := strings.TrimSpace(s)
	if ss == "" {
		return time.Time{}, false
	}
	// ISO
	if t, err := time.Parse("2006-01-02", ss); err == nil {
		return t, true
	}
	// D/M/YYYY, DD/MM/YYYY
	layouts := []string{"2/1/2006", "02/01/2006"}
	for _, ly := range layouts {
		if t, err := time.Parse(ly, ss); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func dateTrunc(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// สถานะที่นับเป็น "จองจริง" ต่อสต็อกวันนั้น
func statusAffects(s string) bool {
	switch strings.TrimSpace(s) {
	case "อนุมัติ", "รับของแล้ว":
		return true
	default:
		return false
	}
}
