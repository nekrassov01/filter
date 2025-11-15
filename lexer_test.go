package filter

import (
	"reflect"
	"testing"
)

func Test_tokenType_String(t *testing.T) {
	tests := []struct {
		name     string
		typ      tokenType
		expected string
	}{
		{
			name:     "error",
			typ:      tokenError,
			expected: "error",
		},
		{
			name:     "eof",
			typ:      tokenEOF,
			expected: "EOF",
		},
		{
			name:     "ident",
			typ:      tokenIdent,
			expected: "identifier",
		},
		{
			name:     "gt",
			typ:      tokenGT,
			expected: "\"greater than\" operator",
		},
		{
			name:     "gte",
			typ:      tokenGTE,
			expected: "\"greater than or equal to\" operator",
		},
		{
			name:     "lt",
			typ:      tokenLT,
			expected: "\"less than\" operator",
		},
		{
			name:     "lte",
			typ:      tokenLTE,
			expected: "\"less than or equal to\" operator",
		},
		{
			name:     "eq",
			typ:      tokenEQ,
			expected: "\"equal to\" operator",
		},
		{
			name:     "eqi",
			typ:      tokenEQI,
			expected: "\"case-insensitive equal to\" operator",
		},
		{
			name:     "neq",
			typ:      tokenNEQ,
			expected: "\"not equal to\" operator",
		},
		{
			name:     "neqi",
			typ:      tokenNEQI,
			expected: "\"case-insensitive not equal to\" operator",
		},
		{
			name:     "req",
			typ:      tokenREQ,
			expected: "regex matching operator",
		},
		{
			name:     "reqi",
			typ:      tokenREQI,
			expected: "case-insensitive regex matching operator",
		},
		{
			name:     "nreq",
			typ:      tokenNREQ,
			expected: "negative regex matching operator",
		},
		{
			name:     "nreqi",
			typ:      tokenNREQI,
			expected: "case-insensitive negative regex matching operator",
		},
		{
			name:     "and",
			typ:      tokenAND,
			expected: "logical AND operator",
		},
		{
			name:     "or",
			typ:      tokenOR,
			expected: "logical OR operator",
		},
		{
			name:     "not",
			typ:      tokenNOT,
			expected: "logical NOT operator",
		},
		{
			name:     "(",
			typ:      tokenLparen,
			expected: "left parenthesis",
		},
		{
			name:     ")",
			typ:      tokenRparen,
			expected: "right parenthesis",
		},
		{
			name:     "string",
			typ:      tokenString,
			expected: "string",
		},
		{
			name:     "raw-string",
			typ:      tokenRawString,
			expected: "raw string",
		},
		{
			name:     "number",
			typ:      tokenNumber,
			expected: "number",
		},
		{
			name:     "duration",
			typ:      tokenDuration,
			expected: "duration",
		},
		{
			name:     "bool",
			typ:      tokenBool,
			expected: "boolean",
		},
		{
			name:     "invalid",
			typ:      256,
			expected: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := test.typ.String(); actual != test.expected {
				t.Errorf("expected %v, actual %v", test.expected, actual)
			}
		})
	}
}

