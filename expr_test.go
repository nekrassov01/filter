package filter

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"testing"
	"time"
)

func Test_eval(t *testing.T) {
	type args struct {
		e *expr
		r Resolver
	}
	type want struct {
		val   bool
		isErr bool
		err   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "single identifier true",
			args: args{
				e: &MustParse(`Int==42`).expr,
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "single identifier false",
			args: args{
				e: &MustParse(`Int==41`).expr,
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "distinct identifiers skip the cache",
			args: args{
				e: &MustParse(`Int==42 && String=="HelloWorld"`).expr,
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "shared identifier within the inline cache",
			args: args{
				e: &MustParse(`Int>40 && Int<50 && Int!=0`).expr,
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "shared identifier resolved once",
			args: args{
				e: &MustParse(`Int==1 || Int==2 || Int==42`).expr,
				r: &countingResolver{
					testResolver: testObject,
					limit:        1,
					seen:         map[string]int{},
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "more identifiers than the inline cache",
			args: args{
				e: &MustParse(func() string {
					var b strings.Builder
					for i := range cacheSize + 1 {
						if i > 0 {
							b.WriteString(" && ")
						}
						fmt.Fprintf(&b, "F%d == %d && F%d != 0", i, i+1, i)
					}
					return b.String()
				}()).expr,
				r: func() testResolver {
					r := testResolver{}
					for i := range cacheSize + 1 {
						r[fmt.Sprintf("F%d", i)] = i + 1
					}
					return r
				}(),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "unknown identifier",
			args: args{
				e: &MustParse(`Unknown==1`).expr,
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:1: unknown identifier "Unknown"`,
			},
		},
		{
			name: "unknown identifier in nested expression",
			args: args{
				e: &MustParse(`Int==42 && (Bool==true || Unknown==1)`).expr,
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "unknown identifier reached after short circuit",
			args: args{
				e: &MustParse(`Int==42 && (Bool==false || Unknown==1)`).expr,
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:28: unknown identifier "Unknown"`,
			},
		},
		{
			name: "zero value from resolver",
			args: args{
				e: &MustParse(`Int==1`).expr,
				r: zeroResolver{},
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:1: unknown identifier "Int"`,
			},
		},
		{
			name: "zero value from resolver with shared identifier",
			args: args{
				e: &MustParse(`Int==1 || Int==2`).expr,
				r: zeroResolver{},
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:1: unknown identifier "Int"`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := eval(test.args.e, test.args.r)
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

func Test_evalNode(t *testing.T) {
	type args struct {
		nodes []node
		i     int32
		r     Resolver
		cache []cached
	}
	type want struct {
		val   bool
		isErr bool
		err   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "invalid logical operator",
			args: args{
				nodes: []node{
					{
						typ: nodeBinary,
						op: token{
							typ: tokenEQ,
						},
					},
				},
				i:     0,
				r:     nil,
				cache: nil,
			},
			want: want{
				val:   false,
				isErr: true,
				err:   `eval error: invalid logical operator "=="`,
			},
		},
		{
			name: "invalid node type",
			args: args{
				nodes: []node{
					{
						typ: nodeType(255),
					},
				},
				i:     0,
				r:     nil,
				cache: nil,
			},
			want: want{
				val:   false,
				isErr: true,
				err:   `eval error: invalid node type "error"`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := evalNode(test.args.nodes, test.args.i, test.args.r, test.args.cache)
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

func Test_evalPredicate(t *testing.T) {
	type args struct {
		n *node
		v Value
	}
	type want struct {
		val   bool
		isErr bool
		err   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "string",
			args: args{
				n: &node{
					typ: nodePredicate,
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenString,
						v:   "a",
					},
				},
				v: String("a"),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "number",
			args: args{
				n: &node{
					typ: nodePredicate,
					op: token{
						typ: tokenGT,
					},
					valFloat: 1,
					hasFloat: true,
				},
				v: Float64(2),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time",
			args: args{
				n: &node{
					typ: nodePredicate,
					op: token{
						typ: tokenLT,
					},
					valTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					hasTime: true,
				},
				v: Time(time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC)),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "duration",
			args: args{
				n: &node{
					typ: nodePredicate,
					op: token{
						typ: tokenNEQ,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: Duration(time.Second),
			},
			want: want{
				val: false,
			},
		},
		{
			name: "bool compares as string",
			args: args{
				n: &node{
					typ: nodePredicate,
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenBool,
						v:   "true",
					},
				},
				v: Bool(true),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "zero value",
			args: args{
				n: &node{
					typ: nodePredicate,
					ident: token{
						typ:  tokenIdent,
						v:    "HP",
						line: 1,
						col:  1,
					},
					op: token{
						typ: tokenEQ,
					},
				},
				v: Value{},
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:1: unknown identifier "HP"`,
			},
		},
		{
			name: "unexpected kind",
			args: args{
				n: &node{
					typ: nodePredicate,
					ident: token{
						typ:  tokenIdent,
						v:    "HP",
						line: 2,
						col:  3,
					},
					op: token{
						typ: tokenEQ,
					},
				},
				v: Value{
					kind: kind(255),
				},
			},
			want: want{
				isErr: true,
				err:   `eval error at 2:3: unknown identifier "HP"`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := evalPredicate(test.args.n, test.args.v)
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

func Test_evalString(t *testing.T) {
	type args struct {
		n *node
		v string
	}
	type want struct {
		val   bool
		isErr bool
		err   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "eq",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenString,
						v:   "Knight",
					},
				},
				v: "Knight",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "eq is case sensitive",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenString,
						v:   "Knight",
					},
				},
				v: "knight",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "eq empty strings",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenString,
					},
				},
				v: "",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "eq non-ascii",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenString,
						v:   "軍師",
					},
				},
				v: "軍師",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "neq",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					val: token{
						typ: tokenString,
						v:   "Knight",
					},
				},
				v: "Mage",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "neq same value",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					val: token{
						typ: tokenString,
						v:   "",
					},
				},
				v: "",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "regex match",
			args: args{
				n: &node{
					op: token{
						typ: tokenREQ,
					},
					re: regexp.MustCompile(`^K.*t$`),
				},
				v: "Knight",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "regex no match",
			args: args{
				n: &node{
					op: token{
						typ: tokenREQ,
					},
					re: regexp.MustCompile(`^K.*t$`),
				},
				v: "Mage",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "regex non-ascii",
			args: args{
				n: &node{
					op: token{
						typ: tokenREQ,
					},
					re: regexp.MustCompile(`^(諸葛亮|龐統)`),
				},
				v: "諸葛亮 孔明",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "regex empty input",
			args: args{
				n: &node{
					op: token{
						typ: tokenREQ,
					},
					re: regexp.MustCompile(`^$`),
				},
				v: "",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "not regex match",
			args: args{
				n: &node{
					op: token{
						typ: tokenNREQ,
					},
					re: regexp.MustCompile(`^K`),
				},
				v: "Knight",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "not regex no match",
			args: args{
				n: &node{
					op: token{
						typ: tokenNREQ,
					},
					re: regexp.MustCompile(`^K`),
				},
				v: "Mage",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "invalid operator",
			args: args{
				n: &node{
					op: token{
						typ:  tokenGT,
						line: 1,
						col:  5,
					},
					val: token{
						typ: tokenString,
						v:   "a",
					},
				},
				v: "b",
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:5: invalid operator for string value ">"`,
			},
		},
		{
			name: "logical operator",
			args: args{
				n: &node{
					op: token{
						typ:  tokenAND,
						line: 1,
						col:  3,
					},
				},
				v: "a",
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:3: invalid operator for string value "&&"`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := evalString(test.args.n, test.args.v)
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

func Test_evalNumber(t *testing.T) {
	type args struct {
		n    *node
		eval func(*node) (bool, error)
	}
	type want struct {
		val   bool
		isErr bool
		err   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "int64/gt",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valInt: 1,
					hasInt: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 2)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/gte equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenGTE,
					},
					valInt: 1,
					hasInt: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/lt",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valInt: 2,
					hasInt: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/lte equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenLTE,
					},
					valInt: 1,
					hasInt: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/eq",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valInt: 1,
					hasInt: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/neq",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					valInt: 1,
					hasInt: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 2)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/mixed fraction",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valFloat: 1.5,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/mixed equality is exact",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valFloat: 1.0000000001,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "int64/same value across float syntax",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/nan equality",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valFloat: math.NaN(),
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "int64/nan inequality",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					valFloat: math.NaN(),
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/nan ordering",
			args: args{
				n: &node{
					op: token{
						typ: tokenGTE,
					},
					valFloat: math.NaN(),
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "int64/negative versus unsigned",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valUint: 0,
					hasUint: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, -1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/positive versus unsigned",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valUint: 1,
					hasUint: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 2)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/equal versus unsigned",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valUint: 1,
					hasUint: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/unsigned boundary",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valUint: 1 << 63,
					hasUint: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, math.MaxInt64)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/large integer versus float",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valFloat: 9007199254740992,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 9007199254740993)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/uncached integer",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenString,
						v:   "42",
					},
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 42)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/uncached float",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenString,
						v:   "42.0",
					},
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 42)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64/invalid literal",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ:  tokenString,
						v:    "bad",
						line: 1,
						col:  6,
					},
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:6: invalid number "bad"`,
			},
		},
		{
			name: "int64/invalid operator",
			args: args{
				n: &node{
					op: token{
						typ:  tokenREQ,
						line: 1,
						col:  3,
					},
					valInt: 1,
					hasInt: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[int64](n, 1)
				},
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:3: invalid operator for number value "=~"`,
			},
		},
		{
			name: "uint64/gt",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valUint: 1,
					hasUint: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 2)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/gte equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenGTE,
					},
					valUint: 1,
					hasUint: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/lt",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valUint: 2,
					hasUint: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/lte equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenLTE,
					},
					valUint: 1,
					hasUint: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/eq",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valUint: 1,
					hasUint: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/neq",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					valUint: 1,
					hasUint: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 2)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/mixed fraction",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valFloat: 1.5,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/mixed equality is exact",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valFloat: 1.0000000001,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "uint64/same value across float syntax",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/nan equality",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valFloat: math.NaN(),
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "uint64/nan inequality",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					valFloat: math.NaN(),
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/nan ordering",
			args: args{
				n: &node{
					op: token{
						typ: tokenGTE,
					},
					valFloat: math.NaN(),
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "uint64/unsigned versus negative",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valInt: -1,
					hasInt: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 0)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/unsigned versus positive",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valInt: 1,
					hasInt: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 0)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/equal versus signed",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valInt: 1,
					hasInt: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/unsigned boundary",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valInt: math.MaxInt64,
					hasInt: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1<<63)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/unsigned maximum versus float",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valFloat: 0x1p64,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, math.MaxUint64)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/uncached integer",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenString,
						v:   "42",
					},
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 42)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/uncached float",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenString,
						v:   "42.0",
					},
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 42)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64/invalid literal",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ:  tokenString,
						v:    "bad",
						line: 1,
						col:  6,
					},
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:6: invalid number "bad"`,
			},
		},
		{
			name: "uint64/invalid operator",
			args: args{
				n: &node{
					op: token{
						typ:  tokenREQ,
						line: 1,
						col:  3,
					},
					valInt: 1,
					hasInt: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[uint64](n, 1)
				},
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:3: invalid operator for number value "=~"`,
			},
		},
		{
			name: "float64/cached gt",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 2)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/cached gt equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 1)
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "float64/gte equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenGTE,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/gte less",
			args: args{
				n: &node{
					op: token{
						typ: tokenGTE,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber(n, 0.5)
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "float64/lt",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, -1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/lt equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 1)
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "float64/lte equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenLTE,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/lte greater",
			args: args{
				n: &node{
					op: token{
						typ: tokenLTE,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber(n, 1.1)
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "float64/uncached literal parsed",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenString,
						v:   "1.5",
					},
				},
				eval: func(n *node) (bool, error) {
					return evalNumber(n, 1.5)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/uncached literal with exponent",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					val: token{
						typ: tokenString,
						v:   "1e3",
					},
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 1001)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/uncached literal invalid",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ:  tokenString,
						v:    "abc",
						line: 1,
						col:  6,
					},
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 1)
				},
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:6: invalid number "abc"`,
			},
		},
		{
			name: "float64/uncached literal empty",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ:  tokenString,
						line: 1,
						col:  6,
					},
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 0)
				},
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:6: invalid number ""`,
			},
		},
		{
			name: "float64/eq within epsilon",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valFloat: 1 + 1e-10,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/eq at epsilon",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valFloat: 0,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber(n, Epsilon)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/eq beyond epsilon",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valFloat: 1 + 1e-8,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 1)
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "float64/neq within epsilon",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					valFloat: 1 + 1e-10,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 1)
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "float64/neq beyond epsilon",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					valFloat: 1 + 1e-8,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 1)
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/negative zero equals zero",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valFloat: 0,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber(n, math.Copysign(0, -1))
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/nan eq",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber(n, math.NaN())
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "float64/nan neq is true",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber(n, math.NaN())
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/nan gt",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber(n, math.NaN())
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "float64/nan lte",
			args: args{
				n: &node{
					op: token{
						typ: tokenLTE,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber(n, math.NaN())
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "float64/positive infinity gt max",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valFloat: math.MaxFloat64,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber(n, math.Inf(1))
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/positive infinity eq itself",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valFloat: math.Inf(1),
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber(n, math.Inf(1))
				},
			},
			want: want{
				val: false,
			},
		},
		{
			name: "float64/negative infinity lt min",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valFloat: -math.MaxFloat64,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber(n, math.Inf(-1))
				},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float64/invalid operator",
			args: args{
				n: &node{
					op: token{
						typ:  tokenREQ,
						line: 1,
						col:  4,
					},
					valFloat: 1,
					hasFloat: true,
				},
				eval: func(n *node) (bool, error) {
					return evalNumber[float64](n, 1)
				},
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:4: invalid operator for number value "=~"`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.args.eval(test.args.n)
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

func Test_evalTime(t *testing.T) {
	type args struct {
		n *node
		v time.Time
	}
	type want struct {
		val   bool
		isErr bool
		err   string
	}
	epoch := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "not cached is parsed at evaluation",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenString,
						v:   "2025-01-01T00:00:00Z",
					},
				},
				v: epoch,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "not cached and invalid",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ:  tokenString,
						v:    "yesterday",
						line: 1,
						col:  8,
					},
				},
				v: epoch,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:8: invalid time "yesterday"`,
			},
		},
		{
			name: "gt after",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch.Add(time.Nanosecond),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "gt equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "gte equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenGTE,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "gte before",
			args: args{
				n: &node{
					op: token{
						typ: tokenGTE,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch.Add(-time.Nanosecond),
			},
			want: want{
				val: false,
			},
		},
		{
			name: "lt before",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch.Add(-time.Nanosecond),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "lt equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "lte equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenLTE,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "lte after",
			args: args{
				n: &node{
					op: token{
						typ: tokenLTE,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch.Add(time.Nanosecond),
			},
			want: want{
				val: false,
			},
		},
		{
			name: "eq same instant in another location",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch.In(time.FixedZone("JST", 9*3600)),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "eq nanosecond apart",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch.Add(time.Nanosecond),
			},
			want: want{
				val: false,
			},
		},
		{
			name: "neq nanosecond apart",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch.Add(time.Nanosecond),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "neq same instant",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "before unix epoch",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valTime: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
					hasTime: true,
				},
				v: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "far future",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "zero time",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: time.Time{},
			},
			want: want{
				val: true,
			},
		},
		{
			name: "invalid operator",
			args: args{
				n: &node{
					op: token{
						typ:  tokenREQ,
						line: 1,
						col:  5,
					},
					valTime: epoch,
					hasTime: true,
				},
				v: epoch,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:5: invalid operator for time value "=~"`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := evalTime(test.args.n, test.args.v)
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

func Test_evalDuration(t *testing.T) {
	type args struct {
		n *node
		v time.Duration
	}
	type want struct {
		val   bool
		isErr bool
		err   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "cached gt",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: time.Second + time.Nanosecond,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "cached gt equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: time.Second,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "gte equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenGTE,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: time.Second,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "gte less",
			args: args{
				n: &node{
					op: token{
						typ: tokenGTE,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: time.Second - time.Nanosecond,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "lt",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: time.Millisecond,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "lt equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: time.Second,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "lte equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenLTE,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: time.Second,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "lte greater",
			args: args{
				n: &node{
					op: token{
						typ: tokenLTE,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: time.Minute,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "eq",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valDuration: 1500 * time.Millisecond,
					hasDuration: true,
				},
				v: 1500 * time.Millisecond,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "eq nanosecond apart",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: time.Second + time.Nanosecond,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "neq",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: time.Second + time.Nanosecond,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "neq equal",
			args: args{
				n: &node{
					op: token{
						typ: tokenNEQ,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: time.Second,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "zero duration",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					valDuration: 0,
					hasDuration: true,
				},
				v: 0,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "negative duration lt zero",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					valDuration: 0,
					hasDuration: true,
				},
				v: -time.Second,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "negative literal",
			args: args{
				n: &node{
					op: token{
						typ: tokenGT,
					},
					valDuration: -time.Minute,
					hasDuration: true,
				},
				v: -time.Second,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uncached literal parsed",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ: tokenString,
						v:   "1h30m",
					},
				},
				v: 90 * time.Minute,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uncached literal fraction",
			args: args{
				n: &node{
					op: token{
						typ: tokenLT,
					},
					val: token{
						typ: tokenString,
						v:   "1.5s",
					},
				},
				v: time.Second,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uncached literal invalid",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ:  tokenString,
						v:    "bad",
						line: 1,
						col:  10,
					},
				},
				v: time.Second,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:10: invalid duration "bad"`,
			},
		},
		{
			name: "uncached literal without unit",
			args: args{
				n: &node{
					op: token{
						typ: tokenEQ,
					},
					val: token{
						typ:  tokenString,
						v:    "10",
						line: 1,
						col:  10,
					},
				},
				v: 10,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:10: invalid duration "10"`,
			},
		},
		{
			name: "invalid operator",
			args: args{
				n: &node{
					op: token{
						typ:  tokenNREQ,
						line: 1,
						col:  9,
					},
					valDuration: time.Second,
					hasDuration: true,
				},
				v: time.Second,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:9: invalid operator for duration value "!~"`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := evalDuration(test.args.n, test.args.v)
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

// countingResolver fails once an identifier is resolved more often than limit,
// proving that eval serves repeated identifiers from its cache.
type countingResolver struct {
	testResolver

	limit int
	seen  map[string]int
}

func (r *countingResolver) Resolve(name string) (Value, bool) {
	r.seen[name]++
	if r.seen[name] > r.limit {
		return Value{}, false
	}
	return r.testResolver.Resolve(name)
}
