package filter

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// timeLayouts are the layouts a time literal may use, most common first.
var timeLayouts = [...]string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	time.DateTime,
	time.DateOnly,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC850,
	time.RFC822,
	time.RFC822Z,
}

// parseNumber parses a literal of the requested numeric type. Integer literals
// are never accepted as floats, so out-of-range integers cannot be rounded.
func parseNumber[T int64 | uint64 | float64](s string) (T, error) {
	var zero T
	if _, ok := any(zero).(float64); ok {
		digits := s
		if len(digits) > 0 && (digits[0] == '+' || digits[0] == '-') {
			digits = digits[1:]
		}
		integer := digits != ""
		for _, c := range digits {
			if (c < '0' || c > '9') && c != '_' {
				integer = false
				break
			}
		}
		if integer {
			return zero, fmt.Errorf("invalid floating-point literal %q", s)
		}
		v, err := strconv.ParseFloat(s, 64)
		return T(v), err
	}
	if strings.Contains(s, "_") {
		// Base 10 rejects separators, while base 0 treats leading zeros as octal.
		// Validate separators before removing them; ignore the floating-point value.
		if _, err := strconv.ParseFloat(s, 64); err != nil {
			return zero, err
		}
		s = strings.ReplaceAll(s, "_", "")
	}
	if _, ok := any(zero).(int64); ok {
		v, err := strconv.ParseInt(s, 10, 64)
		return T(v), err
	}
	v, err := strconv.ParseUint(strings.TrimPrefix(s, "+"), 10, 64)
	return T(v), err
}

// parseTime converts Unix seconds or a literal in one of timeLayouts to a
// UTC time. A zone abbreviation other than UTC or GMT is rejected, since
// time.Parse resolves no other abbreviation to an offset.
func parseTime(s string) (time.Time, error) {
	// Unix seconds.
	digits := s
	if digits != "" && (digits[0] == '-' || digits[0] == '+') {
		digits = digits[1:]
	}
	if digits != "" && digits[0] != '_' && digits[len(digits)-1] != '_' {
		var sec int64
		integer := true
		for i := 0; i < len(digits) && integer; i++ {
			switch c := digits[i]; {
			case c == '_':
			case '0' <= c && c <= '9':
				d := int64(c - '0')
				if sec > (math.MaxInt64-d)/10 {
					return time.Time{}, fmt.Errorf("unix seconds out of range %q", s)
				}
				sec = sec*10 + d
			default:
				integer = false
			}
		}
		if integer {
			if s[0] == '-' {
				sec = -sec
			}
			return time.Unix(sec, 0).UTC(), nil
		}
	}
	for _, layout := range timeLayouts {
		t, err := time.ParseInLocation(layout, s, time.UTC)
		if err != nil {
			continue
		}
		if name, _ := t.Zone(); name != "" && name != "UTC" && name != "GMT" {
			return time.Time{}, fmt.Errorf("unknown time zone %q", name)
		}
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("unrecognized time %q", s)
}
