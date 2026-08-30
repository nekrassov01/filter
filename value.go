package filter

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// kind identifies how a Value's fields are read.
type kind uint8

// Kinds of resolved values.
const (
	kindNone     kind = iota // zero Value: nothing resolved
	kindString               // s
	kindNumber               // a holds the float64 bits
	kindDuration             // a holds the duration
	kindTime                 // a holds Unix seconds, b the nanoseconds
)

// Value is the value of an identifier, as returned by a Resolver.
type Value struct {
	kind kind
	s    string
	a, b int64
}

// String returns a Value holding s.
func String(s string) Value {
	return Value{
		kind: kindString,
		s:    s,
	}
}

// Number returns a Value holding n.
func Number(n float64) Value {
	//nolint:gosec // bit pattern conversion
	return Value{
		kind: kindNumber,
		a:    int64(math.Float64bits(n)),
	}
}

// Duration returns a Value holding d.
func Duration(d time.Duration) Value {
	return Value{
		kind: kindDuration,
		a:    int64(d),
	}
}

// Time returns a Value holding the instant t.
func Time(t time.Time) Value {
	return Value{
		kind: kindTime,
		a:    t.Unix(),
		b:    int64(t.Nanosecond()),
	}
}

// Bool returns a Value holding b, which compares as the string "true" or "false".
func Bool(b bool) Value {
	return Value{
		kind: kindString,
		s:    strconv.FormatBool(b),
	}
}

// ValueOf converts a Go value to a Value. Strings, integer and float types,
// time.Time, time.Duration, and bool keep their kind; any other value is
// formatted with fmt.Sprint and compared as a string.
func ValueOf(v any) Value {
	switch v := v.(type) {
	case string:
		return String(v)
	case int:
		return Number(float64(v))
	case int8:
		return Number(float64(v))
	case int16:
		return Number(float64(v))
	case int32:
		return Number(float64(v))
	case int64:
		return Number(float64(v))
	case uint:
		return Number(float64(v))
	case uint8:
		return Number(float64(v))
	case uint16:
		return Number(float64(v))
	case uint32:
		return Number(float64(v))
	case uint64:
		return Number(float64(v))
	case float32:
		return Number(float64(v))
	case float64:
		return Number(v)
	case time.Time:
		return Time(v)
	case time.Duration:
		return Duration(v)
	case bool:
		return Bool(v)
	default:
		return String(fmt.Sprint(v))
	}
}
