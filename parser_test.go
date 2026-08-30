package filter

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

func Test_parse(t *testing.T) {
	type args struct {
		input string
	}
	type want struct {
		val    string
		nodes  int
		nident int
		shared bool
		isErr  bool
		err    string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "single comparison",
			args: args{
				input: `A==1`,
			},
			want: want{
				val:    `(A == 1)`,
				nodes:  1,
				nident: 1,
			},
		},
		{
			name: "same identifier twice sets shared",
			args: args{
				input: `A==1 && A==2`,
			},
			want: want{
				val:    `((A == 1) && (A == 2))`,
				nodes:  3,
				nident: 1,
				shared: true,
			},
		},
		{
			name: "distinct identifiers do not set shared",
			args: args{
				input: `A==1 && B==2`,
			},
			want: want{
				val:    `((A == 1) && (B == 2))`,
				nodes:  3,
				nident: 2,
			},
		},
		{
			name: "more nodes and identifiers than the inline buffers",
			args: args{
				input: func() string {
					var b strings.Builder
					b.WriteString("A0==0")
					for i := 1; i < 17; i++ {
						b.WriteString(" && A")
						b.WriteString(string(rune('0' + i/10)))
						b.WriteString(string(rune('0' + i%10)))
						b.WriteString("==")
						b.WriteString(string(rune('0' + i/10)))
						b.WriteString(string(rune('0' + i%10)))
					}
					return b.String()
				}(),
			},
			want: want{
				val: func() string {
					s := "(A0 == 0)"
					for i := 1; i < 17; i++ {
						n := string(rune('0'+i/10)) + string(rune('0'+i%10))
						s = "(" + s + " && (A" + n + " == " + n + "))"
					}
					return s
				}(),
				nodes:  33,
				nident: 17,
			},
		},
		{
			name: "input of exactly MaxInput bytes",
			args: args{
				input: `A=="` + strings.Repeat("x", MaxInput-5) + `"`,
			},
			want: want{
				val:    `(A == "` + strings.Repeat("x", MaxInput-5) + `")`,
				nodes:  1,
				nident: 1,
			},
		},
		{
			name: "empty input",
			args: args{
				input: "",
			},
			want: want{
				isErr: true,
				err:   `parse error: empty input`,
			},
		},
		{
			name: "input one byte over MaxInput",
			args: args{
				input: `A=="` + strings.Repeat("x", MaxInput-4) + `"`,
			},
			want: want{
				isErr: true,
				err:   `parse error: input too long: 1048577 bytes exceeds limit 1048576`,
			},
		},
		{
			name: "error from the expression",
			args: args{
				input: `A==`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:4: expected value, got EOF: ""`,
			},
		},
		{
			name: "trailing token",
			args: args{
				input: `A==1 B`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:6: unexpected token after parsing: B`,
			},
		},
		{
			name: "trailing lexer error",
			args: args{
				input: `A==1 $`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:6: unexpected token after parsing: unexpected character U+0024 '$'`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parse(test.args.input)
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
			if val := repr(&Expr{expr: got}); val != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", val, test.want.val)
			}
			if len(got.nodes) != test.want.nodes {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", len(got.nodes), test.want.nodes)
			}
			if got.nident != test.want.nident {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got.nident, test.want.nident)
			}
			if got.shared != test.want.shared {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got.shared, test.want.shared)
			}
		})
	}
}

func Test_newParser(t *testing.T) {
	type args struct {
		input string
	}
	type want struct {
		val parser
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty input",
			args: args{
				input: "",
			},
			want: want{
				val: parser{
					lexer: newLexer(""),
				},
			},
		},
		{
			name: "input is handed to the lexer untouched",
			args: args{
				input: `A == 1`,
			},
			want: want{
				val: parser{
					lexer: newLexer(`A == 1`),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newParser(test.args.input)
			if !reflect.DeepEqual(got, test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_parser_parseExpr(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val   string
		isErr bool
		err   string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "single comparison",
			fields: fields{
				input: `A==1`,
			},
			want: want{
				val: `(A == 1)`,
			},
		},
		{
			name: "or",
			fields: fields{
				input: `A==1 || B==2`,
			},
			want: want{
				val: `((A == 1) || (B == 2))`,
			},
		},
		{
			name: "or is left associative",
			fields: fields{
				input: `A==1 || B==2 || C==3`,
			},
			want: want{
				val: `(((A == 1) || (B == 2)) || (C == 3))`,
			},
		},
		{
			name: "and binds tighter than or",
			fields: fields{
				input: `A==1 || B==2 && C==3`,
			},
			want: want{
				val: `((A == 1) || ((B == 2) && (C == 3)))`,
			},
		},
		{
			name: "not binds tighter than or",
			fields: fields{
				input: `!A==1 || B==2`,
			},
			want: want{
				val: `((! (A == 1)) || (B == 2))`,
			},
		},
		{
			name: "stops before a trailing token",
			fields: fields{
				input: `A==1 B==2`,
			},
			want: want{
				val: `(A == 1)`,
			},
		},
		{
			name: "error on the left side",
			fields: fields{
				input: `==1 || B==2`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:1: expected left parenthesis or identifier, got "equal to" operator: "=="`,
			},
		},
		{
			name: "error on the right side",
			fields: fields{
				input: `A==1 || ==2`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:9: expected left parenthesis or identifier, got "equal to" operator: "=="`,
			},
		},
		{
			name: "missing right side",
			fields: fields{
				input: `A==1 ||`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:8: expected left parenthesis or identifier, got EOF: ""`,
			},
		},
		{
			name: "empty input",
			fields: fields{
				input: ``,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:1: expected left parenthesis or identifier, got EOF: ""`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := newParser(test.fields.input)
			got, err := p.parseExpr()
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
			if val := reprAt(&p, got); val != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", val, test.want.val)
			}
		})
	}
}

func Test_parser_parseAND(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val   string
		isErr bool
		err   string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "single comparison",
			fields: fields{
				input: `A==1`,
			},
			want: want{
				val: `(A == 1)`,
			},
		},
		{
			name: "and",
			fields: fields{
				input: `A==1 && B==2`,
			},
			want: want{
				val: `((A == 1) && (B == 2))`,
			},
		},
		{
			name: "and is left associative",
			fields: fields{
				input: `A==1 && B==2 && C==3`,
			},
			want: want{
				val: `(((A == 1) && (B == 2)) && (C == 3))`,
			},
		},
		{
			name: "not on the right side",
			fields: fields{
				input: `A==1 && !B==2`,
			},
			want: want{
				val: `((A == 1) && (! (B == 2)))`,
			},
		},
		{
			name: "stops before or",
			fields: fields{
				input: `A==1 || B==2`,
			},
			want: want{
				val: `(A == 1)`,
			},
		},
		{
			name: "error on the left side",
			fields: fields{
				input: `1 && B==2`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:1: expected left parenthesis or identifier, got number: "1"`,
			},
		},
		{
			name: "error on the right side",
			fields: fields{
				input: `A==1 && 2`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:9: expected left parenthesis or identifier, got number: "2"`,
			},
		},
		{
			name: "missing right side",
			fields: fields{
				input: `A==1 &&`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:8: expected left parenthesis or identifier, got EOF: ""`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := newParser(test.fields.input)
			got, err := p.parseAND()
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
			if val := reprAt(&p, got); val != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", val, test.want.val)
			}
		})
	}
}

