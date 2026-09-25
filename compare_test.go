package filter

import (
	"math"
	"testing"
)

func Test_compareNumber(t *testing.T) {
	type args struct {
		compare func() (int, bool, bool)
	}
	type want struct {
		order   int
		equal   bool
		ordered bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "int64_int64/less",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, int64](-1, 0)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_int64/equal",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, int64](0, 0)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "int64_int64/greater",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, int64](1, 0)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_int64/full range",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, int64](math.MinInt64, math.MaxInt64)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_int64/large neighbors",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, int64](9007199254740993, 9007199254740992)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_uint64/negative signed",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, uint64](-1, 0)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_uint64/less",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, uint64](1, 2)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_uint64/equal",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, uint64](1, 1)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "int64_uint64/greater",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, uint64](2, 1)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_uint64/signed boundary",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, uint64](math.MaxInt64, 1<<63)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_uint64/full range",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, uint64](math.MinInt64, math.MaxUint64)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_float64/less than fraction",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64](1, 1.5)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_float64/greater than fraction",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64](2, 1.5)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_float64/equal",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, float64](1, 1)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "int64_float64/large integer retains low bit",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64, float64](9007199254740993, 9007199254740992)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_float64/signed minimum",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64](math.MinInt64, -0x1p63)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "int64_float64/exclusive upper bound",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64](math.MaxInt64, 0x1p63)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_float64/no epsilon for mixed types",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64](1, 1+1e-10)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_float64/negative zero",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64](0, math.Copysign(0, -1))
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "int64_float64/positive infinity",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64](math.MaxInt64, math.Inf(1))
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_float64/negative infinity",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64](math.MinInt64, math.Inf(-1))
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "int64_float64/nan",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[int64](0, math.NaN())
				},
			},
			want: want{
				order:   0,
				equal:   false,
				ordered: false,
			},
		},
		{
			name: "uint64_int64/negative signed",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64, int64](0, -1)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_int64/less",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64, int64](1, 2)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_int64/equal",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64, int64](1, 1)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "uint64_int64/greater",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64, int64](2, 1)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_int64/signed boundary",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64, int64](1<<63, math.MaxInt64)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_int64/full range",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64, int64](math.MaxUint64, math.MinInt64)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_uint64/less",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64, uint64](0, 1)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_uint64/equal",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64, uint64](math.MaxUint64, math.MaxUint64)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "uint64_uint64/greater",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64, uint64](1, 0)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_uint64/large neighbors",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64, uint64](math.MaxUint64, math.MaxUint64-1)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_float64/less than fraction",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64](1, 1.5)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_float64/greater than fraction",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64](2, 1.5)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_float64/equal",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64, float64](1, 1)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "uint64_float64/negative fraction",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64](0, -0.5)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_float64/large integer retains low bit",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64, float64](9007199254740993, 9007199254740992)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_float64/exclusive upper bound",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64](math.MaxUint64, 0x1p64)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_float64/no epsilon for mixed types",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64](1, 1+1e-10)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_float64/negative zero",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64](0, math.Copysign(0, -1))
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "uint64_float64/positive infinity",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64](math.MaxUint64, math.Inf(1))
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_float64/negative infinity",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64](0, math.Inf(-1))
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "uint64_float64/nan",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[uint64](0, math.NaN())
				},
			},
			want: want{
				order:   0,
				equal:   false,
				ordered: false,
			},
		},
		{
			name: "float64_int64/less than fraction",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, int64](1.5, 1)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_int64/greater than fraction",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, int64](1.5, 2)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_int64/equal",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, int64](1, 1)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "float64_int64/large integer retains low bit",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, int64](9007199254740992, 9007199254740993)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_int64/signed minimum",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, int64](-0x1p63, math.MinInt64)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "float64_int64/exclusive upper bound",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, int64](0x1p63, math.MaxInt64)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_int64/no epsilon for mixed types",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, int64](1+1e-10, 1)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_int64/negative zero",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, int64](math.Copysign(0, -1), 0)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "float64_int64/positive infinity",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, int64](math.Inf(1), math.MaxInt64)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_int64/negative infinity",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, int64](math.Inf(-1), math.MinInt64)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_int64/nan",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, int64](math.NaN(), 0)
				},
			},
			want: want{
				order:   0,
				equal:   false,
				ordered: false,
			},
		},
		{
			name: "float64_uint64/less than fraction",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, uint64](1.5, 1)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_uint64/greater than fraction",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, uint64](1.5, 2)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_uint64/equal",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, uint64](1, 1)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "float64_uint64/negative fraction",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, uint64](-0.5, 0)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_uint64/large integer retains low bit",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, uint64](9007199254740992, 9007199254740993)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_uint64/exclusive upper bound",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, uint64](0x1p64, math.MaxUint64)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_uint64/no epsilon for mixed types",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, uint64](1+1e-10, 1)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_uint64/negative zero",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, uint64](math.Copysign(0, -1), 0)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "float64_uint64/positive infinity",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, uint64](math.Inf(1), math.MaxUint64)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_uint64/negative infinity",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, uint64](math.Inf(-1), 0)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_uint64/nan",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, uint64](math.NaN(), 0)
				},
			},
			want: want{
				order:   0,
				equal:   false,
				ordered: false,
			},
		},
		{
			name: "float64_float64/less",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, float64](1, 2)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_float64/equal",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, float64](1, 1)
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "float64_float64/greater",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, float64](2, 1)
				},
			},
			want: want{
				order:   1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_float64/equal within epsilon",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64](1, 1+1e-10)
				},
			},
			want: want{
				order:   -1,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "float64_float64/outside epsilon",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64](1, 1+1e-8)
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_float64/signed zeros",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64](0, math.Copysign(0, -1))
				},
			},
			want: want{
				order:   0,
				equal:   true,
				ordered: true,
			},
		},
		{
			name: "float64_float64/left nan",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64, float64](math.NaN(), 1)
				},
			},
			want: want{
				order:   0,
				equal:   false,
				ordered: false,
			},
		},
		{
			name: "float64_float64/right nan",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber[float64](1, math.NaN())
				},
			},
			want: want{
				order:   0,
				equal:   false,
				ordered: false,
			},
		},
		{
			name: "float64_float64/both nan",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber(math.NaN(), math.NaN())
				},
			},
			want: want{
				order:   0,
				equal:   false,
				ordered: false,
			},
		},
		{
			name: "float64_float64/positive infinities remain unequal",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber(math.Inf(1), math.Inf(1))
				},
			},
			want: want{
				order:   0,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_float64/negative infinities remain unequal",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber(math.Inf(-1), math.Inf(-1))
				},
			},
			want: want{
				order:   0,
				equal:   false,
				ordered: true,
			},
		},
		{
			name: "float64_float64/opposite infinities",
			args: args{
				compare: func() (int, bool, bool) {
					return compareNumber(math.Inf(-1), math.Inf(1))
				},
			},
			want: want{
				order:   -1,
				equal:   false,
				ordered: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, equal, ordered := test.args.compare()
			got := want{
				order:   c,
				equal:   equal,
				ordered: ordered,
			}
			if got != test.want {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want)
			}
		})
	}
}

