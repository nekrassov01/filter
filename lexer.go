package filter

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

// eof is the rune returned by next at the end of input.
const eof = -1

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

// isComparisonOperatorType reports whether the token is a comparison operator.
func (t tokenType) isComparisonOperatorType() bool {
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

// isValueType reports whether the token can be the right-hand side of a comparison.
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

// state represents the state of the scanner between tokens.
// Dispatch is by value rather than by function pointer so that the lexer can
// stay on the caller's stack (an indirect call would force it to escape).
type state uint8

// Scanner states.
const (
	stateDone state = iota // scanning finished; only EOF tokens follow
	stateStmt              // ready to scan the next token
)

// mark is a position in the input that the lexer can return to.
type mark struct {
	pos  int32
	line int32
	col  int32
}

// lexer scans an input string into tokens on demand.
type lexer struct {
	input      string // the string being scanned
	state      state  // current state
	token      token  // last emitted token waiting to be consumed
	hasNext    bool   // flag there is a pending token
	prev       mark   // position before the last next; backup returns here
	parenDepth int    // nesting depth of ( ) exprs
	pos        int32  // current byte offset in the input
	startPos   int32  // byte offset where the current token starts
	line       int32  // 1+number of newlines seen
	startLine  int32  // line where the current token starts
	col        int32  // 1+display width of the runes since the last newline
	startCol   int32  // column where the current token starts
}

// newLexer creates a new lexer for the input string.
func newLexer(input string) lexer {
	return lexer{
		input:     input,
		state:     stateStmt,
		line:      1,
		startLine: 1,
		col:       1,
		startCol:  1,
	}
}

// lexStmt dispatches on the next rune to the state that scans the token it starts.
func (l *lexer) lexStmt() state {
	switch r := l.next(); {
	case r == eof:
		return l.lexEOF()
	case isSpace(r):
		return l.lexSpace()
	case r == '"':
		return l.lexDoubleQuotedString()
	case r == '\'':
		return l.lexSingleQuotedString()
	case r == '`':
		return l.lexRawString()
	case r == '(':
		return l.lexLparen()
	case r == ')':
		return l.lexRparen()
	case r == '=':
		return l.lexEQ()
	case r == '!':
		return l.lexNOT()
	case r == '<':
		return l.lexLT()
	case r == '>':
		return l.lexGT()
	case r == '&':
		return l.lexAND()
	case r == '|':
		return l.lexOR()
	case isNumberStart(r):
		return l.lexNumber()
	case unicode.IsLetter(r) || r == '_':
		return l.lexKeywordOrIdent()
	default:
		l.backup()
		return l.errorf("unexpected character %#U", r)
	}
}

// lexEOF emits the EOF token once the input is consumed, or an error when
// parentheses are unbalanced.
func (l *lexer) lexEOF() state {
	if l.parenDepth < 0 {
		return l.errorf("unexpected right parenthesis")
	}
	if l.parenDepth > 0 {
		return l.errorf("unclosed left parenthesis")
	}
	l.emit(tokenEOF)
	return stateDone
}

// lexSpace scans a run of space characters.
// One space has already been seen.
func (l *lexer) lexSpace() state {
	for int(l.pos) < len(l.input) {
		switch l.input[l.pos] {
		case ' ', '\t', '\r':
			l.pos++
			l.col++
		case '\n':
			l.pos++
			l.line++
			l.col = 1
		default:
			l.ignore()
			return stateStmt
		}
	}
	l.ignore()
	return stateStmt
}

// lexLparen emits a left parenthesis.
func (l *lexer) lexLparen() state {
	l.emit(tokenLparen)
	l.parenDepth++
	return stateStmt
}

// lexRparen emits a right parenthesis.
func (l *lexer) lexRparen() state {
	l.emit(tokenRparen)
	l.parenDepth--
	return stateStmt
}

// lexEQ scans the operators starting with '=': == and =~.
// The leading '=' has already been seen.
func (l *lexer) lexEQ() state {
	switch r := l.peek(); r {
	case '=':
		l.next()
		l.emit(tokenEQ)
	case '~':
		l.next()
		l.emit(tokenREQ)
	case eof:
		return l.errorf("unexpected end of input after '='")
	default:
		return l.errorf("unexpected character %q after '='", r)
	}
	return stateStmt
}

// lexNOT scans the operators starting with '!': !=, !~, and the unary NOT
// when no operator character follows. The leading '!' has already been seen.
func (l *lexer) lexNOT() state {
	switch l.peek() {
	case '=':
		l.next()
		l.emit(tokenNEQ)
	case '~':
		l.next()
		l.emit(tokenNREQ)
	default:
		l.emit(tokenNOT)
	}
	return stateStmt
}

// lexLT scans the < and <= operators.
// The leading '<' has already been seen.
func (l *lexer) lexLT() state {
	if l.peek() == '=' {
		l.next()
		l.emit(tokenLTE)
	} else {
		l.emit(tokenLT)
	}
	return stateStmt
}

// lexGT scans the > and >= operators.
// The leading '>' has already been seen.
func (l *lexer) lexGT() state {
	if l.peek() == '=' {
		l.next()
		l.emit(tokenGTE)
	} else {
		l.emit(tokenGT)
	}
	return stateStmt
}

// lexAND scans the && operator.
// The leading '&' has already been seen.
func (l *lexer) lexAND() state {
	switch r := l.peek(); r {
	case '&':
		l.next()
		l.emit(tokenAND)
	case eof:
		return l.errorf("unexpected end of input after '&'")
	default:
		return l.errorf("unexpected character %q after '&'", r)
	}
	return stateStmt
}

// lexOR scans the || operator.
// The leading '|' has already been seen.
func (l *lexer) lexOR() state {
	switch r := l.peek(); r {
	case '|':
		l.next()
		l.emit(tokenOR)
	case eof:
		return l.errorf("unexpected end of input after '|'")
	default:
		return l.errorf("unexpected character %q after '|'", r)
	}
	return stateStmt
}

// lexDoubleQuotedString scans a double-quoted string.
// One double quote has already been seen.
func (l *lexer) lexDoubleQuotedString() state {
	return l.lexString('"')
}

// lexSingleQuotedString scans a single-quoted string.
// One single quote has already been seen.
func (l *lexer) lexSingleQuotedString() state {
	return l.lexString('\'')
}

// lexString scans a quoted string up to the matching closing quote,
// validating escape sequences. The opening quote has already been seen.
func (l *lexer) lexString(quote rune) state {
Loop:
	for {
		switch l.next() {
		case utf8.RuneError:
			return l.errorf("invalid utf8 encoding in string")
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '\\':
			if !l.scanEscape() {
				return l.errorf("invalid escape sequence in string")
			}
		case quote:
			break Loop
		}
	}
	l.emit(tokenString)
	return stateStmt
}

// lexRawString scans a backtick quoted string.
// One backtick has already been seen.
func (l *lexer) lexRawString() state {
Loop:
	for {
		switch l.next() {
		case utf8.RuneError:
			return l.errorf("invalid utf8 encoding in raw string")
		case eof:
			return l.errorf("unterminated raw string")
		case '`':
			break Loop
		}
	}
	l.emit(tokenRawString)
	return stateStmt
}

// lexNumber scans a time, duration, or number literal, trying them in that
// order, so that 2023-01-02 is a date rather than three numbers.
// The leading digit, sign, or dot has already been seen.
func (l *lexer) lexNumber() state {
	l.backup() // rescan the leading character consumed by lexStmt
	start := l.mark()
	if l.scanTime() {
		l.emit(tokenTime)
		return stateStmt
	}
	l.reset(start)
	if l.scanDuration() {
		l.emit(tokenDuration)
		return stateStmt
	}
	l.reset(start)
	l.scanNumber()
	l.emit(tokenNumber)
	return stateStmt
}

// lexKeywordOrIdent scans an identifier and emits it as a boolean literal
// when it spells true or false. The leading character has already been seen.
func (l *lexer) lexKeywordOrIdent() state {
	// ASCII bytes advance without decoding; anything else takes the rune path.
	for int(l.pos) < len(l.input) {
		b := l.input[l.pos]
		if b >= utf8.RuneSelf {
			break
		}
		if !('a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || '0' <= b && b <= '9' || b == '_') {
			goto emit
		}
		l.pos++
		l.col++
	}
	for {
		r := l.next()
		if !isAlphaNumeric(r) && r != '_' {
			l.backup()
			break
		}
	}
emit:
	if isBoolLiteral(l.input[l.startPos:l.pos]) {
		l.emit(tokenBool)
		return stateStmt
	}
	l.emit(tokenIdent)
	return stateStmt
}

// scanEscape consumes the escape sequence following a backslash and reports
// whether it is valid.
func (l *lexer) scanEscape() bool {
	r := l.next()
	switch r {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\':
		// These are valid escape sequences
		return true
	case '"', '\'':
		// escaped quotes are valid in strings
		return true
	case '0':
		// Simple \0 for null character
		return true
	case 'x':
		// \xHH - 2 digit hex
		return l.scanHexEscape(2)
	case 'u':
		// \uHHHH - 4 digit unicode
		return l.scanHexEscape(4)
	case eof:
		// Error if we reach EOF in an escape sequence
		return false
	default:
		// Error for any other escape sequence
		return false
	}
}

// scanHexEscape consumes the given number of hexadecimal digits and reports
// whether all of them were present.
func (l *lexer) scanHexEscape(digits int) bool {
	for range digits {
		r := l.next()
		if !(unicode.IsDigit(r) || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')) {
			return false
		}
	}
	return true
}

// scanTime scans a time literal and reports whether one was found: a date
// (YYYY-MM-DD) optionally followed by 'T' and an RFC 3339 clock time whose
// zone may be omitted.
func (l *lexer) scanTime() bool {
	// Date: YYYY-MM-DD
	if !l.acceptDigits(4) || !l.accept("-") || !l.acceptDigits(2) || !l.accept("-") || !l.acceptDigits(2) {
		return false
	}
	// A date alone is a complete literal. Each further part is taken only
	// when it is complete; otherwise the literal ends before it.
	date := l.mark()
	if !l.accept("T") {
		return true
	}
	// Time: HH:MM:SS
	if !l.acceptDigits(2) || !l.accept(":") || !l.acceptDigits(2) || !l.accept(":") || !l.acceptDigits(2) {
		l.reset(date)
		return true
	}
	// Optional fractional seconds: '.' 1+DIGIT
	clock := l.mark()
	if l.accept(".") {
		if !unicode.IsDigit(l.next()) {
			l.reset(clock)
			return true
		}
		l.acceptRun("0123456789")
	}
	// Optional zone: 'Z' or (+|-)HH:MM
	if l.accept("Z") {
		return true
	}
	zone := l.mark()
	if l.accept("+-") {
		if !l.acceptDigits(2) || !l.accept(":") || !l.acceptDigits(2) {
			l.reset(zone)
		}
	}
	return true
}

// scanDuration scans a duration literal made of number and unit pairs and
// reports whether one was found. It takes the longest match; the remainder
// becomes the next token.
func (l *lexer) scanDuration() bool {
	valid := false
	for {
		start := l.mark()
		if !l.scanDurationNumber() {
			break
		}
		found := false
		switch r := l.next(); r {
		case 'n':
			if l.accept("s") {
				found = true
			}
		case 'u':
			if l.accept("s") {
				found = true
			}
		case 'μ':
			if l.accept("s") {
				found = true
			}
		case 'm':
			l.accept("s")
			found = true
		case 's':
			found = true
		case 'h':
			found = true
		default:
			l.reset(start)
		}
		if !found {
			break
		}
		valid = true
		r := l.peek()
		if r == eof || (!unicode.IsDigit(r) && r != '.') {
			break
		}
	}
	if !valid {
		return false
	}
	return true
}

// scanDurationNumber scans the signed number before a unit in a duration
// literal and reports whether any digits were found.
func (l *lexer) scanDurationNumber() bool {
	start := l.mark()
	l.accept("+-")
	if n := l.acceptRun("0123456789."); n > 0 {
		return true
	}
	l.reset(start)
	return false
}

// scanNumber scans the longest run that can form a number literal.
// It is purely lexical; the parser validates the text with strconv.
// See https://github.com/golang/go/blob/master/src/text/template/parse/lex.go
func (l *lexer) scanNumber() {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789_"
	if l.accept("0") {
		// Note: Leading 0 does not mean octal in floats.
		if l.accept("xX") {
			digits = "0123456789abcdefABCDEF_"
		} else if l.accept("oO") {
			digits = "01234567_"
		} else if l.accept("bB") {
			digits = "01_"
		}
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if len(digits) == 10+1 && l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	if len(digits) == 16+6+1 && l.accept("pP") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
}

// nextToken returns the next token, advancing the state machine until one is
// emitted. After the input is exhausted it keeps returning EOF tokens.
func (l *lexer) nextToken() token {
	for {
		if l.hasNext {
			l.hasNext = false
			return l.token
		}
		if l.state == stateDone {
			return token{
				typ:  tokenEOF,
				pos:  l.pos,
				line: l.line,
				col:  l.col,
			}
		}
		l.state = l.lexStmt()
	}
}

// next consumes and returns the next rune, or eof at the end of input,
// updating the line and column.
func (l *lexer) next() rune {
	l.prev = l.mark()
	if int(l.pos) >= len(l.input) {
		return eof
	}
	if b := l.input[l.pos]; b < utf8.RuneSelf {
		l.pos++
		if b == '\n' {
			l.line++
			l.col = 1
		} else {
			l.col++
		}
		return rune(b)
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	//nolint:gosec // a rune is at most utf8.UTFMax bytes
	l.pos += int32(w)
	l.col += width(r)
	return r
}

// peek returns the next rune without consuming it.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back to the position before the last next.
// At end of input next does not advance, so backup does not move;
// calling it twice is the same as calling it once.
func (l *lexer) backup() {
	l.reset(l.prev)
}

// mark returns the current position for a later reset.
func (l *lexer) mark() mark {
	return mark{
		pos:  l.pos,
		line: l.line,
		col:  l.col,
	}
}

// reset moves the lexer back to a marked position.
func (l *lexer) reset(m mark) {
	l.pos = m.pos
	l.line = m.line
	l.col = m.col
	l.prev = m
}

// emit records a token of the given type spanning the pending input and
// starts the next token after it.
func (l *lexer) emit(typ tokenType) {
	l.token = token{
		typ:  typ,
		v:    l.input[l.startPos:l.pos],
		pos:  l.startPos,
		line: l.startLine,
		col:  l.startCol,
	}
	l.hasNext = true
	l.startPos = l.pos
	l.startLine = l.line
	l.startCol = l.col
}

// ignore discards the pending input without emitting a token.
func (l *lexer) ignore() {
	l.startPos = l.pos
	l.startLine = l.line
	l.startCol = l.col
}

// accept consumes the next byte if it is in valid and reports whether it did.
func (l *lexer) accept(valid string) bool {
	if int(l.pos) < len(l.input) && strings.IndexByte(valid, l.input[l.pos]) >= 0 {
		l.prev = l.mark()
		l.pos++
		l.col++
		return true
	}
	return false
}

// acceptRun consumes a run of bytes from valid and returns how many it consumed.
func (l *lexer) acceptRun(valid string) int {
	n := 0
	for int(l.pos) < len(l.input) && strings.IndexByte(valid, l.input[l.pos]) >= 0 {
		l.pos++
		l.col++
		n++
	}
	return n
}

// acceptDigits consumes exactly n decimal digits and reports whether all of
// them were present.
func (l *lexer) acceptDigits(n int) bool {
	for range n {
		if int(l.pos) >= len(l.input) || l.input[l.pos] < '0' || l.input[l.pos] > '9' {
			return false
		}
		l.pos++
		l.col++
	}
	return true
}

// errorf emits an error token and terminates the scan.
// The token is positioned where the error was detected, not where the
// current token started.
func (l *lexer) errorf(format string, args ...any) state {
	l.token = token{
		typ:  tokenError,
		v:    fmt.Sprintf(format, args...),
		pos:  l.pos,
		line: l.line,
		col:  l.col,
	}
	l.hasNext = true
	return stateDone
}

// width returns the display width of a non-ASCII rune used for column tracking.
func width(r rune) int32 {
	//nolint:gosec // display width is 1 or 2
	return int32(max(runewidth.RuneWidth(r), 1))
}

// isNumberStart reports whether the rune can begin a number, duration, or time literal.
func isNumberStart(r rune) bool {
	return unicode.IsDigit(r) || r == '.' || r == '+' || r == '-'
}

// isSpace reports whether the rune is a space, tab, carriage return, or newline.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether the rune is a letter, a digit, or an underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isBoolLiteral reports whether s spells true or false in lower, title, or upper case.
func isBoolLiteral(s string) bool {
	switch s {
	case "false", "False", "FALSE", "true", "True", "TRUE":
		return true
	default:
		return false
	}
}
