package filter

import (
	"math"
	"testing"
	"time"
)

func Test_parseNumber_int64(t *testing.T) {
	type args struct {
		s string
	}
	type want struct {
		val   int64
		isErr bool
		err   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				s: "0",
			},
			want: want{
				val: 0,
			},
		},
		{
			name: "negative zero",
			args: args{
				s: "-0",
			},
			want: want{
				val: 0,
			},
		},
		{
			name: "leading plus",
			args: args{
				s: "+42",
			},
			want: want{
				val: 42,
			},
		},
		{
			name: "negative",
			args: args{
				s: "-42",
			},
			want: want{
				val: -42,
			},
		},
		{
			name: "leading zeros are decimal",
			args: args{
				s: "08",
			},
			want: want{
				val: 8,
			},
		},
		{
			name: "separators",
			args: args{
				s: "9_007_199_254_740_993",
			},
			want: want{
				val: 9007199254740993,
			},
		},
		{
			name: "negative separators",
			args: args{
				s: "-9_007_199_254_740_993",
			},
			want: want{
				val: -9007199254740993,
			},
		},
		{
			name: "minimum",
			args: args{
				s: "-9223372036854775808",
			},
			want: want{
				val: math.MinInt64,
			},
		},
		{
			name: "maximum",
			args: args{
				s: "9223372036854775807",
			},
			want: want{
				val: math.MaxInt64,
			},
		},
		{
			name: "below minimum",
			args: args{
				s: "-9223372036854775809",
			},
			want: want{
				isErr: true,
				err:   "value out of range",
			},
		},
		{
			name: "above maximum",
			args: args{
				s: "9223372036854775808",
			},
			want: want{
				isErr: true,
				err:   "value out of range",
			},
		},
		{
			name: "fraction rejected",
			args: args{
				s: "1.5",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "exponent rejected",
			args: args{
				s: "1e3",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "sign only",
			args: args{
				s: "-",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "text",
			args: args{
				s: "abc",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "hex integer rejected",
			args: args{
				s: "0x1f",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "binary integer rejected",
			args: args{
				s: "0b10",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "octal integer rejected",
			args: args{
				s: "0o10",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "repeated separator",
			args: args{
				s: "1__0",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "leading separator",
			args: args{
				s: "_1",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "trailing separator",
			args: args{
				s: "1_",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "minimum separators",
			args: args{
				s: "-9_223_372_036_854_775_808",
			},
			want: want{
				val: math.MinInt64,
			},
		},
		{
			name: "leading zeros and separators",
			args: args{
				s: "0_0_8",
			},
			want: want{
				val: 8,
			},
		},
		{
			name: "positive sign only",
			args: args{
				s: "+",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "only separator",
			args: args{
				s: "_",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "separator after sign",
			args: args{
				s: "+_1",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "overflow cannot wrap",
			args: args{
				s: "184467440737095516160000",
			},
			want: want{
				isErr: true,
				err:   "value out of range",
			},
		},
		{
			name: "unsigned maximum rejected",
			args: args{
				s: "18446744073709551615",
			},
			want: want{
				isErr: true,
				err:   "value out of range",
			},
		},
		{
			name: "unicode digit",
			args: args{
				s: "１",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "invalid separator before exponent",
			args: args{
				s: "1__0e2",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseNumber[int64](test.args.s)
			isErr := err != nil
			if isErr != test.want.isErr {
				t.Errorf("error mismatch\ngot=%v\nwant=%v\n", isErr, test.want.isErr)
				return
			}
			if isErr {
				if err.Error() != test.want.err {
					t.Errorf("error mismatch\ngot=%v\nwant=%v\n", err, test.want.err)
				}
				return
			}
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_parseNumber_uint64(t *testing.T) {
	type args struct {
		s string
	}
	type want struct {
		val   uint64
		isErr bool
		err   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				s: "0",
			},
			want: want{
				val: 0,
			},
		},
		{
			name: "leading plus",
			args: args{
				s: "+42",
			},
			want: want{
				val: 42,
			},
		},
		{
			name: "leading zeros are decimal",
			args: args{
				s: "08",
			},
			want: want{
				val: 8,
			},
		},
		{
			name: "separators",
			args: args{
				s: "18_446_744_073_709_551_615",
			},
			want: want{
				val: math.MaxUint64,
			},
		},
		{
			name: "signed boundary",
			args: args{
				s: "9223372036854775808",
			},
			want: want{
				val: 1 << 63,
			},
		},
		{
			name: "maximum",
			args: args{
				s: "18446744073709551615",
			},
			want: want{
				val: math.MaxUint64,
			},
		},
		{
			name: "negative rejected",
			args: args{
				s: "-1",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "negative zero rejected",
			args: args{
				s: "-0",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "above maximum",
			args: args{
				s: "18446744073709551616",
			},
			want: want{
				isErr: true,
				err:   "value out of range",
			},
		},
		{
			name: "fraction rejected",
			args: args{
				s: "1.5",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "exponent rejected",
			args: args{
				s: "1e3",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "sign only",
			args: args{
				s: "-",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "text",
			args: args{
				s: "abc",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "hex integer rejected",
			args: args{
				s: "0x1f",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "binary integer rejected",
			args: args{
				s: "0b10",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "octal integer rejected",
			args: args{
				s: "0o10",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "repeated separator",
			args: args{
				s: "1__0",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "leading separator",
			args: args{
				s: "_1",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "trailing separator",
			args: args{
				s: "1_",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "unsigned sign and separators",
			args: args{
				s: "+18_446_744_073_709_551_615",
			},
			want: want{
				val: math.MaxUint64,
			},
		},
		{
			name: "leading zeros and separators",
			args: args{
				s: "0_0_8",
			},
			want: want{
				val: 8,
			},
		},
		{
			name: "overflow cannot wrap",
			args: args{
				s: "184467440737095516160000",
			},
			want: want{
				isErr: true,
				err:   "value out of range",
			},
		},
		{
			name: "overflow separators",
			args: args{
				s: "18_446_744_073_709_551_616",
			},
			want: want{
				isErr: true,
				err:   "value out of range",
			},
		},
		{
			name: "positive sign only",
			args: args{
				s: "+",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "only separator",
			args: args{
				s: "_",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "separator after sign",
			args: args{
				s: "+_1",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "unicode digit",
			args: args{
				s: "１",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseNumber[uint64](test.args.s)
			isErr := err != nil
			if isErr != test.want.isErr {
				t.Errorf("error mismatch\ngot=%v\nwant=%v\n", isErr, test.want.isErr)
				return
			}
			if isErr {
				if err.Error() != test.want.err {
					t.Errorf("error mismatch\ngot=%v\nwant=%v\n", err, test.want.err)
				}
				return
			}
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_parseNumber_float64(t *testing.T) {
	type args struct {
		s string
	}
	type want struct {
		val   float64
		isErr bool
		err   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				s: "0.0",
			},
			want: want{
				val: 0,
			},
		},
		{
			name: "negative zero",
			args: args{
				s: "-0.0",
			},
			want: want{
				val: math.Copysign(0, -1),
			},
		},
		{
			name: "fraction",
			args: args{
				s: "1.5",
			},
			want: want{
				val: 1.5,
			},
		},
		{
			name: "negative fraction",
			args: args{
				s: "-1.5",
			},
			want: want{
				val: -1.5,
			},
		},
		{
			name: "leading plus",
			args: args{
				s: "+1.5",
			},
			want: want{
				val: 1.5,
			},
		},
		{
			name: "exponent",
			args: args{
				s: "1e3",
			},
			want: want{
				val: 1000,
			},
		},
		{
			name: "hex float",
			args: args{
				s: "0x1.fp3",
			},
			want: want{
				val: 15.5,
			},
		},
		{
			name: "separators",
			args: args{
				s: "1_000.5",
			},
			want: want{
				val: 1000.5,
			},
		},
		{
			name: "smallest positive",
			args: args{
				s: "5e-324",
			},
			want: want{
				val: math.SmallestNonzeroFloat64,
			},
		},
		{
			name: "maximum",
			args: args{
				s: "1.7976931348623157e308",
			},
			want: want{
				val: math.MaxFloat64,
			},
		},
		{
			name: "positive infinity",
			args: args{
				s: "Inf",
			},
			want: want{
				val: math.Inf(1),
			},
		},
		{
			name: "negative infinity",
			args: args{
				s: "-Inf",
			},
			want: want{
				val: math.Inf(-1),
			},
		},
		{
			name: "nan",
			args: args{
				s: "NaN",
			},
			want: want{
				val: math.NaN(),
			},
		},
		{
			name: "overflow",
			args: args{
				s: "1e400",
			},
			want: want{
				isErr: true,
				err:   "strconv.ParseFloat: parsing \"1e400\": value out of range",
			},
		},
		{
			name: "incomplete exponent",
			args: args{
				s: "1e",
			},
			want: want{
				isErr: true,
				err:   "strconv.ParseFloat: parsing \"1e\": invalid syntax",
			},
		},
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: want{
				isErr: true,
				err:   "strconv.ParseFloat: parsing \"\": invalid syntax",
			},
		},
		{
			name: "sign only",
			args: args{
				s: "-",
			},
			want: want{
				isErr: true,
				err:   "strconv.ParseFloat: parsing \"-\": invalid syntax",
			},
		},
		{
			name: "text",
			args: args{
				s: "abc",
			},
			want: want{
				isErr: true,
				err:   "strconv.ParseFloat: parsing \"abc\": invalid syntax",
			},
		},
		{
			name: "hex integer rejected",
			args: args{
				s: "0x1f",
			},
			want: want{
				isErr: true,
				err:   "strconv.ParseFloat: parsing \"0x1f\": invalid syntax",
			},
		},
		{
			name: "binary integer rejected",
			args: args{
				s: "0b10",
			},
			want: want{
				isErr: true,
				err:   "strconv.ParseFloat: parsing \"0b10\": invalid syntax",
			},
		},
		{
			name: "octal integer rejected",
			args: args{
				s: "0o10",
			},
			want: want{
				isErr: true,
				err:   "strconv.ParseFloat: parsing \"0o10\": invalid syntax",
			},
		},
		{
			name: "invalid fraction separator",
			args: args{
				s: "1_.0",
			},
			want: want{
				isErr: true,
				err:   "strconv.ParseFloat: parsing \"1_.0\": invalid syntax",
			},
		},
		{
			name: "integer rejected",
			args: args{
				s: "1",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "negative integer rejected",
			args: args{
				s: "-1",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "integer separators rejected",
			args: args{
				s: "1_000",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "unsigned overflow is not rounded",
			args: args{
				s: "18446744073709551616",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "fraction without integer",
			args: args{
				s: ".5",
			},
			want: want{
				val: 0.5,
			},
		},
		{
			name: "fraction without digits",
			args: args{
				s: "1.",
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "negative exponent",
			args: args{
				s: "1e-3",
			},
			want: want{
				val: 0.001,
			},
		},
		{
			name: "hex separators",
			args: args{
				s: "0x_1.fp3",
			},
			want: want{
				val: 15.5,
			},
		},
		{
			name: "explicit positive infinity",
			args: args{
				s: "+Inf",
			},
			want: want{
				val: math.Inf(1),
			},
		},
		{
			name: "long infinity",
			args: args{
				s: "Infinity",
			},
			want: want{
				val: math.Inf(1),
			},
		},
		{
			name: "lowercase nan",
			args: args{
				s: "nan",
			},
			want: want{
				val: math.NaN(),
			},
		},
		{
			name: "integer overflow followed by exponent",
			args: args{
				s: "18446744073709551616e0",
			},
			want: want{
				val: 0x1p64,
			},
		},
		{
			name: "integer overflow followed by fraction",
			args: args{
				s: "18446744073709551616.0",
			},
			want: want{
				val: 0x1p64,
			},
		},
		{
			name: "long integer followed by exponent",
			args: args{
				s: "1000000000000000000000000000000000000000e-37",
			},
			want: want{
				val: 100,
			},
		},
		{
			name: "negative integer overflow followed by fraction",
			args: args{
				s: "-9223372036854775809.0",
			},
			want: want{
				val: -0x1p63,
			},
		},
		{
			name: "repeated integer separator",
			args: args{
				s: "1__0",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "leading integer separator",
			args: args{
				s: "_1",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "trailing integer separator",
			args: args{
				s: "1_",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "only separator",
			args: args{
				s: "_",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "separator after sign",
			args: args{
				s: "+_1",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "negative integer overflow",
			args: args{
				s: "-9223372036854775809",
			},
			want: want{
				isErr: true,
				err:   "invalid syntax",
			},
		},
		{
			name: "positive sign only",
			args: args{
				s: "+",
			},
			want: want{
				isErr: true,
				err:   "strconv.ParseFloat: parsing \"+\": invalid syntax",
			},
		},
		{
			name: "invalid separator before exponent",
			args: args{
				s: "1__0e2",
			},
			want: want{
				isErr: true,
				err:   "strconv.ParseFloat: parsing \"1__0e2\": invalid syntax",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseNumber[float64](test.args.s)
			isErr := err != nil
			if isErr != test.want.isErr {
				t.Errorf("error mismatch\ngot=%v\nwant=%v\n", isErr, test.want.isErr)
				return
			}
			if isErr {
				if err.Error() != test.want.err {
					t.Errorf("error mismatch\ngot=%v\nwant=%v\n", err, test.want.err)
				}
				return
			}
			if math.Float64bits(got) != math.Float64bits(test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_parseTime(t *testing.T) {
	type args struct {
		s string
	}
	type want struct {
		val   time.Time
		isErr bool
		err   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "rfc3339 utc",
			args: args{
				s: "2025-01-01T00:00:00Z",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "rfc3339 offset",
			args: args{
				s: "2025-01-01T09:00:00+09:00",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "rfc3339 fraction",
			args: args{
				s: "2025-01-01T00:00:00.25Z",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 250000000, time.UTC),
			},
		},
		{
			name: "datetime without zone",
			args: args{
				s: "2025-01-01T00:00:00",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "datetime with space",
			args: args{
				s: "2025-01-01 00:00:00",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "datetime with space and fraction",
			args: args{
				s: "2025-01-01 00:00:00.5",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 500000000, time.UTC),
			},
		},
		{
			name: "date only",
			args: args{
				s: "2025-01-01",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "unix seconds",
			args: args{
				s: "1735689600",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "negative unix seconds",
			args: args{
				s: "-1",
			},
			want: want{
				val: time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name: "rfc1123",
			args: args{
				s: "Wed, 01 Jan 2025 00:00:00 UTC",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "rfc850",
			args: args{
				s: "Wednesday, 01-Jan-25 00:00:00 UTC",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "rfc822",
			args: args{
				s: "01 Jan 25 00:00 UTC",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "rfc1123 with numeric offset",
			args: args{
				s: "Wed, 01 Jan 2025 09:00:00 +0900",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "rfc822 with numeric offset",
			args: args{
				s: "01 Jan 25 09:00 +0900",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "gmt abbreviation",
			args: args{
				s: "Wed, 01 Jan 2025 00:00:00 GMT",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "unix seconds with underscores",
			args: args{
				s: "1_735_689_600",
			},
			want: want{
				val: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "zone abbreviation other than utc or gmt",
			args: args{
				s: "Wed, 01 Jan 2025 00:00:00 EST",
			},
			want: want{
				isErr: true,
				err:   `unknown time zone "EST"`,
			},
		},
		{
			name: "lowercase z",
			args: args{
				s: "2025-01-01T00:00:00z",
			},
			want: want{
				isErr: true,
				err:   `unrecognized time "2025-01-01T00:00:00z"`,
			},
		},
		{
			name: "out of range month",
			args: args{
				s: "2025-13-01",
			},
			want: want{
				isErr: true,
				err:   `unrecognized time "2025-13-01"`,
			},
		},
		{
			name: "clock only",
			args: args{
				s: "12:00:00",
			},
			want: want{
				isErr: true,
				err:   `unrecognized time "12:00:00"`,
			},
		},
		{
			name: "unix with fraction",
			args: args{
				s: "1735689600.5",
			},
			want: want{
				isErr: true,
				err:   `unrecognized time "1735689600.5"`,
			},
		},
		{
			name: "unix seconds overflow",
			args: args{
				s: "99999999999999999999",
			},
			want: want{
				isErr: true,
				err:   `unix seconds out of range "99999999999999999999"`,
			},
		},
		{
			name: "sign only",
			args: args{
				s: "-",
			},
			want: want{
				isErr: true,
				err:   `unrecognized time "-"`,
			},
		},
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: want{
				isErr: true,
				err:   `unrecognized time ""`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseTime(test.args.s)
			isErr := err != nil
			if isErr != test.want.isErr {
				t.Errorf("error mismatch\ngot=%v\nwant=%v\n", isErr, test.want.isErr)
				return
			}
			if isErr {
				if err.Error() != test.want.err {
					t.Errorf("error mismatch\ngot=%v\nwant=%v\n", err, test.want.err)
				}
				return
			}
			if !got.Equal(test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if got.Location() != time.UTC {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got.Location(), time.UTC)
			}
		})
	}
}
