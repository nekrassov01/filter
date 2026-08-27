package filter

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

// eof defines the end of input.
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

const (
	tokenError     tokenType = iota // error
	tokenEOF                        // end of file
	tokenIdent                      // identifier
	tokenGT                         // greater than
	tokenGTE                        // greater than or equal to
	tokenLT                         // less than
	tokenLTE                        // less than or equal to
	tokenEQ                         // equal to
	tokenEQI                        // equal to (case insensitive)
	tokenNEQ                        // not equal to
	tokenNEQI                       // not equal to (case insensitive)
	tokenREQ                        // matches regular expression
	tokenREQI                       // matches regular expression (case insensitive)
	tokenNREQ                       // does not match regular expression
	tokenNREQI                      // does not match regular expression (case insensitive)
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

// String returns a string representation of the token type.
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
	case tokenEQI:
		return "\"case-insensitive equal to\" operator"
	case tokenNEQ:
		return "\"not equal to\" operator"
	case tokenNEQI:
		return "\"case-insensitive not equal to\" operator"
	case tokenREQ:
		return "regex matching operator"
	case tokenREQI:
		return "case-insensitive regex matching operator"
	case tokenNREQ:
		return "negative regex matching operator"
	case tokenNREQI:
		return "case-insensitive negative regex matching operator"
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

// literal returns the literal of a operator token.
// If the literal is not unique, an empty string is returned.
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
	case tokenEQI:
		return "==*"
	case tokenNEQ:
		return "!="
	case tokenNEQI:
		return "!=*"
	case tokenREQ:
		return "=~"
	case tokenREQI:
		return "=~*"
	case tokenNREQ:
		return "!~"
	case tokenNREQI:
		return "!~*"
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
	case tokenEQ, tokenEQI, tokenNEQ, tokenNEQI, tokenGT, tokenGTE, tokenLT, tokenLTE, tokenREQ, tokenREQI, tokenNREQ, tokenNREQI:
		return true
	default:
		return false
	}
}

// isRegexOperatorType reports whether the token is a regex operator.
func (t tokenType) isRegexOperatorType() bool {
	switch t {
	case tokenREQ, tokenREQI, tokenNREQ, tokenNREQI:
		return true
	default:
		return false
	}
}

// isCaseInsensitiveRegexOperatorType reports whether the token is a case insensitive regex operator.
func (t tokenType) isCaseInsensitiveRegexOperatorType() bool {
	switch t {
	case tokenREQI, tokenNREQI:
		return true
	default:
		return false
	}
}

// isValueType reports whether the token is a value type.
func (t tokenType) isValueType() bool {
	switch t {
	case tokenString, tokenRawString, tokenNumber, tokenTime, tokenDuration, tokenBool:
		return true
	default:
		return false
	}
}

// isStringType reports whether the token is a string type.
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

const (
	stateDone state = iota // scanning finished; only EOF tokens follow
	stateStmt              // ready to scan the next token
)

// mark is a position in the input that the lexer can return to.
type mark struct {
	pos  int
	line int
	col  int
}

