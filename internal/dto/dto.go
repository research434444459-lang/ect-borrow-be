package dto

// ===== API Envelope =====

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIResponse[T any] struct {
	Data  *T        `json:"data,omitempty"`
	Meta  any       `json:"meta"` // stays null in this service
	Error *APIError `json:"error,omitempty"`
}

// ===== Domain & View Models =====

type DeviceItem struct {
	Name      string `json:"name"`
	Serial    string `json:"serial"`
	Available int    `json:"available"`
	ImageURL  string `json:"imageUrl"`
}

type DevicesMeta struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type DevicesData struct {
	Meta  DevicesMeta  `json:"meta"`
	Items []DeviceItem `json:"items"`
}

// Row parsed from Google Sheets (Inventory)
// A: machineNo, B: serial, C: group, D: device name, H: imageUrl
// We only map what we need for this endpoint.

type InventoryRow struct {
	Serial   string
	Group    string
	Device   string
	Status   string // ⬅️ ใหม่: คอลัมน์ E
	ImageURL string
}
