package filter

import (
	"encoding/binary"
	"fmt"
	"math"
	"net/netip"
	"strconv"
	"time"
)

// kind identifies how a Value's fields are read.
type kind uint8

// Kinds of resolved values.
const (
	kindNone     kind = iota // zero Value: nothing resolved
	kindString               // s
	kindInt64                // a holds the signed integer
	kindUint64               // a holds the unsigned integer bits
	kindFloat64              // a holds the float64 bits
	kindTime                 // a holds Unix seconds, b the nanoseconds
	kindDuration             // a holds the duration
	kindAddr                 // a and b hold address bits, s the zone, bits the address family
)

// Value is the value of an identifier, as returned by a Resolver.
type Value struct {
	s    string
	a    int64
	b    int64
	kind kind
	bits uint8 // Addr only: 0 = invalid, 32 = IPv4, 128 = IPv6
}

// addr reconstructs an address without conflating IPv4 and mapped IPv6 values.
func (v Value) addr() netip.Addr {
	if v.bits == 0 {
		return netip.Addr{}
	}
	var b [16]byte
	//nolint:gosec // bit pattern conversion
	binary.BigEndian.PutUint64(b[:8], uint64(v.a))
	//nolint:gosec // bit pattern conversion
	binary.BigEndian.PutUint64(b[8:], uint64(v.b))
	ip := netip.AddrFrom16(b)
	if v.bits == 32 {
		return ip.Unmap()
	}
	return ip.WithZone(v.s)
}

// String returns a Value holding s.
func String(s string) Value {
	return Value{
		kind: kindString,
		s:    s,
	}
}

// Int returns a Value holding n.
func Int(n int) Value {
	return Int64(int64(n))
}

// Int64 returns a Value holding n.
func Int64(n int64) Value {
	return Value{
		kind: kindInt64,
		a:    n,
	}
}

// Uint64 returns a Value holding n.
func Uint64(n uint64) Value {
	//nolint:gosec // bit pattern conversion
	return Value{
		kind: kindUint64,
		a:    int64(n),
	}
}

// Float64 returns a Value holding n.
func Float64(n float64) Value {
	//nolint:gosec // bit pattern conversion
	return Value{
		kind: kindFloat64,
		a:    int64(math.Float64bits(n)),
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

// Duration returns a Value holding d.
func Duration(d time.Duration) Value {
	return Value{
		kind: kindDuration,
		a:    int64(d),
	}
}

// Addr returns a Value holding ip, preserving its address family and zone.
// A zero netip.Addr remains invalid and sorts before valid addresses.
func Addr(ip netip.Addr) Value {
	b := ip.As16()
	//nolint:gosec // bit pattern conversion
	return Value{
		kind: kindAddr,
		bits: uint8(ip.BitLen()),
		s:    ip.Zone(),
		a:    int64(binary.BigEndian.Uint64(b[:8])),
		b:    int64(binary.BigEndian.Uint64(b[8:])),
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
// time.Time, time.Duration, and netip.Addr keep their kind. Booleans compare as
// the strings "true" or "false"; any other value is formatted with fmt.Sprint
// and compared as a string.
func ValueOf(v any) Value {
	switch v := v.(type) {
	case string:
		return String(v)
	case int:
		return Int(v)
	case int8:
		return Int64(int64(v))
	case int16:
		return Int64(int64(v))
	case int32:
		return Int64(int64(v))
	case int64:
		return Int64(v)
	case uint:
		return Uint64(uint64(v))
	case uint8:
		return Uint64(uint64(v))
	case uint16:
		return Uint64(uint64(v))
	case uint32:
		return Uint64(uint64(v))
	case uint64:
		return Uint64(v)
	case float32:
		return Float64(float64(v))
	case float64:
		return Float64(v)
	case time.Time:
		return Time(v)
	case time.Duration:
		return Duration(v)
	case netip.Addr:
		return Addr(v)
	case bool:
		return Bool(v)
	default:
		return String(fmt.Sprint(v))
	}
}
