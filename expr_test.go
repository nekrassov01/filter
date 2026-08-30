package filter

import (
	"strings"
	"testing"
)

func Test_evalNode(t *testing.T) {
	tests := []struct {
		name     string
		nodes    []node
		expected string
	}{
		{
			name: "invalid logical operator",
			nodes: []node{
				{
					typ: nodeBinary,
					op: token{
						typ: tokenEQ,
					},
				},
			},
			expected: "invalid logical operator",
		},
		{
			name: "invalid node type",
			nodes: []node{
				{
					typ: nodeType(255),
				},
			},
			expected: "invalid node type",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := evalNode(test.nodes, 0, testObject, nil)
			if err == nil || !strings.Contains(err.Error(), test.expected) {
				t.Errorf(testTemplate, test.nodes, test.expected, err)
			}
		})
	}
}
