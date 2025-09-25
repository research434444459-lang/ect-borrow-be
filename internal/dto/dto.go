package dto

// ===== API Envelope =====

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ListMeta struct {
	TraceID   string `json:"traceId,omitempty"`
	Page      int    `json:"page,omitempty"`
	PerPage   int    `json:"perPage,omitempty"`
	Total     int    `json:"total,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

type APIResponse[T any] struct {
	Data  *T        `json:"data,omitempty"`
	Meta  any       `json:"meta"`
	Error *APIError `json:"error,omitempty"`
}

// ===== Devices (เดิม) =====

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

// Inventory row (A..Z) ที่ใช้ใน /api/devices
// A: machineNo, B: serial, C: group, D: device, E: status, H: imageUrl
type InventoryRow struct {
	Serial   string
	Group    string
	Device   string
	Status   string
	ImageURL string
}

// ===== Requests (list เดิม) =====

type RequestItem struct {
	Name      string `json:"name"`      // H
	StudentID string `json:"studentId"` // G
	Date      string `json:"date"`      // I (fallback D)
	Status    string `json:"status"`    // F
}

type RequestRow struct {
	ColA      string // A: Timestamp
	ColD      string // D
	Status    string // F
	StudentID string // G
	Name      string // H
	ColI      string // I
	RowIndex  int    // fallback sort
}

// ===== Request Detail (เดิม) =====

type ItemLabelValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type RequestDetail struct {
	Name           string           `json:"name"`
	StudentID      string           `json:"studentId"`
	Date           string           `json:"date"`
	Status         string           `json:"status"`
	TS             string           `json:"ts"`
	Year           string           `json:"year"`
	Phone          string           `json:"phone"`
	GroupMembers   []string         `json:"groupMembers"`
	ConfirmedRules bool             `json:"confirmedRules"`
	PickupDate     string           `json:"pickupDate"`
	PickupTime     string           `json:"pickupTime"`
	CourseName     string           `json:"courseName"`
	OtherCourse    string           `json:"otherCourse"`
	Teacher        string           `json:"teacher"`
	Items          []ItemLabelValue `json:"items"`
	AdminNote      string           `json:"adminNote"`
}

type RequestDetailRow struct {
	A_TS           string
	B_Confirmed    string
	C_Year         string
	D_Date         string
	E_PickupTime   string
	F_Status       string
	G_Unknown      string
	H_Name         string
	I_StudentID    string
	J_Phone        string
	K_GroupMembers string
	L_CourseName   string
	M_OtherCourse  string
	N_Teacher      string
	O_Item1        string
	P_Item2        string
	Q_Item3        string
	R_Item4        string
	S_Item5        string
	T_Item6        string
	U_Item7        string
	V_Item8        string
	W_Item9        string
	X_Item10       string
	RowIndex       int
}

// ===== Admin Today (เดิม) =====

type AdminTodayItem struct {
	RequestID  string  `json:"requestId"`
	StudentID  string  `json:"studentId"`  // I
	Name       string  `json:"name"`       // H
	DateBorrow string  `json:"dateBorrow"` // D
	Time       string  `json:"time"`       // E
	DateReturn *string `json:"dateReturn"` // AC
	Giver      *string `json:"giver"`      // AD
	Receiver   *string `json:"receiver"`   // AE
}

type AdminTodayData struct {
	Date    string           `json:"date"`
	Borrow  []AdminTodayItem `json:"borrow"`
	Returns []AdminTodayItem `json:"returns"`
}

type AdminTodayMeta struct {
	Timezone     string `json:"timezone"`
	CountBorrow  int    `json:"countBorrow"`
	CountReturns int    `json:"countReturns"`
	ServerTime   string `json:"serverTime"`
}

type AdminTodayRow struct {
	A_TS        string // A
	D_Date      string // D
	E_Time      string // E
	H_Name      string // H
	I_StudentID string // I
	AC_Return   string // AC
	AD_Giver    string // AD
	AE_Receiver string // AE
	RowIndex    int
}

// ===== Admin Requests (ใหม่) =====

type AdminRequestsItem struct {
	RequestID  string  `json:"requestId"`
	StudentID  string  `json:"studentId"`  // I
	Name       string  `json:"name"`       // H
	DateBorrow string  `json:"dateBorrow"` // D
	Time       string  `json:"time"`       // E
	DateReturn *string `json:"dateReturn"` // AC (nullable)
	Giver      *string `json:"giver"`      // AD (nullable)
	Receiver   *string `json:"receiver"`   // AE (nullable)
	Status     string  `json:"status"`     // F
}

type AdminRequestsData struct {
	Items    []AdminRequestsItem `json:"items"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
	Total    int                 `json:"total"`
}

