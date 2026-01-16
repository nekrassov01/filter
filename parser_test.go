package filter

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	type expected struct {
		ok   bool
		repr string
		err  string
	}
	tests := []struct {
		name     string
		input    string
		expected expected
	}{
		// Strings
		{
			name:  "eq string",
			input: `Class=="軍師"`,
			expected: expected{
				ok:   true,
				repr: `(Class == "軍師")`,
			},
		},
		{
			name:  "neq raw string",
			input: "Name!='孔明'",
			expected: expected{
				ok:   true,
				repr: `(Name != "孔明")`,
			},
		},
		{
			name:  "eqi string",
			input: `Tag==*"Admin"`,
			expected: expected{
				ok:   true,
				repr: `(Tag ==* Admin)`,
			},
		},
		{
			name:  "regex",
			input: `Name=~'A.*'`,
			expected: expected{
				ok:   true,
				repr: `(Name =~ "A.*")`,
			},
		},
		{
			name:  "regex not",
			input: `Name!~'A.*'`,
			expected: expected{
				ok:   true,
				repr: `(Name !~ "A.*")`,
			},
		},
		{
			name:  "regex raw string",
			input: "Name=~`^H.*d$`",
			expected: expected{
				ok:   true,
				repr: `(Name =~ "^H.*d$")`,
			},
		},
		{
			name:  "regex case insensitive",
			input: "Name=~*`A`",
			expected: expected{
				ok:   true,
				repr: `(Name =~* "(?i)A")`,
			},
		},
		{
			name:  "regex case insensitive not",
			input: "Name!~*`A`",
			expected: expected{
				ok:   true,
				repr: `(Name !~* "(?i)A")`,
			},
		},
		// Numbers
		{
			name:  "gt number",
			input: `HP>50`,
			expected: expected{
				ok:   true,
				repr: `(HP > 50)`,
			},
		},
		{
			name:  "gte number",
			input: `MP>=100`,
			expected: expected{
				ok:   true,
				repr: `(MP >= 100)`,
			},
		},
		{
			name:  "lt number float",
			input: `Rate<0.75`,
			expected: expected{
				ok:   true,
				repr: `(Rate < 0.75)`,
			},
		},
		{
			name:  "hex float",
			input: `X==0x1.fp3`,
			expected: expected{
				ok:   true,
				repr: `(X == 0x1.fp3)`,
			},
		},
		// Durations
		{
			name:  "duration gte",
			input: `Delay>=1h30m`,
			expected: expected{
				ok:   true,
				repr: `(Delay >= 1h30m)`,
			},
		},
		{
			name:  "duration lt",
			input: `Timeout<500ms`,
			expected: expected{
				ok:   true,
				repr: `(Timeout < 500ms)`,
			},
		},
		{
			name:  "duration micro",
			input: `Mic==4000μs`,
			expected: expected{
				ok:   true,
				repr: `(Mic == 4000μs)`,
			},
		},
		// Times
		{
			name:  "time before",
			input: `Time>=2023-01-02T15:04:05Z`,
			expected: expected{
				ok:   true,
				repr: `(Time >= "2023-01-02T15:04:05Z")`,
			},
		},
		{
			name:  "time after",
			input: `Time<2023-01-02T15:04:05Z`,
			expected: expected{
				ok:   true,
				repr: `(Time < "2023-01-02T15:04:05Z")`,
			},
		},
		// Booleans
		{
			name:  "bool eq",
			input: `Flag==true`,
			expected: expected{
				ok:   true,
				repr: `(Flag == true)`,
			},
		},
		{
			name:  "bool neq",
			input: `Flag!=False`,
			expected: expected{
				ok:   true,
				repr: `(Flag != False)`,
			},
		},
		// Logic and precedence
		{
			name:  "and or precedence",
			input: `HP>50&&MP>=100||LP==0`,
			expected: expected{
				ok:   true,
				repr: `(((HP > 50) && (MP >= 100)) || (LP == 0))`,
			},
		},
		{
			name:  "paren grouping",
			input: `(HP>50&&MP>=100)||LP==0`,
			expected: expected{
				ok:   true,
				repr: `(((HP > 50) && (MP >= 100)) || (LP == 0))`,
			},
		},
		{
			name:  "not group",
			input: `!(SPD<20)`,
			expected: expected{
				ok:   true,
				repr: `(! (SPD < 20))`,
			},
		},
		{
			name:  "complex",
			input: `Class=="軍師"&&Name=~'孔明'&&(HP>50&&MP>=100&&LP!=0)&&(MAG>=20||!(SPD<20))`,
			expected: expected{
				ok:   true,
				repr: `((((Class == "軍師") && (Name =~ "孔明")) && (((HP > 50) && (MP >= 100)) && (LP != 0))) && ((MAG >= 20) || (! (SPD < 20))))`,
			},
		},
		// Errors
		{
			name:  "regex empty pattern",
			input: `Name=~''`,
			expected: expected{
				ok:  false,
				err: `invalid regex`,
			},
		},
		{
			name:  "invalid regex",
			input: `Name=~'['`,
			expected: expected{
				ok:  false,
				err: `invalid regex`,
			},
		},
		{
			name:  "missing op",
			input: `HP 50`,
			expected: expected{
				ok:  false,
				err: `expected comparison operator`,
			},
		},
		{
			name:  "missing rhs",
			input: `HP>`,
			expected: expected{
				ok:  false,
				err: `expected value`,
			},
		},
		{
			name:  "unexpected trailing",
			input: `HP>50 extra`,
			expected: expected{
				ok:  false,
				err: `unexpected token after parsing`,
			},
		},
		{
			name:  "leading not without operand",
			input: `!`,
			expected: expected{
				ok:  false,
				err: `expected left parenthesis or identifier`,
			},
		},
		{
			name:  "empty",
			input: ``,
			expected: expected{
				ok:  false,
				err: `empty input`,
			},
		},
		{
			name:  "unclosed paren",
			input: `(HP>1`,
			expected: expected{
				ok:  false,
				err: `unclosed left parenthesis`,
			},
		},
		{
			name:  "extra right paren",
			input: `HP>1)`,
			expected: expected{
				ok:  false,
				err: `unexpected token after parsing`,
			},
		},
		{
			name:  "double logical op",
			input: `HP>1&&||MP>2`,
			expected: expected{
				ok:  false,
				err: `expected left parenthesis or identifier`,
			},
		},
		{
			name:  "non ident lhs",
			input: `123==456`,
			expected: expected{
				ok:  false,
				err: `expected left parenthesis or identifier`,
			},
		},
		{
			name:  "unterminated regex string",
			input: `Name=~'abc`,
			expected: expected{
				ok:  false,
				err: `unterminated quoted string`,
			},
		},
		{
			name:  "number then missing op",
			input: `HP50`,
			expected: expected{
				ok:  false,
				err: `expected comparison operator`,
			},
		},
		{
			name:  "duration segment missing unit",
			input: `Delay==1h30`,
			expected: expected{
				ok:  false,
				err: `unexpected token after parsing`,
			},
		},
		{
			name:  "expect mismatch right paren",
			input: `(HP>1 Name==X)`,
			expected: expected{
				ok:  false,
				err: `expected right parenthesis`,
			},
		},
		{
			name:  "expect mismatch nested right paren",
			input: `((HP>1) Name==X)`,
			expected: expected{
				ok:  false,
				err: `expected right parenthesis`,
			},
		},
		{
			name:  "parseExpr initial next failure",
			input: `#&&HP>1`,
			expected: expected{
				ok:  false,
				err: `unexpected character`,
			},
		},
		{
			name:  "parseAND right side next failure",
			input: `HP>1&&#`,
			expected: expected{
				ok:  false,
				err: `unexpected character`,
			},
		},
		{
			name:  "parseAND operator malformed",
			input: `HP>1&X==1`,
			expected: expected{
				ok:  false,
				err: `unexpected character`,
			},
		},
		{
			name:  "parseNOT next failure",
			input: `!#`,
			expected: expected{
				ok:  false,
				err: `unexpected character`,
			},
		},
		{
			name:  "parsePrimary inner expr failure",
			input: `(#)`,
			expected: expected{
				ok:  false,
				err: `unexpected character`,
			},
		},
		{
			name:  "parsePrimary parseExpr failure",
			input: `(##)`,
			expected: expected{
				ok:  false,
				err: `unexpected character`,
			},
		},
		{
			name:  "parseComparison expect ident failure",
			input: `==1`,
			expected: expected{
				ok:  false,
				err: `expected left parenthesis or identifier`,
			},
		},
		{
			name:  "parseComparison operator next failure",
			input: `A$1`,
			expected: expected{
				ok:  false,
				err: `unexpected character`,
			},
		},
		{
			name:  "parseComparison value next failure",
			input: `A==#`,
			expected: expected{
				ok:  false,
				err: `unexpected character`,
			},
		},
		// Parenthesis limit
		{
			name: "paren limit ok (256)",
			input: func() string {
				n := 256
				var b strings.Builder
				b.Grow(n*2 + len(`HP>1`))
				for range n {
					b.WriteByte('(')
				}
				b.WriteString(`HP>1`)
				for range n {
					b.WriteByte(')')
				}
				return b.String()
			}(),
			expected: expected{
				ok:   true,
				repr: `(HP > 1)`,
			},
		},
		{
			name: "paren limit ng (257)",
			input: func() string {
				n := 257
				var b strings.Builder
				b.Grow(n*2 + len(`HP>1`))
				for range n {
					b.WriteByte('(')
				}
				b.WriteString(`HP>1`)
				for range n {
					b.WriteByte(')')
				}
				return b.String()
			}(),
			expected: expected{
				ok:  false,
				err: `too many parentheses`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := Parse(test.input)
			if test.expected.ok {
				if err != nil {
					t.Errorf(testTemplate, test.input, "", err)
					return
				}
				repr := repr(*expr)
				if repr != test.expected.repr {
					t.Errorf(testTemplate, test.input, test.expected.repr, repr)
				}
				return
			}
			if err == nil {
				t.Errorf(testTemplate, test.input, test.expected.err, "")
				return
			}
			if !strings.Contains(err.Error(), test.expected.err) {
				t.Errorf(testTemplate, test.input, test.expected.err, err)
			}
		})
	}
}

