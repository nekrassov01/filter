package filter

import (
	"fmt"
	"strings"
)

// ErrorKind identifies the stage at which an error was detected.
type ErrorKind int

// Error kinds, in the order the stages run at evaluation time.
const (
	// KindEval is the evaluation error kind.
	KindEval ErrorKind = iota

	// KindParse is the parsing error kind.
	KindParse

	// KindLex is the lexical error kind.
	KindLex
)

// Error is an error detected while lexing, parsing, or evaluating an expression.
// Line and Col locate the offending token in the input (1-based, Col counted in
// display width); both are 0 when the error has no position, such as empty input.
type Error struct {
	Kind ErrorKind
	Line int
	Col  int
	Err  error
}

// Error returns the message in the form "<kind> at <line>:<col>: <detail>",
// omitting the position when Line is 0.
func (e *Error) Error() string {
	var prefix string
	switch e.Kind {
	case KindEval:
		prefix = "eval error"
	case KindParse:
		prefix = "parse error"
	case KindLex:
		prefix = "token error"
	default:
		prefix = "unknown error"
	}
	if e.Line > 0 {
		prefix = fmt.Sprintf("%s at %d:%d", prefix, e.Line, e.Col)
	}
	return message(prefix, e.Err.Error())
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.Err
}

// newError returns an Error of the given kind positioned at t.
// The message is formatted as by fmt.Errorf, so %w wraps a cause.
func newError(kind ErrorKind, t token, format string, args ...any) *Error {
	return &Error{
		Kind: kind,
		Line: int(t.line),
		Col:  int(t.col),
		Err:  fmt.Errorf(format, args...),
	}
}

// message joins a prefix and a message with a colon, returning only the
// prefix when the message is empty.
func message(prefix, msg string) string {
	if msg == "" {
		return prefix
	}
	var b strings.Builder
	b.Grow(len(prefix) + 2 + len(msg))
	b.WriteString(prefix)
	b.WriteString(": ")
	b.WriteString(msg)
	return b.String()
}
