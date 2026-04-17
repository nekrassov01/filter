package filter

import (
	"strings"
)

// ErrorKind represents the kind of error.
type ErrorKind int

const (
	// KindEval is the evaluation error kind.
	KindEval ErrorKind = iota

	// KindParse is the parsing error kind.
	KindParse

	// KindLex is the lexical error kind.
	KindLex
)

// Error represents an error in the filter processing.
type Error struct {
	Kind ErrorKind
	Err  error
}

// Error returns the error message.
func (e *Error) Error() string {
	switch e.Kind {
	case KindEval:
		return message("eval error", e.Err.Error())
	case KindParse:
		return message("parse error", e.Err.Error())
	case KindLex:
		return message("token error", e.Err.Error())
	default:
		return message("unknown error", e.Err.Error())
	}
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.Err
}

// message constructs an error message with a prefix and message.
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