// repr converts ast to a string.
func repr(e Expr) string {
	val := func(v string) string {
		if isNumericLike(v) || isDurationLike(v) || isBoolLiteral(v) {
			return v
		}
		return "\"" + v + "\""
	}
	var walk func(int) string
	walk = func(i int) string {
		n := e.parser.nodes[i]
		switch n.typ {
		case nodeBinary:
			return "(" + walk(n.lhs) + " " + n.op.typ.literal() + " " + walk(n.rhs) + ")"
		case nodeNOT:
			return "(! " + walk(n.lhs) + ")"
		case nodeComparison:
			return "(" + n.ident.v + " " + n.op.typ.literal() + " " + val(n.val.v) + ")"
		default:
			return "<unknown>"
		}
	}
	return walk(e.root)
}

func isNumericLike(s string) bool {
	if len(s) == 0 {
		return false
	}
	digit := false
	for _, r := range s {
		if r >= '0' && r <= '9' {
			digit = true
			break
		}
	}
	if !digit {
		return false
	}
	for _, r := range s {
		switch r {
		case '+', '-', '.', '_', 'x', 'X', 'o', 'O', 'b', 'B', 'p', 'P', 'e', 'E':
			continue
		}
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') {
			continue
		}
		return false
	}
	return true
}

func isDurationLike(s string) bool {
	units := []string{"ns", "us", "μs", "ms", "s", "m", "h"}
	for _, unit := range units {
		if strings.Contains(s, unit) {
			return true
		}
	}
	return false
}
