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

var testObject = testTarget{
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
}

type testTarget map[string]any

func (t testTarget) GetField(key string) (any, error) {
	v, ok := t[key]
	if !ok {
		return nil, &FilterError{
			Kind: KindEval,
			Err:  fmt.Errorf("field not found: %q", key),
		}
	}
	return v, nil
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
		target   testTarget
		expected expected
	}{
		// String comparison
		{
			name:   "string eq",
			input:  `String=="HelloWorld"`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "string eq false",
			input:  `String=="X"`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "string neq",
			input:  `String!="X"`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "string eqi true",
			input:  `String==*"helloworld"`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "string eqi false",
			input:  `String==*"hellox"`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "string neqi true",
			input:  `String!=*"hellox"`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "string neqi false",
			input:  `String!=*"helloworld"`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "regex match",
			input:  `String=~"^Hello"`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "regex no match",
			input:  `String=~"world$"`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "regex neg match",
			input:  `String!~"^Hello"`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		// Numeric comparisons
		{
			name:   "int gt",
			input:  `Int>40`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "int gt false",
			input:  `Int>100`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "int eq",
			input:  `Int==42`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "int eq false",
			input:  `Int==41`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "int neq",
			input:  `Int!=41`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "int neq false",
			input:  `Int!=42`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "int gte false",
			input:  `Int>=100`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "int lt false",
			input:  `Int<40`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "int lte false",
			input:  `Int<=41`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "int8 gt",
			input:  `Int8>1`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "int16 gt",
			input:  `Int16>1`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:   "int32 gt",
			input:  `Int32>1`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:   "int64 gt",
			input:  `Int64>1`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:   "uint gt",
			input:  `Uint>1`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:   "uint8 gt",
			input:  `Uint8>1`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:   "uint16 gt",
			input:  `Uint16>1`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:   "uint32 gt",
			input:  `Uint32>1`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:   "uint64 gt",
			input:  `Uint64>1`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:   "float32 gt",
			input:  `Float32>2`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			}},
		{
			name:   "float lt",
			input:  `Float64<3.2`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "float gte",
			input:  `Float64>=3.14`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "float eq epsilon",
			input:  `Float64==3.1400000001`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "float neq epsilon",
			input:  `Float64!=3.1401`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		// Duration
		{
			name:   "duration gt",
			input:  `Duration>1s`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "duration gte false",
			input:  `Duration>=2s`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "duration gt false",
			input:  `Duration>2s`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "duration gte true",
			input:  `Duration>=1500ms`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "duration lt",
			input:  `Duration<2s`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "duration lt false",
			input:  `Duration<1s`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "duration lte true",
			input:  `Duration<=1500ms`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "duration lte false",
			input:  `Duration<=1s`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "duration eq",
			input:  `Duration==1500ms`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "duration eq false",
			input:  `Duration==2s`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "duration neq",
			input:  `Duration!=2s`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "invalid operator duration",
			input:  `Duration=~"1500ms"`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `invalid operator for duration`,
			},
		},
		{
			name:   "duration neq false",
			input:  `Duration!=1500ms`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "duration invalid at eval",
			input:  `Duration>bad`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `expected value, got identifier`,
			},
		},
		// Bool
		{
			name:   "bool eq",
			input:  `Bool==true`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "bool neq",
			input:  `Bool!=false`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "and true",
			input:  `Int>40&&Float64<4`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "and false",
			input:  `Int>40&&Float64>4`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "or true",
			input:  `Int>100||Float64<4`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "or short-circuit left true",
			input:  `Bool==true || InvalidField==1`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "not true->false",
			input:  `!(Int>40)`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "not false->true",
			input:  `!(Int<40)`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		{
			name:   "and short-circuit left false",
			input:  `Int>100 && InvalidField==1`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: false,
			},
		},
		{
			name:   "not inner eval error",
			input:  `!(InvalidField==1)`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		// Mixed
		{
			name:   "precedence",
			input:  `Int>40&&Float64<4||Bool==false`,
			target: testObject,
			expected: expected{
				ok:  true,
				val: true,
			},
		},
		// Errors
		{
			name:   "binary left eval error",
			input:  `UnknownField==1 && Bool==true`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:   "binary right eval error",
			input:  `Bool==true && UnknownField==1`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:   "invalid field",
			input:  `InvalidField==1`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:   "type mismatch 1",
			input:  `Int>"abc"`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:   "type mismatch 2",
			input:  `String>"HelloWorld"`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:   "type mismatch 3",
			input:  `Int=~"42"`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `eval error`,
			},
		},
		{
			name:   "invalid number rhs",
			input:  `Int>1+0`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `parse error`,
			},
		},
		{
			name:   "invalid duration rhs",
			input:  `Duration>1xs`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `parse error`,
			},
		},
		{
			name:   "regex compile error",
			input:  `String=~"["`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `parse error`,
			},
		},
		{
			name:   "regex not found",
			input:  `String=~""`,
			target: testObject,
			expected: expected{
				ok:  false,
				err: `parse error`,
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
				actual, evalError := expr.Eval(test.target)
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
				_, evalError := expr.Eval(test.target)
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
