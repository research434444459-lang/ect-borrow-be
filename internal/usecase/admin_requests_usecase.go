package usecase

import (
	"context"
	"sort"
	"strings"
	"time"

	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
)

type AdminRequestsUsecase interface {
	SearchAdminRequests(ctx context.Context,
		q, student, date, status string,
		page, pageSize int,
	) (items []dto.AdminRequestListItem, total int, err error)
}

type adminRequestsUsecase struct {
	repo repository.Repository
}

func NewAdminRequestsUsecase(repo repository.Repository) AdminRequestsUsecase {
	return &adminRequestsUsecase{repo: repo}
}

func (u *adminRequestsUsecase) SearchAdminRequests(
	ctx context.Context,
	q, student, date, status string,
	page, pageSize int,
) ([]dto.AdminRequestListItem, int, error) {

	rows, err := u.repo.FetchAdminRequestRows(ctx)
	if err != nil {
		return nil, 0, err
	}

	norm := func(s string) string { return strings.TrimSpace(strings.ToLower(s)) }
	has := func(hay, needle string) bool {
		if needle == "" {
			return true
		}
		return strings.Contains(norm(hay), norm(needle))
	}
	dateParamNorm, okDateParam := normalizedISODate(date)

	type wrap struct {
		r  dto.AdminRequestListRow
		t  time.Time
		ok bool
		i  int
	}
	matches := make([]wrap, 0, len(rows))

	for _, r := range rows {
		dBorrow, okD := normalizedISODate(r.D_Date)
		dReturn, okR := normalizedISODate(r.AC_Return)

		if student != "" && !(has(r.I_StudentID, student) || has(r.H_Name, student)) {
			continue
		}
		if status != "" && !has(r.F_Status, status) {
			continue
		}
		if q != "" {
			joined := strings.Join([]string{
				r.I_StudentID, r.H_Name, r.D_Date, r.E_Time, r.F_Status,
				r.AC_Return, r.AD_Giver, r.AE_Receiver, r.A_TS,
			}, " ")
			if !has(joined, q) {
				continue
			}
		}
		if date != "" && okDateParam {
			okRow := (okD && dBorrow == dateParamNorm) || (okR && dReturn == dateParamNorm)
			if !okRow {
				continue
			}
		}

		ts, okTs := parseTimestampA(r.A_TS)
		matches = append(matches, wrap{r: r, t: ts, ok: okTs, i: r.RowIndex})
	}

	// sort ตามคอลัมน์ A (เก่าสุดก่อน)
	sort.SliceStable(matches, func(i, j int) bool {
		ai, aj := matches[i], matches[j]
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

	all := make([]dto.AdminRequestListItem, 0, len(matches))
	for _, w := range matches {
		r := w.r
		toPtr := func(s string) *string {
			if strings.TrimSpace(s) == "" {
				return nil
			}
			v := s
			return &v
		}

		dateBorrow := r.D_Date
		if v, ok := normalizedISODate(r.D_Date); ok {
			dateBorrow = v
		}

		all = append(all, dto.AdminRequestListItem{
			RequestID:  makeRequestID(r.A_TS, r.I_StudentID, r.RowIndex),
			StudentID:  r.I_StudentID,
			Name:       r.H_Name,
			DateBorrow: dateBorrow,
			Time:       r.E_Time,
			DateReturn: toPtr(r.AC_Return),
			Giver:      toPtr(r.AD_Giver),
			Receiver:   toPtr(r.AE_Receiver),
			Status:     r.F_Status,
		})
	}

	total := len(all)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}
	start := (page - 1) * pageSize
	if start > total {
		return []dto.AdminRequestListItem{}, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}
