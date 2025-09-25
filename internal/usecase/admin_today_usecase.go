package usecase

import (
	"context"
	"sort"
	"strings"
	"time"

	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
)

type AdminTodayUsecase interface {
	GetAdminToday(ctx context.Context, date string) (dto.AdminTodayData, dto.AdminTodayMeta, error)
}

type adminTodayUsecase struct {
	repo repository.Repository
}

func NewAdminTodayUsecase(repo repository.Repository) AdminTodayUsecase {
	return &adminTodayUsecase{repo: repo}
}

func (u *adminTodayUsecase) GetAdminToday(ctx context.Context, date string) (dto.AdminTodayData, dto.AdminTodayMeta, error) {
	rows, err := u.repo.FetchAdminTodayRows(ctx)
	if err != nil {
		return dto.AdminTodayData{}, dto.AdminTodayMeta{}, err
	}

	dateParam, ok := normalizedISODate(date)
	if !ok {
		dateParam = date
	}

	type wrap struct {
		it dto.AdminTodayItem
		t  time.Time
		ok bool
		i  int
	}

	bor := make([]wrap, 0)
	ret := make([]wrap, 0)

	for _, r := range rows {
		// normalize วันที่ในแถว
		dBorrow, okD := normalizedISODate(r.D_Date)
		dReturn, okR := normalizedISODate(r.AC_Return)

		// สร้าง pointer ช่วย
		toPtr := func(s string) *string {
			if strings.TrimSpace(s) == "" {
				return nil
			}
			v := s
			return &v
		}

		item := dto.AdminTodayItem{
			RequestID:  makeRequestID(r.A_TS, r.I_StudentID, r.RowIndex),
			StudentID:  r.I_StudentID,
			Name:       r.H_Name,
			DateBorrow: dBorrow,
			Time:       r.E_Time,
			DateReturn: toPtr(r.AC_Return),
			Giver:      toPtr(r.AD_Giver),
			Receiver:   toPtr(r.AE_Receiver),
		}

		// เวลาไว้เรียงคิวตาม A
		ts, okTs := parseTimestampA(r.A_TS)

		// เข้ากลุ่ม borrow/returns ตาม dateParam
		if dateParam != "" && okD && dBorrow == dateParam {
			bor = append(bor, wrap{it: item, t: ts, ok: okTs, i: r.RowIndex})
		}
		if dateParam != "" && okR && dReturn == dateParam {
			ret = append(ret, wrap{it: item, t: ts, ok: okTs, i: r.RowIndex})
		}
	}

	// เรียงเก่าสุดก่อนตาม Timestamp (A)
	sort.SliceStable(bor, func(i, j int) bool {
		ai, aj := bor[i], bor[j]
		if ai.ok && aj.ok {
			return ai.t.Before(aj.t)
		}
		if ai.ok && !aj.ok {
			return true
		}
		if !ai.ok && aj.ok {
			return false
		}
		return ai.i < aj.i
	})
	sort.SliceStable(ret, func(i, j int) bool {
		ai, aj := ret[i], ret[j]
		if ai.ok && aj.ok {
			return ai.t.Before(aj.t)
		}
		if ai.ok && !aj.ok {
			return true
		}
		if !ai.ok && aj.ok {
			return false
		}
		return ai.i < aj.i
	})

	bItems := make([]dto.AdminTodayItem, 0, len(bor))
	for _, w := range bor {
		bItems = append(bItems, w.it)
	}
	rItems := make([]dto.AdminTodayItem, 0, len(ret))
	for _, w := range ret {
		rItems = append(rItems, w.it)
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	meta := dto.AdminTodayMeta{
		Timezone:     "Asia/Bangkok",
		CountBorrow:  len(bItems),
		CountReturns: len(rItems),
		ServerTime:   time.Now().In(loc).Format(time.RFC3339),
	}
	data := dto.AdminTodayData{
		Date:    dateParam,
		Borrow:  bItems,
		Returns: rItems,
	}
	return data, meta, nil
}
