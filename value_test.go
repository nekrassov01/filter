package filter

import (
	"errors"
	"math"
	"net/netip"
	"testing"
	"time"
)

func TestValue_addr(t *testing.T) {
	type want struct {
		val netip.Addr
	}
	tests := []struct {
		name string
		tr   Value
		want want
	}{
		{
			name: "invalid",
			tr: Value{
				kind: kindAddr,
				bits: 0,
				s:    "",
				a:    0,
				b:    0,
			},
			want: want{
				val: netip.Addr{},
			},
		},
		{
			name: "ipv4",
			tr: Value{
				kind: kindAddr,
				bits: 32,
				s:    "",
				a:    0,
				b:    0x0000ffffc0000201,
			},
			want: want{
				val: netip.MustParseAddr("192.0.2.1"),
			},
		},
		{
			name: "ipv4 unspecified",
			tr: Value{
				kind: kindAddr,
				bits: 32,
				s:    "",
				a:    0,
				b:    0x0000ffff00000000,
			},
			want: want{
				val: netip.MustParseAddr("0.0.0.0"),
			},
		},
		{
			name: "ipv6",
			tr: Value{
				kind: kindAddr,
				bits: 128,
				s:    "",
				a:    0x20010db800000000,
				b:    1,
			},
			want: want{
				val: netip.MustParseAddr("2001:db8::1"),
			},
		},
		{
			name: "ipv6 unspecified",
			tr: Value{
				kind: kindAddr,
				bits: 128,
				s:    "",
				a:    0,
				b:    0,
			},
			want: want{
				val: netip.IPv6Unspecified(),
			},
		},
		{
			name: "mapped ipv4",
			tr: Value{
				kind: kindAddr,
				bits: 128,
				s:    "",
				a:    0,
				b:    0x0000ffffc0000201,
			},
			want: want{
				val: netip.MustParseAddr("::ffff:192.0.2.1"),
			},
		},
		{
			name: "zone",
			tr: Value{
				kind: kindAddr,
				bits: 128,
				s:    "eth0",
				a:    -0x0180000000000000,
				b:    1,
			},
			want: want{
				val: netip.MustParseAddr("fe80::1%eth0"),
			},
		},
		{
			name: "high bits",
			tr: Value{
				kind: kindAddr,
				bits: 128,
				s:    "",
				a:    -1,
				b:    -1,
			},
			want: want{
				val: netip.MustParseAddr("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.tr.addr()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestString(t *testing.T) {
	type args struct {
		s string
	}
	type want struct {
		val Value
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "text",
			args: args{
				s: "Knight",
			},
			want: want{
				val: Value{
					kind: kindString,
					s:    "Knight",
				},
			},
		},
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: want{
				val: Value{
					kind: kindString,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := String(test.args.s)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestInt(t *testing.T) {
	type args struct {
		n int
	}
	type want struct {
		val Value
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive",
			args: args{
				n: 42,
			},
			want: want{
				val: Value{
					kind: kindInt64,
					a:    42,
				},
			},
		},
		{
			name: "negative",
			args: args{
				n: -42,
			},
			want: want{
				val: Value{
					kind: kindInt64,
					a:    -42,
				},
			},
		},
		{
			name: "zero",
			args: args{
				n: 0,
			},
			want: want{
				val: Value{
					kind: kindInt64,
					a:    0,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Int(test.args.n)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	type args struct {
		n int64
	}
	type want struct {
		val Value
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive",
			args: args{
				n: 42,
			},
			want: want{
				val: Value{
					kind: kindInt64,
					a:    42,
				},
			},
		},
		{
			name: "negative",
			args: args{
				n: -42,
			},
			want: want{
				val: Value{
					kind: kindInt64,
					a:    -42,
				},
			},
		},
		{
			name: "zero",
			args: args{
				n: 0,
			},
			want: want{
				val: Value{
					kind: kindInt64,
					a:    0,
				},
			},
		},
		{
			name: "above exact float range",
			args: args{
				n: 9007199254740993,
			},
			want: want{
				val: Value{
					kind: kindInt64,
					a:    9007199254740993,
				},
			},
		},
		{
			name: "minimum",
			args: args{
				n: math.MinInt64,
			},
			want: want{
				val: Value{
					kind: kindInt64,
					a:    math.MinInt64,
				},
			},
		},
		{
			name: "maximum",
			args: args{
				n: math.MaxInt64,
			},
			want: want{
				val: Value{
					kind: kindInt64,
					a:    math.MaxInt64,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Int64(test.args.n)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestUint64(t *testing.T) {
	type args struct {
		n uint64
	}
	type want struct {
		val Value
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive",
			args: args{
				n: 42,
			},
			want: want{
				val: Value{
					kind: kindUint64,
					a:    42,
				},
			},
		},
		{
			name: "zero",
			args: args{
				n: 0,
			},
			want: want{
				val: Value{
					kind: kindUint64,
					a:    0,
				},
			},
		},
		{
			name: "above exact float range",
			args: args{
				n: 9007199254740993,
			},
			want: want{
				val: Value{
					kind: kindUint64,
					a:    9007199254740993,
				},
			},
		},
		{
			name: "signed boundary",
			args: args{
				n: 1 << 63,
			},
			want: want{
				val: Value{
					kind: kindUint64,
					a:    math.MinInt64,
				},
			},
		},
		{
			name: "maximum",
			args: args{
				n: math.MaxUint64,
			},
			want: want{
				val: Value{
					kind: kindUint64,
					a:    -1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Uint64(test.args.n)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	type args struct {
		n float64
	}
	type want struct {
		val Value
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "integer",
			args: args{
				n: 42,
			},
			want: want{
				val: Value{
					kind: kindFloat64,
					a:    float64Bits(42),
				},
			},
		},
		{
			name: "fraction",
			args: args{
				n: 3.14,
			},
			want: want{
				val: Value{
					kind: kindFloat64,
					a:    float64Bits(3.14),
				},
			},
		},
		{
			name: "negative",
			args: args{
				n: -1.5,
			},
			want: want{
				val: Value{
					kind: kindFloat64,
					a:    float64Bits(-1.5),
				},
			},
		},
		{
			name: "zero",
			args: args{
				n: 0,
			},
			want: want{
				val: Value{
					kind: kindFloat64,
					a:    0,
				},
			},
		},
		{
			name: "largest",
			args: args{
				n: math.MaxFloat64,
			},
			want: want{
				val: Value{
					kind: kindFloat64,
					a:    float64Bits(math.MaxFloat64),
				},
			},
		},
		{
			name: "negative zero keeps its sign bit",
			args: args{
				n: math.Copysign(0, -1),
			},
			want: want{
				val: Value{
					kind: kindFloat64,
					a:    float64Bits(math.Copysign(0, -1)),
				},
			},
		},
		{
			name: "positive infinity",
			args: args{
				n: math.Inf(1),
			},
			want: want{
				val: Value{
					kind: kindFloat64,
					a:    float64Bits(math.Inf(1)),
				},
			},
		},
		{
			name: "nan keeps its bit pattern",
			args: args{
				n: math.NaN(),
			},
			want: want{
				val: Value{
					kind: kindFloat64,
					a:    float64Bits(math.NaN()),
				},
			},
		},
		{
			name: "smallest positive",
			args: args{
				n: math.SmallestNonzeroFloat64,
			},
			want: want{
				val: Value{
					kind: kindFloat64,
					a:    float64Bits(math.SmallestNonzeroFloat64),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Float64(test.args.n)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	type want struct {
		val Value
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "utc",
			args: args{
				t: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: want{
				val: Value{
					kind: kindTime,
					a:    1735689600,
				},
			},
		},
		{
			name: "nanoseconds",
			args: args{
				t: time.Date(2025, 1, 1, 0, 0, 0, 123456789, time.UTC),
			},
			want: want{
				val: Value{
					kind: kindTime,
					a:    1735689600,
					b:    123456789,
				},
			},
		},
		{
			name: "same instant in another location",
			args: args{
				t: time.Date(2025, 1, 1, 9, 0, 0, 0, time.FixedZone("JST", 9*3600)),
			},
			want: want{
				val: Value{
					kind: kindTime,
					a:    1735689600,
				},
			},
		},
		{
			name: "before the epoch",
			args: args{
				t: time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			want: want{
				val: Value{
					kind: kindTime,
					a:    -1,
				},
			},
		},
		{
			name: "zero time",
			args: args{
				t: time.Time{},
			},
			want: want{
				val: Value{
					kind: kindTime,
					a:    -62135596800,
				},
			},
		},
		{
			name: "far future",
			args: args{
				t: time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
			},
			want: want{
				val: Value{
					kind: kindTime,
					a:    253402300799,
					b:    999999999,
				},
			},
		},
		{
			name: "monotonic reading is dropped",
			args: args{
				t: time.Unix(1735689600, 0).Add(0),
			},
			want: want{
				val: Value{
					kind: kindTime,
					a:    1735689600,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Time(test.args.t)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	type args struct {
		d time.Duration
	}
	type want struct {
		val Value
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive",
			args: args{
				d: 1500 * time.Millisecond,
			},
			want: want{
				val: Value{
					kind: kindDuration,
					a:    int64(1500 * time.Millisecond),
				},
			},
		},
		{
			name: "negative",
			args: args{
				d: -time.Hour,
			},
			want: want{
				val: Value{
					kind: kindDuration,
					a:    int64(-time.Hour),
				},
			},
		},
		{
			name: "zero",
			args: args{
				d: 0,
			},
			want: want{
				val: Value{
					kind: kindDuration,
				},
			},
		},
		{
			name: "largest",
			args: args{
				d: time.Duration(math.MaxInt64),
			},
			want: want{
				val: Value{
					kind: kindDuration,
					a:    math.MaxInt64,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Duration(test.args.d)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestAddr(t *testing.T) {
	type args struct {
		ip netip.Addr
	}
	type want struct {
		val Value
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "invalid",
			args: args{
				ip: netip.Addr{},
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 0,
					s:    "",
					a:    0,
					b:    0,
				},
			},
		},
		{
			name: "ipv4",
			args: args{
				ip: netip.MustParseAddr("192.0.2.1"),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 32,
					s:    "",
					a:    0,
					b:    0x0000ffffc0000201,
				},
			},
		},
		{
			name: "ipv4 unspecified",
			args: args{
				ip: netip.MustParseAddr("0.0.0.0"),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 32,
					s:    "",
					a:    0,
					b:    0x0000ffff00000000,
				},
			},
		},
		{
			name: "ipv6",
			args: args{
				ip: netip.MustParseAddr("2001:db8::1"),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 128,
					s:    "",
					a:    0x20010db800000000,
					b:    1,
				},
			},
		},
		{
			name: "ipv6 unspecified",
			args: args{
				ip: netip.IPv6Unspecified(),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 128,
					s:    "",
					a:    0,
					b:    0,
				},
			},
		},
		{
			name: "mapped ipv4",
			args: args{
				ip: netip.MustParseAddr("::ffff:192.0.2.1"),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 128,
					s:    "",
					a:    0,
					b:    0x0000ffffc0000201,
				},
			},
		},
		{
			name: "zone",
			args: args{
				ip: netip.MustParseAddr("fe80::1%eth0"),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 128,
					s:    "eth0",
					a:    -0x0180000000000000,
					b:    1,
				},
			},
		},
		{
			name: "high bits",
			args: args{
				ip: netip.MustParseAddr("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 128,
					s:    "",
					a:    -1,
					b:    -1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Addr(test.args.ip)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestBool(t *testing.T) {
	type args struct {
		b bool
	}
	type want struct {
		val Value
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "true",
			args: args{
				b: true,
			},
			want: want{
				val: Value{
					kind: kindString,
					s:    "true",
				},
			},
		},
		{
			name: "false",
			args: args{
				b: false,
			},
			want: want{
				val: Value{
					kind: kindString,
					s:    "false",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Bool(test.args.b)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestValueOf(t *testing.T) {
	when := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	type args struct {
		v any
	}
	type want struct {
		val Value
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "string",
			args: args{
				v: "Knight",
			},
			want: want{
				val: String("Knight"),
			},
		},
		{
			name: "int",
			args: args{
				v: int(-7),
			},
			want: want{
				val: Int64(-7),
			},
		},
		{
			name: "int8",
			args: args{
				v: int8(5),
			},
			want: want{
				val: Int64(5),
			},
		},
		{
			name: "int16",
			args: args{
				v: int16(5),
			},
			want: want{
				val: Int64(5),
			},
		},
		{
			name: "int32",
			args: args{
				v: int32(5),
			},
			want: want{
				val: Int64(5),
			},
		},
		{
			name: "int64",
			args: args{
				v: int64(5),
			},
			want: want{
				val: Int64(5),
			},
		},
		{
			name: "uint",
			args: args{
				v: uint(5),
			},
			want: want{
				val: Uint64(5),
			},
		},
		{
			name: "uint8",
			args: args{
				v: uint8(5),
			},
			want: want{
				val: Uint64(5),
			},
		},
		{
			name: "uint16",
			args: args{
				v: uint16(5),
			},
			want: want{
				val: Uint64(5),
			},
		},
		{
			name: "uint32",
			args: args{
				v: uint32(5),
			},
			want: want{
				val: Uint64(5),
			},
		},
		{
			name: "uint64",
			args: args{
				v: uint64(5),
			},
			want: want{
				val: Uint64(5),
			},
		},
		{
			name: "float32",
			args: args{
				v: float32(2.5),
			},
			want: want{
				val: Float64(2.5),
			},
		},
		{
			name: "float64",
			args: args{
				v: 3.14,
			},
			want: want{
				val: Float64(3.14),
			},
		},
		{
			name: "time",
			args: args{
				v: when,
			},
			want: want{
				val: Time(when),
			},
		},
		{
			name: "duration",
			args: args{
				v: 1500 * time.Millisecond,
			},
			want: want{
				val: Duration(1500 * time.Millisecond),
			},
		},
		{
			name: "addr invalid",
			args: args{
				v: netip.Addr{},
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 0,
					s:    "",
					a:    0,
					b:    0,
				},
			},
		},
		{
			name: "addr ipv4",
			args: args{
				v: netip.MustParseAddr("192.0.2.1"),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 32,
					s:    "",
					a:    0,
					b:    0x0000ffffc0000201,
				},
			},
		},
		{
			name: "addr ipv4 unspecified",
			args: args{
				v: netip.MustParseAddr("0.0.0.0"),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 32,
					s:    "",
					a:    0,
					b:    0x0000ffff00000000,
				},
			},
		},
		{
			name: "addr ipv6",
			args: args{
				v: netip.MustParseAddr("2001:db8::1"),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 128,
					s:    "",
					a:    0x20010db800000000,
					b:    1,
				},
			},
		},
		{
			name: "addr ipv6 unspecified",
			args: args{
				v: netip.IPv6Unspecified(),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 128,
					s:    "",
					a:    0,
					b:    0,
				},
			},
		},
		{
			name: "addr mapped ipv4",
			args: args{
				v: netip.MustParseAddr("::ffff:192.0.2.1"),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 128,
					s:    "",
					a:    0,
					b:    0x0000ffffc0000201,
				},
			},
		},
		{
			name: "addr zone",
			args: args{
				v: netip.MustParseAddr("fe80::1%eth0"),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 128,
					s:    "eth0",
					a:    -0x0180000000000000,
					b:    1,
				},
			},
		},
		{
			name: "addr high bits",
			args: args{
				v: netip.MustParseAddr("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
			},
			want: want{
				val: Value{
					kind: kindAddr,
					bits: 128,
					s:    "",
					a:    -1,
					b:    -1,
				},
			},
		},
		{
			name: "bool",
			args: args{
				v: true,
			},
			want: want{
				val: Bool(true),
			},
		},
		{
			name: "other type is formatted",
			args: args{
				v: struct{ X int }{X: 1},
			},
			want: want{
				val: String("{1}"),
			},
		},
		{
			name: "nil is formatted",
			args: args{
				v: nil,
			},
			want: want{
				val: String("<nil>"),
			},
		},
		{
			name: "largest int64",
			args: args{
				v: int64(math.MaxInt64),
			},
			want: want{
				val: Int64(math.MaxInt64),
			},
		},
		{
			name: "largest uint64 stays exact",
			args: args{
				v: uint64(math.MaxUint64),
			},
			want: want{
				val: Uint64(math.MaxUint64),
			},
		},
		{
			name: "empty string",
			args: args{
				v: "",
			},
			want: want{
				val: String(""),
			},
		},
		{
			name: "byte slice is formatted",
			args: args{
				v: []byte("ab"),
			},
			want: want{
				val: String("[97 98]"),
			},
		},
		{
			name: "error is formatted",
			args: args{
				v: errors.New("boom"),
			},
			want: want{
				val: String("boom"),
			},
		},
		{
			name: "stringer is formatted",
			args: args{
				v: time.Month(1),
			},
			want: want{
				val: String("January"),
			},
		},
		{
			name: "typed nil pointer is formatted",
			args: args{
				v: (*int)(nil),
			},
			want: want{
				val: String("<nil>"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ValueOf(test.args.v)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

// float64Bits returns the bit pattern of f as Float64 stores it.
func float64Bits(f float64) int64 {
	//nolint:gosec // bit pattern conversion
	return int64(math.Float64bits(f))
}
