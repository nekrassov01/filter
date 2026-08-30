package filter

import (
	"testing"
	"time"
)

// testObject resolves one identifier per Go type that ValueOf accepts.
var testObject = testResolver{
	"String":       "HelloWorld",
	"StringNumber": "123",
	"Int":          42,
	"Int8":         int8(5),
	"Int16":        int16(5),
	"Int32":        int32(5),
	"Int64":        int64(5),
	"Uint":         uint(5),
	"Uint8":        uint8(5),
	"Uint16":       uint16(5),
	"Uint32":       uint32(5),
	"Uint64":       uint64(5),
	"Float32":      float32(2.5),
	"Float64":      3.14,
	"Time":         time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	"Duration":     1500 * time.Millisecond,
	"Bool":         true,
	"Struct":       struct{ X int }{X: 1},
}

// testResolver resolves identifiers from a map through ValueOf.
type testResolver map[string]any

func (t testResolver) Resolve(name string) (Value, bool) {
	v, ok := t[name]
	if !ok {
		return Value{}, false
	}
	return ValueOf(v), true
}

// zeroResolver reports every identifier as known but resolves it to nothing.
type zeroResolver struct{}

func (zeroResolver) Resolve(string) (Value, bool) {
	return Value{}, true
}

// repr renders the expression tree of e as nested prefix groups for assertions.
func repr(e *Expr) string {
	val := func(t token) string {
		switch t.typ {
		case tokenNumber, tokenDuration, tokenBool:
			return t.v
		default:
			return "\"" + t.v + "\""
		}
	}
	var walk func(int32) string
	walk = func(i int32) string {
		n := e.nodes[i]
		switch n.typ {
		case nodeBinary:
			return "(" + walk(n.left) + " " + n.op.typ.literal() + " " + walk(n.right) + ")"
		case nodeNOT:
			return "(! " + walk(n.left) + ")"
		case nodeComparison:
			return "(" + n.ident.v + " " + n.op.typ.literal() + " " + val(n.val) + ")"
		default:
			return "<unknown>"
		}
	}
	return walk(e.root)
}

func Test_repr(t *testing.T) {
	type args struct {
		e *Expr
	}
	type want struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "comparison with quoted string",
			args: args{
				e: MustParse(`Name == "a"`),
			},
			want: want{
				val: `(Name == "a")`,
			},
		},
		{
			name: "comparison with number",
			args: args{
				e: MustParse(`HP > 50`),
			},
			want: want{
				val: `(HP > 50)`,
			},
		},
		{
			name: "comparison with duration",
			args: args{
				e: MustParse(`Latency <= 1.5s`),
			},
			want: want{
				val: `(Latency <= 1.5s)`,
			},
		},
		{
			name: "comparison with bool",
			args: args{
				e: MustParse(`Enabled != true`),
			},
			want: want{
				val: `(Enabled != true)`,
			},
		},
		{
			name: "comparison with time is quoted",
			args: args{
				e: MustParse(`Birth < 2025-01-01`),
			},
			want: want{
				val: `(Birth < "2025-01-01")`,
			},
		},
		{
			name: "nested binary and not",
			args: args{
				e: MustParse(`A == 1 || !(B == 2 && C == 3)`),
			},
			want: want{
				val: `((A == 1) || (! ((B == 2) && (C == 3))))`,
			},
		},
		{
			name: "unknown node type",
			args: args{
				e: &Expr{
					expr: expr{
						nodes: []node{
							{
								typ: nodeType(255),
							},
						},
					},
				},
			},
			want: want{
				val: "<unknown>",
			},
		},
		{
			name: "unknown node type nested in a binary node",
			args: args{
				e: &Expr{
					expr: expr{
						nodes: []node{
							{
								typ: nodeType(255),
							},
							{
								typ:   nodeBinary,
								left:  0,
								right: 0,
								op: token{
									typ: tokenAND,
								},
							},
						},
						root: 1,
					},
				},
			},
			want: want{
				val: "(<unknown> && <unknown>)",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := repr(test.args.e)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}
