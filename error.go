package filter

import (
	"strings"
)

// ErrorKind represents the kind of error.
type ErrorKind int

const (
	KindEval  ErrorKind = iota // evaluation error kind
	KindParse                  // parsing error kind
	KindLex                    // lexical error kind
)

// FilterError represents an error in the filter processing.
type FilterError struct {
	Kind ErrorKind
	Err  error
}

// Error returns the error message.
func (e *FilterError) Error() string {
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
func (e *FilterError) Unwrap() error {
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