// lexer holds the state of the scanner.
type lexer struct {
	input      string // the string being scanned
	state      state  // current state
	token      token  // last emitted token waiting to be consumed
	hasNext    bool   // flag there is a pending token
	prev       mark   // position before the last next; backup returns here
	parenDepth int    // nesting depth of ( ) exprs
	pos        int    // current position in the input
	startPos   int    // start position of this token
	line       int    // 1+number of newlines seen
	startLine  int    // start line of this token
	col        int    // 1+number of characters since last newline
	startCol   int    // start column of this token
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

// nextToken returns the next token from the input (on-demand state machine advancement).
func (l *lexer) nextToken() token {
	for {
		if l.hasNext {
			l.hasNext = false
			return l.token
		}
		if l.state == stateDone {
			return token{
				typ:  tokenEOF,
				pos:  int32(l.pos),  //nolint:gosec // bounded by MaxInput
				line: int32(l.line), //nolint:gosec // bounded by MaxInput
				col:  int32(l.col),  //nolint:gosec // bounded by MaxInput
			}
		}
		l.state = lexStmt(l)
	}
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	l.prev = l.mark()
	if l.pos >= len(l.input) {
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += w
	if r == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col += width(r)
	}
	return r
}

// peek returns but does not consume the next rune in the input.
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

// mark returns the current position.
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

// emit passes an token back to the parser.
func (l *lexer) emit(typ tokenType) {
	l.token = token{
		typ:  typ,
		v:    l.input[l.startPos:l.pos],
		pos:  int32(l.startPos),  //nolint:gosec // bounded by MaxInput
		line: int32(l.startLine), //nolint:gosec // bounded by MaxInput
		col:  int32(l.startCol),  //nolint:gosec // bounded by MaxInput
	}
	l.hasNext = true
	l.startPos = l.pos
	l.startLine = l.line
	l.startCol = l.col
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.startPos = l.pos
	l.startLine = l.line
	l.startCol = l.col
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) int {
	n := 0
	for strings.ContainsRune(valid, l.next()) {
		n++
	}
	l.backup()
	return n
}

func (l *lexer) acceptDigits(n int) bool {
	for range n {
		if !unicode.IsDigit(l.next()) {
			return false
		}
	}
	return true
}

// errorf emits an error token and terminates the scan by returning stateDone.
func (l *lexer) errorf(format string, args ...any) state {
	l.token = token{
		typ:  tokenError,
		v:    fmt.Sprintf(format, args...),
		pos:  int32(l.startPos),  //nolint:gosec // bounded by MaxInput
		line: int32(l.startLine), //nolint:gosec // bounded by MaxInput
		col:  int32(l.startCol),  //nolint:gosec // bounded by MaxInput
	}
	l.hasNext = true
	return stateDone
}

// scanEscape handles escape sequences in strings
// It consumes the escape character and expects a valid escape sequence.
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

// scanHexEscape handles hexadecimal escape sequences
// It consumes the specified number of hex digits.
func (l *lexer) scanHexEscape(digits int) bool {
	for range digits {
		r := l.next()
		if !(unicode.IsDigit(r) || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')) {
			return false
		}
	}
	return true
}

// scanTime scans a time literal.
func (l *lexer) scanTime() bool {
	// Date: YYYY-MM-DD
	if !l.acceptDigits(4) || !l.accept("-") || !l.acceptDigits(2) || !l.accept("-") || !l.acceptDigits(2) {
		return false
	}
	// 'T' separator
	if !l.accept("T") {
		return false
	}
	// Time: HH:MM:SS
	if !l.acceptDigits(2) || !l.accept(":") || !l.acceptDigits(2) || !l.accept(":") || !l.acceptDigits(2) {
		return false
	}
	// Optional fractional seconds: '.' 1+DIGIT
	if l.accept(".") {
		r := l.next()
		if !unicode.IsDigit(r) {
			return false
		}
		l.acceptRun("0123456789")
	}
	// Optional timezone: 'Z'/'z' or (+|-)HH:MM
	if l.accept("Zz") {
		return true
	}
	if l.accept("+-") {
		if !l.acceptDigits(2) || !l.accept(":") || !l.acceptDigits(2) {
			return false
		}
		return true
	}
	// No timezone provided (allowed by our extension)
	return true
}

// scanDuration scans for duration literals.
// Determines validity by the longest match,
// the remainder is treated as the next token.
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

// scanDurationNumber scans a number in a duration literal.
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

// lexStmt is the top-level state for lexing.
func lexStmt(l *lexer) state {
	switch r := l.next(); {
	case r == eof:
		return lexEOF(l)
	case isSpace(r):
		return lexSpace(l)
	case r == '"':
		return lexDoubleQuotedString(l)
	case r == '\'':
		return lexSingleQuotedString(l)
	case r == '`':
		return lexRawString(l)
	case r == '(':
		return lexLparen(l)
	case r == ')':
		return lexRparen(l)
	case r == '=':
		return lexEQ(l)
	case r == '!':
		return lexNOT(l)
	case r == '<':
		return lexLT(l)
	case r == '>':
		return lexGT(l)
	case r == '&':
		return lexAND(l)
	case r == '|':
		return lexOR(l)
	case unicode.IsDigit(r) || r == '.' || r == '+' || r == '-':
		return lexNumber(l)
	case unicode.IsLetter(r) || r == '_':
		return lexKeywordOrIdent(l)
	default:
		return l.errorf("unexpected character %#U at %d:%d", r, l.line, l.col-width(r))
	}
}

// lexEOF checks for the end of input and emits an EOF token.
// Called when input is completely consumed.
func lexEOF(l *lexer) state {
	if l.parenDepth < 0 {
		return l.errorf("unexpected right parenthesis at %d:%d", l.line, l.col)
	}
	if l.parenDepth > 0 {
		return l.errorf("unclosed left parenthesis at %d:%d", l.line, l.col)
	}
	l.emit(tokenEOF)
	return stateDone
}

// lexSpace scans a run of space characters.
// One space has already been seen.
func lexSpace(l *lexer) state {
	for isSpace(l.peek()) {
		l.next()
	}
	l.ignore()
	return stateStmt
}

// lexDoubleQuotedString scans a double-quoted string.
// One double quote has already been seen.
func lexDoubleQuotedString(l *lexer) state {
	return lexString(l, '"')
}

// lexSingleQuotedString scans a single-quoted string.
// One single quote has already been seen.
func lexSingleQuotedString(l *lexer) state {
	return lexString(l, '\'')
}

// lexString scans a quoted string, handling escape sequences.
// It consumes the opening quote and expects a matching closing quote.
func lexString(l *lexer, quote rune) state {
Loop:
	for {
		switch l.next() {
		case utf8.RuneError:
			return l.errorf("invalid utf8 encoding in string at %d:%d", l.line, l.col)
		case eof, '\n':
			return l.errorf("unterminated quoted string at %d:%d", l.line, l.col)
		case '\\':
			if !l.scanEscape() {
				return l.errorf("invalid escape sequence in string at %d:%d", l.line, l.col)
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
func lexRawString(l *lexer) state {
Loop:
	for {
		switch l.next() {
		case utf8.RuneError:
			return l.errorf("invalid utf8 encoding in raw string at %d:%d", l.line, l.col)
		case eof:
			return l.errorf("unterminated raw string at %d:%d", l.line, l.col)
		case '`':
			break Loop
		}
	}
	l.emit(tokenRawString)
	return stateStmt
}

// lexLparen emits a left parenthesis.
func lexLparen(l *lexer) state {
	l.emit(tokenLparen)
	l.parenDepth++
	return stateStmt
}

// lexRparen emits a right parenthesis.
func lexRparen(l *lexer) state {
	l.emit(tokenRparen)
	l.parenDepth--
	return stateStmt
}

// lexEQ scans for operators starting with an equality sign.
// The leading '=' has already been seen.
func lexEQ(l *lexer) state {
	switch l.peek() {
	case '=':
		l.next()
		if r := l.peek(); r == '*' {
			l.next()
			l.emit(tokenEQI)
		} else {
			l.emit(tokenEQ)
		}
	case '~':
		l.next()
		if r := l.peek(); r == '*' {
			l.next()
			l.emit(tokenREQI)
		} else {
			l.emit(tokenREQ)
		}
	default:
		return l.errorf("unexpected character %q after '=' at %d:%d", l.peek(), l.line, l.col)
	}
	return stateStmt
}

// lexNOT scans for operators starting with a negative sign.
// The leading '!' has already been seen.
// If unary, it emits a negative operator.
func lexNOT(l *lexer) state {
	switch l.peek() {
	case '=':
		l.next()
		if r := l.peek(); r == '*' {
			l.next()
			l.emit(tokenNEQI)
		} else {
			l.emit(tokenNEQ)
		}
	case '~':
		l.next()
		if r := l.peek(); r == '*' {
			l.next()
			l.emit(tokenNREQI)
		} else {
			l.emit(tokenNREQ)
		}
	default:
		l.emit(tokenNOT)
	}
	return stateStmt
}

// lexLT scans for less than operators.
// The leading '<' has already been seen.
func lexLT(l *lexer) state {
	if l.peek() == '=' {
		l.next()
		l.emit(tokenLTE)
	} else {
		l.emit(tokenLT)
	}
	return stateStmt
}

// lexGT scans for greater than operators.
// The leading '>' has already been seen.
func lexGT(l *lexer) state {
	if l.peek() == '=' {
		l.next()
		l.emit(tokenGTE)
	} else {
		l.emit(tokenGT)
	}
	return stateStmt
}

// lexAND scans for the logical AND operator.
// The leading '&' has already been seen.
func lexAND(l *lexer) state {
	r := l.peek()
	if r == '&' {
		l.next()
		l.emit(tokenAND)
	} else {
		return l.errorf("unexpected character %q after '&' at %d:%d", r, l.line, l.col)
	}
	return stateStmt
}

// lexOR scans for the logical OR operator.
// The leading '|' has already been seen.
func lexOR(l *lexer) state {
	r := l.peek()
	if r == '|' {
		l.next()
		l.emit(tokenOR)
	} else {
		return l.errorf("unexpected character %q after '|' at %d:%d", r, l.line, l.col)
	}
	return stateStmt
}

// lexNumber scans for numbers, duration, and time literals.
// The leading digit or sign has already been seen.
func lexNumber(l *lexer) state {
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

// lexKeywordOrIdent scans for keywords or identifiers.
// The leading character has already been seen.
func lexKeywordOrIdent(l *lexer) state {
	for {
		r := l.next()
		if !isAlphaNumeric(r) && r != '_' {
			l.backup()
			break
		}
	}
	if isBoolLiteral(l.input[l.startPos:l.pos]) {
		l.emit(tokenBool)
		return stateStmt
	}
	l.emit(tokenIdent)
	return stateStmt
}

// width returns the display width of the rune used for column tracking.
// ASCII is resolved without consulting the runewidth tables.
func width(r rune) int {
	if r < utf8.RuneSelf {
		return 1
	}
	return max(runewidth.RuneWidth(r), 1)
}

// isSpace reports whether the rune is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether the rune is a valid alphanumeric character.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isBoolLiteral checks if the string is a boolean literal.
func isBoolLiteral(s string) bool {
	switch s {
	case "false", "False", "FALSE", "true", "True", "TRUE":
		return true
	default:
		return false
	}
}
