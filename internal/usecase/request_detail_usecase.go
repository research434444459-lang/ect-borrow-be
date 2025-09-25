package usecase

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
)

// ใช้เช็ค not found จาก handler
var ErrNotFound = errors.New("not found")

// ใช้กับ endpoint: GET /api/requests/{studentId}/{date}
type RequestDetailUsecase interface {
	GetRequestDetail(ctx context.Context, studentId, date string) (dto.RequestDetail, error)
}

type requestDetailUsecase struct {
	repo repository.Repository
}

func NewRequestDetailUsecase(repo repository.Repository) RequestDetailUsecase {
	return &requestDetailUsecase{repo: repo}
}

func (u *requestDetailUsecase) GetRequestDetail(ctx context.Context, studentId, date string) (dto.RequestDetail, error) {
	rows, err := u.repo.FetchRequestDetails(ctx)
	if err != nil {
		return dto.RequestDetail{}, err
	}

	norm := func(s string) string { return strings.TrimSpace(strings.ToLower(s)) }

	toArabicDigits := func(s string) string { // แปลงเลขไทย -> อารบิก
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
	twoDigits := func(n int) string {
		if n < 10 {
			return "0" + strconv.Itoa(n)
		}
		return strconv.Itoa(n)
	}
	formatYMD := func(y, m, d int) string {
		return strconv.Itoa(y) + "-" + twoDigits(m) + "-" + twoDigits(d)
	}

	// map เดือนภาษาไทย (ย่อ/เต็ม)
	thMonths := map[string]int{
		"มกราคม": 1, "ม.ค": 1, "ม.ค.": 1,
		"กุมภาพันธ์": 2, "ก.พ": 2, "ก.พ.": 2,
		"มีนาคม": 3, "มี.ค": 3, "มี.ค.": 3,
		"เมษายน": 4, "เม.ย": 4, "เม.ย.": 4,
		"พฤษภาคม": 5, "พ.ค": 5, "พ.ค.": 5,
		"มิถุนายน": 6, "มิ.ย": 6, "มิ.ย.": 6,
		"กรกฎาคม": 7, "ก.ค": 7, "ก.ค.": 7,
		"สิงหาคม": 8, "ส.ค": 8, "ส.ค.": 8,
		"กันยายน": 9, "ก.ย": 9, "ก.ย.": 9,
		"ตุลาคม": 10, "ต.ค": 10, "ต.ค.": 10,
		"พฤศจิกายน": 11, "พ.ย": 11, "พ.ย.": 11,
		"ธันวาคม": 12, "ธ.ค": 12, "ธ.ค.": 12,
	}

	// แปลงสตริงวันที่ให้เป็น YYYY-MM-DD (รองรับเลขไทย, /.-, และชื่อเดือนภาษาไทย)
	normalizedISO := func(s string) (string, bool) {
		s = strings.TrimSpace(s)
		if s == "" {
			return "", false
		}
		s = toArabicDigits(s)

		// 1) เส้นทาง numeric ปกติ: yyyy/mm/dd หรือ dd/mm/yyyy
		numericTry := func(x string) (string, bool) {
			// ตัดเวลาที่ตามหลังช่องว่าง/T
			if i := strings.IndexAny(x, " T"); i >= 0 {
				x = x[:i]
			}
			// รวม separator ให้เป็น "/"
			x = strings.NewReplacer(".", "/", "-", "/", "—", "/", "–", "/").Replace(x)

			// เก็บเฉพาะเลขและ '/'
			var keep strings.Builder
			for _, r := range x {
				if unicode.IsDigit(r) || r == '/' {
					keep.WriteRune(r)
				}
			}
			parts := strings.Split(keep.String(), "/")
			if len(parts) != 3 {
				return "", false
			}
			atoi := func(z string) (int, bool) {
				z = strings.TrimLeft(z, "0")
				if z == "" {
					z = "0"
				}
				n, err := strconv.Atoi(z)
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
				} // พ.ศ. -> ค.ศ.
				if y < 1000 {
					return "", false
				}
				return formatYMD(y, m, d), true
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
				return formatYMD(y, m, d), true
			}
			return "", false
		}
		if v, ok := numericTry(s); ok {
			return v, true
		}

		// 2) เส้นทางชื่อเดือนภาษาไทย: "12 ก.ย. 2568", "12 กันยายน 2568"
		// ล้างคำ/สัญลักษณ์ที่รบกวน แล้วแทนที่ชื่อเดือนด้วย /MM/
		sm := s
		sm = strings.NewReplacer(",", " ", "เวลา", " ", "น.", " ", "นาฬิกา", " ").Replace(sm)
		// บีบช่องว่าง
		sm = strings.Join(strings.Fields(sm), " ")
		for k, m := range thMonths {
			// เพิ่มทั้งรูปแบบมี/ไม่มีเว้นวรรคข้างๆ
			sm = strings.ReplaceAll(sm, " "+k+" ", " "+"/"+twoDigits(m)+"/"+" ")
			sm = strings.ReplaceAll(sm, k, "/"+twoDigits(m)+"/")
		}
		// ลอง parse แบบ numeric อีกครั้งหลังแทนเดือนแล้ว
		if v, ok := numericTry(sm); ok {
			return v, true
		}
		return "", false
	}

	sameDate := func(a, b string) bool {
		ia, oka := normalizedISO(a)
		ib, okb := normalizedISO(b)
		return oka && okb && ia == ib
	}

	// parse timestamp จากสตริงที่อาจมีทั้งวันที่แบบไทยและเวลาแบบ "HH:MM" หรือ "HH.MM"
	parseTimestamp := func(s string) (time.Time, bool) {
		s = strings.TrimSpace(s)
		if s == "" {
			return time.Time{}, false
		}
		s = toArabicDigits(s)

		// หา date ส่วนที่ parse ได้ (ลองทีละ token และแบบ 3 token เผื่อชื่อเดือนไทย)
		var dateISO string
		parts := strings.Fields(s)
		for i := 0; i < len(parts); i++ {
			if v, ok := normalizedISO(parts[i]); ok {
				dateISO = v
				break
			}
			if i+2 < len(parts) {
				joined := parts[i] + " " + parts[i+1] + " " + parts[i+2]
				if v, ok := normalizedISO(joined); ok {
					dateISO = v
					break
				}
			}
		}
		if dateISO == "" {
			// เผื่อทั้งสตริงเป็นวันที่
			if v, ok := normalizedISO(s); ok {
				dateISO = v
			} else {
				return time.Time{}, false
			}
		}

		// หา time token แบบยืดหยุ่น
		var timePart string
		clean := func(x string) string {
			x = strings.Trim(x, " ,")
			x = strings.ReplaceAll(x, "เวลา", "")
			x = strings.ReplaceAll(x, "นาฬิกา", "")
			x = strings.TrimSuffix(x, "น.")
			x = strings.TrimSuffix(x, "น")
			return strings.TrimSpace(x)
		}
		for _, p := range parts {
			pp := clean(p)
			if strings.Contains(pp, ":") || strings.Contains(pp, ".") {
				pp = strings.ReplaceAll(pp, ".", ":")
				// ตัดกรณีหลงเหลือเช่น "09:30น."
				pp = clean(pp)
				if strings.Count(pp, ":") == 1 {
					pp += ":00"
				}
				// ตรวจดูว่าเป็น HH:MM[:SS] จริง ๆ
				segs := strings.Split(pp, ":")
				if len(segs) >= 2 && len(segs) <= 3 {
					timePart = pp
					break
				}
			}
		}

		layout := "2006-01-02"
		val := dateISO
		if timePart != "" {
			layout = "2006-01-02 15:04:05"
			val = dateISO + " " + timePart
		}
		t, err := time.ParseInLocation(layout, val, time.Local)
		if err != nil {
			return time.Time{}, false
		}
		return t, true
	}

	boolish := func(s string) bool {
		x := norm(s)
		return x == "true" || x == "1" || x == "yes" || x == "y" || x == "ใช่" ||
			strings.Contains(x, "ยืนยัน") || strings.Contains(x, "ตกลง") || strings.Contains(x, "ยอมรับ")
	}
	splitMembers := func(s string) []string {
		s = strings.ReplaceAll(s, "\r\n", "\n")
		s = strings.ReplaceAll(s, "、", ",")
		for _, sep := range []string{",", "\n", " และ "} {
			s = strings.ReplaceAll(s, sep, "|")
		}
		parts := strings.Split(s, "|")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				out = append(out, p)
			}
		}
		return out
	}
	isoOrRaw := func(s string) string {
		if v, ok := normalizedISO(s); ok {
			return v
		}
		return s
	}
	tsOrRaw := func(s string) string {
		if t, ok := parseTimestamp(s); ok {
			return t.Format("2006-01-02 15:04")
		}
		return s
	}

	studentParam := norm(studentId)
	dateParam := strings.TrimSpace(date)

	type mrow struct {
		row dto.RequestDetailRow
		ts  time.Time
		ok  bool
	}
	matches := make([]mrow, 0)
	for _, r := range rows {
		if norm(r.I_StudentID) != studentParam {
			continue
		} // คอลัมน์ I
		if !sameDate(r.D_Date, dateParam) {
			continue
		} // คอลัมน์ D
		t, ok := parseTimestamp(r.A_TS) // คอลัมน์ A (timestamp)
		matches = append(matches, mrow{row: r, ts: t, ok: ok})
	}
	if len(matches) == 0 {
		return dto.RequestDetail{}, ErrNotFound
	}

	// ถ้ามีหลายรายการ เลือก "เก่าสุดก่อน"
	sort.SliceStable(matches, func(i, j int) bool {
		ai, aj := matches[i], matches[j]
		if ai.ok && aj.ok {
			return ai.ts.Before(aj.ts)
		}
		if ai.ok && !aj.ok {
			return true
		}
		if !ai.ok && aj.ok {
			return false
		}
		return matches[i].row.RowIndex < matches[j].row.RowIndex
	})
	r := matches[0].row

	// ส่ง items ครบ 10 ช่องเสมอ (O..X)
	values := []string{
		r.O_Item1, r.P_Item2, r.Q_Item3, r.R_Item4, r.S_Item5,
		r.T_Item6, r.U_Item7, r.V_Item8, r.W_Item9, r.X_Item10,
	}
	items := make([]dto.ItemLabelValue, 0, 10)
	for i := 0; i < 10; i++ {
		label := "อุปกรณ์สำหรับงานผลิตสื่อชิ้นที่ " + strconv.Itoa(i+1)
		val := ""
		if i < len(values) {
			val = values[i]
		}
		items = append(items, dto.ItemLabelValue{Label: label, Value: val})
	}

	detail := dto.RequestDetail{
		Name:           r.H_Name,                       // H
		StudentID:      r.I_StudentID,                  // I
		Date:           isoOrRaw(r.D_Date),             // D
		Status:         r.F_Status,                     // F
		TS:             tsOrRaw(r.A_TS),                // A (แปลงเป็น "YYYY-MM-DD HH:MM")
		Year:           r.C_Year,                       // C
		Phone:          r.J_Phone,                      // J
		GroupMembers:   splitMembers(r.K_GroupMembers), // K
		ConfirmedRules: boolish(r.B_Confirmed),         // B
		PickupDate:     isoOrRaw(r.D_Date),             // D
		PickupTime:     r.E_PickupTime,                 // E (เก็บตามเดิม)
		CourseName:     r.L_CourseName,                 // L
		OtherCourse:    r.M_OtherCourse,                // M
		Teacher:        r.N_Teacher,                    // N
		Items:          items,                          // O..X
		AdminNote:      "",
	}
	return detail, nil
}
