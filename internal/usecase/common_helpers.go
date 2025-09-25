package usecase

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func makeRequestID(ts, studentID string, rowIndex int) string {
	h := sha1.New()
	h.Write([]byte(ts + "|" + studentID + "|" + strconv.Itoa(rowIndex)))
	return "req_" + hex.EncodeToString(h.Sum(nil))[:10]
}

func toArabicDigits(s string) string {
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

func fmtYMD(y, m, d int) string {
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

func normalizedISODate(s string) (string, bool) {
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

func parseTimestampA(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	s = toArabicDigits(s)
	parts := strings.Fields(s) // "30/8/2568 09:30"
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
	iso, ok := normalizedISODate(datePart)
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
