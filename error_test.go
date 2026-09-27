package filter

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Err  error
		Kind ErrorKind
		Line int
		Col  int
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
			name: "eval error with position",
			fields: fields{
				Err:  errors.New(`unknown identifier "Unknown"`),
				Kind: KindEval,
				Line: 1,
				Col:  15,
			},
			want: want{
				val: `eval error at 1:15: unknown identifier "Unknown"`,
			},
		},
		{
			name: "parse error with position",
			fields: fields{
				Err:  errors.New(`invalid number "-"`),
				Kind: KindParse,
				Line: 2,
				Col:  4,
			},
			want: want{
				val: `parse error at 2:4: invalid number "-"`,
			},
		},
		{
			name: "lex error with position",
			fields: fields{
				Err:  errors.New("unexpected end of input after '='"),
				Kind: KindLex,
				Line: 1,
				Col:  2,
			},
			want: want{
				val: "lex error at 1:2: unexpected end of input after '='",
			},
		},
		{
			name: "eval error",
			fields: fields{
				Err:  errors.New("some eval error"),
				Kind: KindEval,
			},
			want: want{
				val: "eval error: some eval error",
			},
		},
		{
			name: "parse error",
			fields: fields{
				Err:  errors.New("some parse error"),
				Kind: KindParse,
			},
			want: want{
				val: "parse error: some parse error",
			},
		},
		{
			name: "lex error",
			fields: fields{
				Err:  errors.New("some lex error"),
				Kind: KindLex,
			},
			want: want{
				val: "lex error: some lex error",
			},
		},
		{
			name: "unknown error",
			fields: fields{
				Err:  errors.New("some unknown error"),
				Kind: ErrorKind(999),
			},
			want: want{
				val: "unknown error: some unknown error",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := &Error{
				Err:  test.fields.Err,
				Kind: test.fields.Kind,
				Line: test.fields.Line,
				Col:  test.fields.Col,
			}
			got := e.Error()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	type fields struct {
		Err  error
		Kind ErrorKind
		Line int
		Col  int
	}
	type want struct {
		val error
	}
	errEval := errors.New("some eval error")
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "unwrap error",
			fields: fields{
				Err:  errEval,
				Kind: KindEval,
			},
			want: want{
				val: errEval,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := &Error{
				Err:  test.fields.Err,
				Kind: test.fields.Kind,
				Line: test.fields.Line,
				Col:  test.fields.Col,
			}
			got := e.Unwrap()
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_newError(t *testing.T) {
	type args struct {
		kind   ErrorKind
		t      token
		format string
		args   []any
	}
	type want struct {
		val   string
		kind  ErrorKind
		line  int
		col   int
		wraps error
	}
	errCause := errors.New("cause")
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "eval error at token",
			args: args{
				kind: KindEval,
				t: token{
					v:    "Unknown",
					pos:  14,
					line: 1,
					col:  15,
					typ:  tokenIdent,
				},
				format: "unknown identifier %q",
				args:   []any{"Unknown"},
			},
			want: want{
				val:  `eval error at 1:15: unknown identifier "Unknown"`,
				kind: KindEval,
				line: 1,
				col:  15,
			},
		},
		{
			name: "parse error on later line",
			args: args{
				kind: KindParse,
				t: token{
					v:    "-",
					line: 2,
					col:  4,
					typ:  tokenNumber,
				},
				format: "invalid number %q",
				args:   []any{"-"},
			},
			want: want{
				val:  `parse error at 2:4: invalid number "-"`,
				kind: KindParse,
				line: 2,
				col:  4,
			},
		},
		{
			name: "lex error without arguments",
			args: args{
				kind: KindLex,
				t: token{
					line: 1,
					col:  6,
					typ:  tokenError,
				},
				format: "unclosed left parenthesis",
			},
			want: want{
				val:  "lex error at 1:6: unclosed left parenthesis",
				kind: KindLex,
				line: 1,
				col:  6,
			},
		},
		{
			name: "zero token has no position",
			args: args{
				kind:   KindParse,
				t:      token{},
				format: "empty input",
			},
			want: want{
				val:  "parse error: empty input",
				kind: KindParse,
			},
		},
		{
			name: "wraps a cause",
			args: args{
				kind: KindParse,
				t: token{
					line: 1,
					col:  7,
				},
				format: "invalid regex %q: %w",
				args:   []any{"[", errCause},
			},
			want: want{
				val:   `parse error at 1:7: invalid regex "[": cause`,
				kind:  KindParse,
				line:  1,
				col:   7,
				wraps: errCause,
			},
		},
		{
			name: "unknown kind",
			args: args{
				kind: ErrorKind(255),
				t: token{
					line: 3,
					col:  1,
				},
				format: "%d %s",
				args:   []any{7, "things"},
			},
			want: want{
				val:  "unknown error at 3:1: 7 things",
				kind: ErrorKind(255),
				line: 3,
				col:  1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newError(test.args.kind, test.args.t, test.args.format, test.args.args...)
			if got.Error() != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got.Error(), test.want.val)
			}
			if got.Kind != test.want.kind {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got.Kind, test.want.kind)
			}
			if got.Line != test.want.line {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got.Line, test.want.line)
			}
			if got.Col != test.want.col {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got.Col, test.want.col)
			}
			if test.want.wraps != nil && !errors.Is(got, test.want.wraps) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got.Unwrap(), test.want.wraps)
			}
		})
	}
}

func Test_message(t *testing.T) {
	type args struct {
		prefix string
		msg    string
	}
	type want struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "normal message",
			args: args{
				prefix: "error occurred",
				msg:    "invalid input",
			},
			want: want{
				val: "error occurred: invalid input",
			},
		},
		{
			name: "empty message",
			args: args{
				prefix: "error occurred",
				msg:    "",
			},
			want: want{
				val: "error occurred",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := message(test.args.prefix, test.args.msg)
			if got != test.want.val {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}
