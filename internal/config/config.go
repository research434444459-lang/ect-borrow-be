package config

import (
	"os"
)

type Config struct {
	Port                 string
	SheetID              string
	SheetTab             string
	GoogleAppCredentials string // for GOOGLE_APPLICATION_CREDENTIALS, optional
	Subtitle             string // for data.meta.subtitle, default: "ดูรายการอุปกรณ์"
	ReadyStatus          string // ค่าที่คอลัมน์ E ต้องเท่ากับ (default: "พร้อมใช้งาน")
}

func Load() Config {
	c := Config{
		Port:                 valueOr("PORT", "8080"),
		SheetID:              valueOr("SHEET_ID", "1VygRdp0QzE3iHBeCwaC1l7crHa8CfDMCWmyq5XRB_NY"),
		SheetTab:             valueOr("SHEET_TAB", "Inventory"),
		Subtitle:             valueOr("DEVICES_SUBTITLE", "ดูรายการอุปกรณ์"),
		GoogleAppCredentials: valueOr("GOOGLE_APPLICATION_CREDENTIALS", "/ect-borrow-be/internal/config/ect-borrow-credential.json"), // no default, optional
		ReadyStatus:          valueOr("READY_STATUS", "พร้อมใช้งาน"),
	}
	return c
}

func valueOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