func Test_parser_parseNOT(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val   string
		isErr bool
		err   string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "without not",
			fields: fields{
				input: `A==1`,
			},
			want: want{
				val: `(A == 1)`,
			},
		},
		{
			name: "not",
			fields: fields{
				input: `!A==1`,
			},
			want: want{
				val: `(! (A == 1))`,
			},
		},
		{
			name: "double not through parentheses",
			fields: fields{
				input: `!(!A==1)`,
			},
			want: want{
				val: `(! (! (A == 1)))`,
			},
		},
		{
			name: "not applies to the parenthesized group",
			fields: fields{
				input: `!(A==1 || B==2)`,
			},
			want: want{
				val: `(! ((A == 1) || (B == 2)))`,
			},
		},
		{
			name: "not applies to one comparison only",
			fields: fields{
				input: `!A==1 && B==2`,
			},
			want: want{
				val: `(! (A == 1))`,
			},
		},
		{
			name: "not without operand",
			fields: fields{
				input: `!`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:2: expected left parenthesis or identifier, got EOF: ""`,
			},
		},
		{
			name: "double not without parentheses",
			fields: fields{
				input: `!!A==1`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:2: expected left parenthesis or identifier, got logical NOT operator: "!"`,
			},
		},
		{
			name: "not before an operator",
			fields: fields{
				input: `! ==1`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:3: expected left parenthesis or identifier, got "equal to" operator: "=="`,
			},
		},
		{
			name: "not glued to equals lexes as not equal",
			fields: fields{
				input: `!==1`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:1: expected left parenthesis or identifier, got "not equal to" operator: "!="`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := newParser(test.fields.input)
			got, err := p.parseNOT()
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
			if val := reprAt(&p, got); val != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", val, test.want.val)
			}
		})
	}
}

