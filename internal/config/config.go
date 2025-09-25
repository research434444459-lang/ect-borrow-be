package config

import "os"

type Config struct {
	Port string

	SheetID           string
	SheetTabInventory string // "Inventory"
	SheetTabRequests  string // "การตอบแบบฟอร์ม 1"

	Subtitle    string // /api/devices
	ReadyStatus string // /api/devices
}

func Load() Config {
	return Config{
		Port:              valueOr("PORT", "8080"),
		SheetID:           valueOr("SHEET_ID", ""),
		SheetTabInventory: valueOr("SHEET_TAB_INVENTORY", "Inventory"),
		SheetTabRequests:  valueOr("SHEET_TAB_REQUESTS", "การตอบแบบฟอร์ม 1"),
		Subtitle:          valueOr("DEVICES_SUBTITLE", "ดูรายการอุปกรณ์"),
		ReadyStatus:       valueOr("READY_STATUS", "พร้อมใช้งาน"),
	}
}

func valueOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
