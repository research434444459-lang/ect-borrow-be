package config

import "os"

type Config struct {
	Port string

	SheetID           string
	SheetTabInventory string // "Inventory"
	SheetTabRequests  string // "การตอบแบบฟอร์ม 1"

	Subtitle    string // /api/devices
	ReadyStatus string // /api/devices

	FrontendOrigin string // <— เพิ่มบรรทัดนี้

}

func Load() Config {
	return Config{
		Port:              valueOr("PORT", "8080"),
		SheetID:           valueOr("SHEET_ID", ""),
		SheetTabInventory: valueOr("SHEET_TAB_INVENTORY", "Inventory"),
		SheetTabRequests:  valueOr("SHEET_TAB_REQUESTS", "การตอบแบบฟอร์ม 1"),
		Subtitle:          valueOr("DEVICES_SUBTITLE", "ดูรายการอุปกรณ์"),
		ReadyStatus:       valueOr("READY_STATUS", "พร้อมใช้งาน"),
		FrontendOrigin:    valueOr("FRONTEND_ORIGIN", "http://localhost:4200"),
	}
}

func valueOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