func Test_parser_parsePrimary(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val        string
		parenCount int
		isErr      bool
		err        string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "comparison",
			fields: fields{
				input: `A==1`,
			},
			want: want{
				val: `(A == 1)`,
			},
		},
		{
			name: "parenthesized comparison",
			fields: fields{
				input: `(A==1)`,
			},
			want: want{
				val:        `(A == 1)`,
				parenCount: 1,
			},
		},
		{
			name: "nested parentheses",
			fields: fields{
				input: `((A==1))`,
			},
			want: want{
				val:        `(A == 1)`,
				parenCount: 2,
			},
		},
		{
			name: "parentheses group an or",
			fields: fields{
				input: `(A==1 || B==2) && C==3`,
			},
			want: want{
				val:        `((A == 1) || (B == 2))`,
				parenCount: 1,
			},
		},
		{
			name: "exactly MaxParen nested parentheses",
			fields: fields{
				input: strings.Repeat("(", MaxParen) + `A==1` + strings.Repeat(")", MaxParen),
			},
			want: want{
				val:        `(A == 1)`,
				parenCount: MaxParen,
			},
		},
		{
			name: "one more than MaxParen nested parentheses",
			fields: fields{
				input: strings.Repeat("(", MaxParen+1) + `A==1` + strings.Repeat(")", MaxParen+1),
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:257: too many parentheses: exceeded limit 256`,
			},
		},
		{
			name: "missing right parenthesis",
			fields: fields{
				input: `(A==1 B`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:7: expected right parenthesis, got identifier: "B"`,
			},
		},
		{
			name: "unclosed left parenthesis",
			fields: fields{
				input: `(A==1`,
			},
			want: want{
				isErr: true,
				err:   `token error at 1:6: unclosed left parenthesis`,
			},
		},
		{
			name: "left parenthesis only",
			fields: fields{
				input: `(`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:2: expected left parenthesis or identifier, got error: "unclosed left parenthesis"`,
			},
		},
		{
			name: "error inside parentheses",
			fields: fields{
				input: `(1)`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:2: expected left parenthesis or identifier, got number: "1"`,
			},
		},
		{
			name: "number instead of identifier",
			fields: fields{
				input: `1`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:1: expected left parenthesis or identifier, got number: "1"`,
			},
		},
		{
			name: "empty input",
			fields: fields{
				input: ``,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:1: expected left parenthesis or identifier, got EOF: ""`,
			},
		},
		{
			name: "lexer error",
			fields: fields{
				input: `$`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:1: expected left parenthesis or identifier, got error: "unexpected character U+0024 '$'"`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := newParser(test.fields.input)
			got, err := p.parsePrimary()
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
			if val := reprAt(&p, got); val != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", val, test.want.val)
			}
			if p.parenCount != test.want.parenCount {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", p.parenCount, test.want.parenCount)
			}
		})
	}
}

func Test_parser_parseComparison(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val     string
		hasNum  bool
		hasTime bool
		hasDur  bool
		regex   bool
		isErr   bool
		err     string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "number caches the number and the unix time",
			fields: fields{
				input: `A==1`,
			},
			want: want{
				val:     `(A == 1)`,
				hasNum:  true,
				hasTime: true,
			},
		},
		{
			name: "fractional number caches the number only",
			fields: fields{
				input: `A>=1.5`,
			},
			want: want{
				val:    `(A >= 1.5)`,
				hasNum: true,
			},
		},
		{
			name: "every comparison operator",
			fields: fields{
				input: `A<=1`,
			},
			want: want{
				val:     `(A <= 1)`,
				hasNum:  true,
				hasTime: true,
			},
		},
		{
			name: "string literal is unquoted",
			fields: fields{
				input: `A!="abc"`,
			},
			want: want{
				val: `(A != "abc")`,
			},
		},
		{
			name: "raw string literal is unquoted",
			fields: fields{
				input: `A=='abc'`,
			},
			want: want{
				val: `(A == "abc")`,
			},
		},
		{
			name: "string that spells a number is converted",
			fields: fields{
				input: `A>"50"`,
			},
			want: want{
				val:     `(A > "50")`,
				hasNum:  true,
				hasTime: true,
			},
		},
		{
			name: "string that spells a time is converted",
			fields: fields{
				input: `A<'2025-01-01T00:00:00Z'`,
			},
			want: want{
				val:     `(A < "2025-01-01T00:00:00Z")`,
				hasTime: true,
			},
		},
		{
			name: "string that spells a duration is converted",
			fields: fields{
				input: `A>"1h30m"`,
			},
			want: want{
				val:    `(A > "1h30m")`,
				hasDur: true,
			},
		},
		{
			name: "bare time",
			fields: fields{
				input: `A<2025-01-01`,
			},
			want: want{
				val:     `(A < "2025-01-01")`,
				hasTime: true,
			},
		},
		{
			name: "duration",
			fields: fields{
				input: `A>10s`,
			},
			want: want{
				val:    `(A > 10s)`,
				hasDur: true,
			},
		},
		{
			name: "bool",
			fields: fields{
				input: `A==true`,
			},
			want: want{
				val: `(A == true)`,
			},
		},
		{
			name: "regex",
			fields: fields{
				input: `A=~'^a.*'`,
			},
			want: want{
				val:   `(A =~ "^a.*")`,
				regex: true,
			},
		},
		{
			name: "negated regex",
			fields: fields{
				input: `A!~"b$"`,
			},
			want: want{
				val:   `(A !~ "b$")`,
				regex: true,
			},
		},
		{
			name: "regex against a number literal",
			fields: fields{
				input: `A=~1`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:4: expected string pattern, got number: "1"`,
			},
		},
		{
			name: "stops after the literal",
			fields: fields{
				input: `A==1 && B==2`,
			},
			want: want{
				val:     `(A == 1)`,
				hasNum:  true,
				hasTime: true,
			},
		},
		{
			name: "empty regex",
			fields: fields{
				input: `A=~''`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:4: invalid regex "": empty pattern`,
			},
		},
		{
			name: "invalid regex",
			fields: fields{
				input: `A=~'['`,
			},
			want: want{
				isErr: true,
				err:   "parse error at 1:4: invalid regex \"[\": error parsing regexp: missing closing ]: `[`",
			},
		},
		{
			name: "invalid time",
			fields: fields{
				input: `A<2025-13-01`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:3: invalid time "2025-13-01"`,
			},
		},
		{
			name: "invalid duration",
			fields: fields{
				input: `A>1.5.5s`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:3: invalid duration "1.5.5s"`,
			},
		},
		{
			name: "invalid number",
			fields: fields{
				input: `A==1e`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:4: invalid number "1e"`,
			},
		},
		{
			name: "missing identifier",
			fields: fields{
				input: `==1`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:1: expected identifier, got "equal to" operator: "=="`,
			},
		},
		{
			name: "missing operator",
			fields: fields{
				input: `A 1`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:3: expected comparison operator, got number: "1"`,
			},
		},
		{
			name: "logical operator instead of comparison operator",
			fields: fields{
				input: `A && B`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:3: expected comparison operator, got logical AND operator: "&&"`,
			},
		},
		{
			name: "missing value",
			fields: fields{
				input: `A==`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:4: expected value, got EOF: ""`,
			},
		},
		{
			name: "identifier instead of value",
			fields: fields{
				input: `A== B`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:5: expected value, got identifier: "B"`,
			},
		},
		{
			name: "lexer error at the operator",
			fields: fields{
				input: `A $`,
			},
			want: want{
				isErr: true,
				err:   `token error at 1:3: unexpected character U+0024 '$'`,
			},
		},
		{
			name: "lexer error at the value",
			fields: fields{
				input: `A==$`,
			},
			want: want{
				isErr: true,
				err:   `token error at 1:4: unexpected character U+0024 '$'`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := newParser(test.fields.input)
			got, err := p.parseComparison()
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
			if val := reprAt(&p, got); val != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", val, test.want.val)
			}
			n := p.node(got)
			if n.hasNum != test.want.hasNum {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", n.hasNum, test.want.hasNum)
			}
			if n.hasTime != test.want.hasTime {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", n.hasTime, test.want.hasTime)
			}
			if n.hasDur != test.want.hasDur {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", n.hasDur, test.want.hasDur)
			}
			if regex := n.re != nil; regex != test.want.regex {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", regex, test.want.regex)
			}
		})
	}
}

func Test_parser_cacheRegex(t *testing.T) {
	type fields struct {
		cached string
	}
	type args struct {
		i int32
		t token
	}
	type want struct {
		val    string
		cached bool
		isErr  bool
		err    string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "compiles and stores the pattern",
			args: args{
				i: 0,
				t: token{
					typ:  tokenRawString,
					v:    `^Test_[a-z]+$`,
					line: 1,
					col:  4,
				},
			},
			want: want{
				val: `^Test_[a-z]+$`,
			},
		},
		{
			name: "reuses the compiled pattern from regexMap",
			fields: fields{
				cached: `^cached$`,
			},
			args: args{
				i: 0,
				t: token{
					typ:  tokenRawString,
					v:    `^cached$`,
					line: 1,
					col:  4,
				},
			},
			want: want{
				val:    `^cached$`,
				cached: true,
			},
		},
		{
			name: "stores on the requested node",
			args: args{
				i: 3,
				t: token{
					typ:  tokenString,
					v:    `abc`,
					line: 1,
					col:  4,
				},
			},
			want: want{
				val: `abc`,
			},
		},
		{
			name: "empty pattern",
			args: args{
				i: 0,
				t: token{
					typ:  tokenString,
					v:    ``,
					line: 1,
					col:  4,
				},
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:4: invalid regex "": empty pattern`,
			},
		},
		{
			name: "invalid pattern",
			args: args{
				i: 0,
				t: token{
					typ:  tokenString,
					v:    `(`,
					line: 2,
					col:  7,
				},
			},
			want: want{
				isErr: true,
				err:   "parse error at 2:7: invalid regex \"(\": error parsing regexp: missing closing ): `(`",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stored *regexp.Regexp
			if test.fields.cached != "" {
				stored = regexp.MustCompile(test.fields.cached)
				regexMap.Store(test.fields.cached, stored)
			}
			p := newParser("")
			err := p.cacheRegex(test.args.i, test.args.t)
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
			got := p.node(test.args.i).re
			if got.String() != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got.String(), test.want.val)
			}
			if cached := got == stored; cached != test.want.cached {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", cached, test.want.cached)
			}
		})
	}
}

