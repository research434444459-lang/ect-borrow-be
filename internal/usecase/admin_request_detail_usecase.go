package usecase

import (
	"context"
	"sort"
	"strings"
	"time"

	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
)

type AdminRequestDetailUsecase interface {
	GetByID(ctx context.Context, requestID string) (dto.AdminRequestDetail, error)
	UpdateByID(ctx context.Context, requestID string, payload dto.AdminRequestUpdatePayload) (dto.AdminRequestDetail, error)
}

type adminRequestDetailUsecase struct {
	repo repository.Repository
}

func NewAdminRequestDetailUsecase(repo repository.Repository) AdminRequestDetailUsecase {
	return &adminRequestDetailUsecase{repo: repo}
}

func normalizeDate(s string) string {
	t, err := time.Parse("2006-01-02", strings.TrimSpace(s))
	if err == nil {
		return t.Format("2006-01-02")
	}
	// ลองแบบไทยนิยม 12/9/2025 หรือ 12/09/2025
	t2, err2 := time.Parse("2/1/2006", strings.TrimSpace(s))
	if err2 == nil {
		return t2.Format("2006-01-02")
	}
	t3, err3 := time.Parse("02/01/2006", strings.TrimSpace(s))
	if err3 == nil {
		return t3.Format("2006-01-02")
	}
	// ถ้า parse ไม่ได้ ให้คืนค่าที่รับมา (กันข้อมูลเดิมหาย)
	return s
}

func (u *adminRequestDetailUsecase) findRowByRequestID(ctx context.Context, requestID string) (dto.AdminRequestDetailRow, bool, error) {
	rows, err := u.repo.FetchAdminRequestDetailRows(ctx)
	if err != nil {
		return dto.AdminRequestDetailRow{}, false, err
	}

	type wrap struct {
		row dto.AdminRequestDetailRow
		t   time.Time
		ok  bool
	}
	ws := make([]wrap, 0, len(rows))
	for _, r := range rows {
		ts, ok := parseTimestampA(r.A_TS)
		ws = append(ws, wrap{row: r, t: ts, ok: ok})
	}
	// เรียงตามคอลัมน์ A: เก่าสุดก่อน
	sort.SliceStable(ws, func(i, j int) bool {
		ai, aj := ws[i], ws[j]
		if ai.ok && aj.ok {
			return ai.t.Before(aj.t)
		}
		if ai.ok && !aj.ok {
			return true
		}
		if !ai.ok && aj.ok {
			return false
		}
		return ws[i].row.RowIndex < ws[j].row.RowIndex
	})

	for _, w := range ws {
		id := makeRequestID(w.row.A_TS, w.row.I_StudentID, w.row.RowIndex)
		if id == requestID {
			return w.row, true, nil
		}
	}
	return dto.AdminRequestDetailRow{}, false, nil
}

