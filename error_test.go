package filter

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Kind ErrorKind
		Line int
		Col  int
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "eval error with position",
			fields: fields{
				Kind: KindEval,
				Line: 1,
				Col:  15,
				Err:  errors.New(`unknown identifier "Unknown"`),
			},
			want: `eval error at 1:15: unknown identifier "Unknown"`,
		},
		{
			name: "parse error with position",
			fields: fields{
				Kind: KindParse,
				Line: 2,
				Col:  4,
				Err:  errors.New(`invalid number "-"`),
			},
			want: `parse error at 2:4: invalid number "-"`,
		},
		{
			name: "lex error with position",
			fields: fields{
				Kind: KindLex,
				Line: 1,
				Col:  2,
				Err:  errors.New("unexpected end of input after '='"),
			},
			want: "token error at 1:2: unexpected end of input after '='",
		},
		{
			name: "eval error",
			fields: fields{
				Kind: KindEval,
				Err:  errors.New("some eval error"),
			},
			want: "eval error: some eval error",
		},
		{
			name: "parse error",
			fields: fields{
				Kind: KindParse,
				Err:  errors.New("some parse error"),
			},
			want: "parse error: some parse error",
		},
		{
			name: "lex error",
			fields: fields{
				Kind: KindLex,
				Err:  errors.New("some lex error"),
			},
			want: "token error: some lex error",
		},
		{
			name: "unknown error",
			fields: fields{
				Kind: ErrorKind(999),
				Err:  errors.New("some unknown error"),
			},
			want: "unknown error: some unknown error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Kind: tt.fields.Kind,
				Line: tt.fields.Line,
				Col:  tt.fields.Col,
				Err:  tt.fields.Err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	type fields struct {
		Kind ErrorKind
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		err    error
	}{
		{
			name: "unwrap error",
			fields: fields{
				Kind: KindEval,
				Err:  errors.New("some eval error"),
			},
			err: errors.New("some eval error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Kind: tt.fields.Kind,
				Err:  tt.fields.Err,
			}
			if err := e.Unwrap(); err.Error() != tt.err.Error() {
				t.Errorf("Error.Unwrap() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func Test_message(t *testing.T) {
	type args struct {
		prefix string
		msg    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal message",
			args: args{
				prefix: "error occurred",
				msg:    "invalid input",
			},
			want: "error occurred: invalid input",
		},
		{
			name: "empty message",
			args: args{
				prefix: "error occurred",
				msg:    "",
			},
			want: "error occurred",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := message(tt.args.prefix, tt.args.msg); got != tt.want {
				t.Errorf("message() = %v, want %v", got, tt.want)
			}
		})
	}
}