func Test_parser_cacheValues(t *testing.T) {
	type args struct {
		i int32
		s string
	}
	type want struct {
		hasNum  bool
		num     float64
		hasDur  bool
		dur     time.Duration
		hasTime bool
		time    time.Time
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "integer spells a number and unix seconds",
			args: args{
				i: 0,
				s: "123",
			},
			want: want{
				hasNum:  true,
				num:     123,
				hasTime: true,
				time:    time.Unix(123, 0).UTC(),
			},
		},
		{
			name: "fraction spells a number only",
			args: args{
				i: 0,
				s: "1.5",
			},
			want: want{
				hasNum: true,
				num:    1.5,
			},
		},
		{
			name: "leading dot",
			args: args{
				i: 0,
				s: ".5",
			},
			want: want{
				hasNum: true,
				num:    0.5,
			},
		},
		{
			name: "leading plus",
			args: args{
				i: 0,
				s: "+1",
			},
			want: want{
				hasNum:  true,
				num:     1,
				hasTime: true,
				time:    time.Unix(1, 0).UTC(),
			},
		},
		{
			name: "leading minus",
			args: args{
				i: 0,
				s: "-1",
			},
			want: want{
				hasNum:  true,
				num:     -1,
				hasTime: true,
				time:    time.Unix(-1, 0).UTC(),
			},
		},
		{
			name: "underscore separators",
			args: args{
				i: 0,
				s: "1_000",
			},
			want: want{
				hasNum:  true,
				num:     1000,
				hasTime: true,
				time:    time.Unix(1000, 0).UTC(),
			},
		},
		{
			name: "sign only",
			args: args{
				i: 0,
				s: "-",
			},
			want: want{},
		},
		{
			name: "duration",
			args: args{
				i: 0,
				s: "1h30m",
			},
			want: want{
				hasDur: true,
				dur:    90 * time.Minute,
			},
		},
		{
			name: "rfc3339 time",
			args: args{
				i: 0,
				s: "2025-01-01T00:00:00Z",
			},
			want: want{
				hasTime: true,
				time:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "time whose layout contains spaces",
			args: args{
				i: 0,
				s: "2025-01-01 00:00:00",
			},
			want: want{
				hasTime: true,
				time:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "time whose layout starts with a weekday",
			args: args{
				i: 0,
				s: "Wed, 01 Jan 2025 00:00:00 UTC",
			},
			want: want{
				hasTime: true,
				time:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "weekday layout that is not a time",
			args: args{
				i: 0,
				s: "Hello, World",
			},
			want: want{},
		},
		{
			name: "number followed by letters",
			args: args{
				i: 0,
				s: "1x",
			},
			want: want{},
		},
		{
			name: "leading space",
			args: args{
				i: 0,
				s: " 1",
			},
			want: want{},
		},
		{
			name: "trailing space",
			args: args{
				i: 0,
				s: "1 ",
			},
			want: want{},
		},
		{
			name: "plain text",
			args: args{
				i: 0,
				s: "abc",
			},
			want: want{},
		},
		{
			name: "empty",
			args: args{
				i: 0,
				s: "",
			},
			want: want{},
		},
		{
			name: "stores on the requested node",
			args: args{
				i: 5,
				s: "42",
			},
			want: want{
				hasNum:  true,
				num:     42,
				hasTime: true,
				time:    time.Unix(42, 0).UTC(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := newParser("")
			p.cacheValues(test.args.i, test.args.s)
			n := p.node(test.args.i)
			got := want{
				hasNum:  n.hasNum,
				num:     n.num,
				hasDur:  n.hasDur,
				dur:     n.dur,
				hasTime: n.hasTime,
				time:    n.time,
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want)
			}
		})
	}
}

func Test_parser_cacheTime(t *testing.T) {
	type args struct {
		i int32
		s string
	}
	type want struct {
		val     bool
		hasTime bool
		time    time.Time
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "rfc3339",
			args: args{
				i: 0,
				s: "2025-01-01T09:00:00+09:00",
			},
			want: want{
				val:     true,
				hasTime: true,
				time:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "unix seconds",
			args: args{
				i: 0,
				s: "0",
			},
			want: want{
				val:     true,
				hasTime: true,
				time:    time.Unix(0, 0).UTC(),
			},
		},
		{
			name: "stores on the requested node",
			args: args{
				i: 2,
				s: "2025-01-01",
			},
			want: want{
				val:     true,
				hasTime: true,
				time:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "invalid",
			args: args{
				i: 0,
				s: "2025-13-01",
			},
			want: want{},
		},
		{
			name: "empty",
			args: args{
				i: 0,
				s: "",
			},
			want: want{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := newParser("")
			got := p.cacheTime(test.args.i, test.args.s)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			n := p.node(test.args.i)
			if n.hasTime != test.want.hasTime {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", n.hasTime, test.want.hasTime)
			}
			if !n.time.Equal(test.want.time) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", n.time, test.want.time)
			}
		})
	}
}

func Test_parser_cacheDuration(t *testing.T) {
	type args struct {
		i int32
		s string
	}
	type want struct {
		val    bool
		hasDur bool
		dur    time.Duration
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "seconds",
			args: args{
				i: 0,
				s: "10s",
			},
			want: want{
				val:    true,
				hasDur: true,
				dur:    10 * time.Second,
			},
		},
		{
			name: "compound with fraction",
			args: args{
				i: 0,
				s: "1h30.5m",
			},
			want: want{
				val:    true,
				hasDur: true,
				dur:    time.Hour + 30*time.Minute + 30*time.Second,
			},
		},
		{
			name: "negative",
			args: args{
				i: 0,
				s: "-1ms",
			},
			want: want{
				val:    true,
				hasDur: true,
				dur:    -time.Millisecond,
			},
		},
		{
			name: "zero",
			args: args{
				i: 0,
				s: "0",
			},
			want: want{
				val:    true,
				hasDur: true,
			},
		},
		{
			name: "stores on the requested node",
			args: args{
				i: 2,
				s: "1ns",
			},
			want: want{
				val:    true,
				hasDur: true,
				dur:    time.Nanosecond,
			},
		},
		{
			name: "missing unit",
			args: args{
				i: 0,
				s: "10",
			},
			want: want{},
		},
		{
			name: "unknown unit",
			args: args{
				i: 0,
				s: "10d",
			},
			want: want{},
		},
		{
			name: "empty",
			args: args{
				i: 0,
				s: "",
			},
			want: want{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := newParser("")
			got := p.cacheDuration(test.args.i, test.args.s)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			n := p.node(test.args.i)
			if n.hasDur != test.want.hasDur {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", n.hasDur, test.want.hasDur)
			}
			if n.dur != test.want.dur {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", n.dur, test.want.dur)
			}
		})
	}
}

func Test_parser_cacheNumber(t *testing.T) {
	type args struct {
		i int32
		s string
	}
	type want struct {
		val    bool
		hasNum bool
		num    float64
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "integer",
			args: args{
				i: 0,
				s: "42",
			},
			want: want{
				val:    true,
				hasNum: true,
				num:    42,
			},
		},
		{
			name: "negative fraction",
			args: args{
				i: 0,
				s: "-1.5",
			},
			want: want{
				val:    true,
				hasNum: true,
				num:    -1.5,
			},
		},
		{
			name: "exponent",
			args: args{
				i: 0,
				s: "1e3",
			},
			want: want{
				val:    true,
				hasNum: true,
				num:    1000,
			},
		},
		{
			name: "underscore separators",
			args: args{
				i: 0,
				s: "1_000",
			},
			want: want{
				val:    true,
				hasNum: true,
				num:    1000,
			},
		},
		{
			name: "stores on the requested node",
			args: args{
				i: 2,
				s: "7",
			},
			want: want{
				val:    true,
				hasNum: true,
				num:    7,
			},
		},
		{
			name: "hexadecimal without exponent",
			args: args{
				i: 0,
				s: "0x1f",
			},
			want: want{},
		},
		{
			name: "out of range",
			args: args{
				i: 0,
				s: "1e400",
			},
			want: want{},
		},
		{
			name: "text",
			args: args{
				i: 0,
				s: "abc",
			},
			want: want{},
		},
		{
			name: "empty",
			args: args{
				i: 0,
				s: "",
			},
			want: want{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := newParser("")
			got := p.cacheNumber(test.args.i, test.args.s)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			n := p.node(test.args.i)
			if n.hasNum != test.want.hasNum {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", n.hasNum, test.want.hasNum)
			}
			if n.num != test.want.num {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", n.num, test.want.num)
			}
		})
	}
}

