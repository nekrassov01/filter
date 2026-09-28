package filter

import (
	"reflect"
	"testing"
)

func Test_nodeType_String(t *testing.T) {
	type want struct {
		val string
	}
	tests := []struct {
		name string
		tr   nodeType
		want want
	}{
		{
			name: "binary",
			tr:   nodeBinary,
			want: want{
				val: "binary node",
			},
		},
		{
			name: "not",
			tr:   nodeUnary,
			want: want{
				val: "unary node",
			},
		},
		{
			name: "predicate",
			tr:   nodePredicate,
			want: want{
				val: "predicate node",
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

func Test_newNodeBinary(t *testing.T) {
	type args struct {
		left  int32
		op    token
		right int32
	}
	type want struct {
		val node
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "and with children",
			args: args{
				left: 0,
				op: token{
					v:    "&&",
					pos:  8,
					line: 1,
					col:  9,
					typ:  tokenAND,
				},
				right: 1,
			},
			want: want{
				val: node{
					op: token{
						v:    "&&",
						pos:  8,
						line: 1,
						col:  9,
						typ:  tokenAND,
					},
					typ:   nodeBinary,
					left:  0,
					right: 1,
				},
			},
		},
		{
			name: "or with later children",
			args: args{
				left: 3,
				op: token{
					v:   "||",
					typ: tokenOR,
				},
				right: 7,
			},
			want: want{
				val: node{
					op: token{
						v:   "||",
						typ: tokenOR,
					},
					typ:   nodeBinary,
					left:  3,
					right: 7,
				},
			},
		},
		{
			name: "zero token",
			args: args{
				left:  0,
				op:    token{},
				right: 0,
			},
			want: want{
				val: node{
					typ: nodeBinary,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newNodeBinary(test.args.left, test.args.op, test.args.right)
			if !reflect.DeepEqual(got, test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_newNodeUnary(t *testing.T) {
	type args struct {
		child int32
		op    token
	}
	type want struct {
		val node
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "not with child",
			args: args{
				child: 2,
				op: token{
					v:    "!",
					pos:  4,
					line: 1,
					col:  5,
					typ:  tokenNOT,
				},
			},
			want: want{
				val: node{
					op: token{
						v:    "!",
						pos:  4,
						line: 1,
						col:  5,
						typ:  tokenNOT,
					},
					typ:  nodeUnary,
					left: 2,
				},
			},
		},
		{
			name: "child zero",
			args: args{
				child: 0,
				op: token{
					typ: tokenNOT,
				},
			},
			want: want{
				val: node{
					op: token{
						typ: tokenNOT,
					},
					typ: nodeUnary,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newNodeUnary(test.args.child, test.args.op)
			if !reflect.DeepEqual(got, test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_newNodePredicate(t *testing.T) {
	type args struct {
		ident    token
		op       token
		val      token
		identIdx int32
	}
	type want struct {
		val node
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "string predicate",
			args: args{
				ident: token{
					v:    "Name",
					line: 1,
					col:  1,
					typ:  tokenIdent,
				},
				op: token{
					v:    "==",
					pos:  4,
					line: 1,
					col:  5,
					typ:  tokenEQ,
				},
				val: token{
					v:    `"a"`,
					pos:  6,
					line: 1,
					col:  7,
					typ:  tokenString,
				},
			},
			want: want{
				val: node{
					ident: token{
						v:    "Name",
						line: 1,
						col:  1,
						typ:  tokenIdent,
					},
					op: token{
						v:    "==",
						pos:  4,
						line: 1,
						col:  5,
						typ:  tokenEQ,
					},
					val: token{
						v:    `"a"`,
						pos:  6,
						line: 1,
						col:  7,
						typ:  tokenString,
					},
					s:   "a",
					typ: nodePredicate,
				},
			},
		},
		{
			name: "identifier index kept",
			args: args{
				ident: token{
					v:   "HP",
					typ: tokenIdent,
				},
				op: token{
					v:   ">",
					typ: tokenGT,
				},
				val: token{
					v:   "1",
					typ: tokenNumber,
				},
				identIdx: 5,
			},
			want: want{
				val: node{
					ident: token{
						v:   "HP",
						typ: tokenIdent,
					},
					op: token{
						v:   ">",
						typ: tokenGT,
					},
					val: token{
						v:   "1",
						typ: tokenNumber,
					},
					s:        "1",
					typ:      nodePredicate,
					identIdx: 5,
				},
			},
		},
		{
			name: "zero tokens leave caches empty",
			args: args{
				ident: token{},
				op:    token{},
				val:   token{},
			},
			want: want{
				val: node{
					typ: nodePredicate,
				},
			},
		},
		{
			name: "single quoted",
			args: args{
				val: token{
					v:   "'text'",
					typ: tokenString,
				},
			},
			want: want{
				val: node{
					val: token{
						v:   "'text'",
						typ: tokenString,
					},
					s:   "text",
					typ: nodePredicate,
				},
			},
		},
		{
			name: "raw string",
			args: args{
				val: token{
					v:   "`text`",
					typ: tokenRawString,
				},
			},
			want: want{
				val: node{
					val: token{
						v:   "`text`",
						typ: tokenRawString,
					},
					s:   "text",
					typ: nodePredicate,
				},
			},
		},
		{
			name: "title case boolean",
			args: args{
				val: token{
					v:   "True",
					typ: tokenBool,
				},
			},
			want: want{
				val: node{
					val: token{
						v:   "True",
						typ: tokenBool,
					},
					s:   "true",
					typ: nodePredicate,
				},
			},
		},
		{
			name: "upper case boolean",
			args: args{
				val: token{
					v:   "FALSE",
					typ: tokenBool,
				},
			},
			want: want{
				val: node{
					val: token{
						v:   "FALSE",
						typ: tokenBool,
					},
					s:   "false",
					typ: nodePredicate,
				},
			},
		},
		{
			name: "empty string",
			args: args{
				val: token{
					v:   "\"\"",
					typ: tokenString,
				},
			},
			want: want{
				val: node{
					val: token{
						v:   "\"\"",
						typ: tokenString,
					},
					s:   "",
					typ: nodePredicate,
				},
			},
		},
		{
			name: "short string token",
			args: args{
				val: token{
					v:   "\"",
					typ: tokenString,
				},
			},
			want: want{
				val: node{
					val: token{
						v:   "\"",
						typ: tokenString,
					},
					s:   "\"",
					typ: nodePredicate,
				},
			},
		},
		{
			name: "escape spelling retained",
			args: args{
				val: token{
					v:   "\"a\\n\"",
					typ: tokenString,
				},
			},
			want: want{
				val: node{
					val: token{
						v:   "\"a\\n\"",
						typ: tokenString,
					},
					s:   "a\\n",
					typ: nodePredicate,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newNodePredicate(test.args.ident, test.args.op, test.args.val, test.args.identIdx)
			if !reflect.DeepEqual(got, test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}
