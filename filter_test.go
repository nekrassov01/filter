package filter

import (
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		input string
	}
	type want struct {
		val   string
		isErr bool
		err   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		// Strings
		{
			name: "eq string",
			args: args{
				input: `Class=="軍師"`,
			},
			want: want{
				val: `(Class == "軍師")`,
			},
		},
		{
			name: "neq raw string",
			args: args{
				input: "Name!='孔明'",
			},
			want: want{
				val: `(Name != "孔明")`,
			},
		},
		{
			name: "eqi string",
			args: args{
				input: `Tag==*"Admin"`,
			},
			want: want{
				isErr: true,
				err:   `token error at 1:6: unexpected character U+002A '*'`,
			},
		},
		{
			name: "regex",
			args: args{
				input: `Name=~'A.*'`,
			},
			want: want{
				val: `(Name =~ "A.*")`,
			},
		},
		{
			name: "regex not",
			args: args{
				input: `Name!~'A.*'`,
			},
			want: want{
				val: `(Name !~ "A.*")`,
			},
		},
		{
			name: "regex raw string",
			args: args{
				input: "Name=~`^H.*d$`",
			},
			want: want{
				val: `(Name =~ "^H.*d$")`,
			},
		},
		{
			name: "regex with inline flag",
			args: args{
				input: "Name=~`(?i)A`",
			},
			want: want{
				val: `(Name =~ "(?i)A")`,
			},
		},
		{
			name: "negative regex with inline flag",
			args: args{
				input: "Name!~`(?i)A`",
			},
			want: want{
				val: `(Name !~ "(?i)A")`,
			},
		},
		// Numbers
		{
			name: "gt number",
			args: args{
				input: `HP>50`,
			},
			want: want{
				val: `(HP > 50)`,
			},
		},
		{
			name: "gte number",
			args: args{
				input: `MP>=100`,
			},
			want: want{
				val: `(MP >= 100)`,
			},
		},
		{
			name: "lt number float",
			args: args{
				input: `Rate<0.75`,
			},
			want: want{
				val: `(Rate < 0.75)`,
			},
		},
		{
			name: "hex float",
			args: args{
				input: `X==0x1.fp3`,
			},
			want: want{
				val: `(X == 0x1.fp3)`,
			},
		},
		// Durations
		{
			name: "duration gte",
			args: args{
				input: `Delay>=1h30m`,
			},
			want: want{
				val: `(Delay >= 1h30m)`,
			},
		},
		{
			name: "duration lt",
			args: args{
				input: `Timeout<500ms`,
			},
			want: want{
				val: `(Timeout < 500ms)`,
			},
		},
		{
			name: "duration micro",
			args: args{
				input: `Mic==4000μs`,
			},
			want: want{
				val: `(Mic == 4000μs)`,
			},
		},
		// Times
		{
			name: "time before",
			args: args{
				input: `Time>=2023-01-02T15:04:05Z`,
			},
			want: want{
				val: `(Time >= "2023-01-02T15:04:05Z")`,
			},
		},
		{
			name: "time after",
			args: args{
				input: `Time<2023-01-02T15:04:05Z`,
			},
			want: want{
				val: `(Time < "2023-01-02T15:04:05Z")`,
			},
		},
		// Booleans
		{
			name: "bool eq",
			args: args{
				input: `Flag==true`,
			},
			want: want{
				val: `(Flag == true)`,
			},
		},
		{
			name: "bool neq",
			args: args{
				input: `Flag!=False`,
			},
			want: want{
				val: `(Flag != false)`,
			},
		},
		// Logic and precedence
		{
			name: "and or precedence",
			args: args{
				input: `HP>50&&MP>=100||LP==0`,
			},
			want: want{
				val: `(((HP > 50) && (MP >= 100)) || (LP == 0))`,
			},
		},
		{
			name: "paren grouping",
			args: args{
				input: `(HP>50&&MP>=100)||LP==0`,
			},
			want: want{
				val: `(((HP > 50) && (MP >= 100)) || (LP == 0))`,
			},
		},
		{
			name: "not group",
			args: args{
				input: `!(SPD<20)`,
			},
			want: want{
				val: `(! (SPD < 20))`,
			},
		},
		{
			name: "complex",
			args: args{
				input: `Class=="軍師"&&Name=~'孔明'&&(HP>50&&MP>=100&&LP!=0)&&(MAG>=20||!(SPD<20))`,
			},
			want: want{
				val: `((((Class == "軍師") && (Name =~ "孔明")) && (((HP > 50) && (MP >= 100)) && (LP != 0))) && ((MAG >= 20) || (! (SPD < 20))))`,
			},
		},
		// Errors
		{
			name: "regex with a number pattern",
			args: args{
				input: `Name=~1`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:7: expected string pattern, got number: "1"`,
			},
		},
		{
			name: "regex empty pattern",
			args: args{
				input: `Name=~''`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:7: invalid regex "": empty pattern`,
			},
		},
		{
			name: "invalid regex",
			args: args{
				input: `Name=~'['`,
			},
			want: want{
				isErr: true,
				err:   "parse error at 1:7: invalid regex \"[\": error parsing regexp: missing closing ]: `[`",
			},
		},
		{
			name: "missing op",
			args: args{
				input: `HP 50`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:4: expected comparison operator, got number: "50"`,
			},
		},
		{
			name: "missing right",
			args: args{
				input: `HP>`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:4: expected value, got EOF: ""`,
			},
		},
		{
			name: "unexpected trailing",
			args: args{
				input: `HP>50 extra`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:7: unexpected token after parsing: extra`,
			},
		},
		{
			name: "leading not without operand",
			args: args{
				input: `!`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:2: expected left parenthesis or identifier, got EOF: ""`,
			},
		},
		{
			name: "empty",
			args: args{
				input: ``,
			},
			want: want{
				isErr: true,
				err:   `parse error: empty input`,
			},
		},
		{
			name: "unclosed paren",
			args: args{
				input: `(HP>1`,
			},
			want: want{
				isErr: true,
				err:   `token error at 1:6: unclosed left parenthesis`,
			},
		},
		{
			name: "extra right paren",
			args: args{
				input: `HP>1)`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:5: unexpected token after parsing: )`,
			},
		},
		{
			name: "double logical op",
			args: args{
				input: `HP>1&&||MP>2`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:7: expected left parenthesis or identifier, got logical OR operator: "||"`,
			},
		},
		{
			name: "non ident left",
			args: args{
				input: `123==456`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:1: expected left parenthesis or identifier, got number: "123"`,
			},
		},
		{
			name: "unterminated regex string",
			args: args{
				input: `Name=~'abc`,
			},
			want: want{
				isErr: true,
				err:   `token error at 1:11: unterminated quoted string`,
			},
		},
		{
			name: "number then missing op",
			args: args{
				input: `HP50`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:5: expected comparison operator, got EOF: ""`,
			},
		},
		{
			name: "duration segment missing unit",
			args: args{
				input: `Delay==1h30`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:10: unexpected token after parsing: 30`,
			},
		},
		{
			name: "expect mismatch right paren",
			args: args{
				input: `(HP>1 Name==X)`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:7: expected right parenthesis, got identifier: "Name"`,
			},
		},
		{
			name: "expect mismatch nested right paren",
			args: args{
				input: `((HP>1) Name==X)`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:9: expected right parenthesis, got identifier: "Name"`,
			},
		},
		{
			name: "bare dot as number",
			args: args{
				input: `HP > .`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:6: invalid number "."`,
			},
		},
		{
			name: "bare sign as number",
			args: args{
				input: `HP > -`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:6: invalid number "-"`,
			},
		},
		{
			name: "bare dot after wide identifier",
			args: args{
				input: `名前 == .`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:9: invalid number "."`,
			},
		},
		{
			name: "bare sign on second line",
			args: args{
				input: "a\n== -",
			},
			want: want{
				isErr: true,
				err:   `parse error at 2:4: invalid number "-"`,
			},
		},
		{
			name: "number with two dots",
			args: args{
				input: `HP > 1.2.3`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:9: unexpected token after parsing: .3`,
			},
		},
		{
			name: "base prefix without digits",
			args: args{
				input: `HP > 0x`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:6: invalid number "0x"`,
			},
		},
		{
			name: "exponent without digits",
			args: args{
				input: `HP > 1e`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:6: invalid number "1e"`,
			},
		},
		{
			name: "date only",
			args: args{
				input: `Time >= 2023-01-02`,
			},
			want: want{
				val: `(Time >= "2023-01-02")`,
			},
		},
		{
			name: "time without zone",
			args: args{
				input: `Time > 2023-01-02T15:04:05`,
			},
			want: want{
				val: `(Time > "2023-01-02T15:04:05")`,
			},
		},
		{
			name: "date followed by T alone",
			args: args{
				input: `Time > 2023-01-02T`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:18: unexpected token after parsing: T`,
			},
		},
		{
			name: "date out of range",
			args: args{
				input: `Time > 2023-13-01`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:8: invalid time "2023-13-01"`,
			},
		},
		{
			name: "time out of range",
			args: args{
				input: `Time > 2023-13-01T00:00:00Z`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:8: invalid time "2023-13-01T00:00:00Z"`,
			},
		},
		{
			name: "duration with repeated fraction",
			args: args{
				input: `Duration > 1.5.5s`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:12: invalid duration "1.5.5s"`,
			},
		},
		{
			name: "parseExpr initial next failure",
			args: args{
				input: `#&&HP>1`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:1: expected left parenthesis or identifier, got error: "unexpected character U+0023 '#'"`,
			},
		},
		{
			name: "parseAND right side next failure",
			args: args{
				input: `HP>1&&#`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:7: expected left parenthesis or identifier, got error: "unexpected character U+0023 '#'"`,
			},
		},
		{
			name: "parseAND operator malformed",
			args: args{
				input: `HP>1&X==1`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:6: unexpected token after parsing: unexpected character 'X' after '&'`,
			},
		},
		{
			name: "parseNOT next failure",
			args: args{
				input: `!#`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:2: expected left parenthesis or identifier, got error: "unexpected character U+0023 '#'"`,
			},
		},
		{
			name: "parsePrimary inner expr failure",
			args: args{
				input: `(#)`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:2: expected left parenthesis or identifier, got error: "unexpected character U+0023 '#'"`,
			},
		},
		{
			name: "parsePrimary parseExpr failure",
			args: args{
				input: `(##)`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:2: expected left parenthesis or identifier, got error: "unexpected character U+0023 '#'"`,
			},
		},
		{
			name: "or missing right operand",
			args: args{
				input: `HP>1 ||`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:8: expected left parenthesis or identifier, got EOF: ""`,
			},
		},
		{
			name: "more identifiers and nodes than the inline buffers",
			args: args{
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
			},
			want: want{
				val: func() string {
					s := "(F0 > 0)"
					for i := 1; i < 20; i++ {
						s = fmt.Sprintf("(%s && (F%d > %d))", s, i, i)
					}
					return s
				}(),
			},
		},
		{
			name: "input too long",
			args: args{
				input: "HP > " + strings.Repeat("1", MaxInput),
			},
			want: want{
				isErr: true,
				err:   `parse error: input too long: 1048581 bytes exceeds limit 1048576`,
			},
		},
		{
			name: "parseComparison expect ident failure",
			args: args{
				input: `==1`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:1: expected left parenthesis or identifier, got "equal to" operator: "=="`,
			},
		},
		{
			name: "parseComparison operator next failure",
			args: args{
				input: `A$1`,
			},
			want: want{
				isErr: true,
				err:   `token error at 1:2: unexpected character U+0024 '$'`,
			},
		},
		{
			name: "parseComparison value next failure",
			args: args{
				input: `A==#`,
			},
			want: want{
				isErr: true,
				err:   `token error at 1:4: unexpected character U+0023 '#'`,
			},
		},
		// Parenthesis limit
		{
			name: "paren limit ok (256)",
			args: args{
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
			},
			want: want{
				val: `(HP > 1)`,
			},
		},
		{
			name: "paren limit ng (257)",
			args: args{
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
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:257: too many parentheses: exceeded limit 256`,
			},
		},
		{
			name: "duration invalid at eval",
			args: args{
				input: `Duration>bad`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:10: expected value, got identifier: "bad"`,
			},
		},
		{
			name: "invalid number right",
			args: args{
				input: `Int>1+0`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:6: unexpected token after parsing: +0`,
			},
		},
		{
			name: "invalid duration right",
			args: args{
				input: `Duration>1xs`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:11: unexpected token after parsing: xs`,
			},
		},
		{
			name: "regex compile error",
			args: args{
				input: `String=~"["`,
			},
			want: want{
				isErr: true,
				err:   "parse error at 1:9: invalid regex \"[\": error parsing regexp: missing closing ]: `[`",
			},
		},
		{
			name: "regex not found",
			args: args{
				input: `String=~""`,
			},
			want: want{
				isErr: true,
				err:   `parse error at 1:9: invalid regex "": empty pattern`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Parse(test.args.input)
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
			if got := repr(got); got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	type args struct {
		input string
	}
	type want struct {
		val    string
		panics bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "valid input",
			args: args{
				input: `Name == "a"`,
			},
			want: want{
				val: `(Name == "a")`,
			},
		},
		{
			name: "invalid input",
			args: args{
				input: `Name ==`,
			},
			want: want{
				panics: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if panics := recover() != nil; panics != test.want.panics {
					t.Errorf("panic mismatch\ngot=%v\nwant=%v\n", panics, test.want.panics)
				}
			}()
			if got := repr(MustParse(test.args.input)); got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestExpr_Eval(t *testing.T) {
	type fields struct {
		expr expr
	}
	type args struct {
		r Resolver
	}
	type want struct {
		val   bool
		isErr bool
		err   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		// String comparison
		{
			name: "string eq",
			fields: fields{
				expr: MustParse(`String=="HelloWorld"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "string eq false",
			fields: fields{
				expr: MustParse(`String=="X"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "string neq",
			fields: fields{
				expr: MustParse(`String!="X"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "regex ignoring case",
			fields: fields{
				expr: MustParse(`String=~"(?i)^helloworld$"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "regex ignoring case false",
			fields: fields{
				expr: MustParse(`String=~"(?i)^hellox$"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "negative regex ignoring case",
			fields: fields{
				expr: MustParse(`String!~"(?i)^hellox$"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "negative regex ignoring case false",
			fields: fields{
				expr: MustParse(`String!~"(?i)^helloworld$"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "regex match",
			fields: fields{
				expr: MustParse(`String=~"^Hello"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "regex no match",
			fields: fields{
				expr: MustParse(`String=~"world$"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "regex neg match",
			fields: fields{
				expr: MustParse(`String!~"^Hello"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		// Numeric comparisons
		{
			name: "int gt",
			fields: fields{
				expr: MustParse(`Int>40`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int gt false",
			fields: fields{
				expr: MustParse(`Int>100`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "int eq",
			fields: fields{
				expr: MustParse(`Int==42`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int eq false",
			fields: fields{
				expr: MustParse(`Int==41`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "int neq",
			fields: fields{
				expr: MustParse(`Int!=41`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int neq false",
			fields: fields{
				expr: MustParse(`Int!=42`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "int gte false",
			fields: fields{
				expr: MustParse(`Int>=100`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "int lt false",
			fields: fields{
				expr: MustParse(`Int<40`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "int lte false",
			fields: fields{
				expr: MustParse(`Int<=41`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "int8 gt",
			fields: fields{
				expr: MustParse(`Int8>1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int16 gt",
			fields: fields{
				expr: MustParse(`Int16>1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int32 gt",
			fields: fields{
				expr: MustParse(`Int32>1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "int64 gt",
			fields: fields{
				expr: MustParse(`Int64>1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint gt",
			fields: fields{
				expr: MustParse(`Uint>1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint8 gt",
			fields: fields{
				expr: MustParse(`Uint8>1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint16 gt",
			fields: fields{
				expr: MustParse(`Uint16>1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint32 gt",
			fields: fields{
				expr: MustParse(`Uint32>1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uint64 gt",
			fields: fields{
				expr: MustParse(`Uint64>1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float32 gt",
			fields: fields{
				expr: MustParse(`Float32>2`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float lt",
			fields: fields{
				expr: MustParse(`Float64<3.2`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float gte",
			fields: fields{
				expr: MustParse(`Float64>=3.14`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float eq epsilon",
			fields: fields{
				expr: MustParse(`Float64==3.1400000001`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "float neq epsilon",
			fields: fields{
				expr: MustParse(`Float64!=3.1401`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		// Duration
		{
			name: "duration gt",
			fields: fields{
				expr: MustParse(`Duration>1s`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "duration gte false",
			fields: fields{
				expr: MustParse(`Duration>=2s`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "duration gt false",
			fields: fields{
				expr: MustParse(`Duration>2s`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "duration gte true",
			fields: fields{
				expr: MustParse(`Duration>=1500ms`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "duration lt",
			fields: fields{
				expr: MustParse(`Duration<2s`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "duration lt false",
			fields: fields{
				expr: MustParse(`Duration<1s`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "duration lte true",
			fields: fields{
				expr: MustParse(`Duration<=1500ms`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "duration lte false",
			fields: fields{
				expr: MustParse(`Duration<=1s`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "duration eq",
			fields: fields{
				expr: MustParse(`Duration==1500ms`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "duration eq false",
			fields: fields{
				expr: MustParse(`Duration==2s`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "duration neq",
			fields: fields{
				expr: MustParse(`Duration!=2s`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "invalid operator duration",
			fields: fields{
				expr: MustParse(`Duration=~"1500ms"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:9: invalid operator for duration value "=~"`,
			},
		},
		{
			name: "duration neq false",
			fields: fields{
				expr: MustParse(`Duration!=1500ms`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "other types compare by their formatted text",
			fields: fields{
				expr: MustParse(`Struct=="{1}"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "string of a lone sign stays a string",
			fields: fields{
				expr: MustParse(`String!="-"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "duration string converted at eval",
			fields: fields{
				expr: MustParse(`Duration>'0'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "number string converted at eval",
			fields: fields{
				expr: MustParse(`Float64<'Inf'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		// Time
		{
			name: "time gt",
			fields: fields{
				expr: MustParse(`Time>'2024-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time gt false",
			fields: fields{
				expr: MustParse(`Time>'2026-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "time gte",
			fields: fields{
				expr: MustParse(`Time>='2025-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time gte false",
			fields: fields{
				expr: MustParse(`Time>='2026-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "time lt",
			fields: fields{
				expr: MustParse(`Time<'2026-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time lt false",
			fields: fields{
				expr: MustParse(`Time<'2024-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "time lte",
			fields: fields{
				expr: MustParse(`Time<='2025-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time lte false",
			fields: fields{
				expr: MustParse(`Time<='2024-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "time eq",
			fields: fields{
				expr: MustParse(`Time=='2025-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time eq false",
			fields: fields{
				expr: MustParse(`Time=='2024-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "time neq",
			fields: fields{
				expr: MustParse(`Time!='2024-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time neq false",
			fields: fields{
				expr: MustParse(`Time!='2025-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		// Combined logicals
		{
			name: "combined logicals",
			fields: fields{
				expr: MustParse(`String=="HelloWorld" && Int==42 || Float64<3.0`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "combined logicals false",
			fields: fields{
				expr: MustParse(`String=="HelloWorld" && Int==41 || Float64<3.0`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "combined logicals with parens",
			fields: fields{
				expr: MustParse(`(String=="HelloWorld" && Int==41) || (Float64<3.2 && Bool==true)`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "combined logicals with parens false",
			fields: fields{
				expr: MustParse(`(String=="HelloWorld" && Int==41) || (Float64<3.0 && Bool==true)`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		// Bool
		{
			name: "bool eq",
			fields: fields{
				expr: MustParse(`Bool==true`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "bool eq uppercase literal",
			fields: fields{
				expr: MustParse(`Bool==TRUE`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "bool neq capitalized literal",
			fields: fields{
				expr: MustParse(`Bool!=False`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "bool neq",
			fields: fields{
				expr: MustParse(`Bool!=false`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "and true",
			fields: fields{
				expr: MustParse(`Int>40&&Float64<4`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "and false",
			fields: fields{
				expr: MustParse(`Int>40&&Float64>4`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "or true",
			fields: fields{
				expr: MustParse(`Int>100||Float64<4`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "or short-circuit left true",
			fields: fields{
				expr: MustParse(`Bool==true || Invalid==1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "or left error",
			fields: fields{
				expr: MustParse(`Invalid==1 || Bool==true`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:1: unknown identifier "Invalid"`,
			},
		},
		{
			name: "same identifier referenced twice",
			fields: fields{
				expr: MustParse(`Int>40 && Int<50`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "more identifiers than the inline cache",
			fields: fields{
				expr: MustParse(func() string {
					var b strings.Builder
					for i := range 17 {
						if i > 0 {
							b.WriteString(" && ")
						}
						fmt.Fprintf(&b, "F%d == %d", i, i)
					}
					return b.String() + " && F0 == 0"
				}()).expr,
			},
			args: args{
				r: func() testResolver {
					t := testResolver{}
					for i := range 17 {
						t[fmt.Sprintf("F%d", i)] = i
					}
					return t
				}(),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "not true->false",
			fields: fields{
				expr: MustParse(`!(Int>40)`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "not false->true",
			fields: fields{
				expr: MustParse(`!(Int<40)`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "and short-circuit left false",
			fields: fields{
				expr: MustParse(`Int>100 && Invalid==1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: false,
			},
		},
		{
			name: "not inner eval error",
			fields: fields{
				expr: MustParse(`!(Invalid==1)`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:3: unknown identifier "Invalid"`,
			},
		},
		// Mixed
		{
			name: "precedence",
			fields: fields{
				expr: MustParse(`Int>40&&Float64<4||Bool==false`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		// Errors
		{
			name: "binary left eval error",
			fields: fields{
				expr: MustParse(`Unknown==1 && Bool==true`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:1: unknown identifier "Unknown"`,
			},
		},
		{
			name: "binary right eval error",
			fields: fields{
				expr: MustParse(`Bool==true && Unknown==1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:15: unknown identifier "Unknown"`,
			},
		},
		{
			name: "unknown identifier",
			fields: fields{
				expr: MustParse(`Invalid==1`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:1: unknown identifier "Invalid"`,
			},
		},
		{
			name: "zero value from resolver",
			fields: fields{
				expr: MustParse(`Int==1`).expr,
			},
			args: args{
				r: zeroResolver{},
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:1: unknown identifier "Int"`,
			},
		},
		{
			name: "type mismatch 1",
			fields: fields{
				expr: MustParse(`Int>"abc"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:5: invalid number "abc"`,
			},
		},
		{
			name: "type mismatch 2",
			fields: fields{
				expr: MustParse(`String>"HelloWorld"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:7: invalid operator for string value ">"`,
			},
		},
		{
			name: "type mismatch 3",
			fields: fields{
				expr: MustParse(`Int=~"42"`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:4: invalid operator for number value "=~"`,
			},
		},
		{
			name: "time eq bare rfc3339",
			fields: fields{
				expr: MustParse(`Time==2025-01-01T00:00:00Z`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time eq bare datetime without zone",
			fields: fields{
				expr: MustParse(`Time==2025-01-01T00:00:00`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time eq bare date",
			fields: fields{
				expr: MustParse(`Time==2025-01-01`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time eq date only",
			fields: fields{
				expr: MustParse(`Time=='2025-01-01'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time gt date only",
			fields: fields{
				expr: MustParse(`Time>'2024-12-31'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time eq datetime with space",
			fields: fields{
				expr: MustParse(`Time=='2025-01-01 00:00:00'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time eq datetime without zone",
			fields: fields{
				expr: MustParse(`Time=='2025-01-01T00:00:00'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time lt datetime with fraction and no zone",
			fields: fields{
				expr: MustParse(`Time<'2025-01-01 00:00:00.5'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time eq rfc1123 string",
			fields: fields{
				expr: MustParse(`Time=='Wed, 01 Jan 2025 00:00:00 UTC'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time eq rfc1123z string",
			fields: fields{
				expr: MustParse(`Time=='Wed, 01 Jan 2025 09:00:00 +0900'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time eq rfc850 string",
			fields: fields{
				expr: MustParse(`Time=='Wednesday, 01-Jan-25 00:00:00 UTC'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time eq rfc822 string",
			fields: fields{
				expr: MustParse(`Time=='01 Jan 25 00:00 UTC'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time eq unix seconds with underscores",
			fields: fields{
				expr: MustParse(`Time==1_735_689_600`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time with zone abbreviation other than utc",
			fields: fields{
				expr: MustParse(`Time=='Wed, 01 Jan 2025 00:00:00 EST'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:7: invalid time "Wed, 01 Jan 2025 00:00:00 EST"`,
			},
		},
		{
			name: "time eq unix seconds as string",
			fields: fields{
				expr: MustParse(`Time=='1735689600'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time eq unix seconds as number",
			fields: fields{
				expr: MustParse(`Time==1735689600`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time lt unix seconds as number",
			fields: fields{
				expr: MustParse(`Time<1735689601`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "time with fractional unix seconds",
			fields: fields{
				expr: MustParse(`Time==1735689600.5`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:7: invalid time "1735689600.5"`,
			},
		},
		{
			name: "time with out of range date",
			fields: fields{
				expr: MustParse(`Time>'2025-13-01'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:6: invalid time "2025-13-01"`,
			},
		},
		{
			name: "time with clock only",
			fields: fields{
				expr: MustParse(`Time>'12:00:00'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:6: invalid time "12:00:00"`,
			},
		},
		{
			name: "invalid time",
			fields: fields{
				expr: MustParse(`Time>'invalid-time'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:6: invalid time "invalid-time"`,
			},
		},
		{
			name: "time invalid operator",
			fields: fields{
				expr: MustParse(`Time=~'2025-01-01T00:00:00Z'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:5: invalid operator for time value "=~"`,
			},
		},
		{
			name: "invalid duration",
			fields: fields{
				expr: MustParse(`Duration>'bad-duration'`).expr,
			},
			args: args{
				r: testObject,
			},
			want: want{
				isErr: true,
				err:   `eval error at 1:10: invalid duration "bad-duration"`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := &Expr{
				expr: test.fields.expr,
			}
			got, err := e.Eval(test.args.r)
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
