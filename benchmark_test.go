package filter

import (
	"strings"
	"testing"
)

var simpleInput = `String == "HelloWorld"`

var largeInput = `String == "HelloWorld"
&&
StringNumber =~ '^[0-9]+$'
&&
Int > 40
&&
(
	Int8 < 10
	&&
	Int16 <= 5
	&&
	Int32 != 0
)
&&
(
	Float32 >= 2.5
	||
	!
	(
		Float64 < 3.0
	)
)
&&
(
	(
		Time <= 2023-01-01T00:00:00Z
	)
	||
	(
		Duration < 2s30ms100μs1000ns
	)
	||
	(
		Bool == TRUE
	)
)
`

func BenchmarkParseSimple(b *testing.B) {
	for b.Loop() {
		if _, err := Parse(simpleInput); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvalSimple(b *testing.B) {
	expr, err := Parse(simpleInput)
	if err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		if ok, err := expr.Eval(&testObject); !ok || err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseHeavy(b *testing.B) {
	for b.Loop() {
		if _, err := Parse(largeInput); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvalHeavy(b *testing.B) {
	expr, err := Parse(largeInput)
	if err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		if ok, err := expr.Eval(&testObject); !ok || err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseRepeated(b *testing.B) {
	input := repeatInput(largeInput, 30)
	for b.Loop() {
		if _, err := Parse(input); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvalRepeated(b *testing.B) {
	input := repeatInput(largeInput, 30)
	expr, err := Parse(input)
	if err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		if ok, err := expr.Eval(&testObject); !ok || err != nil {
			b.Fatal(err)
		}
	}
}

func repeatInput(input string, n int) string {
	if n <= 0 {
		return input
	}
	var sb strings.Builder
	sb.Grow(len(input) + n*(len(input)+2))
	sb.WriteString(input)
	for range n {
		sb.WriteString("&&")
		sb.WriteString(input)
	}
	return sb.String()
}