type AdminRequestsMeta struct {
	Timezone   string `json:"timezone"`
	ServerTime string `json:"serverTime"`
}

// แถวดิบที่ใช้สำหรับ /api/admin/requests
type AdminListRow struct {
	A_TS        string // A
	D_Date      string // D
	E_Time      string // E
	F_Status    string // F
	H_Name      string // H
	I_StudentID string // I
	AC_Return   string // AC
	AD_Giver    string // AD
	AE_Receiver string // AE
	RowIndex    int
}

// ===== Admin — Requests List =====

type AdminRequestListItem struct {
	RequestID  string  `json:"requestId"`
	StudentID  string  `json:"studentId"`  // I
	Name       string  `json:"name"`       // H
	DateBorrow string  `json:"dateBorrow"` // D (normalize ได้ ถ้า parse สำเร็จ)
	Time       string  `json:"time"`       // E
	DateReturn *string `json:"dateReturn"` // AC (nullable)
	Giver      *string `json:"giver"`      // AD (nullable)
	Receiver   *string `json:"receiver"`   // AE (nullable)
	Status     string  `json:"status"`     // F
}

type AdminRequestListData struct {
	Items    []AdminRequestListItem `json:"items"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"pageSize"`
	Total    int                    `json:"total"`
}

type AdminRequestListMeta struct {
	Timezone   string `json:"timezone"`
	ServerTime string `json:"serverTime"`
}

// แถวดิบสำหรับ admin/requests (ต้องอ่านถึง AE)
type AdminRequestListRow struct {
	A_TS        string // A
	D_Date      string // D
	E_Time      string // E
	F_Status    string // F
	H_Name      string // H
	I_StudentID string // I
	AC_Return   string // AC
	AD_Giver    string // AD
	AE_Receiver string // AE
	RowIndex    int
}

// ===== Admin — Request Detail =====

type AdminAttachment struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

type AdminRequestItem struct {
	Name   string  `json:"name"`             // O..X
	Number *string `json:"number,omitempty"` // AF..AO (nullable)

	// คงพฤติกรรมเดิม: บาง client เก่าจะส่ง Y/Z มาผูกไว้กับ item ตัวแรก
	Non   *string `json:"non,omitempty"`   // Y (เฉพาะ item แรก)
	Other *string `json:"other,omitempty"` // Z (เฉพาะ item แรก)
}

type AdminRequestDetail struct {
	RequestID    string             `json:"requestId"`
	StudentID    string             `json:"studentId"`
	Name         string             `json:"name"`
	Phone        string             `json:"phone"`
	GroupMembers string             `json:"groupMembers"`
	CourseName   string             `json:"courseName"`
	CourseOther  string             `json:"courseOther"`
	Teacher      string             `json:"teacher"`
	Status       string             `json:"status"`
	Timestamp    string             `json:"timestamp"`
	DateBorrow   string             `json:"dateBorrow"` // YYYY-MM-DD
	Time         string             `json:"time"`
	DateReturn   *string            `json:"dateReturn"`
	Giver        *string            `json:"giver"`    // AD
	Receiver     *string            `json:"receiver"` // AE
	Items        []AdminRequestItem `json:"items"`
	Attachments  []AdminAttachment  `json:"attachments"`

	// ✅ เพิ่ม field ใหม่ให้ FE ใช้โดยตรง (มาจากคอลัมน์ Y, Z)
	MissingItems  *string `json:"missingItems,omitempty"`  // Y
	ActivityItems *string `json:"activityItems,omitempty"` // Z
}

