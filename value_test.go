package filter

import (
	"math"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Value
	}{
		{
			name:  "text",
			input: "Knight",
			expected: Value{
				kind: kindString,
				s:    "Knight",
			},
		},
		{
			name:  "empty",
			input: "",
			expected: Value{
				kind: kindString,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := String(test.input); actual != test.expected {
				t.Errorf(testTemplate, test.input, test.expected, actual)
			}
		})
	}
}

func TestNumber(t *testing.T) {
	tests := []struct {
		name  string
		input float64
	}{
		{
			name:  "integer",
			input: 42,
		},
		{
			name:  "fraction",
			input: 3.14,
		},
		{
			name:  "negative",
			input: -1.5,
		},
		{
			name:  "zero",
			input: 0,
		},
		{
			name:  "largest",
			input: math.MaxFloat64,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Number(test.input)
			if actual.kind != kindNumber {
				t.Errorf(testTemplate, test.input, kindNumber, actual.kind)
			}
			//nolint:gosec // bit pattern conversion
			if f := math.Float64frombits(uint64(actual.a)); f != test.input {
				t.Errorf(testTemplate, test.input, test.input, f)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Duration
		expected Value
	}{
		{
			name:  "positive",
			input: 1500 * time.Millisecond,
			expected: Value{
				kind: kindDuration,
				a:    int64(1500 * time.Millisecond),
			},
		},
		{
			name:  "negative",
			input: -time.Hour,
			expected: Value{
				kind: kindDuration,
				a:    int64(-time.Hour),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := Duration(test.input); actual != test.expected {
				t.Errorf(testTemplate, test.input, test.expected, actual)
			}
		})
	}
}

func TestTime(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected Value
	}{
		{
			name:  "utc",
			input: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: Value{
				kind: kindTime,
				a:    1735689600,
			},
		},
		{
			name:  "nanoseconds",
			input: time.Date(2025, 1, 1, 0, 0, 0, 123456789, time.UTC),
			expected: Value{
				kind: kindTime,
				a:    1735689600,
				b:    123456789,
			},
		},
		{
			name:  "same instant in another location",
			input: time.Date(2025, 1, 1, 9, 0, 0, 0, time.FixedZone("JST", 9*3600)),
			expected: Value{
				kind: kindTime,
				a:    1735689600,
			},
		},
		{
			name:  "before the epoch",
			input: time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC),
			expected: Value{
				kind: kindTime,
				a:    -1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := Time(test.input); actual != test.expected {
				t.Errorf(testTemplate, test.input, test.expected, actual)
			}
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected Value
	}{
		{
			name:  "true",
			input: true,
			expected: Value{
				kind: kindString,
				s:    "true",
			},
		},
		{
			name:  "false",
			input: false,
			expected: Value{
				kind: kindString,
				s:    "false",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := Bool(test.input); actual != test.expected {
				t.Errorf(testTemplate, test.input, test.expected, actual)
			}
		})
	}
}

func TestValueOf(t *testing.T) {
	when := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		input    any
		expected Value
	}{
		{
			name:     "string",
			input:    "Knight",
			expected: String("Knight"),
		},
		{
			name:     "int",
			input:    int(-7),
			expected: Number(-7),
		},
		{
			name:     "int8",
			input:    int8(5),
			expected: Number(5),
		},
		{
			name:     "int16",
			input:    int16(5),
			expected: Number(5),
		},
		{
			name:     "int32",
			input:    int32(5),
			expected: Number(5),
		},
		{
			name:     "int64",
			input:    int64(5),
			expected: Number(5),
		},
		{
			name:     "uint",
			input:    uint(5),
			expected: Number(5),
		},
		{
			name:     "uint8",
			input:    uint8(5),
			expected: Number(5),
		},
		{
			name:     "uint16",
			input:    uint16(5),
			expected: Number(5),
		},
		{
			name:     "uint32",
			input:    uint32(5),
			expected: Number(5),
		},
		{
			name:     "uint64",
			input:    uint64(5),
			expected: Number(5),
		},
		{
			name:     "float32",
			input:    float32(2.5),
			expected: Number(2.5),
		},
		{
			name:     "float64",
			input:    3.14,
			expected: Number(3.14),
		},
		{
			name:     "time",
			input:    when,
			expected: Time(when),
		},
		{
			name:     "duration",
			input:    1500 * time.Millisecond,
			expected: Duration(1500 * time.Millisecond),
		},
		{
			name:     "bool",
			input:    true,
			expected: Bool(true),
		},
		{
			name:     "other type is formatted",
			input:    struct{ X int }{X: 1},
			expected: String("{1}"),
		},
		{
			name:     "nil is formatted",
			input:    nil,
			expected: String("<nil>"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := ValueOf(test.input); actual != test.expected {
				t.Errorf(testTemplate, test.input, test.expected, actual)
			}
		})
	}
}