func (u *adminRequestDetailUsecase) rowToDetail(r dto.AdminRequestDetailRow, requestID string) dto.AdminRequestDetail {
	// Attachments (แบบง่าย: 1 รายการ ถ้ามีค่า)
	atts := []dto.AdminAttachment{}
	if strings.TrimSpace(r.AA_File) != "" || strings.TrimSpace(r.AB_Att) != "" {
		atts = append(atts, dto.AdminAttachment{
			ID:       r.AB_Att,
			Label:    r.AB_Att,
			Filename: r.AA_File,
			URL:      r.AA_File,
		})
	}

	// Items (ชื่อ O..X + หมายเลข AF..AO) + ฝังค่า Y/Z ไว้ที่ item แรก (คง behavior เดิม)
	names := []string{r.O1, r.P2, r.Q3, r.R4, r.S5, r.T6, r.U7, r.V8, r.W9, r.X10}
	nums := []string{r.AFnum1, r.AGnum2, r.AHnum3, r.AInum4, r.AJnum5, r.AKnum6, r.ALnum7, r.AMnum8, r.ANnum9, r.AOnum10}
	items := make([]dto.AdminRequestItem, 0, 10)
	for i := 0; i < 10; i++ {
		var numPtr *string
		if strings.TrimSpace(nums[i]) != "" {
			v := nums[i]
			numPtr = &v
		}
		it := dto.AdminRequestItem{
			Name:   names[i],
			Number: numPtr,
		}
		if i == 0 {
			if strings.TrimSpace(r.Y_Non) != "" {
				v := r.Y_Non
				it.Non = &v
			}
			if strings.TrimSpace(r.Z_Other) != "" {
				v := r.Z_Other
				it.Other = &v
			}
		}
		items = append(items, it)
	}

	// Timestamp -> "YYYY-MM-DD HH:MM" ถ้า parse ได้
	ts := r.A_TS
	if t, ok := parseTimestampA(r.A_TS); ok {
		ts = t.Format("2006-01-02 15:04")
	}

	// pointer helpers
	var dateReturn *string
	if strings.TrimSpace(r.AC_Ret) != "" {
		v := r.AC_Ret
		dateReturn = &v
	}
	var giver *string
	if strings.TrimSpace(r.AD_Giver) != "" {
		v := r.AD_Giver
		giver = &v
	}
	var receiver *string
	if strings.TrimSpace(r.AE_Recv) != "" {
		v := r.AE_Recv
		receiver = &v
	}

	// ✅ top-level missing/activity
	var missing *string
	if strings.TrimSpace(r.Y_Non) != "" {
		v := r.Y_Non
		missing = &v
	}
	var activity *string
	if strings.TrimSpace(r.Z_Other) != "" {
		v := r.Z_Other
		activity = &v
	}

	return dto.AdminRequestDetail{
		RequestID:    requestID,
		StudentID:    r.I_StudentID,
		Name:         r.H_Name,
		Phone:        r.J_Phone,
		GroupMembers: r.K_Group,
		CourseName:   r.L_Course,
		CourseOther:  r.M_Other,
		Teacher:      r.N_Teacher,
		Status:       r.F_Status,
		Timestamp:    ts,
		DateBorrow:   normalizeDate(r.D_Date),
		Time:         r.E_Time,
		DateReturn:   dateReturn,
		Giver:        giver,
		Receiver:     receiver,
		Items:        items,
		Attachments:  atts,

		// ✅ ส่งค่าออกไปเป็น top-level
		MissingItems:  missing,  // Y
		ActivityItems: activity, // Z
	}
}

func (u *adminRequestDetailUsecase) GetByID(ctx context.Context, requestID string) (dto.AdminRequestDetail, error) {
	row, ok, err := u.findRowByRequestID(ctx, requestID)
	if err != nil {
		return dto.AdminRequestDetail{}, err
	}
	if !ok {
		return dto.AdminRequestDetail{}, ErrNotFound
	}
	id := makeRequestID(row.A_TS, row.I_StudentID, row.RowIndex)
	return u.rowToDetail(row, id), nil
}

func (u *adminRequestDetailUsecase) UpdateByID(ctx context.Context, requestID string, payload dto.AdminRequestUpdatePayload) (dto.AdminRequestDetail, error) {
	row, ok, err := u.findRowByRequestID(ctx, requestID)
	if err != nil {
		return dto.AdminRequestDetail{}, err
	}
	if !ok {
		return dto.AdminRequestDetail{}, ErrNotFound
	}

	// เขียนกลับลงชีต
	if err := u.repo.UpdateAdminRequestRow(ctx, row.RowIndex, payload); err != nil {
		return dto.AdminRequestDetail{}, err
	}

	// อ่านใหม่เพื่อสะท้อนค่าที่อัปเดต
	row2, _, err := u.findRowByRequestID(ctx, requestID)
	if err != nil {
		return dto.AdminRequestDetail{}, err
	}
	id := makeRequestID(row2.A_TS, row2.I_StudentID, row2.RowIndex)
	return u.rowToDetail(row2, id), nil
}