func Test_tokenType_literal(t *testing.T) {
	tests := []struct {
		name     string
		typ      tokenType
		expected string
	}{
		{
			name:     "error",
			typ:      tokenError,
			expected: "",
		},
		{
			name:     "eof",
			typ:      tokenEOF,
			expected: "",
		},
		{
			name:     "ident",
			typ:      tokenIdent,
			expected: "",
		},
		{
			name:     "gt",
			typ:      tokenGT,
			expected: ">",
		},
		{
			name:     "gte",
			typ:      tokenGTE,
			expected: ">=",
		},
		{
			name:     "lt",
			typ:      tokenLT,
			expected: "<",
		},
		{
			name:     "lte",
			typ:      tokenLTE,
			expected: "<=",
		},
		{
			name:     "eq",
			typ:      tokenEQ,
			expected: "==",
		},
		{
			name:     "eqi",
			typ:      tokenEQI,
			expected: "==*",
		},
		{
			name:     "neq",
			typ:      tokenNEQ,
			expected: "!=",
		},
		{
			name:     "neqi",
			typ:      tokenNEQI,
			expected: "!=*",
		},
		{
			name:     "req",
			typ:      tokenREQ,
			expected: "=~",
		},
		{
			name:     "reqi",
			typ:      tokenREQI,
			expected: "=~*",
		},
		{
			name:     "nreq",
			typ:      tokenNREQ,
			expected: "!~",
		},
		{
			name:     "nreqi",
			typ:      tokenNREQI,
			expected: "!~*",
		},
		{
			name:     "and",
			typ:      tokenAND,
			expected: "&&",
		},
		{
			name:     "or",
			typ:      tokenOR,
			expected: "||",
		},
		{
			name:     "not",
			typ:      tokenNOT,
			expected: "!",
		},
		{
			name:     "left paren",
			typ:      tokenLparen,
			expected: "(",
		},
		{
			name:     "right paren",
			typ:      tokenRparen,
			expected: ")",
		},
		{
			name:     "string",
			typ:      tokenString,
			expected: "",
		},
		{
			name:     "raw-string",
			typ:      tokenRawString,
			expected: "",
		},
		{
			name:     "number",
			typ:      tokenNumber,
			expected: "",
		},
		{
			name:     "duration",
			typ:      tokenDuration,
			expected: "",
		},
		{
			name:     "bool",
			typ:      tokenBool,
			expected: "",
		},
		{
			name:     "invalid",
			typ:      256,
			expected: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := test.typ.literal(); actual != test.expected {
				t.Errorf("expected %v, actual %v", test.expected, actual)
			}
		})
	}
}