func Test_parser_identIndex(t *testing.T) {
	type fields struct {
		identBuf [identBufSize]string
		idents   []string
		nident   int32
	}
	type args struct {
		name string
	}
	type want struct {
		val    int32
		nident int32
		shared bool
		idents []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "first identifier",
			args: args{
				name: "A",
			},
			want: want{
				val:    0,
				nident: 1,
			},
		},
		{
			name: "second identifier",
			fields: fields{
				identBuf: [identBufSize]string{"A"},
				nident:   1,
			},
			args: args{
				name: "B",
			},
			want: want{
				val:    1,
				nident: 2,
			},
		},
		{
			name: "repeated identifier sets shared",
			fields: fields{
				identBuf: [identBufSize]string{"A", "B"},
				nident:   2,
			},
			args: args{
				name: "A",
			},
			want: want{
				val:    0,
				nident: 2,
				shared: true,
			},
		},
		{
			name: "identifier that only shares a prefix is new",
			fields: fields{
				identBuf: [identBufSize]string{"Ab"},
				nident:   1,
			},
			args: args{
				name: "A",
			},
			want: want{
				val:    1,
				nident: 2,
			},
		},
		{
			name: "inline buffer full moves identifiers to the heap",
			fields: fields{
				identBuf: [identBufSize]string{"A", "B", "C", "D", "E", "F", "G", "H"},
				nident:   identBufSize,
			},
			args: args{
				name: "I",
			},
			want: want{
				val:    identBufSize,
				nident: identBufSize + 1,
				idents: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"},
			},
		},
		{
			name: "heap identifiers grow",
			fields: fields{
				idents: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"},
				nident: identBufSize + 1,
			},
			args: args{
				name: "J",
			},
			want: want{
				val:    identBufSize + 1,
				nident: identBufSize + 2,
				idents: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"},
			},
		},
		{
			name: "repeated identifier on the heap sets shared",
			fields: fields{
				idents: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"},
				nident: identBufSize + 1,
			},
			args: args{
				name: "I",
			},
			want: want{
				val:    identBufSize,
				nident: identBufSize + 1,
				shared: true,
				idents: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &parser{
				identBuf: test.fields.identBuf,
				idents:   test.fields.idents,
				nident:   test.fields.nident,
			}
			got := p.identIndex(test.args.name)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if p.nident != test.want.nident {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", p.nident, test.want.nident)
			}
			if p.shared != test.want.shared {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", p.shared, test.want.shared)
			}
			if !reflect.DeepEqual(p.idents, test.want.idents) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", p.idents, test.want.idents)
			}
		})
	}
}

