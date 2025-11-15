package filter

import "testing"

func Test_nodeType_String(t *testing.T) {
	tests := []struct {
		name     string
		typ      nodeType
		expected string
	}{
		{
			name:     "binary",
			typ:      nodeBinary,
			expected: "binary node",
		},
		{
			name:     "not",
			typ:      nodeNOT,
			expected: "not node",
		},
		{
			name:     "comparison",
			typ:      nodeComparison,
			expected: "comparison node",
		},
		{
			name:     "invalid",
			typ:      256,
			expected: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := test.typ.String(); actual != test.expected {
				t.Errorf("expected %v, actual %v", test.expected, actual)
			}
		})
	}
}
