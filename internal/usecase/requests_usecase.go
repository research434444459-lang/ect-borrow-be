package usecase

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
)

// ใช้กับ endpoint: GET /api/requests?student=&date=
type RequestUsecase interface {
	SearchRequests(ctx context.Context, student, date string) ([]dto.RequestItem, int, error)
}

type requestUsecase struct {
	repo repository.Repository
}

func NewRequestUsecase(repo repository.Repository) RequestUsecase {
	return &requestUsecase{repo: repo}
}

func (u *requestUsecase) SearchRequests(ctx context.Context, student, date string) ([]dto.RequestItem, int, error) {
	rows, err := u.repo.FetchRequests(ctx)
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

	// ----- helpers สำหรับ normalize date และ parse timestamp A -----
	toArabicDigits := func(s string) string {
		var b strings.Builder
		for _, r := range s {
			switch r {
			case '๐':
				b.WriteRune('0')
			case '๑':
				b.WriteRune('1')
			case '๒':
				b.WriteRune('2')
			case '๓':
				b.WriteRune('3')
			case '๔':
				b.WriteRune('4')
			case '๕':
				b.WriteRune('5')
			case '๖':
				b.WriteRune('6')
			case '๗':
				b.WriteRune('7')
			case '๘':
				b.WriteRune('8')
			case '๙':
				b.WriteRune('9')
			default:
				b.WriteRune(r)
			}
		}
		return b.String()
	}
	fmtYMD := func(y, m, d int) string {
		yy := strconv.Itoa(y)
		mm := strconv.Itoa(m)
		dd := strconv.Itoa(d)
		if m < 10 {
			mm = "0" + mm
		}
		if d < 10 {
			dd = "0" + dd
		}
		return yy + "-" + mm + "-" + dd
	}
	normalizedISO := func(s string) (string, bool) {
		s = strings.TrimSpace(s)
		if s == "" {
			return "", false
		}
		s = toArabicDigits(s)
		if i := strings.IndexAny(s, " T"); i >= 0 {
			s = s[:i]
		}
		s = strings.NewReplacer(".", "/", "-", "/", "—", "/", "–", "/").Replace(s)
		var keep strings.Builder
		for _, r := range s {
			if unicode.IsDigit(r) || r == '/' {
				keep.WriteRune(r)
			}
		}
		parts := strings.Split(keep.String(), "/")
		if len(parts) != 3 {
			return "", false
		}
		atoi := func(x string) (int, bool) {
			x = strings.TrimLeft(x, "0")
			if x == "" {
				x = "0"
			}
			n, err := strconv.Atoi(x)
			return n, err == nil
		}
		// yyyy/mm/dd
		if len(parts[0]) == 4 {
			y, ok1 := atoi(parts[0])
			m, ok2 := atoi(parts[1])
			d, ok3 := atoi(parts[2])
			if !(ok1 && ok2 && ok3) {
				return "", false
			}
			if y >= 2400 {
				y -= 543
			}
			if y < 1000 {
				return "", false
			}
			return fmtYMD(y, m, d), true
		}
		// dd/mm/yyyy
		if len(parts[2]) == 4 {
			d, ok1 := atoi(parts[0])
			m, ok2 := atoi(parts[1])
			y, ok3 := atoi(parts[2])
			if !(ok1 && ok2 && ok3) {
				return "", false
			}
			if y >= 2400 {
				y -= 543
			}
			if y < 1000 {
				return "", false
			}
			return fmtYMD(y, m, d), true
		}
		return "", false
	}
	sameDate := func(a, b string) bool {
		ia, oka := normalizedISO(a)
		ib, okb := normalizedISO(b)
		return oka && okb && ia == ib
	}
	parseTimestampA := func(s string) (time.Time, bool) {
		s = strings.TrimSpace(s)
		if s == "" {
			return time.Time{}, false
		}
		s = toArabicDigits(s)
		parts := strings.Fields(s)
		var datePart, timePart string
		if len(parts) >= 1 {
			datePart = parts[0]
		}
		if len(parts) >= 2 {
			timePart = parts[1]
			if strings.Count(timePart, ":") == 1 {
				timePart += ":00"
			}
		}
		iso, ok := normalizedISO(datePart)
		if !ok {
			return time.Time{}, false
		}
		layout := "2006-01-02"
		val := iso
		if timePart != "" {
			layout = "2006-01-02 15:04:05"
			val = iso + " " + timePart
		}
		t, err := time.ParseInLocation(layout, val, time.Local)
		if err != nil {
			return time.Time{}, false
		}
		return t, true
	}
	// ---------------------------------------------------------------

	studentParam := strings.TrimSpace(student)
	dateParam := strings.TrimSpace(date)

	type rowWrap struct {
		item dto.RequestItem
		ts   time.Time
		ok   bool
		idx  int
	}

	buf := make([]rowWrap, 0, len(rows))
	for i, r := range rows {
		// match student: contains ใน D/H/I และรวม G (studentId)
		matchStudent := (studentParam == "") ||
			has(r.ColD, studentParam) || has(r.Name, studentParam) || has(r.ColI, studentParam) || has(r.StudentID, studentParam)

		// match date: contains เดิม + เทียบ normalized กับ D/I
		matchDate := (dateParam == "") ||
			has(r.ColD, dateParam) || has(r.Name, dateParam) || has(r.ColI, dateParam) ||
			sameDate(r.ColD, dateParam) || sameDate(r.ColI, dateParam)

		if !(matchStudent && matchDate) {
			continue
		}

		dateVal := strings.TrimSpace(r.ColI)
		if dateVal == "" {
			dateVal = r.ColD
		}

		ts, ok := parseTimestampA(r.ColA)

		buf = append(buf, rowWrap{
			item: dto.RequestItem{
				Name:      r.Name,      // H
				StudentID: r.StudentID, // G
				Date:      dateVal,     // I -> D
				Status:    r.Status,    // F
			},
			ts:  ts,
			ok:  ok,
			idx: i,
		})
	}

	// เรียงตาม Timestamp A: เก่าสุดก่อน
	sort.SliceStable(buf, func(i, j int) bool {
		ai, aj := buf[i], buf[j]
		if ai.ok && aj.ok {
			return ai.ts.Before(aj.ts)
		}
		if ai.ok && !aj.ok {
			return true
		}
		if !ai.ok && aj.ok {
			return false
		}
		return ai.idx < aj.idx
	})

	items := make([]dto.RequestItem, 0, len(buf))
	for _, w := range buf {
		items = append(items, w.item)
	}

	return items, len(items), nil
}