func Test_parser_addNode(t *testing.T) {
	type fields struct {
		nnode int32
		nodes []node
	}
	type args struct {
		n node
	}
	type want struct {
		val   int32
		nnode int32
		nodes int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "first node",
			args: args{
				n: node{
					typ: nodeComparison,
				},
			},
			want: want{
				val:   0,
				nnode: 1,
			},
		},
		{
			name: "last inline slot",
			fields: fields{
				nnode: nodeBufSize - 1,
			},
			args: args{
				n: node{
					typ: nodeComparison,
				},
			},
			want: want{
				val:   nodeBufSize - 1,
				nnode: nodeBufSize,
			},
		},
		{
			name: "inline buffer full moves nodes to the heap",
			fields: fields{
				nnode: nodeBufSize,
			},
			args: args{
				n: node{
					typ: nodeComparison,
				},
			},
			want: want{
				val:   nodeBufSize,
				nnode: nodeBufSize + 1,
				nodes: nodeBufSize + 1,
			},
		},
		{
			name: "heap nodes grow",
			fields: fields{
				nnode: nodeBufSize + 1,
				nodes: make([]node, nodeBufSize+1),
			},
			args: args{
				n: node{
					typ: nodeComparison,
				},
			},
			want: want{
				val:   nodeBufSize + 1,
				nnode: nodeBufSize + 2,
				nodes: nodeBufSize + 2,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &parser{
				nnode: test.fields.nnode,
				nodes: test.fields.nodes,
			}
			got := p.addNode(test.args.n)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if p.nnode != test.want.nnode {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", p.nnode, test.want.nnode)
			}
			if len(p.nodes) != test.want.nodes {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", len(p.nodes), test.want.nodes)
			}
			if !reflect.DeepEqual(*p.node(got), test.args.n) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", *p.node(got), test.args.n)
			}
		})
	}
}

