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
			tr:   nodeNOT,
			want: want{
				val: "not node",
			},
		},
		{
			name: "comparison",
			tr:   nodeComparison,
			want: want{
				val: "comparison node",
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
					typ:  tokenAND,
					v:    "&&",
					pos:  8,
					line: 1,
					col:  9,
				},
				right: 1,
			},
			want: want{
				val: node{
					typ:  nodeBinary,
					left: 0,
					op: token{
						typ:  tokenAND,
						v:    "&&",
						pos:  8,
						line: 1,
						col:  9,
					},
					right: 1,
				},
			},
		},
		{
			name: "or with later children",
			args: args{
				left: 3,
				op: token{
					typ: tokenOR,
					v:   "||",
				},
				right: 7,
			},
			want: want{
				val: node{
					typ:  nodeBinary,
					left: 3,
					op: token{
						typ: tokenOR,
						v:   "||",
					},
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

func Test_newNodeNOT(t *testing.T) {
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
					typ:  tokenNOT,
					v:    "!",
					pos:  4,
					line: 1,
					col:  5,
				},
			},
			want: want{
				val: node{
					typ:  nodeNOT,
					left: 2,
					op: token{
						typ:  tokenNOT,
						v:    "!",
						pos:  4,
						line: 1,
						col:  5,
					},
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
					typ: nodeNOT,
					op: token{
						typ: tokenNOT,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newNodeNOT(test.args.child, test.args.op)
			if !reflect.DeepEqual(got, test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}

func Test_newNodeComparison(t *testing.T) {
	type args struct {
		ident token
		op    token
		val   token
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
			name: "string comparison",
			args: args{
				ident: token{
					typ:  tokenIdent,
					v:    "Name",
					line: 1,
					col:  1,
				},
				op: token{
					typ:  tokenEQ,
					v:    "==",
					pos:  4,
					line: 1,
					col:  5,
				},
				val: token{
					typ:  tokenString,
					v:    "a",
					pos:  6,
					line: 1,
					col:  7,
				},
			},
			want: want{
				val: node{
					typ: nodeComparison,
					ident: token{
						typ:  tokenIdent,
						v:    "Name",
						line: 1,
						col:  1,
					},
					op: token{
						typ:  tokenEQ,
						v:    "==",
						pos:  4,
						line: 1,
						col:  5,
					},
					val: token{
						typ:  tokenString,
						v:    "a",
						pos:  6,
						line: 1,
						col:  7,
					},
				},
			},
		},
		{
			name: "identifier index kept",
			args: args{
				ident: token{
					typ: tokenIdent,
					v:   "HP",
					idx: 5,
				},
				op: token{
					typ: tokenGT,
					v:   ">",
				},
				val: token{
					typ: tokenNumber,
					v:   "1",
				},
			},
			want: want{
				val: node{
					typ: nodeComparison,
					ident: token{
						typ: tokenIdent,
						v:   "HP",
						idx: 5,
					},
					op: token{
						typ: tokenGT,
						v:   ">",
					},
					val: token{
						typ: tokenNumber,
						v:   "1",
					},
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
					typ: nodeComparison,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newNodeComparison(test.args.ident, test.args.op, test.args.val)
			if !reflect.DeepEqual(got, test.want.val) {
				t.Errorf("value mismatch\ngot=%v\nwant=%v\n", got, test.want.val)
			}
		})
	}
}
