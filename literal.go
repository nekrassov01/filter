package filter

import (
	"fmt"
	"math"
	"strconv"
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

// parseNumber parses a literal of the requested numeric type. Decimal integers
// are never accepted as floats, so out-of-range integers cannot be rounded.
// Integer failures reuse strconv errors; callers attach the input and position.
func parseNumber[T int64 | uint64 | float64](s string) (T, error) {
	var zero T
	digits := s
	negative := false
	if len(digits) > 0 && (digits[0] == '+' || digits[0] == '-') {
		negative = digits[0] == '-'
		digits = digits[1:]
	}
	if _, ok := any(zero).(float64); ok {
		integer := digits != ""
		for i := 0; i < len(digits); i++ {
			if c := digits[i]; (c < '0' || c > '9') && c != '_' {
				integer = false
				break
			}
		}
		if integer {
			return zero, strconv.ErrSyntax
		}
		v, err := strconv.ParseFloat(s, 64)
		return T(v), err
	}
	if digits == "" {
		return zero, strconv.ErrSyntax
	}
	limit := uint64(math.MaxInt64)
	switch any(zero).(type) {
	case int64:
		if negative {
			limit++
		}
	case uint64:
		if negative {
			return zero, strconv.ErrSyntax
		}
		limit = math.MaxUint64
	}
	var n uint64
	for i := 0; i < len(digits); i++ {
		switch c := digits[i]; {
		case '0' <= c && c <= '9':
			d := uint64(c - '0')
			if n > limit/10 || (n == limit/10 && d > limit%10) {
				return zero, strconv.ErrRange
			}
			n = n*10 + d
		case c == '_':
			if i == 0 || i+1 == len(digits) || digits[i-1] < '0' || digits[i-1] > '9' || digits[i+1] < '0' || digits[i+1] > '9' {
				return zero, strconv.ErrSyntax
			}
		default:
			return zero, strconv.ErrSyntax
		}
	}
	if negative {
		// The magnitude 1<<63 represents MinInt64 after negation.
		return -T(n), nil
	}
	return T(n), nil
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