func Test_parser_node(t *testing.T) {
	type fields struct {
		nodeBuf [nodeBufSize]node
		nodes   []node
	}
	type args struct {
		i int32
	}
	type want struct {
		val node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "first inline node",
			fields: fields{
				nodeBuf: [nodeBufSize]node{
					{
						typ: nodeComparison,
					},
					{
						typ: nodeNOT,
					},
				},
			},
			args: args{
				i: 0,
			},
			want: want{
				val: node{
					typ: nodeComparison,
				},
			},
		},
		{
			name: "last inline node",
			fields: fields{
				nodeBuf: func() [nodeBufSize]node {
					var buf [nodeBufSize]node
					buf[nodeBufSize-1] = node{
						typ: nodeBinary,
					}
					return buf
				}(),
			},
			args: args{
				i: nodeBufSize - 1,
			},
			want: want{
				val: node{
					typ: nodeBinary,
				},
			},
		},
		{
			name: "heap nodes take precedence over the inline buffer",
			fields: fields{
				nodeBuf: [nodeBufSize]node{
					{
						typ: nodeComparison,
					},
				},
				nodes: []node{
					{
						typ: nodeNOT,
					},
				},
			},
			args: args{
				i: 0,
			},
			want: want{
				val: node{
					typ: nodeNOT,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &parser{
				nodeBuf: test.fields.nodeBuf,
				nodes:   test.fields.nodes,
			}
			got := p.node(test.args.i)
			if !reflect.DeepEqual(*got, test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", *got, test.want.val)
			}
			got.left = 99
			if p.node(test.args.i).left != 99 {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", p.node(test.args.i).left, 99)
			}
		})
	}
}

