package filter

import "testing"

func Test_tokenType_String(t *testing.T) {
	type want struct {
		val string
	}
	tests := []struct {
		name string
		tr   tokenType
		want want
	}{
		{
			name: "error",
			tr:   tokenError,
			want: want{
				val: "error",
			},
		},
		{
			name: "eof",
			tr:   tokenEOF,
			want: want{
				val: "EOF",
			},
		},
		{
			name: "ident",
			tr:   tokenIdent,
			want: want{
				val: "identifier",
			},
		},
		{
			name: "gt",
			tr:   tokenGT,
			want: want{
				val: "\"greater than\" operator",
			},
		},
		{
			name: "gte",
			tr:   tokenGTE,
			want: want{
				val: "\"greater than or equal to\" operator",
			},
		},
		{
			name: "lt",
			tr:   tokenLT,
			want: want{
				val: "\"less than\" operator",
			},
		},
		{
			name: "lte",
			tr:   tokenLTE,
			want: want{
				val: "\"less than or equal to\" operator",
			},
		},
		{
			name: "eq",
			tr:   tokenEQ,
			want: want{
				val: "\"equal to\" operator",
			},
		},
		{
			name: "neq",
			tr:   tokenNEQ,
			want: want{
				val: "\"not equal to\" operator",
			},
		},
		{
			name: "req",
			tr:   tokenREQ,
			want: want{
				val: "regex matching operator",
			},
		},
		{
			name: "nreq",
			tr:   tokenNREQ,
			want: want{
				val: "negative regex matching operator",
			},
		},
		{
			name: "and",
			tr:   tokenAND,
			want: want{
				val: "logical AND operator",
			},
		},
		{
			name: "or",
			tr:   tokenOR,
			want: want{
				val: "logical OR operator",
			},
		},
		{
			name: "not",
			tr:   tokenNOT,
			want: want{
				val: "logical NOT operator",
			},
		},
		{
			name: "(",
			tr:   tokenLparen,
			want: want{
				val: "left parenthesis",
			},
		},
		{
			name: ")",
			tr:   tokenRparen,
			want: want{
				val: "right parenthesis",
			},
		},
		{
			name: "string",
			tr:   tokenString,
			want: want{
				val: "string",
			},
		},
		{
			name: "raw-string",
			tr:   tokenRawString,
			want: want{
				val: "raw string",
			},
		},
		{
			name: "number",
			tr:   tokenNumber,
			want: want{
				val: "number",
			},
		},
		{
			name: "duration",
			tr:   tokenDuration,
			want: want{
				val: "duration",
			},
		},
		{
			name: "time",
			tr:   tokenTime,
			want: want{
				val: "time",
			},
		},
		{
			name: "bool",
			tr:   tokenBool,
			want: want{
				val: "boolean",
			},
		},
		{
			name: "invalid",
			tr:   255,
			want: want{
				val: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.tr.String()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_tokenType_literal(t *testing.T) {
	type want struct {
		val string
	}
	tests := []struct {
		name string
		tr   tokenType
		want want
	}{
		{
			name: "error",
			tr:   tokenError,
			want: want{
				val: "",
			},
		},
		{
			name: "eof",
			tr:   tokenEOF,
			want: want{
				val: "",
			},
		},
		{
			name: "ident",
			tr:   tokenIdent,
			want: want{
				val: "",
			},
		},
		{
			name: "gt",
			tr:   tokenGT,
			want: want{
				val: ">",
			},
		},
		{
			name: "gte",
			tr:   tokenGTE,
			want: want{
				val: ">=",
			},
		},
		{
			name: "lt",
			tr:   tokenLT,
			want: want{
				val: "<",
			},
		},
		{
			name: "lte",
			tr:   tokenLTE,
			want: want{
				val: "<=",
			},
		},
		{
			name: "eq",
			tr:   tokenEQ,
			want: want{
				val: "==",
			},
		},
		{
			name: "neq",
			tr:   tokenNEQ,
			want: want{
				val: "!=",
			},
		},
		{
			name: "req",
			tr:   tokenREQ,
			want: want{
				val: "=~",
			},
		},
		{
			name: "nreq",
			tr:   tokenNREQ,
			want: want{
				val: "!~",
			},
		},
		{
			name: "and",
			tr:   tokenAND,
			want: want{
				val: "&&",
			},
		},
		{
			name: "or",
			tr:   tokenOR,
			want: want{
				val: "||",
			},
		},
		{
			name: "not",
			tr:   tokenNOT,
			want: want{
				val: "!",
			},
		},
		{
			name: "left paren",
			tr:   tokenLparen,
			want: want{
				val: "(",
			},
		},
		{
			name: "right paren",
			tr:   tokenRparen,
			want: want{
				val: ")",
			},
		},
		{
			name: "string",
			tr:   tokenString,
			want: want{
				val: "",
			},
		},
		{
			name: "raw-string",
			tr:   tokenRawString,
			want: want{
				val: "",
			},
		},
		{
			name: "number",
			tr:   tokenNumber,
			want: want{
				val: "",
			},
		},
		{
			name: "duration",
			tr:   tokenDuration,
			want: want{
				val: "",
			},
		},
		{
			name: "bool",
			tr:   tokenBool,
			want: want{
				val: "",
			},
		},
		{
			name: "invalid",
			tr:   255,
			want: want{
				val: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.tr.literal()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_tokenType_isComparisonOperatorType(t *testing.T) {
	type want struct {
		val bool
	}
	tests := []struct {
		name string
		tr   tokenType
		want want
	}{
		{
			name: "error",
			tr:   tokenError,
			want: want{
				val: false,
			},
		},
		{
			name: "eof",
			tr:   tokenEOF,
			want: want{
				val: false,
			},
		},
		{
			name: "ident",
			tr:   tokenIdent,
			want: want{
				val: false,
			},
		},
		{
			name: "gt",
			tr:   tokenGT,
			want: want{
				val: true,
			},
		},
		{
			name: "gte",
			tr:   tokenGTE,
			want: want{
				val: true,
			},
		},
		{
			name: "lt",
			tr:   tokenLT,
			want: want{
				val: true,
			},
		},
		{
			name: "lte",
			tr:   tokenLTE,
			want: want{
				val: true,
			},
		},
		{
			name: "eq",
			tr:   tokenEQ,
			want: want{
				val: true,
			},
		},
		{
			name: "neq",
			tr:   tokenNEQ,
			want: want{
				val: true,
			},
		},
		{
			name: "req",
			tr:   tokenREQ,
			want: want{
				val: true,
			},
		},
		{
			name: "nreq",
			tr:   tokenNREQ,
			want: want{
				val: true,
			},
		},
		{
			name: "and",
			tr:   tokenAND,
			want: want{
				val: false,
			},
		},
		{
			name: "or",
			tr:   tokenOR,
			want: want{
				val: false,
			},
		},
		{
			name: "not",
			tr:   tokenNOT,
			want: want{
				val: false,
			},
		},
		{
			name: "lparen",
			tr:   tokenLparen,
			want: want{
				val: false,
			},
		},
		{
			name: "rparen",
			tr:   tokenRparen,
			want: want{
				val: false,
			},
		},
		{
			name: "string",
			tr:   tokenString,
			want: want{
				val: false,
			},
		},
		{
			name: "raw string",
			tr:   tokenRawString,
			want: want{
				val: false,
			},
		},
		{
			name: "number",
			tr:   tokenNumber,
			want: want{
				val: false,
			},
		},
		{
			name: "duration",
			tr:   tokenDuration,
			want: want{
				val: false,
			},
		},
		{
			name: "time",
			tr:   tokenTime,
			want: want{
				val: false,
			},
		},
		{
			name: "bool",
			tr:   tokenBool,
			want: want{
				val: false,
			},
		},
		{
			name: "out of range",
			tr:   tokenType(255),
			want: want{
				val: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.tr.isComparisonOperatorType()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_tokenType_isRegexOperatorType(t *testing.T) {
	type want struct {
		val bool
	}
	tests := []struct {
		name string
		tr   tokenType
		want want
	}{
		{
			name: "error",
			tr:   tokenError,
			want: want{
				val: false,
			},
		},
		{
			name: "eof",
			tr:   tokenEOF,
			want: want{
				val: false,
			},
		},
		{
			name: "ident",
			tr:   tokenIdent,
			want: want{
				val: false,
			},
		},
		{
			name: "gt",
			tr:   tokenGT,
			want: want{
				val: false,
			},
		},
		{
			name: "gte",
			tr:   tokenGTE,
			want: want{
				val: false,
			},
		},
		{
			name: "lt",
			tr:   tokenLT,
			want: want{
				val: false,
			},
		},
		{
			name: "lte",
			tr:   tokenLTE,
			want: want{
				val: false,
			},
		},
		{
			name: "eq",
			tr:   tokenEQ,
			want: want{
				val: false,
			},
		},
		{
			name: "neq",
			tr:   tokenNEQ,
			want: want{
				val: false,
			},
		},
		{
			name: "req",
			tr:   tokenREQ,
			want: want{
				val: true,
			},
		},
		{
			name: "nreq",
			tr:   tokenNREQ,
			want: want{
				val: true,
			},
		},
		{
			name: "and",
			tr:   tokenAND,
			want: want{
				val: false,
			},
		},
		{
			name: "or",
			tr:   tokenOR,
			want: want{
				val: false,
			},
		},
		{
			name: "not",
			tr:   tokenNOT,
			want: want{
				val: false,
			},
		},
		{
			name: "lparen",
			tr:   tokenLparen,
			want: want{
				val: false,
			},
		},
		{
			name: "rparen",
			tr:   tokenRparen,
			want: want{
				val: false,
			},
		},
		{
			name: "string",
			tr:   tokenString,
			want: want{
				val: false,
			},
		},
		{
			name: "raw string",
			tr:   tokenRawString,
			want: want{
				val: false,
			},
		},
		{
			name: "number",
			tr:   tokenNumber,
			want: want{
				val: false,
			},
		},
		{
			name: "duration",
			tr:   tokenDuration,
			want: want{
				val: false,
			},
		},
		{
			name: "time",
			tr:   tokenTime,
			want: want{
				val: false,
			},
		},
		{
			name: "bool",
			tr:   tokenBool,
			want: want{
				val: false,
			},
		},
		{
			name: "out of range",
			tr:   tokenType(255),
			want: want{
				val: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.tr.isRegexOperatorType()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_tokenType_isValueType(t *testing.T) {
	type want struct {
		val bool
	}
	tests := []struct {
		name string
		tr   tokenType
		want want
	}{
		{
			name: "error",
			tr:   tokenError,
			want: want{
				val: false,
			},
		},
		{
			name: "eof",
			tr:   tokenEOF,
			want: want{
				val: false,
			},
		},
		{
			name: "ident",
			tr:   tokenIdent,
			want: want{
				val: false,
			},
		},
		{
			name: "gt",
			tr:   tokenGT,
			want: want{
				val: false,
			},
		},
		{
			name: "gte",
			tr:   tokenGTE,
			want: want{
				val: false,
			},
		},
		{
			name: "lt",
			tr:   tokenLT,
			want: want{
				val: false,
			},
		},
		{
			name: "lte",
			tr:   tokenLTE,
			want: want{
				val: false,
			},
		},
		{
			name: "eq",
			tr:   tokenEQ,
			want: want{
				val: false,
			},
		},
		{
			name: "neq",
			tr:   tokenNEQ,
			want: want{
				val: false,
			},
		},
		{
			name: "req",
			tr:   tokenREQ,
			want: want{
				val: false,
			},
		},
		{
			name: "nreq",
			tr:   tokenNREQ,
			want: want{
				val: false,
			},
		},
		{
			name: "and",
			tr:   tokenAND,
			want: want{
				val: false,
			},
		},
		{
			name: "or",
			tr:   tokenOR,
			want: want{
				val: false,
			},
		},
		{
			name: "not",
			tr:   tokenNOT,
			want: want{
				val: false,
			},
		},
		{
			name: "lparen",
			tr:   tokenLparen,
			want: want{
				val: false,
			},
		},
		{
			name: "rparen",
			tr:   tokenRparen,
			want: want{
				val: false,
			},
		},
		{
			name: "string",
			tr:   tokenString,
			want: want{
				val: true,
			},
		},
		{
			name: "raw string",
			tr:   tokenRawString,
			want: want{
				val: true,
			},
		},
		{
			name: "number",
			tr:   tokenNumber,
			want: want{
				val: true,
			},
		},
		{
			name: "duration",
			tr:   tokenDuration,
			want: want{
				val: true,
			},
		},
		{
			name: "time",
			tr:   tokenTime,
			want: want{
				val: true,
			},
		},
		{
			name: "bool",
			tr:   tokenBool,
			want: want{
				val: true,
			},
		},
		{
			name: "out of range",
			tr:   tokenType(255),
			want: want{
				val: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.tr.isValueType()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_tokenType_isStringType(t *testing.T) {
	type want struct {
		val bool
	}
	tests := []struct {
		name string
		tr   tokenType
		want want
	}{
		{
			name: "string",
			tr:   tokenString,
			want: want{
				val: true,
			},
		},
		{
			name: "raw string",
			tr:   tokenRawString,
			want: want{
				val: true,
			},
		},
		{
			name: "number",
			tr:   tokenNumber,
			want: want{
				val: false,
			},
		},
		{
			name: "ident",
			tr:   tokenIdent,
			want: want{
				val: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.tr.isStringType()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}