type AdminRequestDetailMeta struct {
	Timezone   string `json:"timezone"`
	ServerTime string `json:"serverTime"`
}

// แถวดิบเต็มสำหรับ detail: A..AO + AA + AB + AC..AE + Y/Z
type AdminRequestDetailRow struct {
	RowIndex    int
	A_TS        string // A
	B_Confirm   string // B
	C_Year      string // C
	D_Date      string // D
	E_Time      string // E
	F_Status    string // F
	H_Name      string // H
	I_StudentID string // I
	J_Phone     string // J
	K_Group     string // K
	L_Course    string // L
	M_Other     string // M
	N_Teacher   string // N
	// Items O..X (10)
	O1       string // O
	P2       string // P
	Q3       string // Q
	R4       string // R
	S5       string // S
	T6       string // T
	U7       string // U
	V8       string // V
	W9       string // W
	X10      string // X
	Y_Non    string // Y
	Z_Other  string // Z
	AA_File  string // AA (filename/url)
	AB_Att   string // AB (id/label)
	AC_Ret   string // AC
	AD_Giver string // AD
	AE_Recv  string // AE
	// Numbers AF..AO (10)
	AFnum1  string // AF
	AGnum2  string // AG
	AHnum3  string // AH
	AInum4  string // AI
	AJnum5  string // AJ
	AKnum6  string // AK
	ALnum7  string // AL
	AMnum8  string // AM
	ANnum9  string // AN
	AOnum10 string // AO
}

// ===== Update Payload (PATCH) =====

type AdminRequestUpdatePayload struct {
	// เดิม
	Status             *string  `json:"status,omitempty"`             // F
	DateReturn         *string  `json:"dateReturn,omitempty"`         // AC
	Giver              *string  `json:"giver,omitempty"`              // AD
	Receiver           *string  `json:"receiver,omitempty"`           // AE
	ItemsNumbers       []string `json:"itemsNumbers,omitempty"`       // AF..AO
	AttachmentID       *string  `json:"attachmentId,omitempty"`       // AB
	AttachmentLabel    *string  `json:"attachmentLabel,omitempty"`    // AB
	AttachmentFilename *string  `json:"attachmentFilename,omitempty"` // AA
	AttachmentURL      *string  `json:"attachmentUrl,omitempty"`      // AA
	Non                *string  `json:"non,omitempty"`                // Y (compat เดิม)
	Other              *string  `json:"other,omitempty"`              // Z (compat เดิม)

	// ✅ เพิ่มรองรับรูปแบบใหม่จาก FE
	Items         []AdminRequestItem `json:"items,omitempty"`         // {name, number, (optional non/other ใน item[0])}
	ItemsNames    []string           `json:"itemsNames,omitempty"`    // O..X (ถ้าส่งชื่ออย่างเดียว)
	MissingItems  *string            `json:"missingItems,omitempty"`  // Y (ชื่อใหม่)
	ActivityItems *string            `json:"activityItems,omitempty"` // Z (ชื่อใหม่)
}

// ForecastItem ใช้แสดงแถวในหน้า forecast
type ForecastItem struct {
	Name      string `json:"name"`
	Serial    string `json:"serial"`
	ImageURL  string `json:"imageUrl"`
	Remaining int    `json:"remaining"`
}

// ForecastPayload / Response
type ForecastData struct {
	Date  string         `json:"date"`
	Items []ForecastItem `json:"items"`
}