func Test_parser_expect(t *testing.T) {
	type fields struct {
		input string
	}
	type args struct {
		typ tokenType
	}
	type want struct {
		val   token
		isErr bool
		err   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "matching token is consumed",
			fields: fields{
				input: `A`,
			},
			args: args{
				typ: tokenIdent,
			},
			want: want{
				val: token{
					typ:  tokenIdent,
					v:    "A",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "matching token after whitespace",
			fields: fields{
				input: "\n  )",
			},
			args: args{
				typ: tokenRparen,
			},
			want: want{
				val: token{
					typ:  tokenRparen,
					v:    ")",
					pos:  3,
					line: 2,
					col:  3,
				},
			},
		},
		{
			name: "mismatch",
			fields: fields{
				input: `1`,
			},
			args: args{
				typ: tokenIdent,
			},
			want: want{
				val: token{
					typ:  tokenNumber,
					v:    "1",
					pos:  0,
					line: 1,
					col:  1,
				},
				isErr: true,
				err:   `parse error at 1:1: expected identifier, got number: "1"`,
			},
		},
		{
			name: "eof",
			fields: fields{
				input: ``,
			},
			args: args{
				typ: tokenIdent,
			},
			want: want{
				val: token{
					typ:  tokenEOF,
					line: 1,
					col:  1,
				},
				isErr: true,
				err:   `parse error at 1:1: expected identifier, got EOF: ""`,
			},
		},
		{
			name: "lexer error",
			fields: fields{
				input: `$`,
			},
			args: args{
				typ: tokenIdent,
			},
			want: want{
				val: token{
					typ:  tokenError,
					v:    "unexpected character U+0024 '$'",
					line: 1,
					col:  1,
				},
				isErr: true,
				err:   `token error at 1:1: unexpected character U+0024 '$'`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := newParser(test.fields.input)
			got, err := p.expect(test.args.typ)
			isErr := err != nil
			if isErr != test.want.isErr {
				t.Errorf("error mismatch\ngot=%v\nwant=%v\n", isErr, test.want.isErr)
				return
			}
			if isErr && err.Error() != test.want.err {
				t.Errorf("error mismatch\ngot=%v\nwant=%v\n", err, test.want.err)
			}
			if !reflect.DeepEqual(got, test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_parser_next(t *testing.T) {
	type fields struct {
		input  string
		peeked bool
	}
	type want struct {
		val    token
		peeked bool
		isErr  bool
		err    string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "first token",
			fields: fields{
				input: `A == 1`,
			},
			want: want{
				val: token{
					typ:  tokenIdent,
					v:    "A",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "peeked token is returned and the peek is cleared",
			fields: fields{
				input:  `A == 1`,
				peeked: true,
			},
			want: want{
				val: token{
					typ:  tokenIdent,
					v:    "A",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "eof",
			fields: fields{
				input: ``,
			},
			want: want{
				val: token{
					typ:  tokenEOF,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "lexer error",
			fields: fields{
				input: `$`,
			},
			want: want{
				val: token{
					typ:  tokenError,
					v:    "unexpected character U+0024 '$'",
					line: 1,
					col:  1,
				},
				isErr: true,
				err:   `token error at 1:1: unexpected character U+0024 '$'`,
			},
		},
		{
			name: "peeked lexer error",
			fields: fields{
				input:  `$`,
				peeked: true,
			},
			want: want{
				val: token{
					typ:  tokenError,
					v:    "unexpected character U+0024 '$'",
					line: 1,
					col:  1,
				},
				isErr: true,
				err:   `token error at 1:1: unexpected character U+0024 '$'`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := newParser(test.fields.input)
			if test.fields.peeked {
				p.peek()
			}
			got, err := p.next()
			isErr := err != nil
			if isErr != test.want.isErr {
				t.Errorf("error mismatch\ngot=%v\nwant=%v\n", isErr, test.want.isErr)
				return
			}
			if isErr && err.Error() != test.want.err {
				t.Errorf("error mismatch\ngot=%v\nwant=%v\n", err, test.want.err)
			}
			if !reflect.DeepEqual(got, test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if p.peeked != test.want.peeked {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", p.peeked, test.want.peeked)
			}
		})
	}
}

func Test_parser_peek(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val    token
		peeked bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "first token",
			fields: fields{
				input: `A == 1`,
			},
			want: want{
				val: token{
					typ:  tokenIdent,
					v:    "A",
					pos:  0,
					line: 1,
					col:  1,
				},
				peeked: true,
			},
		},
		{
			name: "eof",
			fields: fields{
				input: ``,
			},
			want: want{
				val: token{
					typ:  tokenEOF,
					line: 1,
					col:  1,
				},
				peeked: true,
			},
		},
		{
			name: "lexer error is returned as a token",
			fields: fields{
				input: `$`,
			},
			want: want{
				val: token{
					typ:  tokenError,
					v:    "unexpected character U+0024 '$'",
					line: 1,
					col:  1,
				},
				peeked: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := newParser(test.fields.input)
			got := p.peek()
			if !reflect.DeepEqual(got, test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if p.peeked != test.want.peeked {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", p.peeked, test.want.peeked)
			}
			if again := p.peek(); !reflect.DeepEqual(again, got) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", again, got)
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

func Test_unquote(t *testing.T) {
	type args struct {
		t token
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
			name: "string",
			args: args{
				t: token{
					typ: tokenString,
					v:   `"abc"`,
				},
			},
			want: want{
				val: "abc",
			},
		},
		{
			name: "raw string",
			args: args{
				t: token{
					typ: tokenRawString,
					v:   "`abc`",
				},
			},
			want: want{
				val: "abc",
			},
		},
		{
			name: "empty string",
			args: args{
				t: token{
					typ: tokenString,
					v:   `""`,
				},
			},
			want: want{
				val: "",
			},
		},
		{
			name: "too short",
			args: args{
				t: token{
					typ: tokenString,
					v:   `"`,
				},
			},
			want: want{
				val: `"`,
			},
		},
		{
			name: "number",
			args: args{
				t: token{
					typ: tokenNumber,
					v:   "42",
				},
			},
			want: want{
				val: "42",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := unquote(test.args.t)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

// reprAt renders the subtree rooted at node i of p.
func reprAt(p *parser, i int32) string {
	nodes := p.nodes
	if nodes == nil {
		nodes = p.nodeBuf[:p.nnode]
	}
	return repr(&Expr{
		expr: expr{
			nodes: nodes,
			root:  i,
		},
	})
}
