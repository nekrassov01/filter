package filter

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var testTemplate = `ERROR:
input:    %v
expected: %v
actual:   %v
`

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

// zeroResolver reports every identifier as known but resolves it to nothing.
type zeroResolver struct{}

func (zeroResolver) Resolve(string) (Value, bool) {
	return Value{}, true
}

type testResolver map[string]any

func (t testResolver) Resolve(name string) (Value, bool) {
	v, ok := t[name]
	if !ok {
		return Value{}, false
	}
	return ValueOf(v), true
}

func TestEval(t *testing.T) {
	type expected struct {
		ok  bool
		val bool
		err string
	}
	tests := []struct {
		name     string
		input    string
		resolver testResolver
		expected expected
	}{
		// String comparison
		{
			name:     "string eq",
			input:    `String=="HelloWorld"`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "string eq false",
			input:    `String=="X"`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "string neq",
			input:    `String!="X"`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "regex ignoring case",
			input:    `String=~"(?i)^helloworld$"`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "regex ignoring case false",
			input:    `String=~"(?i)^hellox$"`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "negative regex ignoring case",
			input:    `String!~"(?i)^hellox$"`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "negative regex ignoring case false",
			input:    `String!~"(?i)^helloworld$"`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "regex match",
			input:    `String=~"^Hello"`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "regex no match",
			input:    `String=~"world$"`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "regex neg match",
			input:    `String!~"^Hello"`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		// Numeric comparisons
		{
			name:     "int gt",
			input:    `Int>40`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "int gt false",
			input:    `Int>100`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "int eq",
			input:    `Int==42`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "int eq false",
			input:    `Int==41`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "int neq",
			input:    `Int!=41`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "int neq false",
			input:    `Int!=42`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "int gte false",
			input:    `Int>=100`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "int lt false",
			input:    `Int<40`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "int lte false",
			input:    `Int<=41`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "int8 gt",
			input:    `Int8>1`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "int16 gt",
			input:    `Int16>1`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:     "int32 gt",
			input:    `Int32>1`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:     "int64 gt",
			input:    `Int64>1`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:     "uint gt",
			input:    `Uint>1`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:     "uint8 gt",
			input:    `Uint8>1`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:     "uint16 gt",
			input:    `Uint16>1`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:     "uint32 gt",
			input:    `Uint32>1`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:     "uint64 gt",
			input:    `Uint64>1`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:     "float32 gt",
			input:    `Float32>2`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:     "float lt",
			input:    `Float64<3.2`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "float gte",
			input:    `Float64>=3.14`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "float eq epsilon",
			input:    `Float64==3.1400000001`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "float neq epsilon",
			input:    `Float64!=3.1401`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		// Duration
		{
			name:     "duration gt",
			input:    `Duration>1s`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "duration gte false",
			input:    `Duration>=2s`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "duration gt false",
			input:    `Duration>2s`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "duration gte true",
			input:    `Duration>=1500ms`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "duration lt",
			input:    `Duration<2s`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "duration lt false",
			input:    `Duration<1s`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "duration lte true",
			input:    `Duration<=1500ms`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "duration lte false",
			input:    `Duration<=1s`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "duration eq",
			input:    `Duration==1500ms`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "duration eq false",
			input:    `Duration==2s`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "duration neq",
			input:    `Duration!=2s`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "invalid operator duration",
			input:    `Duration=~"1500ms"`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `invalid operator for duration`,
			},
		},
		{
			name:     "duration neq false",
			input:    `Duration!=1500ms`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "other types compare by their formatted text",
			input:    `Struct=="{1}"`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "string of a lone sign stays a string",
			input:    `String!="-"`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "duration string converted at eval",
			input:    `Duration>'0'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "number string converted at eval",
			input:    `Float64<'Inf'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "duration invalid at eval",
			input:    `Duration>bad`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `expected value, got identifier`,
			},
		},
		// Time
		{
			name:     "time gt",
			input:    `Time>'2024-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time gt false",
			input:    `Time>'2026-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "time gte",
			input:    `Time>='2025-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time gte false",
			input:    `Time>='2026-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "time lt",
			input:    `Time<'2026-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time lt false",
			input:    `Time<'2024-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "time lte",
			input:    `Time<='2025-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time lte false",
			input:    `Time<='2024-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "time eq",
			input:    `Time=='2025-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time eq false",
			input:    `Time=='2024-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "time neq",
			input:    `Time!='2024-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time neq false",
			input:    `Time!='2025-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		// Combined logicals
		{
			name:     "combined logicals",
			input:    `String=="HelloWorld" && Int==42 || Float64<3.0`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "combined logicals false",
			input:    `String=="HelloWorld" && Int==41 || Float64<3.0`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "combined logicals with parens",
			input:    `(String=="HelloWorld" && Int==41) || (Float64<3.2 && Bool==true)`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "combined logicals with parens false",
			input:    `(String=="HelloWorld" && Int==41) || (Float64<3.0 && Bool==true)`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		// Bool
		{
			name:     "bool eq",
			input:    `Bool==true`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "bool neq",
			input:    `Bool!=false`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "and true",
			input:    `Int>40&&Float64<4`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "and false",
			input:    `Int>40&&Float64>4`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "or true",
			input:    `Int>100||Float64<4`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "or short-circuit left true",
			input:    `Bool==true || Invalid==1`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "or left error",
			input:    `Invalid==1 || Bool==true`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error at 1:1: unknown identifier "Invalid"`,
			},
		},
		{
			name:     "same identifier referenced twice",
			input:    `Int>40 && Int<50`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name: "more identifiers than the inline cache",
			input: func() string {
				var b strings.Builder
				for i := range 17 {
					if i > 0 {
						b.WriteString(" && ")
					}
					fmt.Fprintf(&b, "F%d == %d", i, i)
				}
				return b.String() + " && F0 == 0"
			}(),
			resolver: func() testResolver {
				t := testResolver{}
				for i := range 17 {
					t[fmt.Sprintf("F%d", i)] = i
				}
				return t
			}(),
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "not true->false",
			input:    `!(Int>40)`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "not false->true",
			input:    `!(Int<40)`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "and short-circuit left false",
			input:    `Int>100 && Invalid==1`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:     "not inner eval error",
			input:    `!(Invalid==1)`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		// Mixed
		{
			name:     "precedence",
			input:    `Int>40&&Float64<4||Bool==false`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		// Errors
		{
			name:     "binary left eval error",
			input:    `Unknown==1 && Bool==true`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:     "binary right eval error",
			input:    `Bool==true && Unknown==1`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error at 1:15: unknown identifier "Unknown"`,
			},
		},
		{
			name:     "unknown identifier",
			input:    `Invalid==1`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:     "type mismatch 1",
			input:    `Int>"abc"`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:     "type mismatch 2",
			input:    `String>"HelloWorld"`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:     "type mismatch 3",
			input:    `Int=~"42"`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:     "invalid number right",
			input:    `Int>1+0`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `parse error`,
			},
		},
		{
			name:     "invalid duration right",
			input:    `Duration>1xs`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `parse error`,
			},
		},
		{
			name:     "regex compile error",
			input:    `String=~"["`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `parse error`,
			},
		},
		{
			name:     "regex not found",
			input:    `String=~""`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `parse error`,
			},
		},
		{
			name:     "time eq bare rfc3339",
			input:    `Time==2025-01-01T00:00:00Z`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time eq bare datetime without zone",
			input:    `Time==2025-01-01T00:00:00`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time eq bare date",
			input:    `Time==2025-01-01`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time eq date only",
			input:    `Time=='2025-01-01'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time gt date only",
			input:    `Time>'2024-12-31'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time eq datetime with space",
			input:    `Time=='2025-01-01 00:00:00'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time eq datetime without zone",
			input:    `Time=='2025-01-01T00:00:00'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time lt datetime with fraction and no zone",
			input:    `Time<'2025-01-01 00:00:00.5'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time eq rfc1123 string",
			input:    `Time=='Wed, 01 Jan 2025 00:00:00 UTC'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time eq rfc1123z string",
			input:    `Time=='Wed, 01 Jan 2025 09:00:00 +0900'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time eq rfc850 string",
			input:    `Time=='Wednesday, 01-Jan-25 00:00:00 UTC'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time eq rfc822 string",
			input:    `Time=='01 Jan 25 00:00 UTC'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time eq unix seconds with underscores",
			input:    `Time==1_735_689_600`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time with zone abbreviation other than utc",
			input:    `Time=='Wed, 01 Jan 2025 00:00:00 EST'`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error at 1:7: invalid time "Wed, 01 Jan 2025 00:00:00 EST"`,
			},
		},
		{
			name:     "time eq unix seconds as string",
			input:    `Time=='1735689600'`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time eq unix seconds as number",
			input:    `Time==1735689600`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time lt unix seconds as number",
			input:    `Time<1735689601`,
			resolver: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:     "time with fractional unix seconds",
			input:    `Time==1735689600.5`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error at 1:7: invalid time "1735689600.5"`,
			},
		},
		{
			name:     "time with out of range date",
			input:    `Time>'2025-13-01'`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error at 1:6: invalid time "2025-13-01"`,
			},
		},
		{
			name:     "time with clock only",
			input:    `Time>'12:00:00'`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `invalid time "12:00:00"`,
			},
		},
		{
			name:     "invalid time",
			input:    `Time>'invalid-time'`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:     "time invalid operator",
			input:    `Time=~'2025-01-01T00:00:00Z'`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:     "invalid duration",
			input:    `Duration>'bad-duration'`,
			resolver: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, parseError := Parse(test.input)
			if test.expected.ok {
				if parseError != nil {
					t.Errorf(testTemplate, test.input, "", parseError)
					return
				}
				actual, evalError := expr.Eval(test.resolver)
				if evalError != nil {
					t.Errorf(testTemplate, test.input, test.expected.val, evalError)
					return
				}
				if actual != test.expected.val {
					t.Errorf(testTemplate, test.input, test.expected.val, actual)
				}
				return
			}
			if parseError == nil {
				_, evalError := expr.Eval(test.resolver)
				if evalError == nil || !strings.Contains(evalError.Error(), test.expected.err) {
					t.Errorf(testTemplate, test.input, test.expected.err, evalError)
				}
				return
			}
			if !strings.Contains(parseError.Error(), test.expected.err) {
				t.Errorf(testTemplate, test.input, test.expected.err, parseError)
			}
		})
	}
}

func Test_eval(t *testing.T) {
	tests := []struct {
		name     string
		nodes    []node
		expected string
	}{
		{
			name: "invalid logical operator",
			nodes: []node{
				{
					typ: nodeBinary,
					op: token{
						typ: tokenEQ,
					},
				},
			},
			expected: "invalid logical operator",
		},
		{
			name: "invalid node type",
			nodes: []node{
				{
					typ: nodeType(255),
				},
			},
			expected: "invalid node type",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := eval(test.nodes, 0, testObject, nil)
			if err == nil || !strings.Contains(err.Error(), test.expected) {
				t.Errorf(testTemplate, test.nodes, test.expected, err)
			}
		})
	}
}

func TestEval_zeroValue(t *testing.T) {
	expr, err := Parse(`Int==1`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = expr.Eval(zeroResolver{})
	if err == nil || !strings.Contains(err.Error(), `unknown identifier "Int"`) {
		t.Errorf(testTemplate, `Int==1`, `unknown identifier "Int"`, err)
	}
}
