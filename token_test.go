package filter

import (
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
			name:     "neq",
			typ:      tokenNEQ,
			expected: "\"not equal to\" operator",
		},
		{
			name:     "req",
			typ:      tokenREQ,
			expected: "regex matching operator",
		},
		{
			name:     "nreq",
			typ:      tokenNREQ,
			expected: "negative regex matching operator",
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
			name:     "time",
			typ:      tokenTime,
			expected: "time",
		},
		{
			name:     "bool",
			typ:      tokenBool,
			expected: "boolean",
		},
		{
			name:     "invalid",
			typ:      255,
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
			name:     "neq",
			typ:      tokenNEQ,
			expected: "!=",
		},
		{
			name:     "req",
			typ:      tokenREQ,
			expected: "=~",
		},
		{
			name:     "nreq",
			typ:      tokenNREQ,
			expected: "!~",
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
			typ:      255,
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

func Test_tokenType_isStringType(t *testing.T) {
	tests := []struct {
		name     string
		typ      tokenType
		expected bool
	}{
		{
			name:     "string",
			typ:      tokenString,
			expected: true,
		},
		{
			name:     "raw string",
			typ:      tokenRawString,
			expected: true,
		},
		{
			name:     "number",
			typ:      tokenNumber,
			expected: false,
		},
		{
			name:     "ident",
			typ:      tokenIdent,
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.typ.isStringType()
			if actual != test.expected {
				t.Errorf(testTemplate, test.typ, test.expected, actual)
			}
		})
	}
}