func Test_lex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token
	}{
		{
			name:  "simple number 1",
			input: "1",
			expected: []token{
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
		{
			name:  "simple number 2",
			input: "+1",
			expected: []token{
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
		{
			name:  "simple number 3",
			input: "-1",
			expected: []token{
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
		{
			name:  "simple number 4",
			input: ".1",
			expected: []token{
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
		{
			name:  "simple number 5",
			input: "0.1",
			expected: []token{
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
		{
			name:  "simple number 6",
			input: "0x1.fp3",
			expected: []token{
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
		{
			name:  "simple duration",
			input: "1h",
			expected: []token{
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
		{
			name:  "ident",
			input: "id IDENT_1 あいうえお",
			expected: []token{
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
		{
			name:  "comparison operators 1",
			input: "> >= < <=",
			expected: []token{
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
		{
			name:  "comparison operators 2",
			input: "== ==* != !=*",
			expected: []token{
				{
					typ:  tokenEQ,
					v:    "==",
					pos:  0,
					line: 1,
					col:  1,
				},
				{
					typ:  tokenEQI,
					v:    "==*",
					pos:  3,
					line: 1,
					col:  4,
				},
				{
					typ:  tokenNEQ,
					v:    "!=",
					pos:  7,
					line: 1,
					col:  8,
				},
				{
					typ:  tokenNEQI,
					v:    "!=*",
					pos:  10,
					line: 1,
					col:  11,
				},
				{
					typ:  tokenEOF,
					v:    "",
					pos:  13,
					line: 1,
					col:  14,
				},
			},
		},
		{
			name:  "comparison operators 3",
			input: "=~ !~",
			expected: []token{
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
		{
			name:  "logical operators",
			input: "&& || !",
			expected: []token{
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
		{
			name:  "parentheses",
			input: "()",
			expected: []token{
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
		{
			name:  "string",
			input: "\"abc\"",
			expected: []token{
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
		{
			name:  "string with escape",
			input: "\"\\n\\t\\\\\\\"\\'\\0\\a\\b\\f\\r\\v\\x41\\u0041\"",
			expected: []token{
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
		{
			name:  "string with wrong hex",
			input: "'\\xG'",
			expected: []token{
				{
					typ:  tokenError,
					v:    "invalid escape sequence in string at 1:5",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name:  "string with eof",
			input: "\"",
			expected: []token{
				{
					typ:  tokenError,
					v:    "unterminated quoted string at 1:2",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name:  "string with line break",
			input: "\"abc\\ndef\"",
			expected: []token{
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
		{
			name:  "raw string",
			input: "`abc`",
			expected: []token{
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
		{
			name:  "number",
			input: "0 1 +2 -3 0.4 .5 +0.6 -0.7 +.8 -.9 1.23e4 1.23E4 1.23e+4 1.23e-4 0x1A2b 0x1.fp3 0x1.fp+3 0x1.fp-3 0o755 0b1011",
			expected: []token{
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
		{
			name:  "duration",
			input: "1h30m+100s+1h+30m+15s-3000ms-4000us-5000ns 0.1h.5m 1y2m3w4d",
			expected: []token{
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
		{
			name:  "duration/number/ident",
			input: "1h1x",
			expected: []token{
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
		{
			name:  "bool",
			input: "true True TRUE false False FALSE tRue",
			expected: []token{
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
		{
			name:  "invalid character 1",
			input: "\\",
			expected: []token{
				{
					typ:  tokenError,
					v:    "unexpected character U+005C '\\' at 1:1",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name:  "invalid paren depth 1",
			input: "((",
			expected: []token{
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
					v:    "unclosed left parenthesis at 1:3",
					pos:  2,
					line: 1,
					col:  3,
				},
			},
		},
		{
			name:  "invalid paren depth 2",
			input: "))",
			expected: []token{
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
					v:    "unexpected right parenthesis at 1:3",
					pos:  2,
					line: 1,
					col:  3,
				},
			},
		},
		{
			name:  "invalid paren depth 3",
			input: "((())",
			expected: []token{
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
					v:    "unclosed left parenthesis at 1:6",
					pos:  5,
					line: 1,
					col:  6,
				},
			},
		},
		{
			name:  "invalid paren depth 4",
			input: "(()))",
			expected: []token{
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
					v:    "unexpected right parenthesis at 1:6",
					pos:  5,
					line: 1,
					col:  6,
				},
			},
		},
		{
			name:  "rune error in string",
			input: "\"\uFFFD\"",
			expected: []token{
				{
					typ:  tokenError,
					v:    "invalid utf8 encoding in string at 1:3",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name:  "unterminated string 1",
			input: "\"aaa bbb ccc",
			expected: []token{
				{
					typ:  tokenError,
					v:    "unterminated quoted string at 1:13",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name:  "unterminated string 2",
			input: "'aaa bbb ccc",
			expected: []token{
				{
					typ:  tokenError,
					v:    "unterminated quoted string at 1:13",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name:  "invalid escape sequence in string",
			input: "\"aaa\\zbbb\"",
			expected: []token{
				{
					typ:  tokenError,
					v:    "invalid escape sequence in string at 1:7",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name:  "rune error in raw string",
			input: "`\uFFFD`",
			expected: []token{
				{
					typ:  tokenError,
					v:    "invalid utf8 encoding in raw string at 1:3",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name:  "unterminated raw string",
			input: "`aaa bbb ccc",
			expected: []token{
				{
					typ:  tokenError,
					v:    "unterminated raw string at 1:13",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name:  "unexpected operator 1",
			input: "=!",
			expected: []token{
				{
					typ:  tokenError,
					v:    "unexpected character '!' after '=' at 1:2",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name:  "unexpected operator 2",
			input: "&|",
			expected: []token{
				{
					typ:  tokenError,
					v:    "unexpected character '|' after '&' at 1:2",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name:  "unexpected operator 3",
			input: "|&",
			expected: []token{
				{
					typ:  tokenError,
					v:    "unexpected character '&' after '|' at 1:2",
					pos:  0,
					line: 1,
					col:  1,
				},
			},
		},
		{
			name:  "bad number syntax 1",
			input: "10abc",
			expected: []token{
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
		{
			name:  "bad number syntax 2",
			input: "_",
			expected: []token{
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
		{
			name:  "multibyte",
			input: "一二三四五六七八九十",
			expected: []token{
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
		{
			name:  "mixed 1",
			input: `Class=="軍師"&&Name=~'孔明'&&(HP>50&&MP>=100&&LP!=0)&&(MAG>=20||!(SPD<20))`,
			expected: []token{
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
		{
			name: "mixed 2",
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
			expected: []token{
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
		{
			name: "newline in input",
			input: `

test1
test2



		test3



`,
			expected: []token{
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.input)
			actual := make([]token, 0, len(test.input))
			t.Logf("input: %v", test.input)
			for {
				token := l.nextToken()
				actual = append(actual, token)
				t.Logf("token: %v", token.v)
				if token.typ == tokenEOF || token.typ == tokenError {
					break
				}
			}
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf(testTemplate, test.input, test.expected, actual)
			}
		})
	}
}

func Test_lexer_scanEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "newline", input: "n", expected: true},
		{name: "tab", input: "t", expected: true},
		{name: "backslash", input: "\\", expected: true},
		{name: "quote_double", input: "\"", expected: true},
		{name: "quote_single", input: "'", expected: true},
		{name: "null", input: "0", expected: true},
		{name: "bell", input: "a", expected: true},
		{name: "backspace", input: "b", expected: true},
		{name: "formfeed", input: "f", expected: true},
		{name: "carriage_return", input: "r", expected: true},
		{name: "vertical_tab", input: "v", expected: true},
		{name: "hex", input: "x41", expected: true},
		{name: "unicode", input: "u0041", expected: true},
		{name: "invalid_char", input: "z", expected: false},
		{name: "empty", input: "", expected: false},
		{name: "eof", input: string([]byte{0}), expected: false},
		{name: "backtick", input: "`", expected: false},
		{name: "hex_short", input: "x4", expected: false},
		{name: "hex_nonhex", input: "x4G", expected: false},
		{name: "unicode_short", input: "u041", expected: false},
		{name: "unicode_nonhex", input: "u004G", expected: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.input,
				pos:   0,
			}
			actual := l.scanEscape()
			if actual != test.expected {
				t.Errorf(testTemplate, test.input, test.expected, actual)
			}
		})
	}
}

func Test_lexer_scanDuration(t *testing.T) {
	type expected struct {
		valid   bool
		matched string
	}
	tests := []struct {
		name     string
		input    string
		expected expected
	}{
		{name: "hour", input: "1h", expected: expected{valid: true, matched: "1h"}},
		{name: "minute", input: "1m", expected: expected{valid: true, matched: "1m"}},
		{name: "second", input: "1s", expected: expected{valid: true, matched: "1s"}},
		{name: "millisecond", input: "1ms", expected: expected{valid: true, matched: "1ms"}},
		{name: "microsecond 1", input: "1us", expected: expected{valid: true, matched: "1us"}},
		{name: "microsecond 2", input: "1μs", expected: expected{valid: true, matched: "1μs"}},
		{name: "nanosecond", input: "1ns", expected: expected{valid: true, matched: "1ns"}},
		{name: "sign 1", input: "+1h", expected: expected{valid: true, matched: "+1h"}},
		{name: "sign 2", input: "-1h", expected: expected{valid: true, matched: "-1h"}},
		{name: "float 1", input: "0.1h", expected: expected{valid: true, matched: "0.1h"}},
		{name: "float 2", input: "1.1h", expected: expected{valid: true, matched: "1.1h"}},
		{name: "float 3", input: ".1h", expected: expected{valid: true, matched: ".1h"}},
		{name: "float 4", input: "1.h", expected: expected{valid: true, matched: "1.h"}},
		{name: "mixed 1", input: "1h5000ns", expected: expected{valid: true, matched: "1h5000ns"}},
		{name: "mixed 2", input: "5000ns1h", expected: expected{valid: true, matched: "5000ns1h"}},
		{name: "mixed 3", input: "+1h5000ns", expected: expected{valid: true, matched: "+1h5000ns"}},
		{name: "mixed 4", input: "-5000ns1h", expected: expected{valid: true, matched: "-5000ns1h"}},
		{name: "mixed 5", input: "0.1h0.30m", expected: expected{valid: true, matched: "0.1h0.30m"}},
		{name: "mixed 6", input: ".1m.30s", expected: expected{valid: true, matched: ".1m.30s"}},
		{name: "mixed 7", input: "-1.1h.30m", expected: expected{valid: true, matched: "-1.1h.30m"}},
		{name: "mixed 8", input: "+0.1h.30m", expected: expected{valid: true, matched: "+0.1h.30m"}},
		{name: "mixed 9", input: "-.1h.30m", expected: expected{valid: true, matched: "-.1h.30m"}},
		{name: "mixed 10", input: "+.1h.30m", expected: expected{valid: true, matched: "+.1h.30m"}},
		{name: "mixed 11", input: "+1.h30.m", expected: expected{valid: true, matched: "+1.h30.m"}},
		{name: "full", input: "1h30m15s3000ms4000us5000ns", expected: expected{valid: true, matched: "1h30m15s3000ms4000us5000ns"}},
		{name: "duplicated", input: "1h1h", expected: expected{valid: true, matched: "1h1h"}},
		{name: "longest match 1", input: "1h+30m", expected: expected{valid: true, matched: "1h"}},
		{name: "longest match 2", input: "1h-30m", expected: expected{valid: true, matched: "1h"}},
		{name: "longest match 3", input: "+1h+30m+15s+3000ms+4000us+5000ns", expected: expected{valid: true, matched: "+1h"}},
		{name: "longest match 4", input: "-1h-30m-15s-3000ms-4000us-5000ns", expected: expected{valid: true, matched: "-1h"}},
		{name: "longest match 5", input: "1hm", expected: expected{valid: true, matched: "1h"}},
		{name: "longest match 6", input: "1hms", expected: expected{valid: true, matched: "1h"}},
		{name: "longest match 7", input: "1hd", expected: expected{valid: true, matched: "1h"}},
		{name: "longest match 8", input: "1h30m1d", expected: expected{valid: true, matched: "1h30m"}},
		{name: "longest match 9", input: "1h30md", expected: expected{valid: true, matched: "1h30m"}},
		{name: "longest match 10", input: "1h_", expected: expected{valid: true, matched: "1h"}},
		{name: "invalid multiple dot but passed 1", input: "0..1h", expected: expected{valid: true, matched: "0..1h"}},
		{name: "invalid multiple dot but passed 2", input: "..1h", expected: expected{valid: true, matched: "..1h"}},
		{name: "number 1", input: "1", expected: expected{valid: false, matched: ""}},
		{name: "number 2", input: "+1", expected: expected{valid: false, matched: ""}},
		{name: "number 3", input: "-1", expected: expected{valid: false, matched: ""}},
		{name: "invalid unit 1", input: "365d", expected: expected{valid: false, matched: ""}},
		{name: "invalid unit 4", input: "1d30m", expected: expected{valid: false, matched: ""}},
		{name: "only unit 1", input: "h", expected: expected{valid: false, matched: ""}},
		{name: "only unit 2", input: "ms", expected: expected{valid: false, matched: ""}},
		{name: "only sign 1", input: "+", expected: expected{valid: false, matched: ""}},
		{name: "only sign 2", input: "-", expected: expected{valid: false, matched: ""}},
		{name: "sign and unit 1", input: "+ms", expected: expected{valid: false, matched: ""}},
		{name: "sign and unit 2", input: "-ms", expected: expected{valid: false, matched: ""}},
		{name: "unary operator 1", input: "- 1ms", expected: expected{valid: false, matched: ""}},
		{name: "unary operator 2", input: "+ ms", expected: expected{valid: false, matched: ""}},
		{name: "unit repeat 1", input: "msms", expected: expected{valid: false, matched: ""}},
		{name: "empty", input: "", expected: expected{valid: false, matched: ""}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer{
				input: test.input,
				pos:   0,
			}
			actual := l.scanDuration()
			if actual != test.expected.valid {
				t.Errorf(testTemplate, test.input, test.expected.valid, actual)
			}
			if test.input[l.startPos:l.pos] != test.expected.matched {
				t.Errorf(testTemplate, test.input, test.expected.matched, test.input[l.startPos:l.pos])
			}
		})
	}
}
