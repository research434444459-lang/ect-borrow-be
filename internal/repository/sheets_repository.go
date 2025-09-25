package repository

import (
	"context"
	"errors"
	"fmt"
	"strings" // ⬅️ เพิ่มบรรทัดนี้

	"ect-borrow-be/internal/config"
	"ect-borrow-be/internal/dto"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type sheetsRepository struct {
	svc     *sheets.Service
	sheetID string

	invTab string
	reqTab string
}

func newSheetsRepository(cfg config.Config) (*sheetsRepository, error) {
	if cfg.SheetID == "" {
		return nil, errors.New("SHEET_ID is required")
	}
	ctx := context.Background()
	svc, err := sheets.NewService(ctx, option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		return nil, fmt.Errorf("sheets.NewService: %w", err)
	}
	inv := cfg.SheetTabInventory
	if inv == "" {
		inv = "Inventory"
	}
	req := cfg.SheetTabRequests
	if req == "" {
		req = "การตอบแบบฟอร์ม 1"
	}
	return &sheetsRepository{
		svc:     svc,
		sheetID: cfg.SheetID,
		invTab:  inv,
		reqTab:  req,
	}, nil
}

// ===== Inventory (เดิม) =====
func (r *sheetsRepository) FetchInventory(ctx context.Context) ([]dto.InventoryRow, error) {
	rng := fmt.Sprintf("%s!A2:Z", r.invTab)
	resp, err := r.svc.Spreadsheets.Values.Get(r.sheetID, rng).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("read sheet: %w", err)
	}
	rows := make([]dto.InventoryRow, 0, len(resp.Values))
	for _, v := range resp.Values {
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
			Status:   get(4), // E
			ImageURL: get(7), // H
		}
		if row.Device == "" {
			continue
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// ===== Requests (list เดิม) =====
func (r *sheetsRepository) FetchRequests(ctx context.Context) ([]dto.RequestRow, error) {
	rng := fmt.Sprintf("%s!A2:Z", r.reqTab)
	resp, err := r.svc.Spreadsheets.Values.
		Get(r.sheetID, rng).
		Context(ctx).
		ValueRenderOption("FORMATTED_VALUE").
		Do()
	if err != nil {
		return nil, fmt.Errorf("read sheet: %w", err)
	}
	out := make([]dto.RequestRow, 0, len(resp.Values))
	for i, v := range resp.Values {
		get := func(ix int) string {
			if ix < len(v) {
				if s, ok := v[ix].(string); ok {
					return s
				}
			}
			return ""
		}
		out = append(out, dto.RequestRow{
			ColA:      get(0), // A
			ColD:      get(3), // D
			Status:    get(5), // F
			StudentID: get(8), // G
			Name:      get(7), // H
			ColI:      get(3), // I
			RowIndex:  i,
		})
	}
	return out, nil
}

// ===== Request Details (เดิม) =====
func (r *sheetsRepository) FetchRequestDetails(ctx context.Context) ([]dto.RequestDetailRow, error) {
	rng := fmt.Sprintf("%s!A2:X", r.reqTab)
	resp, err := r.svc.Spreadsheets.Values.
		Get(r.sheetID, rng).
		Context(ctx).
		ValueRenderOption("FORMATTED_VALUE").
		Do()
	if err != nil {
		return nil, fmt.Errorf("read sheet: %w", err)
	}
	out := make([]dto.RequestDetailRow, 0, len(resp.Values))
	for i, v := range resp.Values {
		get := func(ix int) string {
			if ix < len(v) {
				if s, ok := v[ix].(string); ok {
					return s
				}
			}
			return ""
		}
		out = append(out, dto.RequestDetailRow{
			A_TS:           get(0),
			B_Confirmed:    get(1),
			C_Year:         get(2),
			D_Date:         get(3),
			E_PickupTime:   get(4),
			F_Status:       get(5),
			G_Unknown:      get(6),
			H_Name:         get(7),
			I_StudentID:    get(8),
			J_Phone:        get(9),
			K_GroupMembers: get(10),
			L_CourseName:   get(11),
			M_OtherCourse:  get(12),
			N_Teacher:      get(13),
			O_Item1:        get(14),
			P_Item2:        get(15),
			Q_Item3:        get(16),
			R_Item4:        get(17),
			S_Item5:        get(18),
			T_Item6:        get(19),
			U_Item7:        get(20),
			V_Item8:        get(21),
			W_Item9:        get(22),
			X_Item10:       get(23),
			RowIndex:       i,
		})
	}
	return out, nil
}

// ===== Admin Today (เดิม) =====
func (r *sheetsRepository) FetchAdminTodayRows(ctx context.Context) ([]dto.AdminTodayRow, error) {
	rng := fmt.Sprintf("%s!A2:AE", r.reqTab)
	resp, err := r.svc.Spreadsheets.Values.
		Get(r.sheetID, rng).
		Context(ctx).
		ValueRenderOption("FORMATTED_VALUE").
		Do()
	if err != nil {
		return nil, fmt.Errorf("read sheet: %w", err)
	}
	out := make([]dto.AdminTodayRow, 0, len(resp.Values))
	for i, v := range resp.Values {
		get := func(ix int) string {
			if ix < len(v) {
				if s, ok := v[ix].(string); ok {
					return s
				}
			}
			return ""
		}
		out = append(out, dto.AdminTodayRow{
			A_TS:        get(0),  // A
			D_Date:      get(3),  // D
			E_Time:      get(4),  // E
			H_Name:      get(7),  // H
			I_StudentID: get(8),  // I
			AC_Return:   get(28), // AC
			AD_Giver:    get(29), // AD
			AE_Receiver: get(30), // AE
			RowIndex:    i,
		})
	}
	return out, nil
}

// ===== Admin — Requests List (A..AE) =====
func (r *sheetsRepository) FetchAdminRequestRows(ctx context.Context) ([]dto.AdminRequestListRow, error) {
	rng := fmt.Sprintf("%s!A2:AE", r.reqTab)
	resp, err := r.svc.Spreadsheets.Values.
		Get(r.sheetID, rng).
		Context(ctx).
		ValueRenderOption("FORMATTED_VALUE").
		Do()
	if err != nil {
		return nil, fmt.Errorf("read sheet: %w", err)
	}
	out := make([]dto.AdminRequestListRow, 0, len(resp.Values))
	for i, v := range resp.Values {
		get := func(ix int) string {
			if ix < len(v) {
				if s, ok := v[ix].(string); ok {
					return s
				}
			}
			return ""
		}
		out = append(out, dto.AdminRequestListRow{
			A_TS:        get(0),  // A
			D_Date:      get(3),  // D
			E_Time:      get(4),  // E
			F_Status:    get(5),  // F
			H_Name:      get(7),  // H
			I_StudentID: get(8),  // I
			AC_Return:   get(28), // AC
			AD_Giver:    get(29), // AD
			AE_Receiver: get(30), // AE
			RowIndex:    i,
		})
	}
	return out, nil
}

// ===== Admin Requests (ใหม่) =====
// ใช้คอลัมน์ A, D, E, F, H, I, AC, AD, AE
func (r *sheetsRepository) FetchAdminListRows(ctx context.Context) ([]dto.AdminListRow, error) {
	rng := fmt.Sprintf("%s!A2:AE", r.reqTab)
	resp, err := r.svc.Spreadsheets.Values.
		Get(r.sheetID, rng).
		Context(ctx).
		ValueRenderOption("FORMATTED_VALUE").
		Do()
	if err != nil {
		return nil, fmt.Errorf("read sheet: %w", err)
	}
	out := make([]dto.AdminListRow, 0, len(resp.Values))
	for i, v := range resp.Values {
		get := func(ix int) string {
			if ix < len(v) {
				if s, ok := v[ix].(string); ok {
					return s
				}
			}
			return ""
		}
		out = append(out, dto.AdminListRow{
			A_TS:        get(0),  // A
			D_Date:      get(3),  // D
			E_Time:      get(4),  // E
			F_Status:    get(5),  // F
			H_Name:      get(7),  // H
			I_StudentID: get(8),  // I
			AC_Return:   get(28), // AC
			AD_Giver:    get(29), // AD
			AE_Receiver: get(30), // AE
			RowIndex:    i,
		})
	}
	return out, nil

}
func (r *sheetsRepository) FetchAdminRequestDetailRows(ctx context.Context) ([]dto.AdminRequestDetailRow, error) {
	// ครอบคลุม A..AO (รวม AA/AB/AC/AD/AE และ Y/Z)
	rng := fmt.Sprintf("%s!A2:AO", r.reqTab)
	resp, err := r.svc.Spreadsheets.Values.
		Get(r.sheetID, rng).
		Context(ctx).
		ValueRenderOption("FORMATTED_VALUE").
		Do()
	if err != nil {
		return nil, fmt.Errorf("read sheet: %w", err)
	}

	get := func(v []interface{}, ix int) string {
		if ix < len(v) {
			if s, ok := v[ix].(string); ok {
				return s
			}
		}
		return ""
	}

	out := make([]dto.AdminRequestDetailRow, 0, len(resp.Values))
	for i, v := range resp.Values {
		out = append(out, dto.AdminRequestDetailRow{
			RowIndex:    i,
			A_TS:        get(v, 0),
			B_Confirm:   get(v, 1),
			C_Year:      get(v, 2),
			D_Date:      get(v, 3),
			E_Time:      get(v, 4),
			F_Status:    get(v, 5),
			H_Name:      get(v, 7),
			I_StudentID: get(v, 8),
			J_Phone:     get(v, 9),
			K_Group:     get(v, 10),
			L_Course:    get(v, 11),
			M_Other:     get(v, 12),
			N_Teacher:   get(v, 13),
			O1:          get(v, 14),
			P2:          get(v, 15),
			Q3:          get(v, 16),
			R4:          get(v, 17),
			S5:          get(v, 18),
			T6:          get(v, 19),
			U7:          get(v, 20),
			V8:          get(v, 21),
			W9:          get(v, 22),
			X10:         get(v, 23),
			Y_Non:       get(v, 24), // Y
			Z_Other:     get(v, 25), // Z
			AA_File:     get(v, 26),
			AB_Att:      get(v, 27),
			AC_Ret:      get(v, 28),
			AD_Giver:    get(v, 29),
			AE_Recv:     get(v, 30),
			AFnum1:      get(v, 31),
			AGnum2:      get(v, 32),
			AHnum3:      get(v, 33),
			AInum4:      get(v, 34),
			AJnum5:      get(v, 35),
			AKnum6:      get(v, 36),
			ALnum7:      get(v, 37),
			AMnum8:      get(v, 38),
			ANnum9:      get(v, 39),
			AOnum10:     get(v, 40),
		})
	}
	return out, nil
}

// ===== Admin — Request Detail: Update Row =====
func (r *sheetsRepository) UpdateAdminRequestRow(ctx context.Context, rowIndex int, payload dto.AdminRequestUpdatePayload) error {
	row := rowIndex + 2 // A2 = rowIndex 0
	tab := r.reqTab

	vranges := []*sheets.ValueRange{}
	makeVR := func(a1 string, v interface{}) {
		vranges = append(vranges, &sheets.ValueRange{
			Range:  a1,
			Values: [][]interface{}{{v}},
		})
	}

	// F (status)
	if payload.Status != nil {
		makeVR(fmt.Sprintf("%s!F%d", tab, row), *payload.Status)
	}
	// AC (dateReturn)
	if payload.DateReturn != nil {
		makeVR(fmt.Sprintf("%s!AC%d", tab, row), *payload.DateReturn)
	}
	// AD/AE (giver/receiver)
	if payload.Giver != nil {
		makeVR(fmt.Sprintf("%s!AD%d", tab, row), *payload.Giver)
	}
	if payload.Receiver != nil {
		makeVR(fmt.Sprintf("%s!AE%d", tab, row), *payload.Receiver)
	}

	// ===== Items (ชื่อ/หมายเลข) =====
	// 1) payload.Items → แตกเป็น O..X (names) และ AF..AO (numbers)
	if len(payload.Items) > 0 {
		names := make([]string, 10)
		numbers := make([]string, 10)
		for i := 0; i < 10; i++ {
			if i < len(payload.Items) {
				names[i] = strings.TrimSpace(payload.Items[i].Name)
				if payload.Items[i].Number != nil {
					numbers[i] = strings.TrimSpace(*payload.Items[i].Number)
				}
				// เก็บ Non/Other จาก item[0] ถ้า top-level ยังไม่ส่งมา
				if i == 0 {
					if payload.Non == nil && payload.Items[i].Non != nil {
						payload.Non = payload.Items[i].Non
					}
					if payload.Other == nil && payload.Items[i].Other != nil {
						payload.Other = payload.Items[i].Other
					}
				}
			}
		}
		colsNames := []string{"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X"}
		colsNums := []string{"AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO"}
		for i, c := range colsNames {
			makeVR(fmt.Sprintf("%s!%s%d", tab, c, row), names[i])
		}
		for i, c := range colsNums {
			makeVR(fmt.Sprintf("%s!%s%d", tab, c, row), numbers[i])
		}
	}

	// 2) itemsNames (อัปเดตเฉพาะ O..X)
	if len(payload.ItemsNames) > 0 {
		cols := []string{"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X"}
		for i := 0; i < len(payload.ItemsNames) && i < len(cols); i++ {
			makeVR(fmt.Sprintf("%s!%s%d", tab, cols[i], row), payload.ItemsNames[i])
		}
	}

	// 3) itemsNumbers (อัปเดตเฉพาะ AF..AO)
	if len(payload.ItemsNumbers) > 0 {
		cols := []string{"AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO"}
		for i := 0; i < len(payload.ItemsNumbers) && i < len(cols); i++ {
			makeVR(fmt.Sprintf("%s!%s%d", tab, cols[i], row), payload.ItemsNumbers[i])
		}
	}

	// ===== Y/Z: missingItems / activityItems (ใหม่) หรือ Non/Other (เดิม) =====
	if payload.MissingItems != nil {
		makeVR(fmt.Sprintf("%s!Y%d", tab, row), *payload.MissingItems)
	} else if payload.Non != nil {
		makeVR(fmt.Sprintf("%s!Y%d", tab, row), *payload.Non)
	}
	if payload.ActivityItems != nil {
		makeVR(fmt.Sprintf("%s!Z%d", tab, row), *payload.ActivityItems)
	} else if payload.Other != nil {
		makeVR(fmt.Sprintf("%s!Z%d", tab, row), *payload.Other)
	}

	// แนบไฟล์ AA/AB (เดิม)
	if payload.AttachmentFilename != nil || payload.AttachmentURL != nil {
		val := ""
		if payload.AttachmentURL != nil && *payload.AttachmentURL != "" {
			val = *payload.AttachmentURL
		} else if payload.AttachmentFilename != nil {
			val = *payload.AttachmentFilename
		}
		makeVR(fmt.Sprintf("%s!AA%d", tab, row), val)
	}
	if payload.AttachmentID != nil || payload.AttachmentLabel != nil {
		val := ""
		if payload.AttachmentLabel != nil && *payload.AttachmentLabel != "" {
			val = *payload.AttachmentLabel
		} else if payload.AttachmentID != nil {
			val = *payload.AttachmentID
		}
		makeVR(fmt.Sprintf("%s!AB%d", tab, row), val)
	}

	if len(vranges) == 0 {
		return nil
	}
	req := &sheets.BatchUpdateValuesRequest{
		Data:             vranges,
		ValueInputOption: "USER_ENTERED",
	}
	_, err := r.svc.Spreadsheets.Values.BatchUpdate(r.sheetID, req).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("batch update: %w", err)
	}
	return nil
}
func (r *sheetsRepository) Service() *sheets.Service {
	return r.svc
}

