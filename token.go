package filter

// token represents a token produced by the lexer.
type token struct {
	typ  tokenType
	v    string
	pos  int32
	line int32
	col  int32
	idx  int32
}

// tokenType represents the type of token produced by the lexer.
type tokenType uint8

// Token types produced by the lexer.
const (
	tokenError     tokenType = iota // error
	tokenEOF                        // end of file
	tokenIdent                      // identifier
	tokenGT                         // greater than
	tokenGTE                        // greater than or equal to
	tokenLT                         // less than
	tokenLTE                        // less than or equal to
	tokenEQ                         // equal to
	tokenNEQ                        // not equal to
	tokenREQ                        // matches regular expression
	tokenNREQ                       // does not match regular expression
	tokenAND                        // logical AND
	tokenOR                         // logical OR
	tokenNOT                        // logical NOT
	tokenLparen                     // left parenthesis
	tokenRparen                     // right parenthesis
	tokenString                     // string literal
	tokenRawString                  // raw string literal
	tokenNumber                     // number literal
	tokenDuration                   // duration literal
	tokenTime                       // time literal
	tokenBool                       // boolean literal
)

// String returns a human-readable name for the token type.
func (t tokenType) String() string {
	switch t {
	case tokenError:
		return "error"
	case tokenEOF:
		return "EOF"
	case tokenIdent:
		return "identifier"
	case tokenGT:
		return "\"greater than\" operator"
	case tokenGTE:
		return "\"greater than or equal to\" operator"
	case tokenLT:
		return "\"less than\" operator"
	case tokenLTE:
		return "\"less than or equal to\" operator"
	case tokenEQ:
		return "\"equal to\" operator"
	case tokenNEQ:
		return "\"not equal to\" operator"
	case tokenREQ:
		return "regex matching operator"
	case tokenNREQ:
		return "negative regex matching operator"
	case tokenAND:
		return "logical AND operator"
	case tokenOR:
		return "logical OR operator"
	case tokenNOT:
		return "logical NOT operator"
	case tokenLparen:
		return "left parenthesis"
	case tokenRparen:
		return "right parenthesis"
	case tokenString:
		return "string"
	case tokenRawString:
		return "raw string"
	case tokenNumber:
		return "number"
	case tokenDuration:
		return "duration"
	case tokenTime:
		return "time"
	case tokenBool:
		return "boolean"
	default:
		return ""
	}
}

// literal returns the source text of an operator or delimiter token, or an
// empty string for tokens whose text is not fixed.
func (t tokenType) literal() string {
	switch t {
	case tokenGT:
		return ">"
	case tokenGTE:
		return ">="
	case tokenLT:
		return "<"
	case tokenLTE:
		return "<="
	case tokenEQ:
		return "=="
	case tokenNEQ:
		return "!="
	case tokenREQ:
		return "=~"
	case tokenNREQ:
		return "!~"
	case tokenAND:
		return "&&"
	case tokenOR:
		return "||"
	case tokenNOT:
		return "!"
	case tokenLparen:
		return "("
	case tokenRparen:
		return ")"
	default:
		return ""
	}
}

// isPredicateOperatorType reports whether the token is a predicate operator.
func (t tokenType) isPredicateOperatorType() bool {
	switch t {
	case tokenEQ, tokenNEQ, tokenGT, tokenGTE, tokenLT, tokenLTE, tokenREQ, tokenNREQ:
		return true
	default:
		return false
	}
}

// isRegexOperatorType reports whether the token is a regex operator.
func (t tokenType) isRegexOperatorType() bool {
	switch t {
	case tokenREQ, tokenNREQ:
		return true
	default:
		return false
	}
}

// isValueType reports whether the token can be the right-hand side of a predicate.
func (t tokenType) isValueType() bool {
	switch t {
	case tokenString, tokenRawString, tokenNumber, tokenTime, tokenDuration, tokenBool:
		return true
	default:
		return false
	}
}

// isStringType reports whether the token is a quoted or raw string literal.
func (t tokenType) isStringType() bool {
	switch t {
	case tokenString, tokenRawString:
		return true
	default:
		return false
	}
}
