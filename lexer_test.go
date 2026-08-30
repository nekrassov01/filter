package filter

import (
	"reflect"
	"testing"
)

func Test_newLexer(t *testing.T) {
	type args struct {
		input string
	}
	type want struct {
		val lexer
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
				val: lexer{
					input:     "",
					state:     stateStmt,
					line:      1,
					startLine: 1,
					col:       1,
					startCol:  1,
				},
			},
		},
		{
			name: "ascii input",
			args: args{
				input: `Name == "a"`,
			},
			want: want{
				val: lexer{
					input:     `Name == "a"`,
					state:     stateStmt,
					line:      1,
					startLine: 1,
					col:       1,
					startCol:  1,
				},
			},
		},
		{
			name: "unicode input starts at the same position",
			args: args{
				input: "軍師",
			},
			want: want{
				val: lexer{
					input:     "軍師",
					state:     stateStmt,
					line:      1,
					startLine: 1,
					col:       1,
					startCol:  1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newLexer(test.args.input)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_lexer_lexStmt(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty input emits EOF",
			fields: fields{
				input: "",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenEOF,
					v:    "",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "space emits nothing",
			fields: fields{
				input: " ",
			},
			want: want{
				val: stateStmt,
				tok: token{},
			},
		},
		{
			name: "double quoted string",
			fields: fields{
				input: `"a"`,
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `"a"`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "single quoted string",
			fields: fields{
				input: `'a'`,
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `'a'`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "raw string",
			fields: fields{
				input: "`a`",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenRawString,
					v:    "`a`",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "left parenthesis",
			fields: fields{
				input: "(",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenLparen,
					v:    "(",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "right parenthesis",
			fields: fields{
				input: ")",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenRparen,
					v:    ")",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "equal",
			fields: fields{
				input: "==",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenEQ,
					v:    "==",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "not",
			fields: fields{
				input: "!",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNOT,
					v:    "!",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "less than",
			fields: fields{
				input: "<",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenLT,
					v:    "<",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "greater than or equal",
			fields: fields{
				input: ">=",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenGTE,
					v:    ">=",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "and",
			fields: fields{
				input: "&&",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenAND,
					v:    "&&",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "or",
			fields: fields{
				input: "||",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenOR,
					v:    "||",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "digit starts a number",
			fields: fields{
				input: "1",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "1",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "sign starts a number",
			fields: fields{
				input: "-1",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "-1",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "dot starts a number",
			fields: fields{
				input: ".5",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    ".5",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "lone sign is still a number token",
			fields: fields{
				input: "+",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "+",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "letter starts an identifier",
			fields: fields{
				input: "abc",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "abc",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "underscore starts an identifier",
			fields: fields{
				input: "_x",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "_x",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "non-ascii letter starts an identifier",
			fields: fields{
				input: "軍師",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "軍師",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "keyword becomes a bool",
			fields: fields{
				input: "true",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenBool,
					v:    "true",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "unexpected ascii character",
			fields: fields{
				input: "#",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected character U+0023 '#'",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "unexpected non-ascii character",
			fields: fields{
				input: "→",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected character U+2192 '→'",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "digit after a non-ascii identifier is part of it",
			fields: fields{
				input: "名前1 ",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "名前1",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			got := l.lexStmt()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexEOF(t *testing.T) {
	type fields struct {
		input      string
		parenDepth int
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "balanced parentheses emit EOF",
			fields: fields{
				input:      "",
				parenDepth: 0,
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenEOF,
					v:    "",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "EOF follows the last token",
			fields: fields{
				input:      "a",
				parenDepth: 0,
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenEOF,
					v:    "",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "EOF follows a wide identifier",
			fields: fields{
				input:      "軍師",
				parenDepth: 0,
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenEOF,
					v:    "",
					pos:  6,
					line: 1,
					col:  5,
				},
			},
		},
		{
			name: "unclosed left parenthesis",
			fields: fields{
				input:      "",
				parenDepth: 1,
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unclosed left parenthesis",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "deeply unclosed left parentheses",
			fields: fields{
				input:      "",
				parenDepth: 3,
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unclosed left parenthesis",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "unexpected right parenthesis",
			fields: fields{
				input:      "",
				parenDepth: -1,
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected right parenthesis",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			if test.fields.input != "" {
				l.nextToken()
			}
			l.parenDepth = test.fields.parenDepth
			got := l.lexEOF()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexSpace(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val  state
		pos  int32
		line int32
		col  int32
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "single space at end of input",
			fields: fields{
				input: " ",
			},
			want: want{
				val:  stateStmt,
				pos:  1,
				line: 1,
				col:  2,
			},
		},
		{
			name: "spaces before a token",
			fields: fields{
				input: "  a",
			},
			want: want{
				val:  stateStmt,
				pos:  2,
				line: 1,
				col:  3,
			},
		},
		{
			name: "tab and carriage return advance the column",
			fields: fields{
				input: " \t\r",
			},
			want: want{
				val:  stateStmt,
				pos:  3,
				line: 1,
				col:  4,
			},
		},
		{
			name: "newline advances the line and resets the column",
			fields: fields{
				input: " \n",
			},
			want: want{
				val:  stateStmt,
				pos:  2,
				line: 2,
				col:  1,
			},
		},
		{
			name: "leading newline consumed by next",
			fields: fields{
				input: "\n",
			},
			want: want{
				val:  stateStmt,
				pos:  1,
				line: 2,
				col:  1,
			},
		},
		{
			name: "several newlines then a token",
			fields: fields{
				input: " \n\n x",
			},
			want: want{
				val:  stateStmt,
				pos:  4,
				line: 3,
				col:  2,
			},
		},
		{
			name: "stops before a non-ascii character",
			fields: fields{
				input: " 軍",
			},
			want: want{
				val:  stateStmt,
				pos:  1,
				line: 1,
				col:  2,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexSpace()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.pos != test.want.pos || l.line != test.want.line || l.col != test.want.col {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.mark(), mark{pos: test.want.pos, line: test.want.line, col: test.want.col})
			}
			if l.startPos != l.pos || l.startLine != l.line || l.startCol != l.col {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", mark{pos: l.startPos, line: l.startLine, col: l.startCol}, l.mark())
			}
			if l.hasNext {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.hasNext, false)
			}
		})
	}
}

func Test_lexer_lexLparen(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val        state
		tok        token
		parenDepth int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "single left parenthesis",
			fields: fields{
				input: "(",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenLparen,
					v:    "(",
					pos:  0,
					line: 1,
					col:  1,
				},
				parenDepth: 1,
			},
		},
		{
			name: "token spans only the parenthesis",
			fields: fields{
				input: "(a",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenLparen,
					v:    "(",
					pos:  0,
					line: 1,
					col:  1,
				},
				parenDepth: 1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexLparen()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
			if l.parenDepth != test.want.parenDepth {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.parenDepth, test.want.parenDepth)
			}
		})
	}
}

func Test_lexer_lexRparen(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val        state
		tok        token
		parenDepth int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "single right parenthesis",
			fields: fields{
				input: ")",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenRparen,
					v:    ")",
					pos:  0,
					line: 1,
					col:  1,
				},
				parenDepth: -1,
			},
		},
		{
			name: "token spans only the parenthesis",
			fields: fields{
				input: ")&&",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenRparen,
					v:    ")",
					pos:  0,
					line: 1,
					col:  1,
				},
				parenDepth: -1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexRparen()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
			if l.parenDepth != test.want.parenDepth {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.parenDepth, test.want.parenDepth)
			}
		})
	}
}

func Test_lexer_lexEQ(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "equal",
			fields: fields{
				input: "==",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenEQ,
					v:    "==",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "regex equal",
			fields: fields{
				input: "=~",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenREQ,
					v:    "=~",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "token stops after the operator",
			fields: fields{
				input: "==x",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenEQ,
					v:    "==",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "end of input after equal",
			fields: fields{
				input: "=",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected end of input after '='",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "ascii character after equal",
			fields: fields{
				input: "=x",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected character 'x' after '='",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "space after equal",
			fields: fields{
				input: "= =",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected character ' ' after '='",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "non-ascii character after equal",
			fields: fields{
				input: "=軍",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected character '軍' after '='",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexEQ()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexNOT(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "not equal",
			fields: fields{
				input: "!=",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNEQ,
					v:    "!=",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "regex not equal",
			fields: fields{
				input: "!~",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNREQ,
					v:    "!~",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "not at end of input",
			fields: fields{
				input: "!",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNOT,
					v:    "!",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "not before an identifier",
			fields: fields{
				input: "!a",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNOT,
					v:    "!",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "not before a parenthesis",
			fields: fields{
				input: "!(",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNOT,
					v:    "!",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "double not emits one token",
			fields: fields{
				input: "!!",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNOT,
					v:    "!",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexNOT()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexLT(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "less than or equal",
			fields: fields{
				input: "<=",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenLTE,
					v:    "<=",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "less than at end of input",
			fields: fields{
				input: "<",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenLT,
					v:    "<",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "less than before a number",
			fields: fields{
				input: "<1",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenLT,
					v:    "<",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "double less than emits one token",
			fields: fields{
				input: "<<",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenLT,
					v:    "<",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexLT()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexGT(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "greater than or equal",
			fields: fields{
				input: ">=",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenGTE,
					v:    ">=",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "greater than at end of input",
			fields: fields{
				input: ">",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenGT,
					v:    ">",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "greater than before a number",
			fields: fields{
				input: ">1",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenGT,
					v:    ">",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "double greater than emits one token",
			fields: fields{
				input: ">>",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenGT,
					v:    ">",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexGT()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexAND(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "and",
			fields: fields{
				input: "&&",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenAND,
					v:    "&&",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "token stops after the operator",
			fields: fields{
				input: "&&&",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenAND,
					v:    "&&",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "end of input after ampersand",
			fields: fields{
				input: "&",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected end of input after '&'",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "ascii character after ampersand",
			fields: fields{
				input: "&x",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected character 'x' after '&'",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "pipe after ampersand",
			fields: fields{
				input: "&|",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected character '|' after '&'",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "non-ascii character after ampersand",
			fields: fields{
				input: "&軍",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected character '軍' after '&'",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexAND()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexOR(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "or",
			fields: fields{
				input: "||",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenOR,
					v:    "||",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "token stops after the operator",
			fields: fields{
				input: "|||",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenOR,
					v:    "||",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "end of input after pipe",
			fields: fields{
				input: "|",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected end of input after '|'",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "ascii character after pipe",
			fields: fields{
				input: "|x",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected character 'x' after '|'",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "ampersand after pipe",
			fields: fields{
				input: "|&",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected character '&' after '|'",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "non-ascii character after pipe",
			fields: fields{
				input: "|軍",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unexpected character '軍' after '|'",
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexOR()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexDoubleQuotedString(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "string",
			fields: fields{
				input: `"abc"`,
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `"abc"`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "single quote inside is literal",
			fields: fields{
				input: `"a'b"`,
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `"a'b"`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "escaped double quote does not close",
			fields: fields{
				input: `"a\"b"`,
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `"a\"b"`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "unterminated",
			fields: fields{
				input: `"abc`,
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unterminated quoted string",
					pos:  4,
					line: 1,
					col:  5,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexDoubleQuotedString()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexSingleQuotedString(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "string",
			fields: fields{
				input: `'abc'`,
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `'abc'`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "double quote inside is literal",
			fields: fields{
				input: `'a"b'`,
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `'a"b'`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "escaped single quote does not close",
			fields: fields{
				input: `'a\'b'`,
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `'a\'b'`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "unterminated",
			fields: fields{
				input: `'abc`,
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unterminated quoted string",
					pos:  4,
					line: 1,
					col:  5,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexSingleQuotedString()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexString(t *testing.T) {
	type fields struct {
		input string
	}
	type args struct {
		quote rune
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "double quoted",
			fields: fields{
				input: `"abc"`,
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `"abc"`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "single quoted",
			fields: fields{
				input: `'abc'`,
			},
			args: args{
				quote: '\'',
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `'abc'`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "empty string",
			fields: fields{
				input: `""`,
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `""`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "token stops at the closing quote",
			fields: fields{
				input: `"a" == b`,
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `"a"`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "non-ascii content",
			fields: fields{
				input: `"軍師"`,
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `"軍師"`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "valid escapes",
			fields: fields{
				input: `"\n\t\\\"\x41A\0"`,
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `"\n\t\\\"\x41A\0"`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "other quote is literal",
			fields: fields{
				input: `"a'b"`,
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenString,
					v:    `"a'b"`,
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "unterminated at end of input",
			fields: fields{
				input: `"abc`,
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unterminated quoted string",
					pos:  4,
					line: 1,
					col:  5,
				},
			},
		},
		{
			name: "unterminated at end of input after a wide rune",
			fields: fields{
				input: `"軍`,
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unterminated quoted string",
					pos:  4,
					line: 1,
					col:  4,
				},
			},
		},
		{
			name: "newline terminates the string",
			fields: fields{
				input: "\"a\nb\"",
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unterminated quoted string",
					pos:  3,
					line: 2,
					col:  1,
				},
			},
		},
		{
			name: "invalid escape",
			fields: fields{
				input: `"\z"`,
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "invalid escape sequence in string",
					pos:  3,
					line: 1,
					col:  4,
				},
			},
		},
		{
			name: "short hex escape",
			fields: fields{
				input: `"\x4"`,
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "invalid escape sequence in string",
					pos:  5,
					line: 1,
					col:  6,
				},
			},
		},
		{
			name: "backslash at end of input",
			fields: fields{
				input: `"\`,
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "invalid escape sequence in string",
					pos:  2,
					line: 1,
					col:  3,
				},
			},
		},
		{
			name: "invalid utf8",
			fields: fields{
				input: "\"\xff\"",
			},
			args: args{
				quote: '"',
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "invalid utf8 encoding in string",
					pos:  2,
					line: 1,
					col:  3,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexString(test.args.quote)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexRawString(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "raw string",
			fields: fields{
				input: "`abc`",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenRawString,
					v:    "`abc`",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "empty raw string",
			fields: fields{
				input: "``",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenRawString,
					v:    "``",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "quotes inside are literal",
			fields: fields{
				input: "`\"'`",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenRawString,
					v:    "`\"'`",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "backslash is literal",
			fields: fields{
				input: "`a\\z`",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenRawString,
					v:    "`a\\z`",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "newline is allowed",
			fields: fields{
				input: "`a\nb`",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenRawString,
					v:    "`a\nb`",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "token stops at the closing backtick",
			fields: fields{
				input: "`a` == b",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenRawString,
					v:    "`a`",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "non-ascii content",
			fields: fields{
				input: "`軍師`",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenRawString,
					v:    "`軍師`",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "unterminated",
			fields: fields{
				input: "`abc",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unterminated raw string",
					pos:  4,
					line: 1,
					col:  5,
				},
			},
		},
		{
			name: "unterminated after a newline",
			fields: fields{
				input: "`a\n",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "unterminated raw string",
					pos:  3,
					line: 2,
					col:  1,
				},
			},
		},
		{
			name: "invalid utf8",
			fields: fields{
				input: "`\xff`",
			},
			want: want{
				val: stateDone,
				tok: token{
					typ:  tokenError,
					v:    "invalid utf8 encoding in raw string",
					pos:  2,
					line: 1,
					col:  3,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexRawString()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexNumber(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "integer",
			fields: fields{
				input: "1",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "1",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "positive sign",
			fields: fields{
				input: "+1",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "+1",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "negative float",
			fields: fields{
				input: "-1.5",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "-1.5",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "leading dot",
			fields: fields{
				input: ".5",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    ".5",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "exponent",
			fields: fields{
				input: "1e3",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "1e3",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "hex",
			fields: fields{
				input: "0x1f",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "0x1f",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "underscore separator",
			fields: fields{
				input: "1_000",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "1_000",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "lone sign",
			fields: fields{
				input: "+",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "+",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "lone dot",
			fields: fields{
				input: ".",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    ".",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "number stops before a letter",
			fields: fields{
				input: "1x",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "1",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "second dot ends the number",
			fields: fields{
				input: "1.2.3",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "1.2",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "duration",
			fields: fields{
				input: "10s",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenDuration,
					v:    "10s",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "compound duration",
			fields: fields{
				input: "1h30m",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenDuration,
					v:    "1h30m",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "fractional duration",
			fields: fields{
				input: "1.5s",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenDuration,
					v:    "1.5s",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "negative duration",
			fields: fields{
				input: "-5m",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenDuration,
					v:    "-5m",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "millisecond duration",
			fields: fields{
				input: "1ms",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenDuration,
					v:    "1ms",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "duration stops before a non-unit",
			fields: fields{
				input: "10s)",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenDuration,
					v:    "10s",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "duration with two dots is left to the parser",
			fields: fields{
				input: "1.5.5s",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenDuration,
					v:    "1.5.5s",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "date",
			fields: fields{
				input: "2025-01-01",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenTime,
					v:    "2025-01-01",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "date time",
			fields: fields{
				input: "2025-01-01T00:00:00Z",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenTime,
					v:    "2025-01-01T00:00:00Z",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "date time with offset",
			fields: fields{
				input: "2025-01-01T09:00:00+09:00",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenTime,
					v:    "2025-01-01T09:00:00+09:00",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "incomplete clock keeps the date",
			fields: fields{
				input: "2025-01-01T00",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenTime,
					v:    "2025-01-01",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "partial date is a number",
			fields: fields{
				input: "2025-01",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "2025",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "position after earlier tokens",
			fields: fields{
				input: "軍 1",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenNumber,
					v:    "1",
					pos:  4,
					line: 1,
					col:  4,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			for !isNumberStart(l.next()) {
				l.ignore()
			}
			got := l.lexNumber()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_lexKeywordOrIdent(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val state
		tok token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "identifier",
			fields: fields{
				input: "abc",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "abc",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "digits and underscores",
			fields: fields{
				input: "a1_b2",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "a1_b2",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "single underscore",
			fields: fields{
				input: "_",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "_",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "lower true",
			fields: fields{
				input: "true",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenBool,
					v:    "true",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "title false",
			fields: fields{
				input: "False",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenBool,
					v:    "False",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "upper true",
			fields: fields{
				input: "TRUE",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenBool,
					v:    "TRUE",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "mixed case is an identifier",
			fields: fields{
				input: "tRUE",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "tRUE",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "keyword prefix is an identifier",
			fields: fields{
				input: "trueish",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "trueish",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "stops before a space",
			fields: fields{
				input: "abc def",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "abc",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "stops before an operator",
			fields: fields{
				input: "abc==",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "abc",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "keyword before an operator is a bool",
			fields: fields{
				input: "true)",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenBool,
					v:    "true",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "hyphen ends the identifier",
			fields: fields{
				input: "a-b",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "a",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "non-ascii letters",
			fields: fields{
				input: "軍師",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "軍師",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "non-ascii letters then digit",
			fields: fields{
				input: "名前1",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "名前1",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "ascii then non-ascii",
			fields: fields{
				input: "a軍",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "a軍",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "non-ascii stops before an operator",
			fields: fields{
				input: "軍師==",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "軍師",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "non-ascii stops before a wide space",
			fields: fields{
				input: "軍　師",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "軍",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "fullwidth letters",
			fields: fields{
				input: "Ａ1",
			},
			want: want{
				val: stateStmt,
				tok: token{
					typ:  tokenIdent,
					v:    "Ａ1",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			l.next()
			got := l.lexKeywordOrIdent()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.token != test.want.tok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.token, test.want.tok)
			}
		})
	}
}

func Test_lexer_scanEscape(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "newline",
			fields: fields{
				input: "n",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "tab",
			fields: fields{
				input: "t",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "backslash",
			fields: fields{
				input: "\\",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "quote_double",
			fields: fields{
				input: "\"",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "quote_single",
			fields: fields{
				input: "'",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "null",
			fields: fields{
				input: "0",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "bell",
			fields: fields{
				input: "a",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "backspace",
			fields: fields{
				input: "b",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "formfeed",
			fields: fields{
				input: "f",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "carriage_return",
			fields: fields{
				input: "r",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "vertical_tab",
			fields: fields{
				input: "v",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "hex",
			fields: fields{
				input: "x41",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "unicode",
			fields: fields{
				input: "u0041",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "invalid_char",
			fields: fields{
				input: "z",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "empty",
			fields: fields{
				input: "",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "eof",
			fields: fields{
				input: string([]byte{0}),
			},
			want: want{
				val: false,
			},
		},
		{
			name: "backtick",
			fields: fields{
				input: "`",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "hex_short",
			fields: fields{
				input: "x4",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "hex_nonhex",
			fields: fields{
				input: "x4G",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "unicode_short",
			fields: fields{
				input: "u041",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "unicode_nonhex",
			fields: fields{
				input: "u004G",
			},
			want: want{
				val: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.fields.input,
			}
			got := l.scanEscape()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_lexer_scanHexEscape(t *testing.T) {
	type fields struct {
		input string
	}
	type args struct {
		digits int
	}
	type want struct {
		val bool
		pos int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "two digits",
			fields: fields{
				input: "41",
			},
			args: args{
				digits: 2,
			},
			want: want{
				val: true,
				pos: 2,
			},
		},
		{
			name: "two digits stop before the rest",
			fields: fields{
				input: "41zz",
			},
			args: args{
				digits: 2,
			},
			want: want{
				val: true,
				pos: 2,
			},
		},
		{
			name: "four digits",
			fields: fields{
				input: "0041",
			},
			args: args{
				digits: 4,
			},
			want: want{
				val: true,
				pos: 4,
			},
		},
		{
			name: "lower case hex",
			fields: fields{
				input: "abcd",
			},
			args: args{
				digits: 4,
			},
			want: want{
				val: true,
				pos: 4,
			},
		},
		{
			name: "upper case hex",
			fields: fields{
				input: "ABCD",
			},
			args: args{
				digits: 4,
			},
			want: want{
				val: true,
				pos: 4,
			},
		},
		{
			name: "zero digits",
			fields: fields{
				input: "",
			},
			args: args{
				digits: 0,
			},
			want: want{
				val: true,
				pos: 0,
			},
		},
		{
			name: "one digit short",
			fields: fields{
				input: "4",
			},
			args: args{
				digits: 2,
			},
			want: want{
				val: false,
				pos: 1,
			},
		},
		{
			name: "three digits short",
			fields: fields{
				input: "004",
			},
			args: args{
				digits: 4,
			},
			want: want{
				val: false,
				pos: 3,
			},
		},
		{
			name: "empty input",
			fields: fields{
				input: "",
			},
			args: args{
				digits: 2,
			},
			want: want{
				val: false,
				pos: 0,
			},
		},
		{
			name: "non-hex letter",
			fields: fields{
				input: "4G",
			},
			args: args{
				digits: 2,
			},
			want: want{
				val: false,
				pos: 2,
			},
		},
		{
			name: "letter past f",
			fields: fields{
				input: "g0",
			},
			args: args{
				digits: 2,
			},
			want: want{
				val: false,
				pos: 1,
			},
		},
		{
			name: "non-ascii digit",
			fields: fields{
				input: "４1",
			},
			args: args{
				digits: 2,
			},
			want: want{
				val: true,
				pos: 4,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.fields.input,
			}
			got := l.scanHexEscape(test.args.digits)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.pos != test.want.pos {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.pos, test.want.pos)
			}
		})
	}
}

func Test_lexer_scanTime(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		ok      bool
		matched string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "utc",
			fields: fields{
				input: "2023-01-02T15:04:05Z",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02T15:04:05Z",
			},
		},
		{
			name: "fraction",
			fields: fields{
				input: "2023-01-02T15:04:05.123Z",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02T15:04:05.123Z",
			},
		},
		{
			name: "offset plus",
			fields: fields{
				input: "2023-01-02T15:04:05+09:00",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02T15:04:05+09:00",
			},
		},
		{
			name: "offset minus",
			fields: fields{
				input: "2023-01-02T15:04:05-05:00",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02T15:04:05-05:00",
			},
		},
		{
			name: "no timezone",
			fields: fields{
				input: "2023-01-02T15:04:05",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02T15:04:05",
			},
		},
		{
			name: "fraction no timezone",
			fields: fields{
				input: "2023-01-02T15:04:05.5",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02T15:04:05.5",
			},
		},
		{
			name: "date only",
			fields: fields{
				input: "2023-01-02",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02",
			},
		},
		{
			name: "date followed by space",
			fields: fields{
				input: "2023-01-02 15:04:05",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02",
			},
		},
		{
			name: "date followed by T alone",
			fields: fields{
				input: "2023-01-02T",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02",
			},
		},
		{
			name: "date followed by T and letters",
			fields: fields{
				input: "2023-01-02Tx",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02",
			},
		},
		{
			name: "missing seconds stops at the date",
			fields: fields{
				input: "2023-01-02T15:04",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02",
			},
		},
		{
			name: "bad clock stops at the date",
			fields: fields{
				input: "2023-01-02T15-04-05Z",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02",
			},
		},
		{
			name: "lowercase z stops before it",
			fields: fields{
				input: "2023-01-02T15:04:05z",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02T15:04:05",
			},
		},
		{
			name: "empty fraction stops before the dot",
			fields: fields{
				input: "2023-01-02T15:04:05.Z",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02T15:04:05",
			},
		},
		{
			name: "bad offset stops before the sign",
			fields: fields{
				input: "2023-01-02T15:04:05+0900",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02T15:04:05",
			},
		},
		{
			name: "short offset stops before the sign",
			fields: fields{
				input: "2023-01-02T15:04:05+09",
			},
			want: want{
				ok:      true,
				matched: "2023-01-02T15:04:05",
			},
		},
		{
			name: "bad date",
			fields: fields{
				input: "2023/01/02T15:04:05Z",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "empty",
			fields: fields{
				input: "",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.fields.input,
			}
			got := l.scanTime()
			if got != test.want.ok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.ok)
			}
			if got {
				if matched := test.fields.input[l.startPos:l.pos]; matched != test.want.matched {
					t.Errorf("value mismatch\ngot=%v\nwant=%v\n", matched, test.want.matched)
				}
			}
		})
	}
}

func Test_lexer_scanDuration(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		ok      bool
		matched string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "hour",
			fields: fields{
				input: "1h",
			},
			want: want{
				ok:      true,
				matched: "1h",
			},
		},
		{
			name: "minute",
			fields: fields{
				input: "1m",
			},
			want: want{
				ok:      true,
				matched: "1m",
			},
		},
		{
			name: "second",
			fields: fields{
				input: "1s",
			},
			want: want{
				ok:      true,
				matched: "1s",
			},
		},
		{
			name: "millisecond",
			fields: fields{
				input: "1ms",
			},
			want: want{
				ok:      true,
				matched: "1ms",
			},
		},
		{
			name: "microsecond 1",
			fields: fields{
				input: "1us",
			},
			want: want{
				ok:      true,
				matched: "1us",
			},
		},
		{
			name: "microsecond 2",
			fields: fields{
				input: "1μs",
			},
			want: want{
				ok:      true,
				matched: "1μs",
			},
		},
		{
			name: "microsecond 3",
			fields: fields{
				input: "1µs",
			},
			want: want{
				ok:      true,
				matched: "1µs",
			},
		},
		{
			name: "nanosecond",
			fields: fields{
				input: "1ns",
			},
			want: want{
				ok:      true,
				matched: "1ns",
			},
		},
		{
			name: "sign 1",
			fields: fields{
				input: "+1h",
			},
			want: want{
				ok:      true,
				matched: "+1h",
			},
		},
		{
			name: "sign 2",
			fields: fields{
				input: "-1h",
			},
			want: want{
				ok:      true,
				matched: "-1h",
			},
		},
		{
			name: "float 1",
			fields: fields{
				input: "0.1h",
			},
			want: want{
				ok:      true,
				matched: "0.1h",
			},
		},
		{
			name: "float 2",
			fields: fields{
				input: "1.1h",
			},
			want: want{
				ok:      true,
				matched: "1.1h",
			},
		},
		{
			name: "float 3",
			fields: fields{
				input: ".1h",
			},
			want: want{
				ok:      true,
				matched: ".1h",
			},
		},
		{
			name: "float 4",
			fields: fields{
				input: "1.h",
			},
			want: want{
				ok:      true,
				matched: "1.h",
			},
		},
		{
			name: "mixed 1",
			fields: fields{
				input: "1h5000ns",
			},
			want: want{
				ok:      true,
				matched: "1h5000ns",
			},
		},
		{
			name: "mixed 2",
			fields: fields{
				input: "5000ns1h",
			},
			want: want{
				ok:      true,
				matched: "5000ns1h",
			},
		},
		{
			name: "mixed 3",
			fields: fields{
				input: "+1h5000ns",
			},
			want: want{
				ok:      true,
				matched: "+1h5000ns",
			},
		},
		{
			name: "mixed 4",
			fields: fields{
				input: "-5000ns1h",
			},
			want: want{
				ok:      true,
				matched: "-5000ns1h",
			},
		},
		{
			name: "mixed 5",
			fields: fields{
				input: "0.1h0.30m",
			},
			want: want{
				ok:      true,
				matched: "0.1h0.30m",
			},
		},
		{
			name: "mixed 6",
			fields: fields{
				input: ".1m.30s",
			},
			want: want{
				ok:      true,
				matched: ".1m.30s",
			},
		},
		{
			name: "mixed 7",
			fields: fields{
				input: "-1.1h.30m",
			},
			want: want{
				ok:      true,
				matched: "-1.1h.30m",
			},
		},
		{
			name: "mixed 8",
			fields: fields{
				input: "+0.1h.30m",
			},
			want: want{
				ok:      true,
				matched: "+0.1h.30m",
			},
		},
		{
			name: "mixed 9",
			fields: fields{
				input: "-.1h.30m",
			},
			want: want{
				ok:      true,
				matched: "-.1h.30m",
			},
		},
		{
			name: "mixed 10",
			fields: fields{
				input: "+.1h.30m",
			},
			want: want{
				ok:      true,
				matched: "+.1h.30m",
			},
		},
		{
			name: "mixed 11",
			fields: fields{
				input: "+1.h30.m",
			},
			want: want{
				ok:      true,
				matched: "+1.h30.m",
			},
		},
		{
			name: "full",
			fields: fields{
				input: "1h30m15s3000ms4000us5000ns",
			},
			want: want{
				ok:      true,
				matched: "1h30m15s3000ms4000us5000ns",
			},
		},
		{
			name: "duplicated",
			fields: fields{
				input: "1h1h",
			},
			want: want{
				ok:      true,
				matched: "1h1h",
			},
		},
		{
			name: "longest match 1",
			fields: fields{
				input: "1h+30m",
			},
			want: want{
				ok:      true,
				matched: "1h",
			},
		},
		{
			name: "longest match 2",
			fields: fields{
				input: "1h-30m",
			},
			want: want{
				ok:      true,
				matched: "1h",
			},
		},
		{
			name: "longest match 3",
			fields: fields{
				input: "+1h+30m+15s+3000ms+4000us+5000ns",
			},
			want: want{
				ok:      true,
				matched: "+1h",
			},
		},
		{
			name: "longest match 4",
			fields: fields{
				input: "-1h-30m-15s-3000ms-4000us-5000ns",
			},
			want: want{
				ok:      true,
				matched: "-1h",
			},
		},
		{
			name: "longest match 5",
			fields: fields{
				input: "1hm",
			},
			want: want{
				ok:      true,
				matched: "1h",
			},
		},
		{
			name: "longest match 6",
			fields: fields{
				input: "1hms",
			},
			want: want{
				ok:      true,
				matched: "1h",
			},
		},
		{
			name: "longest match 7",
			fields: fields{
				input: "1hd",
			},
			want: want{
				ok:      true,
				matched: "1h",
			},
		},
		{
			name: "longest match 8",
			fields: fields{
				input: "1h30m1d",
			},
			want: want{
				ok:      true,
				matched: "1h30m",
			},
		},
		{
			name: "longest match 9",
			fields: fields{
				input: "1h30md",
			},
			want: want{
				ok:      true,
				matched: "1h30m",
			},
		},
		{
			name: "longest match 10",
			fields: fields{
				input: "1h_",
			},
			want: want{
				ok:      true,
				matched: "1h",
			},
		},
		{
			name: "invalid multiple dot but passed 1",
			fields: fields{
				input: "0..1h",
			},
			want: want{
				ok:      true,
				matched: "0..1h",
			},
		},
		{
			name: "invalid multiple dot but passed 2",
			fields: fields{
				input: "..1h",
			},
			want: want{
				ok:      true,
				matched: "..1h",
			},
		},
		{
			name: "number 1",
			fields: fields{
				input: "1",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "number 2",
			fields: fields{
				input: "+1",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "number 3",
			fields: fields{
				input: "-1",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "invalid unit 1",
			fields: fields{
				input: "365d",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "invalid unit 4",
			fields: fields{
				input: "1d30m",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "only unit 1",
			fields: fields{
				input: "h",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "only unit 2",
			fields: fields{
				input: "ms",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "only sign 1",
			fields: fields{
				input: "+",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "only sign 2",
			fields: fields{
				input: "-",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "sign and unit 1",
			fields: fields{
				input: "+ms",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "sign and unit 2",
			fields: fields{
				input: "-ms",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "unary operator 1",
			fields: fields{
				input: "- 1ms",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "unary operator 2",
			fields: fields{
				input: "+ ms",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "unit repeat 1",
			fields: fields{
				input: "msms",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
		{
			name: "empty",
			fields: fields{
				input: "",
			},
			want: want{
				ok:      false,
				matched: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.fields.input,
			}
			got := l.scanDuration()
			if got != test.want.ok {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.ok)
			}
			if matched := test.fields.input[l.startPos:l.pos]; matched != test.want.matched {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", matched, test.want.matched)
			}
		})
	}
}

func Test_lexer_scanDurationNumber(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val bool
		pos int32
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "integer",
			fields: fields{
				input: "1",
			},
			want: want{
				val: true,
				pos: 1,
			},
		},
		{
			name: "fraction",
			fields: fields{
				input: "1.5",
			},
			want: want{
				val: true,
				pos: 3,
			},
		},
		{
			name: "leading dot",
			fields: fields{
				input: ".5",
			},
			want: want{
				val: true,
				pos: 2,
			},
		},
		{
			name: "trailing dot",
			fields: fields{
				input: "1.",
			},
			want: want{
				val: true,
				pos: 2,
			},
		},
		{
			name: "positive sign",
			fields: fields{
				input: "+1",
			},
			want: want{
				val: true,
				pos: 2,
			},
		},
		{
			name: "negative sign",
			fields: fields{
				input: "-1",
			},
			want: want{
				val: true,
				pos: 2,
			},
		},
		{
			name: "stops before the unit",
			fields: fields{
				input: "10s",
			},
			want: want{
				val: true,
				pos: 2,
			},
		},
		{
			name: "dots only",
			fields: fields{
				input: "..",
			},
			want: want{
				val: true,
				pos: 2,
			},
		},
		{
			name: "empty input",
			fields: fields{
				input: "",
			},
			want: want{
				val: false,
				pos: 0,
			},
		},
		{
			name: "sign without digits resets",
			fields: fields{
				input: "+",
			},
			want: want{
				val: false,
				pos: 0,
			},
		},
		{
			name: "sign before a unit resets",
			fields: fields{
				input: "-s",
			},
			want: want{
				val: false,
				pos: 0,
			},
		},
		{
			name: "unit without a number",
			fields: fields{
				input: "s",
			},
			want: want{
				val: false,
				pos: 0,
			},
		},
		{
			name: "double sign resets",
			fields: fields{
				input: "+-1",
			},
			want: want{
				val: false,
				pos: 0,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.fields.input,
			}
			got := l.scanDurationNumber()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
			if l.pos != test.want.pos {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", l.pos, test.want.pos)
			}
		})
	}
}

func Test_lexer_scanNumber(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "integer",
			fields: fields{
				input: "1",
			},
			want: want{
				val: "1",
			},
		},
		{
			name: "zero",
			fields: fields{
				input: "0",
			},
			want: want{
				val: "0",
			},
		},
		{
			name: "float",
			fields: fields{
				input: "1.5",
			},
			want: want{
				val: "1.5",
			},
		},
		{
			name: "leading dot",
			fields: fields{
				input: ".5",
			},
			want: want{
				val: ".5",
			},
		},
		{
			name: "trailing dot",
			fields: fields{
				input: "1.",
			},
			want: want{
				val: "1.",
			},
		},
		{
			name: "sign plus",
			fields: fields{
				input: "+1",
			},
			want: want{
				val: "+1",
			},
		},
		{
			name: "sign minus",
			fields: fields{
				input: "-1",
			},
			want: want{
				val: "-1",
			},
		},
		{
			name: "underscore",
			fields: fields{
				input: "1_000",
			},
			want: want{
				val: "1_000",
			},
		},
		{
			name: "hex",
			fields: fields{
				input: "0x1f",
			},
			want: want{
				val: "0x1f",
			},
		},
		{
			name: "hex float",
			fields: fields{
				input: "0x.8p1",
			},
			want: want{
				val: "0x.8p1",
			},
		},
		{
			name: "octal",
			fields: fields{
				input: "0o17",
			},
			want: want{
				val: "0o17",
			},
		},
		{
			name: "binary",
			fields: fields{
				input: "0b101",
			},
			want: want{
				val: "0b101",
			},
		},
		{
			name: "exponent",
			fields: fields{
				input: "1e3",
			},
			want: want{
				val: "1e3",
			},
		},
		{
			name: "signed exponent",
			fields: fields{
				input: "1e-3",
			},
			want: want{
				val: "1e-3",
			},
		},
		{
			name: "dot only",
			fields: fields{
				input: ".",
			},
			want: want{
				val: ".",
			},
		},
		{
			name: "sign only",
			fields: fields{
				input: "-",
			},
			want: want{
				val: "-",
			},
		},
		{
			name: "base prefix only",
			fields: fields{
				input: "0x",
			},
			want: want{
				val: "0x",
			},
		},
		{
			name: "stops at second dot",
			fields: fields{
				input: "1.2.3",
			},
			want: want{
				val: "1.2",
			},
		},
		{
			name: "stops at operator",
			fields: fields{
				input: "1-2",
			},
			want: want{
				val: "1",
			},
		},
		{
			name: "stops at unit",
			fields: fields{
				input: "1h",
			},
			want: want{
				val: "1",
			},
		},
		{
			name: "empty",
			fields: fields{
				input: "",
			},
			want: want{
				val: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.fields.input,
			}
			l.scanNumber()
			if got := test.fields.input[l.startPos:l.pos]; got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_lexer_nextToken(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		val []token
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "simple number 1",
			fields: fields{
				input: "1",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "1",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  1,
						line: 1,
						col:  2,
					},
				},
			},
		},
		{
			name: "simple number 2",
			fields: fields{
				input: "+1",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "+1",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  2,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "simple number 3",
			fields: fields{
				input: "-1",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "-1",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  2,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "simple number 4",
			fields: fields{
				input: ".1",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    ".1",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  2,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "simple number 5",
			fields: fields{
				input: "0.1",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "0.1",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  3,
						line: 1,
						col:  4,
					},
				},
			},
		},
		{
			name: "simple number 6",
			fields: fields{
				input: "0x1.fp3",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "0x1.fp3",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  7,
						line: 1,
						col:  8,
					},
				},
			},
		},
		{
			name: "simple duration",
			fields: fields{
				input: "1h",
			},
			want: want{
				val: []token{
					{
						typ:  tokenDuration,
						v:    "1h",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  2,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "ident",
			fields: fields{
				input: "id IDENT_1 あいうえお",
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "id",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "IDENT_1",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenIdent,
						v:    "あいうえお",
						pos:  11,
						line: 1,
						col:  12,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  26,
						line: 1,
						col:  22,
					},
				},
			},
		},
		{
			name: "comparison operators 1",
			fields: fields{
				input: "> >= < <=",
			},
			want: want{
				val: []token{
					{
						typ:  tokenGT,
						v:    ">",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenGTE,
						v:    ">=",
						pos:  2,
						line: 1,
						col:  3,
					},
					{
						typ:  tokenLT,
						v:    "<",
						pos:  5,
						line: 1,
						col:  6,
					},
					{
						typ:  tokenLTE,
						v:    "<=",
						pos:  7,
						line: 1,
						col:  8,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  9,
						line: 1,
						col:  10,
					},
				},
			},
		},
		{
			name: "comparison operators 2",
			fields: fields{
				input: "== !=",
			},
			want: want{
				val: []token{
					{
						typ:  tokenEQ,
						v:    "==",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenNEQ,
						v:    "!=",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenEOF,
						pos:  5,
						line: 1,
						col:  6,
					},
				},
			},
		},
		{
			name: "comparison operators 3",
			fields: fields{
				input: "=~ !~",
			},
			want: want{
				val: []token{
					{
						typ:  tokenREQ,
						v:    "=~",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenNREQ,
						v:    "!~",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  5,
						line: 1,
						col:  6,
					},
				},
			},
		},
		{
			name: "logical operators",
			fields: fields{
				input: "&& || !",
			},
			want: want{
				val: []token{
					{
						typ:  tokenAND,
						v:    "&&",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenOR,
						v:    "||",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenNOT,
						v:    "!",
						pos:  6,
						line: 1,
						col:  7,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  7,
						line: 1,
						col:  8,
					},
				},
			},
		},
		{
			name: "parentheses",
			fields: fields{
				input: "()",
			},
			want: want{
				val: []token{
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  1,
						line: 1,
						col:  2,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  2,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "string",
			fields: fields{
				input: "\"abc\"",
			},
			want: want{
				val: []token{
					{
						typ:  tokenString,
						v:    "\"abc\"",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  5,
						line: 1,
						col:  6,
					},
				},
			},
		},
		{
			name: "string with escape",
			fields: fields{
				input: "\"\\n\\t\\\\\\\"\\'\\0\\a\\b\\f\\r\\v\\x41\\u0041\"",
			},
			want: want{
				val: []token{
					{
						typ:  tokenString,
						v:    "\"\\n\\t\\\\\\\"\\'\\0\\a\\b\\f\\r\\v\\x41\\u0041\"",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  34,
						line: 1,
						col:  35,
					},
				},
			},
		},
		{
			name: "string with wrong hex",
			fields: fields{
				input: "'\\xG'",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "invalid escape sequence in string",
						pos:  4,
						line: 1,
						col:  5,
					},
				},
			},
		},
		{
			name: "string with eof",
			fields: fields{
				input: "\"",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "unterminated quoted string",
						pos:  1,
						line: 1,
						col:  2,
					},
				},
			},
		},
		{
			name: "string with line break",
			fields: fields{
				input: "\"abc\\ndef\"",
			},
			want: want{
				val: []token{
					{
						typ:  tokenString,
						v:    "\"abc\\ndef\"",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  10,
						line: 1,
						col:  11,
					},
				},
			},
		},
		{
			name: "raw string",
			fields: fields{
				input: "`abc`",
			},
			want: want{
				val: []token{
					{
						typ:  tokenRawString,
						v:    "`abc`",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  5,
						line: 1,
						col:  6,
					},
				},
			},
		},
		{
			name: "number",
			fields: fields{
				input: "0 1 +2 -3 0.4 .5 +0.6 -0.7 +.8 -.9 1.23e4 1.23E4 1.23e+4 1.23e-4 0x1A2b 0x1.fp3 0x1.fp+3 0x1.fp-3 0o755 0b1011",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "0",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenNumber,
						v:    "1",
						pos:  2,
						line: 1,
						col:  3,
					},
					{
						typ:  tokenNumber,
						v:    "+2",
						pos:  4,
						line: 1,
						col:  5,
					},
					{
						typ:  tokenNumber,
						v:    "-3",
						pos:  7,
						line: 1,
						col:  8,
					},
					{
						typ:  tokenNumber,
						v:    "0.4",
						pos:  10,
						line: 1,
						col:  11,
					},
					{
						typ:  tokenNumber,
						v:    ".5",
						pos:  14,
						line: 1,
						col:  15,
					},
					{
						typ:  tokenNumber,
						v:    "+0.6",
						pos:  17,
						line: 1,
						col:  18,
					},
					{
						typ:  tokenNumber,
						v:    "-0.7",
						pos:  22,
						line: 1,
						col:  23,
					},
					{
						typ:  tokenNumber,
						v:    "+.8",
						pos:  27,
						line: 1,
						col:  28,
					},
					{
						typ:  tokenNumber,
						v:    "-.9",
						pos:  31,
						line: 1,
						col:  32,
					},
					{
						typ:  tokenNumber,
						v:    "1.23e4",
						pos:  35,
						line: 1,
						col:  36,
					},
					{
						typ:  tokenNumber,
						v:    "1.23E4",
						pos:  42,
						line: 1,
						col:  43,
					},
					{
						typ:  tokenNumber,
						v:    "1.23e+4",
						pos:  49,
						line: 1,
						col:  50,
					},
					{
						typ:  tokenNumber,
						v:    "1.23e-4",
						pos:  57,
						line: 1,
						col:  58,
					},
					{
						typ:  tokenNumber,
						v:    "0x1A2b",
						pos:  65,
						line: 1,
						col:  66,
					},
					{
						typ:  tokenNumber,
						v:    "0x1.fp3",
						pos:  72,
						line: 1,
						col:  73,
					},
					{
						typ:  tokenNumber,
						v:    "0x1.fp+3",
						pos:  80,
						line: 1,
						col:  81,
					},
					{
						typ:  tokenNumber,
						v:    "0x1.fp-3",
						pos:  89,
						line: 1,
						col:  90,
					},
					{
						typ:  tokenNumber,
						v:    "0o755",
						pos:  98,
						line: 1,
						col:  99,
					},
					{
						typ:  tokenNumber,
						v:    "0b1011",
						pos:  104,
						line: 1,
						col:  105,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  110,
						line: 1,
						col:  111,
					},
				},
			},
		},
		{
			name: "duration",
			fields: fields{
				input: "1h30m+100s+1h+30m+15s-3000ms-4000us-5000ns 0.1h.5m 1y2m3w4d",
			},
			want: want{
				val: []token{
					{
						typ:  tokenDuration,
						v:    "1h30m",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenDuration,
						v:    "+100s",
						pos:  5,
						line: 1,
						col:  6,
					},
					{
						typ:  tokenDuration,
						v:    "+1h",
						pos:  10,
						line: 1,
						col:  11,
					},
					{
						typ:  tokenDuration,
						v:    "+30m",
						pos:  13,
						line: 1,
						col:  14,
					},
					{
						typ:  tokenDuration,
						v:    "+15s",
						pos:  17,
						line: 1,
						col:  18,
					},
					{
						typ:  tokenDuration,
						v:    "-3000ms",
						pos:  21,
						line: 1,
						col:  22,
					},
					{
						typ:  tokenDuration,
						v:    "-4000us",
						pos:  28,
						line: 1,
						col:  29,
					},
					{
						typ:  tokenDuration,
						v:    "-5000ns",
						pos:  35,
						line: 1,
						col:  36,
					},
					{
						typ:  tokenDuration,
						v:    "0.1h.5m",
						pos:  43,
						line: 1,
						col:  44,
					},
					{
						typ:  tokenNumber,
						v:    "1",
						pos:  51,
						line: 1,
						col:  52,
					},
					{
						typ:  tokenIdent,
						v:    "y2m3w4d",
						pos:  52,
						line: 1,
						col:  53,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  59,
						line: 1,
						col:  60,
					},
				},
			},
		},
		{
			name: "duration/number/ident",
			fields: fields{
				input: "1h1x",
			},
			want: want{
				val: []token{
					{
						typ:  tokenDuration,
						v:    "1h",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenNumber,
						v:    "1",
						pos:  2,
						line: 1,
						col:  3,
					},
					{
						typ:  tokenIdent,
						v:    "x",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  4,
						line: 1,
						col:  5,
					},
				},
			},
		},
		{
			name: "bool",
			fields: fields{
				input: "true True TRUE false False FALSE tRue",
			},
			want: want{
				val: []token{
					{
						typ:  tokenBool,
						v:    "true",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenBool,
						v:    "True",
						pos:  5,
						line: 1,
						col:  6,
					},
					{
						typ:  tokenBool,
						v:    "TRUE",
						pos:  10,
						line: 1,
						col:  11,
					},
					{
						typ:  tokenBool,
						v:    "false",
						pos:  15,
						line: 1,
						col:  16,
					},
					{
						typ:  tokenBool,
						v:    "False",
						pos:  21,
						line: 1,
						col:  22,
					},
					{
						typ:  tokenBool,
						v:    "FALSE",
						pos:  27,
						line: 1,
						col:  28,
					},
					{
						typ:  tokenIdent,
						v:    "tRue",
						pos:  33,
						line: 1,
						col:  34,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  37,
						line: 1,
						col:  38,
					},
				},
			},
		},
		{
			name: "two digit number at end of input",
			fields: fields{
				input: "40",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "40",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						pos:  2,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "number with two dots splits at the second dot",
			fields: fields{
				input: "1.2.3",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "1.2",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenNumber,
						v:    ".3",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenEOF,
						pos:  5,
						line: 1,
						col:  6,
					},
				},
			},
		},
		{
			name: "wide identifier then bare dot",
			fields: fields{
				input: "名前 == .",
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "名前",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEQ,
						v:    "==",
						pos:  7,
						line: 1,
						col:  6,
					},
					{
						typ:  tokenNumber,
						v:    ".",
						pos:  10,
						line: 1,
						col:  9,
					},
					{
						typ:  tokenEOF,
						pos:  11,
						line: 1,
						col:  10,
					},
				},
			},
		},
		{
			name: "bare sign on second line",
			fields: fields{
				input: "a\n== -",
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "a",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEQ,
						v:    "==",
						pos:  2,
						line: 2,
						col:  1,
					},
					{
						typ:  tokenNumber,
						v:    "-",
						pos:  5,
						line: 2,
						col:  4,
					},
					{
						typ:  tokenEOF,
						pos:  6,
						line: 2,
						col:  5,
					},
				},
			},
		},
		{
			name: "bare dot",
			fields: fields{
				input: ".",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    ".",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						pos:  1,
						line: 1,
						col:  2,
					},
				},
			},
		},
		{
			name: "bare plus",
			fields: fields{
				input: "+",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "+",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						pos:  1,
						line: 1,
						col:  2,
					},
				},
			},
		},
		{
			name: "bare minus after value",
			fields: fields{
				input: "1 -",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "1",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenNumber,
						v:    "-",
						pos:  2,
						line: 1,
						col:  3,
					},
					{
						typ:  tokenEOF,
						pos:  3,
						line: 1,
						col:  4,
					},
				},
			},
		},
		{
			name: "invalid character 1",
			fields: fields{
				input: "\\",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "unexpected character U+005C '\\'",
						pos:  0,
						line: 1,
						col:  1,
					},
				},
			},
		},
		{
			name: "invalid paren depth 1",
			fields: fields{
				input: "((",
			},
			want: want{
				val: []token{
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  1,
						line: 1,
						col:  2,
					},
					{
						typ:  tokenError,
						v:    "unclosed left parenthesis",
						pos:  2,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "invalid paren depth 2",
			fields: fields{
				input: "))",
			},
			want: want{
				val: []token{
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  1,
						line: 1,
						col:  2,
					},
					{
						typ:  tokenError,
						v:    "unexpected right parenthesis",
						pos:  2,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "invalid paren depth 3",
			fields: fields{
				input: "((())",
			},
			want: want{
				val: []token{
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  1,
						line: 1,
						col:  2,
					},
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  2,
						line: 1,
						col:  3,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  4,
						line: 1,
						col:  5,
					},
					{
						typ:  tokenError,
						v:    "unclosed left parenthesis",
						pos:  5,
						line: 1,
						col:  6,
					},
				},
			},
		},
		{
			name: "invalid paren depth 4",
			fields: fields{
				input: "(()))",
			},
			want: want{
				val: []token{
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  1,
						line: 1,
						col:  2,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  2,
						line: 1,
						col:  3,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  4,
						line: 1,
						col:  5,
					},
					{
						typ:  tokenError,
						v:    "unexpected right parenthesis",
						pos:  5,
						line: 1,
						col:  6,
					},
				},
			},
		},
		{
			name: "rune error in string",
			fields: fields{
				input: "\"\uFFFD\"",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "invalid utf8 encoding in string",
						pos:  4,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "unterminated string 1",
			fields: fields{
				input: "\"aaa bbb ccc",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "unterminated quoted string",
						pos:  12,
						line: 1,
						col:  13,
					},
				},
			},
		},
		{
			name: "unterminated string 2",
			fields: fields{
				input: "'aaa bbb ccc",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "unterminated quoted string",
						pos:  12,
						line: 1,
						col:  13,
					},
				},
			},
		},
		{
			name: "invalid escape sequence in string",
			fields: fields{
				input: "\"aaa\\zbbb\"",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "invalid escape sequence in string",
						pos:  6,
						line: 1,
						col:  7,
					},
				},
			},
		},
		{
			name: "rune error in raw string",
			fields: fields{
				input: "`\uFFFD`",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "invalid utf8 encoding in raw string",
						pos:  4,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "unterminated raw string",
			fields: fields{
				input: "`aaa bbb ccc",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "unterminated raw string",
						pos:  12,
						line: 1,
						col:  13,
					},
				},
			},
		},
		{
			name: "unexpected operator 1",
			fields: fields{
				input: "=!",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "unexpected character '!' after '='",
						pos:  1,
						line: 1,
						col:  2,
					},
				},
			},
		},
		{
			name: "unexpected operator 2",
			fields: fields{
				input: "&|",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "unexpected character '|' after '&'",
						pos:  1,
						line: 1,
						col:  2,
					},
				},
			},
		},
		{
			name: "unexpected operator 3",
			fields: fields{
				input: "|&",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "unexpected character '&' after '|'",
						pos:  1,
						line: 1,
						col:  2,
					},
				},
			},
		},
		{
			name: "bad number syntax 1",
			fields: fields{
				input: "10abc",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "10",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "abc",
						pos:  2,
						line: 1,
						col:  3,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  5,
						line: 1,
						col:  6,
					},
				},
			},
		},
		{
			name: "bad number syntax 2",
			fields: fields{
				input: "_",
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "_",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  1,
						line: 1,
						col:  2,
					},
				},
			},
		},
		{
			name: "multibyte",
			fields: fields{
				input: "一二三四五六七八九十",
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "一二三四五六七八九十",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  30,
						line: 1,
						col:  21,
					},
				},
			},
		},
		{
			name: "mixed 1",
			fields: fields{
				input: `Class=="軍師"&&Name=~'孔明'&&(HP>50&&MP>=100&&LP!=0)&&(MAG>=20||!(SPD<20))`,
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "Class",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEQ,
						v:    "==",
						pos:  5,
						line: 1,
						col:  6,
					},
					{
						typ:  tokenString,
						v:    "\"軍師\"",
						pos:  7,
						line: 1,
						col:  8,
					},
					{
						typ:  tokenAND,
						v:    "&&",
						pos:  15,
						line: 1,
						col:  14,
					},
					{
						typ:  tokenIdent,
						v:    "Name",
						pos:  17,
						line: 1,
						col:  16,
					},
					{
						typ:  tokenREQ,
						v:    "=~",
						pos:  21,
						line: 1,
						col:  20,
					},
					{
						typ:  tokenString,
						v:    "'孔明'",
						pos:  23,
						line: 1,
						col:  22,
					},
					{
						typ:  tokenAND,
						v:    "&&",
						pos:  31,
						line: 1,
						col:  28,
					},
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  33,
						line: 1,
						col:  30,
					},
					{
						typ:  tokenIdent,
						v:    "HP",
						pos:  34,
						line: 1,
						col:  31,
					},
					{
						typ:  tokenGT,
						v:    ">",
						pos:  36,
						line: 1,
						col:  33,
					},
					{
						typ:  tokenNumber,
						v:    "50",
						pos:  37,
						line: 1,
						col:  34,
					},
					{
						typ:  tokenAND,
						v:    "&&",
						pos:  39,
						line: 1,
						col:  36,
					},
					{
						typ:  tokenIdent,
						v:    "MP",
						pos:  41,
						line: 1,
						col:  38,
					},
					{
						typ:  tokenGTE,
						v:    ">=",
						pos:  43,
						line: 1,
						col:  40,
					},
					{
						typ:  tokenNumber,
						v:    "100",
						pos:  45,
						line: 1,
						col:  42,
					},
					{
						typ:  tokenAND,
						v:    "&&",
						pos:  48,
						line: 1,
						col:  45,
					},
					{
						typ:  tokenIdent,
						v:    "LP",
						pos:  50,
						line: 1,
						col:  47,
					},
					{
						typ:  tokenNEQ,
						v:    "!=",
						pos:  52,
						line: 1,
						col:  49,
					},
					{
						typ:  tokenNumber,
						v:    "0",
						pos:  54,
						line: 1,
						col:  51,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  55,
						line: 1,
						col:  52,
					},
					{
						typ:  tokenAND,
						v:    "&&",
						pos:  56,
						line: 1,
						col:  53,
					},
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  58,
						line: 1,
						col:  55,
					},
					{
						typ:  tokenIdent,
						v:    "MAG",
						pos:  59,
						line: 1,
						col:  56,
					},
					{
						typ:  tokenGTE,
						v:    ">=",
						pos:  62,
						line: 1,
						col:  59,
					},
					{
						typ:  tokenNumber,
						v:    "20",
						pos:  64,
						line: 1,
						col:  61,
					},
					{
						typ:  tokenOR,
						v:    "||",
						pos:  66,
						line: 1,
						col:  63,
					},
					{
						typ:  tokenNOT,
						v:    "!",
						pos:  68,
						line: 1,
						col:  65,
					},
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  69,
						line: 1,
						col:  66,
					},
					{
						typ:  tokenIdent,
						v:    "SPD",
						pos:  70,
						line: 1,
						col:  67,
					},
					{
						typ:  tokenLT,
						v:    "<",
						pos:  73,
						line: 1,
						col:  70,
					},
					{
						typ:  tokenNumber,
						v:    "20",
						pos:  74,
						line: 1,
						col:  71,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  76,
						line: 1,
						col:  73,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  77,
						line: 1,
						col:  74,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  78,
						line: 1,
						col:  75,
					},
				},
			},
		},
		{
			name: "mixed 2",
			fields: fields{
				input: `Class=="軍師"
&&
Name=~'孔明'
&&
(
	HP>50
	&&
	MP>=100
	&&
	LP!=0
)
&&
(
	MAG>=20
	||
	!
	(
		SPD<20
	)
)
`,
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "Class",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEQ,
						v:    "==",
						pos:  5,
						line: 1,
						col:  6,
					},
					{
						typ:  tokenString,
						v:    "\"軍師\"",
						pos:  7,
						line: 1,
						col:  8,
					},
					{
						typ:  tokenAND,
						v:    "&&",
						pos:  16,
						line: 2,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "Name",
						pos:  19,
						line: 3,
						col:  1,
					},
					{
						typ:  tokenREQ,
						v:    "=~",
						pos:  23,
						line: 3,
						col:  5,
					},
					{
						typ:  tokenString,
						v:    "'孔明'",
						pos:  25,
						line: 3,
						col:  7,
					},
					{
						typ:  tokenAND,
						v:    "&&",
						pos:  34,
						line: 4,
						col:  1,
					},
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  37,
						line: 5,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "HP",
						pos:  40,
						line: 6,
						col:  2,
					},
					{
						typ:  tokenGT,
						v:    ">",
						pos:  42,
						line: 6,
						col:  4,
					},
					{
						typ:  tokenNumber,
						v:    "50",
						pos:  43,
						line: 6,
						col:  5,
					},
					{
						typ:  tokenAND,
						v:    "&&",
						pos:  47,
						line: 7,
						col:  2,
					},
					{
						typ:  tokenIdent,
						v:    "MP",
						pos:  51,
						line: 8,
						col:  2,
					},
					{
						typ:  tokenGTE,
						v:    ">=",
						pos:  53,
						line: 8,
						col:  4,
					},
					{
						typ:  tokenNumber,
						v:    "100",
						pos:  55,
						line: 8,
						col:  6,
					},
					{
						typ:  tokenAND,
						v:    "&&",
						pos:  60,
						line: 9,
						col:  2,
					},
					{
						typ:  tokenIdent,
						v:    "LP",
						pos:  64,
						line: 10,
						col:  2,
					},
					{
						typ:  tokenNEQ,
						v:    "!=",
						pos:  66,
						line: 10,
						col:  4,
					},
					{
						typ:  tokenNumber,
						v:    "0",
						pos:  68,
						line: 10,
						col:  6,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  70,
						line: 11,
						col:  1,
					},
					{
						typ:  tokenAND,
						v:    "&&",
						pos:  72,
						line: 12,
						col:  1,
					},
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  75,
						line: 13,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "MAG",
						pos:  78,
						line: 14,
						col:  2,
					},
					{
						typ:  tokenGTE,
						v:    ">=",
						pos:  81,
						line: 14,
						col:  5,
					},
					{
						typ:  tokenNumber,
						v:    "20",
						pos:  83,
						line: 14,
						col:  7,
					},
					{
						typ:  tokenOR,
						v:    "||",
						pos:  87,
						line: 15,
						col:  2,
					},
					{
						typ:  tokenNOT,
						v:    "!",
						pos:  91,
						line: 16,
						col:  2,
					},
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  94,
						line: 17,
						col:  2,
					},
					{
						typ:  tokenIdent,
						v:    "SPD",
						pos:  98,
						line: 18,
						col:  3,
					},
					{
						typ:  tokenLT,
						v:    "<",
						pos:  101,
						line: 18,
						col:  6,
					},
					{
						typ:  tokenNumber,
						v:    "20",
						pos:  102,
						line: 18,
						col:  7,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  106,
						line: 19,
						col:  2,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  108,
						line: 20,
						col:  1,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  110,
						line: 21,
						col:  1,
					},
				},
			},
		},
		{
			name: "newline in input",
			fields: fields{
				input: `

test1
test2



		test3



`,
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "test1",
						pos:  2,
						line: 3,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "test2",
						pos:  8,
						line: 4,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "test3",
						pos:  19,
						line: 8,
						col:  3,
					},
					{
						typ:  tokenEOF,
						v:    "",
						pos:  28,
						line: 12,
						col:  1,
					},
				},
			},
		},
		{
			name: "asterisk after an operator is not an operator",
			fields: fields{
				input: "==*",
			},
			want: want{
				val: []token{
					{
						typ:  tokenEQ,
						v:    "==",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenError,
						v:    "unexpected character U+002A '*'",
						pos:  2,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "eq at end of input",
			fields: fields{
				input: "=",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "unexpected end of input after '='",
						pos:  1,
						line: 1,
						col:  2,
					},
				},
			},
		},
		{
			name: "and at end of input",
			fields: fields{
				input: "&",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "unexpected end of input after '&'",
						pos:  1,
						line: 1,
						col:  2,
					},
				},
			},
		},
		{
			name: "or at end of input",
			fields: fields{
				input: "|",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "unexpected end of input after '|'",
						pos:  1,
						line: 1,
						col:  2,
					},
				},
			},
		},
		{
			name: "eq followed by lone eq",
			fields: fields{
				input: "===",
			},
			want: want{
				val: []token{
					{
						typ:  tokenEQ,
						v:    "==",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenError,
						v:    "unexpected end of input after '='",
						pos:  3,
						line: 1,
						col:  4,
					},
				},
			},
		},
		{
			name: "neq followed by tilde",
			fields: fields{
				input: "!=~",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNEQ,
						v:    "!=",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenError,
						v:    "unexpected character U+007E '~'",
						pos:  2,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "lte followed by negative number",
			fields: fields{
				input: "a<=-1",
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "a",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenLTE,
						v:    "<=",
						pos:  1,
						line: 1,
						col:  2,
					},
					{
						typ:  tokenNumber,
						v:    "-1",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenEOF,
						pos:  5,
						line: 1,
						col:  6,
					},
				},
			},
		},
		{
			name: "escape at end of input",
			fields: fields{
				input: "\"abc\\",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "invalid escape sequence in string",
						pos:  5,
						line: 1,
						col:  6,
					},
				},
			},
		},
		{
			name: "short hex escape",
			fields: fields{
				input: "\"\\x4\"",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "invalid escape sequence in string",
						pos:  5,
						line: 1,
						col:  6,
					},
				},
			},
		},
		{
			name: "short unicode escape",
			fields: fields{
				input: "\"\\u00\"",
			},
			want: want{
				val: []token{
					{
						typ:  tokenError,
						v:    "invalid escape sequence in string",
						pos:  6,
						line: 1,
						col:  7,
					},
				},
			},
		},
		{
			name: "empty strings",
			fields: fields{
				input: "\"\" '' ``",
			},
			want: want{
				val: []token{
					{
						typ:  tokenString,
						v:    "\"\"",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenString,
						v:    "''",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenRawString,
						v:    "``",
						pos:  6,
						line: 1,
						col:  7,
					},
					{
						typ:  tokenEOF,
						pos:  8,
						line: 1,
						col:  9,
					},
				},
			},
		},
		{
			name: "other quote kind inside string",
			fields: fields{
				input: "'it\"s' \"it's\"",
			},
			want: want{
				val: []token{
					{
						typ:  tokenString,
						v:    "'it\"s'",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenString,
						v:    "\"it's\"",
						pos:  7,
						line: 1,
						col:  8,
					},
					{
						typ:  tokenEOF,
						pos:  13,
						line: 1,
						col:  14,
					},
				},
			},
		},
		{
			name: "raw string keeps backslash",
			fields: fields{
				input: "`a\\nb`",
			},
			want: want{
				val: []token{
					{
						typ:  tokenRawString,
						v:    "`a\\nb`",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						pos:  6,
						line: 1,
						col:  7,
					},
				},
			},
		},
		{
			name: "raw string spanning lines",
			fields: fields{
				input: "`a\nb` x",
			},
			want: want{
				val: []token{
					{
						typ:  tokenRawString,
						v:    "`a\nb`",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "x",
						pos:  6,
						line: 2,
						col:  4,
					},
					{
						typ:  tokenEOF,
						pos:  7,
						line: 2,
						col:  5,
					},
				},
			},
		},
		{
			name: "tab and crlf whitespace",
			fields: fields{
				input: "a\t\r\nb",
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "a",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "b",
						pos:  4,
						line: 2,
						col:  1,
					},
					{
						typ:  tokenEOF,
						pos:  5,
						line: 2,
						col:  2,
					},
				},
			},
		},
		{
			name: "wide unexpected character",
			fields: fields{
				input: "a ＃",
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "a",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenError,
						v:    "unexpected character U+FF03 '＃'",
						pos:  2,
						line: 1,
						col:  3,
					},
				},
			},
		},
		{
			name: "time",
			fields: fields{
				input: "2023-01-02T15:04:05Z",
			},
			want: want{
				val: []token{
					{
						typ:  tokenTime,
						v:    "2023-01-02T15:04:05Z",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						pos:  20,
						line: 1,
						col:  21,
					},
				},
			},
		},
		{
			name: "time with fraction and offset",
			fields: fields{
				input: "2023-01-02T15:04:05.123+09:00 x",
			},
			want: want{
				val: []token{
					{
						typ:  tokenTime,
						v:    "2023-01-02T15:04:05.123+09:00",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "x",
						pos:  30,
						line: 1,
						col:  31,
					},
					{
						typ:  tokenEOF,
						pos:  31,
						line: 1,
						col:  32,
					},
				},
			},
		},
		{
			name: "time without zone",
			fields: fields{
				input: "2023-01-02T15:04:05",
			},
			want: want{
				val: []token{
					{
						typ:  tokenTime,
						v:    "2023-01-02T15:04:05",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						pos:  19,
						line: 1,
						col:  20,
					},
				},
			},
		},
		{
			name: "time after operator",
			fields: fields{
				input: "T>2023-01-02T15:04:05Z",
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "T",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenGT,
						v:    ">",
						pos:  1,
						line: 1,
						col:  2,
					},
					{
						typ:  tokenTime,
						v:    "2023-01-02T15:04:05Z",
						pos:  2,
						line: 1,
						col:  3,
					},
					{
						typ:  tokenEOF,
						pos:  22,
						line: 1,
						col:  23,
					},
				},
			},
		},
		{
			name: "date only in comparison",
			fields: fields{
				input: "T>=2023-01-02&&x",
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "T",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenGTE,
						v:    ">=",
						pos:  1,
						line: 1,
						col:  2,
					},
					{
						typ:  tokenTime,
						v:    "2023-01-02",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenAND,
						v:    "&&",
						pos:  13,
						line: 1,
						col:  14,
					},
					{
						typ:  tokenIdent,
						v:    "x",
						pos:  15,
						line: 1,
						col:  16,
					},
					{
						typ:  tokenEOF,
						pos:  16,
						line: 1,
						col:  17,
					},
				},
			},
		},
		{
			name: "date followed by T alone",
			fields: fields{
				input: "2023-01-02T",
			},
			want: want{
				val: []token{
					{
						typ:  tokenTime,
						v:    "2023-01-02",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "T",
						pos:  10,
						line: 1,
						col:  11,
					},
					{
						typ:  tokenEOF,
						pos:  11,
						line: 1,
						col:  12,
					},
				},
			},
		},
		{
			name: "lowercase zone letter is not part of the time",
			fields: fields{
				input: "2023-01-02T15:04:05z",
			},
			want: want{
				val: []token{
					{
						typ:  tokenTime,
						v:    "2023-01-02T15:04:05",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "z",
						pos:  19,
						line: 1,
						col:  20,
					},
					{
						typ:  tokenEOF,
						pos:  20,
						line: 1,
						col:  21,
					},
				},
			},
		},
		{
			name: "date only",
			fields: fields{
				input: "2023-01-02",
			},
			want: want{
				val: []token{
					{
						typ:  tokenTime,
						v:    "2023-01-02",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						pos:  10,
						line: 1,
						col:  11,
					},
				},
			},
		},
		{
			name: "lowercase time separator stops the time at the date",
			fields: fields{
				input: "2023-01-02t15:04:05Z",
			},
			want: want{
				val: []token{
					{
						typ:  tokenTime,
						v:    "2023-01-02",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "t15",
						pos:  10,
						line: 1,
						col:  11,
					},
					{
						typ:  tokenError,
						v:    "unexpected character U+003A ':'",
						pos:  13,
						line: 1,
						col:  14,
					},
				},
			},
		},
		{
			name: "duration with micro sign",
			fields: fields{
				input: "1μs",
			},
			want: want{
				val: []token{
					{
						typ:  tokenDuration,
						v:    "1μs",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						pos:  4,
						line: 1,
						col:  4,
					},
				},
			},
		},
		{
			name: "number with underscore",
			fields: fields{
				input: "1_000",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "1_000",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						pos:  5,
						line: 1,
						col:  6,
					},
				},
			},
		},
		{
			name: "exponent then unit is not a duration",
			fields: fields{
				input: "1e5s",
			},
			want: want{
				val: []token{
					{
						typ:  tokenNumber,
						v:    "1e5",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "s",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenEOF,
						pos:  4,
						line: 1,
						col:  5,
					},
				},
			},
		},
		{
			name: "identifiers with underscore and digits",
			fields: fields{
				input: "_x a1 x_1",
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "_x",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenIdent,
						v:    "a1",
						pos:  3,
						line: 1,
						col:  4,
					},
					{
						typ:  tokenIdent,
						v:    "x_1",
						pos:  6,
						line: 1,
						col:  7,
					},
					{
						typ:  tokenEOF,
						pos:  9,
						line: 1,
						col:  10,
					},
				},
			},
		},
		{
			name: "number in parentheses",
			fields: fields{
				input: "(1)",
			},
			want: want{
				val: []token{
					{
						typ:  tokenLparen,
						v:    "(",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenNumber,
						v:    "1",
						pos:  1,
						line: 1,
						col:  2,
					},
					{
						typ:  tokenRparen,
						v:    ")",
						pos:  2,
						line: 1,
						col:  3,
					},
					{
						typ:  tokenEOF,
						pos:  3,
						line: 1,
						col:  4,
					},
				},
			},
		},
		{
			name: "eof repeats after the input is consumed",
			fields: fields{
				input: "a",
			},
			want: want{
				val: []token{
					{
						typ:  tokenIdent,
						v:    "a",
						pos:  0,
						line: 1,
						col:  1,
					},
					{
						typ:  tokenEOF,
						pos:  1,
						line: 1,
						col:  2,
					},
					{
						typ:  tokenEOF,
						pos:  1,
						line: 1,
						col:  2,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			got := make([]token, 0, len(test.want.val))
			for range test.want.val {
				got = append(got, l.nextToken())
			}
			if !reflect.DeepEqual(got, test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_lexer_next(t *testing.T) {
	type fields struct {
		input string
		pos   int32
		line  int32
		col   int32
	}
	type want struct {
		val  rune
		pos  int32
		line int32
		col  int32
		prev mark
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "ascii byte advances one byte and one column",
			fields: fields{
				input: "ab",
				line:  1,
				col:   1,
			},
			want: want{
				val:  'a',
				pos:  1,
				line: 1,
				col:  2,
				prev: mark{
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "empty input returns eof without moving",
			fields: fields{
				input: "",
				line:  1,
				col:   1,
			},
			want: want{
				val:  eof,
				pos:  0,
				line: 1,
				col:  1,
				prev: mark{
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "end of input returns eof and keeps the position",
			fields: fields{
				input: "a",
				pos:   1,
				line:  1,
				col:   2,
			},
			want: want{
				val:  eof,
				pos:  1,
				line: 1,
				col:  2,
				prev: mark{
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "newline advances the line and resets the column",
			fields: fields{
				input: "a\nb",
				pos:   1,
				line:  1,
				col:   2,
			},
			want: want{
				val:  '\n',
				pos:  2,
				line: 2,
				col:  1,
				prev: mark{
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "carriage return is an ordinary column",
			fields: fields{
				input: "\r\n",
				line:  1,
				col:   1,
			},
			want: want{
				val:  '\r',
				pos:  1,
				line: 1,
				col:  2,
				prev: mark{
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "two byte rune advances two bytes and one column",
			fields: fields{
				input: "é",
				line:  1,
				col:   1,
			},
			want: want{
				val:  'é',
				pos:  2,
				line: 1,
				col:  2,
				prev: mark{
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "three byte wide rune advances three bytes and two columns",
			fields: fields{
				input: "軍師",
				line:  1,
				col:   1,
			},
			want: want{
				val:  '軍',
				pos:  3,
				line: 1,
				col:  3,
				prev: mark{
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "fullwidth letter advances two columns",
			fields: fields{
				input: "Ａ",
				line:  1,
				col:   1,
			},
			want: want{
				val:  'Ａ',
				pos:  3,
				line: 1,
				col:  3,
				prev: mark{
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "four byte emoji advances four bytes and two columns",
			fields: fields{
				input: "😀",
				line:  1,
				col:   1,
			},
			want: want{
				val:  '😀',
				pos:  4,
				line: 1,
				col:  3,
				prev: mark{
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "invalid utf-8 byte yields the replacement rune and one column",
			fields: fields{
				input: "\xffa",
				line:  1,
				col:   1,
			},
			want: want{
				val:  '\ufffd',
				pos:  1,
				line: 1,
				col:  2,
				prev: mark{
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "truncated multibyte sequence yields the replacement rune",
			fields: fields{
				input: "\xe8\xbb",
				line:  1,
				col:   1,
			},
			want: want{
				val:  '\ufffd',
				pos:  1,
				line: 1,
				col:  2,
				prev: mark{
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "reads from the current offset in the middle of the input",
			fields: fields{
				input: "ab軍",
				pos:   2,
				line:  1,
				col:   3,
			},
			want: want{
				val:  '軍',
				pos:  5,
				line: 1,
				col:  5,
				prev: mark{
					pos:  2,
					line: 1,
					col:  3,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.fields.input,
				pos:   test.fields.pos,
				line:  test.fields.line,
				col:   test.fields.col,
			}
			got := want{
				val:  l.next(),
				pos:  l.pos,
				line: l.line,
				col:  l.col,
				prev: l.prev,
			}
			if got != test.want {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want)
			}
		})
	}
}

func Test_lexer_peek(t *testing.T) {
	type fields struct {
		input string
		pos   int32
		line  int32
		col   int32
	}
	type want struct {
		val  rune
		pos  int32
		line int32
		col  int32
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "ascii byte without moving",
			fields: fields{
				input: "ab",
				line:  1,
				col:   1,
			},
			want: want{
				val:  'a',
				pos:  0,
				line: 1,
				col:  1,
			},
		},
		{
			name: "wide rune without moving",
			fields: fields{
				input: "軍",
				line:  1,
				col:   1,
			},
			want: want{
				val:  '軍',
				pos:  0,
				line: 1,
				col:  1,
			},
		},
		{
			name: "newline without changing the line",
			fields: fields{
				input: "\n",
				line:  1,
				col:   1,
			},
			want: want{
				val:  '\n',
				pos:  0,
				line: 1,
				col:  1,
			},
		},
		{
			name: "eof at the end of input",
			fields: fields{
				input: "a",
				pos:   1,
				line:  1,
				col:   2,
			},
			want: want{
				val:  eof,
				pos:  1,
				line: 1,
				col:  2,
			},
		},
		{
			name: "eof on empty input",
			fields: fields{
				input: "",
				line:  1,
				col:   1,
			},
			want: want{
				val:  eof,
				pos:  0,
				line: 1,
				col:  1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.fields.input,
				pos:   test.fields.pos,
				line:  test.fields.line,
				col:   test.fields.col,
			}
			got := want{
				val:  l.peek(),
				pos:  l.pos,
				line: l.line,
				col:  l.col,
			}
			if got != test.want {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want)
			}
		})
	}
}

func Test_lexer_backup(t *testing.T) {
	type fields struct {
		input string
		steps func(l *lexer)
	}
	type want struct {
		pos  int32
		line int32
		col  int32
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "backup is idempotent until the next read",
			fields: fields{
				input: "ab",
				steps: func(l *lexer) {
					l.next()
					l.next()
					l.backup()
					l.backup()
				},
			},
			want: want{
				pos:  1,
				line: 1,
				col:  2,
			},
		},
		{
			name: "backup at end of input does not move",
			fields: fields{
				input: "a",
				steps: func(l *lexer) {
					l.next()
					l.next()
					l.backup()
				},
			},
			want: want{
				pos:  1,
				line: 1,
				col:  2,
			},
		},
		{
			name: "backup across newline restores line and column",
			fields: fields{
				input: "ab\n",
				steps: func(l *lexer) {
					l.next()
					l.next()
					l.next()
					l.backup()
				},
			},
			want: want{
				pos:  2,
				line: 1,
				col:  3,
			},
		},
		{
			name: "backup over wide rune restores column by width",
			fields: fields{
				input: "名前",
				steps: func(l *lexer) {
					l.next()
					l.next()
					l.backup()
				},
			},
			want: want{
				pos:  3,
				line: 1,
				col:  3,
			},
		},
		{
			name: "backup after reset does not move",
			fields: fields{
				input: "abc",
				steps: func(l *lexer) {
					l.next()
					m := l.mark()
					l.next()
					l.next()
					l.reset(m)
					l.backup()
				},
			},
			want: want{
				pos:  1,
				line: 1,
				col:  2,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.fields.input)
			test.fields.steps(&l)
			got := want{
				pos:  l.pos,
				line: l.line,
				col:  l.col,
			}
			if got != test.want {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want)
			}
		})
	}
}

func Test_lexer_mark(t *testing.T) {
	type fields struct {
		pos  int32
		line int32
		col  int32
	}
	type want struct {
		val mark
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "zero position",
			fields: fields{
				pos:  0,
				line: 0,
				col:  0,
			},
			want: want{
				val: mark{
					pos:  0,
					line: 0,
					col:  0,
				},
			},
		},
		{
			name: "initial position",
			fields: fields{
				pos:  0,
				line: 1,
				col:  1,
			},
			want: want{
				val: mark{
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "position on a later line",
			fields: fields{
				pos:  10,
				line: 3,
				col:  4,
			},
			want: want{
				val: mark{
					pos:  10,
					line: 3,
					col:  4,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				pos:  test.fields.pos,
				line: test.fields.line,
				col:  test.fields.col,
			}
			if got := l.mark(); got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_lexer_reset(t *testing.T) {
	type fields struct {
		pos  int32
		line int32
		col  int32
		prev mark
	}
	type args struct {
		m mark
	}
	type want struct {
		pos  int32
		line int32
		col  int32
		prev mark
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "moves back to an earlier position on the same line",
			fields: fields{
				pos:  5,
				line: 1,
				col:  6,
				prev: mark{
					pos:  4,
					line: 1,
					col:  5,
				},
			},
			args: args{
				m: mark{
					pos:  2,
					line: 1,
					col:  3,
				},
			},
			want: want{
				pos:  2,
				line: 1,
				col:  3,
				prev: mark{
					pos:  2,
					line: 1,
					col:  3,
				},
			},
		},
		{
			name: "moves back across a newline",
			fields: fields{
				pos:  4,
				line: 2,
				col:  2,
				prev: mark{
					pos:  3,
					line: 2,
					col:  1,
				},
			},
			args: args{
				m: mark{
					pos:  1,
					line: 1,
					col:  2,
				},
			},
			want: want{
				pos:  1,
				line: 1,
				col:  2,
				prev: mark{
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
		{
			name: "resetting to the current position only pins prev",
			fields: fields{
				pos:  3,
				line: 1,
				col:  4,
				prev: mark{
					pos:  2,
					line: 1,
					col:  3,
				},
			},
			args: args{
				m: mark{
					pos:  3,
					line: 1,
					col:  4,
				},
			},
			want: want{
				pos:  3,
				line: 1,
				col:  4,
				prev: mark{
					pos:  3,
					line: 1,
					col:  4,
				},
			},
		},
		{
			name: "moves forward when the mark is ahead",
			fields: fields{
				pos:  0,
				line: 1,
				col:  1,
			},
			args: args{
				m: mark{
					pos:  7,
					line: 2,
					col:  3,
				},
			},
			want: want{
				pos:  7,
				line: 2,
				col:  3,
				prev: mark{
					pos:  7,
					line: 2,
					col:  3,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				pos:  test.fields.pos,
				line: test.fields.line,
				col:  test.fields.col,
				prev: test.fields.prev,
			}
			l.reset(test.args.m)
			got := want{
				pos:  l.pos,
				line: l.line,
				col:  l.col,
				prev: l.prev,
			}
			if got != test.want {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want)
			}
		})
	}
}

func Test_lexer_emit(t *testing.T) {
	type fields struct {
		input     string
		pos       int32
		startPos  int32
		line      int32
		startLine int32
		col       int32
		startCol  int32
	}
	type args struct {
		typ tokenType
	}
	type want struct {
		token     token
		hasNext   bool
		startPos  int32
		startLine int32
		startCol  int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "token spanning the whole input",
			fields: fields{
				input:     "abc",
				pos:       3,
				startPos:  0,
				line:      1,
				startLine: 1,
				col:       4,
				startCol:  1,
			},
			args: args{
				typ: tokenIdent,
			},
			want: want{
				token: token{
					typ:  tokenIdent,
					v:    "abc",
					pos:  0,
					line: 1,
					col:  1,
				},
				hasNext:   true,
				startPos:  3,
				startLine: 1,
				startCol:  4,
			},
		},
		{
			name: "token in the middle of the input keeps its start position",
			fields: fields{
				input:     "a == b",
				pos:       4,
				startPos:  2,
				line:      1,
				startLine: 1,
				col:       5,
				startCol:  3,
			},
			args: args{
				typ: tokenEQ,
			},
			want: want{
				token: token{
					typ:  tokenEQ,
					v:    "==",
					pos:  2,
					line: 1,
					col:  3,
				},
				hasNext:   true,
				startPos:  4,
				startLine: 1,
				startCol:  5,
			},
		},
		{
			name: "empty span emits an empty value",
			fields: fields{
				input:     "abc",
				pos:       3,
				startPos:  3,
				line:      1,
				startLine: 1,
				col:       4,
				startCol:  4,
			},
			args: args{
				typ: tokenEOF,
			},
			want: want{
				token: token{
					typ:  tokenEOF,
					v:    "",
					pos:  3,
					line: 1,
					col:  4,
				},
				hasNext:   true,
				startPos:  3,
				startLine: 1,
				startCol:  4,
			},
		},
		{
			name: "token starting on an earlier line records the start line",
			fields: fields{
				input:     "a\n\"x\ny\"",
				pos:       7,
				startPos:  2,
				line:      3,
				startLine: 2,
				col:       3,
				startCol:  1,
			},
			args: args{
				typ: tokenString,
			},
			want: want{
				token: token{
					typ:  tokenString,
					v:    "\"x\ny\"",
					pos:  2,
					line: 2,
					col:  1,
				},
				hasNext:   true,
				startPos:  7,
				startLine: 3,
				startCol:  3,
			},
		},
		{
			name: "wide runes are sliced by byte offsets",
			fields: fields{
				input:     "軍師 x",
				pos:       6,
				startPos:  0,
				line:      1,
				startLine: 1,
				col:       5,
				startCol:  1,
			},
			args: args{
				typ: tokenIdent,
			},
			want: want{
				token: token{
					typ:  tokenIdent,
					v:    "軍師",
					pos:  0,
					line: 1,
					col:  1,
				},
				hasNext:   true,
				startPos:  6,
				startLine: 1,
				startCol:  5,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input:     test.fields.input,
				pos:       test.fields.pos,
				startPos:  test.fields.startPos,
				line:      test.fields.line,
				startLine: test.fields.startLine,
				col:       test.fields.col,
				startCol:  test.fields.startCol,
			}
			l.emit(test.args.typ)
			got := want{
				token:     l.token,
				hasNext:   l.hasNext,
				startPos:  l.startPos,
				startLine: l.startLine,
				startCol:  l.startCol,
			}
			if got != test.want {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want)
			}
		})
	}
}

func Test_lexer_ignore(t *testing.T) {
	type fields struct {
		pos       int32
		startPos  int32
		line      int32
		startLine int32
		col       int32
		startCol  int32
	}
	type want struct {
		startPos  int32
		startLine int32
		startCol  int32
		hasNext   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "moves the start to the current position",
			fields: fields{
				pos:       4,
				startPos:  1,
				line:      1,
				startLine: 1,
				col:       5,
				startCol:  2,
			},
			want: want{
				startPos:  4,
				startLine: 1,
				startCol:  5,
				hasNext:   false,
			},
		},
		{
			name: "across a newline the start line follows",
			fields: fields{
				pos:       6,
				startPos:  2,
				line:      3,
				startLine: 1,
				col:       1,
				startCol:  3,
			},
			want: want{
				startPos:  6,
				startLine: 3,
				startCol:  1,
				hasNext:   false,
			},
		},
		{
			name: "no pending input leaves the start unchanged",
			fields: fields{
				pos:       0,
				startPos:  0,
				line:      1,
				startLine: 1,
				col:       1,
				startCol:  1,
			},
			want: want{
				startPos:  0,
				startLine: 1,
				startCol:  1,
				hasNext:   false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				pos:       test.fields.pos,
				startPos:  test.fields.startPos,
				line:      test.fields.line,
				startLine: test.fields.startLine,
				col:       test.fields.col,
				startCol:  test.fields.startCol,
			}
			l.ignore()
			got := want{
				startPos:  l.startPos,
				startLine: l.startLine,
				startCol:  l.startCol,
				hasNext:   l.hasNext,
			}
			if got != test.want {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want)
			}
		})
	}
}

func Test_lexer_accept(t *testing.T) {
	type fields struct {
		input string
		pos   int32
		line  int32
		col   int32
	}
	type args struct {
		valid string
	}
	type want struct {
		val  bool
		pos  int32
		col  int32
		prev mark
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "matching byte is consumed",
			fields: fields{
				input: "+1",
				line:  1,
				col:   1,
			},
			args: args{
				valid: "+-",
			},
			want: want{
				val: true,
				pos: 1,
				col: 2,
				prev: mark{
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "last byte of valid matches",
			fields: fields{
				input: "-1",
				line:  1,
				col:   1,
			},
			args: args{
				valid: "+-",
			},
			want: want{
				val: true,
				pos: 1,
				col: 2,
				prev: mark{
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name: "non-matching byte is left in place",
			fields: fields{
				input: "1",
				line:  1,
				col:   1,
			},
			args: args{
				valid: "+-",
			},
			want: want{
				val: false,
				pos: 0,
				col: 1,
			},
		},
		{
			name: "empty valid set never matches",
			fields: fields{
				input: "a",
				line:  1,
				col:   1,
			},
			args: args{
				valid: "",
			},
			want: want{
				val: false,
				pos: 0,
				col: 1,
			},
		},
		{
			name: "end of input never matches",
			fields: fields{
				input: "a",
				pos:   1,
				line:  1,
				col:   2,
			},
			args: args{
				valid: "a",
			},
			want: want{
				val: false,
				pos: 1,
				col: 2,
			},
		},
		{
			name: "empty input never matches",
			fields: fields{
				input: "",
				line:  1,
				col:   1,
			},
			args: args{
				valid: "a",
			},
			want: want{
				val: false,
				pos: 0,
				col: 1,
			},
		},
		{
			name: "matches at the current offset not the start",
			fields: fields{
				input: "ab",
				pos:   1,
				line:  1,
				col:   2,
			},
			args: args{
				valid: "b",
			},
			want: want{
				val: true,
				pos: 2,
				col: 3,
				prev: mark{
					pos:  1,
					line: 1,
					col:  2,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.fields.input,
				pos:   test.fields.pos,
				line:  test.fields.line,
				col:   test.fields.col,
			}
			got := want{
				val:  l.accept(test.args.valid),
				pos:  l.pos,
				col:  l.col,
				prev: l.prev,
			}
			if got != test.want {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want)
			}
		})
	}
}

func Test_lexer_acceptRun(t *testing.T) {
	type fields struct {
		input string
		pos   int32
		col   int32
	}
	type args struct {
		valid string
	}
	type want struct {
		val int
		pos int32
		col int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "consumes the whole run and stops at the first other byte",
			fields: fields{
				input: "123abc",
				col:   1,
			},
			args: args{
				valid: "0123456789",
			},
			want: want{
				val: 3,
				pos: 3,
				col: 4,
			},
		},
		{
			name: "no matching byte consumes nothing",
			fields: fields{
				input: "abc",
				col:   1,
			},
			args: args{
				valid: "0123456789",
			},
			want: want{
				val: 0,
				pos: 0,
				col: 1,
			},
		},
		{
			name: "run stops at the end of input",
			fields: fields{
				input: "999",
				col:   1,
			},
			args: args{
				valid: "0123456789",
			},
			want: want{
				val: 3,
				pos: 3,
				col: 4,
			},
		},
		{
			name: "empty input consumes nothing",
			fields: fields{
				input: "",
				col:   1,
			},
			args: args{
				valid: "0123456789",
			},
			want: want{
				val: 0,
				pos: 0,
				col: 1,
			},
		},
		{
			name: "empty valid set consumes nothing",
			fields: fields{
				input: "123",
				col:   1,
			},
			args: args{
				valid: "",
			},
			want: want{
				val: 0,
				pos: 0,
				col: 1,
			},
		},
		{
			name: "starts at the current offset",
			fields: fields{
				input: "a12",
				pos:   1,
				col:   2,
			},
			args: args{
				valid: "0123456789_",
			},
			want: want{
				val: 2,
				pos: 3,
				col: 4,
			},
		},
		{
			name: "underscore separators are part of the run",
			fields: fields{
				input: "1_000.5",
				col:   1,
			},
			args: args{
				valid: "0123456789_",
			},
			want: want{
				val: 5,
				pos: 5,
				col: 6,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.fields.input,
				pos:   test.fields.pos,
				col:   test.fields.col,
			}
			got := want{
				val: l.acceptRun(test.args.valid),
				pos: l.pos,
				col: l.col,
			}
			if got != test.want {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want)
			}
		})
	}
}

func Test_lexer_acceptDigits(t *testing.T) {
	type fields struct {
		input string
		pos   int32
		col   int32
	}
	type args struct {
		n int
	}
	type want struct {
		val bool
		pos int32
		col int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "exactly n digits",
			fields: fields{
				input: "2025",
				col:   1,
			},
			args: args{
				n: 4,
			},
			want: want{
				val: true,
				pos: 4,
				col: 5,
			},
		},
		{
			name: "more digits than n consumes only n",
			fields: fields{
				input: "202501",
				col:   1,
			},
			args: args{
				n: 4,
			},
			want: want{
				val: true,
				pos: 4,
				col: 5,
			},
		},
		{
			name: "too few digits before the end of input keeps the consumed ones",
			fields: fields{
				input: "20",
				col:   1,
			},
			args: args{
				n: 4,
			},
			want: want{
				val: false,
				pos: 2,
				col: 3,
			},
		},
		{
			name: "non-digit inside the run stops there",
			fields: fields{
				input: "12a4",
				col:   1,
			},
			args: args{
				n: 4,
			},
			want: want{
				val: false,
				pos: 2,
				col: 3,
			},
		},
		{
			name: "non-digit at the start consumes nothing",
			fields: fields{
				input: "a123",
				col:   1,
			},
			args: args{
				n: 2,
			},
			want: want{
				val: false,
				pos: 0,
				col: 1,
			},
		},
		{
			name: "zero digits always succeeds without moving",
			fields: fields{
				input: "abc",
				col:   1,
			},
			args: args{
				n: 0,
			},
			want: want{
				val: true,
				pos: 0,
				col: 1,
			},
		},
		{
			name: "empty input fails for one digit",
			fields: fields{
				input: "",
				col:   1,
			},
			args: args{
				n: 1,
			},
			want: want{
				val: false,
				pos: 0,
				col: 1,
			},
		},
		{
			name: "digits just below and above the ascii range are rejected",
			fields: fields{
				input: "/:",
				col:   1,
			},
			args: args{
				n: 1,
			},
			want: want{
				val: false,
				pos: 0,
				col: 1,
			},
		},
		{
			name: "fullwidth digits are not digits",
			fields: fields{
				input: "１２",
				col:   1,
			},
			args: args{
				n: 1,
			},
			want: want{
				val: false,
				pos: 0,
				col: 1,
			},
		},
		{
			name: "starts at the current offset",
			fields: fields{
				input: "T15",
				pos:   1,
				col:   2,
			},
			args: args{
				n: 2,
			},
			want: want{
				val: true,
				pos: 3,
				col: 4,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.fields.input,
				pos:   test.fields.pos,
				col:   test.fields.col,
			}
			got := want{
				val: l.acceptDigits(test.args.n),
				pos: l.pos,
				col: l.col,
			}
			if got != test.want {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want)
			}
		})
	}
}

func Test_lexer_errorf(t *testing.T) {
	type fields struct {
		pos       int32
		startPos  int32
		line      int32
		startLine int32
		col       int32
		startCol  int32
	}
	type args struct {
		format string
		args   []any
	}
	type want struct {
		val     state
		token   token
		hasNext bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "plain message at the detection point",
			fields: fields{
				pos:       5,
				startPos:  2,
				line:      1,
				startLine: 1,
				col:       6,
				startCol:  3,
			},
			args: args{
				format: "unterminated quoted string",
			},
			want: want{
				val: stateDone,
				token: token{
					typ:  tokenError,
					v:    "unterminated quoted string",
					pos:  5,
					line: 1,
					col:  6,
				},
				hasNext: true,
			},
		},
		{
			name: "formatted message with a rune argument",
			fields: fields{
				pos:       1,
				startPos:  0,
				line:      1,
				startLine: 1,
				col:       2,
				startCol:  1,
			},
			args: args{
				format: "unexpected character %#U",
				args: []any{
					'#',
				},
			},
			want: want{
				val: stateDone,
				token: token{
					typ:  tokenError,
					v:    "unexpected character U+0023 '#'",
					pos:  1,
					line: 1,
					col:  2,
				},
				hasNext: true,
			},
		},
		{
			name: "position on a later line is reported as is",
			fields: fields{
				pos:       9,
				startPos:  4,
				line:      3,
				startLine: 2,
				col:       2,
				startCol:  1,
			},
			args: args{
				format: "unexpected end of input after '%s'",
				args: []any{
					"&",
				},
			},
			want: want{
				val: stateDone,
				token: token{
					typ:  tokenError,
					v:    "unexpected end of input after '&'",
					pos:  9,
					line: 3,
					col:  2,
				},
				hasNext: true,
			},
		},
		{
			name: "empty format produces an empty message",
			fields: fields{
				pos:       0,
				startPos:  0,
				line:      1,
				startLine: 1,
				col:       1,
				startCol:  1,
			},
			args: args{
				format: "",
			},
			want: want{
				val: stateDone,
				token: token{
					typ:  tokenError,
					v:    "",
					pos:  0,
					line: 1,
					col:  1,
				},
				hasNext: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				pos:       test.fields.pos,
				startPos:  test.fields.startPos,
				line:      test.fields.line,
				startLine: test.fields.startLine,
				col:       test.fields.col,
				startCol:  test.fields.startCol,
			}
			got := want{
				val:     l.errorf(test.args.format, test.args.args...),
				token:   l.token,
				hasNext: l.hasNext,
			}
			if got != test.want {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want)
			}
		})
	}
}

func Test_width(t *testing.T) {
	type args struct {
		r rune
	}
	type want struct {
		val int32
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "latin letter with diacritic",
			args: args{
				r: 'é',
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "cjk ideograph is wide",
			args: args{
				r: '軍',
			},
			want: want{
				val: 2,
			},
		},
		{
			name: "fullwidth latin letter is wide",
			args: args{
				r: 'Ａ',
			},
			want: want{
				val: 2,
			},
		},
		{
			name: "emoji is wide",
			args: args{
				r: '😀',
			},
			want: want{
				val: 2,
			},
		},
		{
			name: "combining mark is counted as one",
			args: args{
				r: '\u0301',
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "zero width space is counted as one",
			args: args{
				r: '\u200b',
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "control character is counted as one",
			args: args{
				r: '\x01',
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "replacement rune is counted as one",
			args: args{
				r: '\ufffd',
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "ascii letter is one",
			args: args{
				r: 'a',
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "halfwidth katakana is one",
			args: args{
				r: 'ｶ',
			},
			want: want{
				val: 1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := width(test.args.r); got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_isNumberStart(t *testing.T) {
	type args struct {
		r rune
	}
	type want struct {
		val bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				r: '0',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "nine",
			args: args{
				r: '9',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "plus",
			args: args{
				r: '+',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "minus",
			args: args{
				r: '-',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "dot",
			args: args{
				r: '.',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "fullwidth digit is a unicode digit",
			args: args{
				r: '１',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "arabic-indic digit is a unicode digit",
			args: args{
				r: '٣',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "letter",
			args: args{
				r: 'a',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "underscore",
			args: args{
				r: '_',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "space",
			args: args{
				r: ' ',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "byte just below zero",
			args: args{
				r: '/',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "byte just above nine",
			args: args{
				r: ':',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "eof",
			args: args{
				r: eof,
			},
			want: want{
				val: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isNumberStart(test.args.r); got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_isSpace(t *testing.T) {
	type args struct {
		r rune
	}
	type want struct {
		val bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "space",
			args: args{
				r: ' ',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "tab",
			args: args{
				r: '\t',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "carriage return",
			args: args{
				r: '\r',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "newline",
			args: args{
				r: '\n',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "vertical tab is not a space",
			args: args{
				r: '\v',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "form feed is not a space",
			args: args{
				r: '\f',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "no-break space is not a space",
			args: args{
				r: '\u00a0',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "ideographic space is not a space",
			args: args{
				r: '\u3000',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "letter",
			args: args{
				r: 'a',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "eof",
			args: args{
				r: eof,
			},
			want: want{
				val: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isSpace(test.args.r); got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_isAlphaNumeric(t *testing.T) {
	type args struct {
		r rune
	}
	type want struct {
		val bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "underscore",
			args: args{
				r: '_',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "lowercase letter",
			args: args{
				r: 'a',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "uppercase letter",
			args: args{
				r: 'Z',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "digit",
			args: args{
				r: '7',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "cjk ideograph is a letter",
			args: args{
				r: '軍',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "hiragana is a letter",
			args: args{
				r: 'あ',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "fullwidth digit is a digit",
			args: args{
				r: '１',
			},
			want: want{
				val: true,
			},
		},
		{
			name: "hyphen",
			args: args{
				r: '-',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "dot",
			args: args{
				r: '.',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "space",
			args: args{
				r: ' ',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "emoji is neither letter nor digit",
			args: args{
				r: '😀',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "combining mark is neither letter nor digit",
			args: args{
				r: '\u0301',
			},
			want: want{
				val: false,
			},
		},
		{
			name: "eof",
			args: args{
				r: eof,
			},
			want: want{
				val: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isAlphaNumeric(test.args.r); got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_isBoolLiteral(t *testing.T) {
	type args struct {
		s string
	}
	type want struct {
		val bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "lower true",
			args: args{
				s: "true",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "title true",
			args: args{
				s: "True",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "upper true",
			args: args{
				s: "TRUE",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "lower false",
			args: args{
				s: "false",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "title false",
			args: args{
				s: "False",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "upper false",
			args: args{
				s: "FALSE",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "mixed case is not a literal",
			args: args{
				s: "tRue",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "prefix with a suffix is not a literal",
			args: args{
				s: "truex",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "prefix alone is not a literal",
			args: args{
				s: "tru",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "surrounding space is not trimmed",
			args: args{
				s: " true",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: want{
				val: false,
			},
		},
		{
			name: "digit literal",
			args: args{
				s: "1",
			},
			want: want{
				val: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isBoolLiteral(test.args.s); got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}
