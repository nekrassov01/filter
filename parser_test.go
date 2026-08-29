package filter

import (
	"fmt"
	"strings"
	"testing"
	"time"
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
				ok:  false,
				err: `unexpected character U+002A '*'`,
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
			name:  "regex with inline flag",
			input: "Name=~`(?i)A`",
			expected: expected{
				ok:   true,
				repr: `(Name =~ "(?i)A")`,
			},
		},
		{
			name:  "negative regex with inline flag",
			input: "Name!~`(?i)A`",
			expected: expected{
				ok:   true,
				repr: `(Name !~ "(?i)A")`,
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
			name:  "missing right",
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
			name:  "non ident left",
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
			name:  "bare dot as number",
			input: `HP > .`,
			expected: expected{
				ok:  false,
				err: `parse error at 1:6: invalid number "."`,
			},
		},
		{
			name:  "bare sign as number",
			input: `HP > -`,
			expected: expected{
				ok:  false,
				err: `parse error at 1:6: invalid number "-"`,
			},
		},
		{
			name:  "bare dot after wide identifier",
			input: `名前 == .`,
			expected: expected{
				ok:  false,
				err: `parse error at 1:9: invalid number "."`,
			},
		},
		{
			name:  "bare sign on second line",
			input: "a\n== -",
			expected: expected{
				ok:  false,
				err: `parse error at 2:4: invalid number "-"`,
			},
		},
		{
			name:  "number with two dots",
			input: `HP > 1.2.3`,
			expected: expected{
				ok:  false,
				err: `unexpected token after parsing: .3`,
			},
		},
		{
			name:  "base prefix without digits",
			input: `HP > 0x`,
			expected: expected{
				ok:  false,
				err: `parse error at 1:6: invalid number "0x"`,
			},
		},
		{
			name:  "exponent without digits",
			input: `HP > 1e`,
			expected: expected{
				ok:  false,
				err: `parse error at 1:6: invalid number "1e"`,
			},
		},
		{
			name:  "date only",
			input: `Time >= 2023-01-02`,
			expected: expected{
				ok:   true,
				repr: `(Time >= "2023-01-02")`,
			},
		},
		{
			name:  "time without zone",
			input: `Time > 2023-01-02T15:04:05`,
			expected: expected{
				ok:   true,
				repr: `(Time > "2023-01-02T15:04:05")`,
			},
		},
		{
			name:  "date followed by T alone",
			input: `Time > 2023-01-02T`,
			expected: expected{
				ok:  false,
				err: `parse error at 1:18: unexpected token after parsing: T`,
			},
		},
		{
			name:  "date out of range",
			input: `Time > 2023-13-01`,
			expected: expected{
				ok:  false,
				err: `parse error at 1:8: invalid time "2023-13-01"`,
			},
		},
		{
			name:  "time out of range",
			input: `Time > 2023-13-01T00:00:00Z`,
			expected: expected{
				ok:  false,
				err: `parse error at 1:8: invalid time "2023-13-01T00:00:00Z"`,
			},
		},
		{
			name:  "duration with repeated fraction",
			input: `Duration > 1.5.5s`,
			expected: expected{
				ok:  false,
				err: `parse error at 1:12: invalid duration "1.5.5s"`,
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
			name:  "or missing right operand",
			input: `HP>1 ||`,
			expected: expected{
				ok:  false,
				err: `expected left parenthesis or identifier`,
			},
		},
		{
			name: "more identifiers and nodes than the inline buffers",
			input: func() string {
				var b strings.Builder
				for i := range 20 {
					if i > 0 {
						b.WriteString(" && ")
					}
					fmt.Fprintf(&b, "F%d > %d", i, i)
				}
				return b.String()
			}(),
			expected: expected{
				ok: true,
				repr: func() string {
					s := "(F0 > 0)"
					for i := 1; i < 20; i++ {
						s = fmt.Sprintf("(%s && (F%d > %d))", s, i, i)
					}
					return s
				}(),
			},
		},
		{
			name:  "input too long",
			input: "HP > " + strings.Repeat("1", MaxInput),
			expected: expected{
				ok:  false,
				err: `input too long`,
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
				repr := repr(expr)
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

func Test_parseTime(t *testing.T) {
	type expected struct {
		ok   bool
		time time.Time
	}
	tests := []struct {
		name     string
		input    string
		expected expected
	}{
		{
			name:  "rfc3339 utc",
			input: "2025-01-01T00:00:00Z",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "rfc3339 offset",
			input: "2025-01-01T09:00:00+09:00",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "rfc3339 fraction",
			input: "2025-01-01T00:00:00.25Z",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 250000000, time.UTC),
			},
		},
		{
			name:  "datetime without zone",
			input: "2025-01-01T00:00:00",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "datetime with space",
			input: "2025-01-01 00:00:00",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "datetime with space and fraction",
			input: "2025-01-01 00:00:00.5",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 500000000, time.UTC),
			},
		},
		{
			name:  "date only",
			input: "2025-01-01",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "unix seconds",
			input: "1735689600",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "negative unix seconds",
			input: "-1",
			expected: expected{
				ok:   true,
				time: time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name:  "rfc1123",
			input: "Wed, 01 Jan 2025 00:00:00 UTC",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "rfc850",
			input: "Wednesday, 01-Jan-25 00:00:00 UTC",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "rfc822",
			input: "01 Jan 25 00:00 UTC",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "rfc1123 with numeric offset",
			input: "Wed, 01 Jan 2025 09:00:00 +0900",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "rfc822 with numeric offset",
			input: "01 Jan 25 09:00 +0900",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "gmt abbreviation",
			input: "Wed, 01 Jan 2025 00:00:00 GMT",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "unix seconds with underscores",
			input: "1_735_689_600",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "zone abbreviation other than utc or gmt",
			input: "Wed, 01 Jan 2025 00:00:00 EST",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "lowercase z",
			input: "2025-01-01T00:00:00z",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "out of range month",
			input: "2025-13-01",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "clock only",
			input: "12:00:00",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "unix with fraction",
			input: "1735689600.5",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "unix seconds overflow",
			input: "99999999999999999999",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "sign only",
			input: "-",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "empty",
			input: "",
			expected: expected{
				ok: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := parseTime(test.input)
			if (err == nil) != test.expected.ok {
				t.Errorf(testTemplate, test.input, test.expected.ok, err)
				return
			}
			if test.expected.ok && !actual.Equal(test.expected.time) {
				t.Errorf(testTemplate, test.input, test.expected.time, actual)
			}
			if test.expected.ok && actual.Location() != time.UTC {
				t.Errorf(testTemplate, test.input, time.UTC, actual.Location())
			}
		})
	}
}

func Test_unquote(t *testing.T) {
	tests := []struct {
		name     string
		token    token
		expected string
	}{
		{
			name: "string",
			token: token{
				typ: tokenString,
				v:   `"abc"`,
			},
			expected: "abc",
		},
		{
			name: "raw string",
			token: token{
				typ: tokenRawString,
				v:   "`abc`",
			},
			expected: "abc",
		},
		{
			name: "empty string",
			token: token{
				typ: tokenString,
				v:   `""`,
			},
			expected: "",
		},
		{
			name: "too short",
			token: token{
				typ: tokenString,
				v:   `"`,
			},
			expected: `"`,
		},
		{
			name: "number",
			token: token{
				typ: tokenNumber,
				v:   "42",
			},
			expected: "42",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := unquote(test.token)
			if actual != test.expected {
				t.Errorf(testTemplate, test.token.v, test.expected, actual)
			}
		})
	}
}

// repr converts ast to a string. String and time literals are quoted;
// numbers, durations, and booleans are printed as written.
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