func Test_compareIntFloat(t *testing.T) {
	type args struct {
		i int64
		f float64
	}
	type want struct {
		val int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "less than fraction",
			args: args{
				i: 1,
				f: 1.5,
			},
			want: want{
				val: -1,
			},
		},
		{
			name: "greater than fraction",
			args: args{
				i: 2,
				f: 1.5,
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "equal",
			args: args{
				i: 1,
				f: 1.0,
			},
			want: want{
				val: 0,
			},
		},
		{
			name: "negative fraction",
			args: args{
				i: -1,
				f: -0.5,
			},
			want: want{
				val: -1,
			},
		},
		{
			name: "negative fractional remainder",
			args: args{
				i: 0,
				f: -0.5,
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "large integer retains low bit",
			args: args{
				i: 9007199254740993,
				f: 9007199254740992.0,
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "large negative integer retains low bit",
			args: args{
				i: -9007199254740993,
				f: -9007199254740992.0,
			},
			want: want{
				val: -1,
			},
		},
		{
			name: "signed minimum equals float",
			args: args{
				i: math.MinInt64,
				f: -0x1p63,
			},
			want: want{
				val: 0,
			},
		},
		{
			name: "float below signed minimum",
			args: args{
				i: math.MinInt64,
				f: math.Nextafter(-0x1p63, math.Inf(-1)),
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "float at exclusive signed upper bound",
			args: args{
				i: math.MaxInt64,
				f: 0x1p63,
			},
			want: want{
				val: -1,
			},
		},
		{
			name: "float below signed upper bound",
			args: args{
				i: math.MaxInt64,
				f: math.Nextafter(0x1p63, 0),
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "positive infinity",
			args: args{
				i: math.MaxInt64,
				f: math.Inf(1),
			},
			want: want{
				val: -1,
			},
		},
		{
			name: "negative infinity",
			args: args{
				i: math.MinInt64,
				f: math.Inf(-1),
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "positive subnormal",
			args: args{
				i: 0,
				f: math.SmallestNonzeroFloat64,
			},
			want: want{
				val: -1,
			},
		},
		{
			name: "negative subnormal",
			args: args{
				i: 0,
				f: -math.SmallestNonzeroFloat64,
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "negative zero",
			args: args{
				i: 0,
				f: math.Copysign(0, -1),
			},
			want: want{
				val: 0,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := compareIntFloat(test.args.i, test.args.f)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_compareUintFloat(t *testing.T) {
	type args struct {
		u uint64
		f float64
	}
	type want struct {
		val int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "less than fraction",
			args: args{
				u: 1,
				f: 1.5,
			},
			want: want{
				val: -1,
			},
		},
		{
			name: "greater than fraction",
			args: args{
				u: 2,
				f: 1.5,
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "equal",
			args: args{
				u: 1,
				f: 1.0,
			},
			want: want{
				val: 0,
			},
		},
		{
			name: "negative fraction",
			args: args{
				u: 0,
				f: -0.5,
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "large integer retains low bit",
			args: args{
				u: 9007199254740993,
				f: 9007199254740992.0,
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "signed boundary equals float",
			args: args{
				u: 1 << 63,
				f: 0x1p63,
			},
			want: want{
				val: 0,
			},
		},
		{
			name: "float at exclusive unsigned upper bound",
			args: args{
				u: math.MaxUint64,
				f: 0x1p64,
			},
			want: want{
				val: -1,
			},
		},
		{
			name: "float below unsigned upper bound",
			args: args{
				u: math.MaxUint64,
				f: math.Nextafter(0x1p64, 0),
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "positive infinity",
			args: args{
				u: math.MaxUint64,
				f: math.Inf(1),
			},
			want: want{
				val: -1,
			},
		},
		{
			name: "negative infinity",
			args: args{
				u: 0,
				f: math.Inf(-1),
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "positive subnormal",
			args: args{
				u: 0,
				f: math.SmallestNonzeroFloat64,
			},
			want: want{
				val: -1,
			},
		},
		{
			name: "negative subnormal",
			args: args{
				u: 0,
				f: -math.SmallestNonzeroFloat64,
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "negative zero",
			args: args{
				u: 0,
				f: math.Copysign(0, -1),
			},
			want: want{
				val: 0,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := compareUintFloat(test.args.u, test.args.f)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}
