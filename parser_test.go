package filter

import (
	"testing"
	"time"
)

func Test_parseTime(t *testing.T) {
	type expected struct {
		ok   bool
		time time.Time
	}
	tests := []struct {
		name     string
		input    string
		expected expected
	}{
		{
			name:  "rfc3339 utc",
			input: "2025-01-01T00:00:00Z",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "rfc3339 offset",
			input: "2025-01-01T09:00:00+09:00",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "rfc3339 fraction",
			input: "2025-01-01T00:00:00.25Z",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 250000000, time.UTC),
			},
		},
		{
			name:  "datetime without zone",
			input: "2025-01-01T00:00:00",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "datetime with space",
			input: "2025-01-01 00:00:00",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "datetime with space and fraction",
			input: "2025-01-01 00:00:00.5",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 500000000, time.UTC),
			},
		},
		{
			name:  "date only",
			input: "2025-01-01",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "unix seconds",
			input: "1735689600",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "negative unix seconds",
			input: "-1",
			expected: expected{
				ok:   true,
				time: time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name:  "rfc1123",
			input: "Wed, 01 Jan 2025 00:00:00 UTC",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "rfc850",
			input: "Wednesday, 01-Jan-25 00:00:00 UTC",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "rfc822",
			input: "01 Jan 25 00:00 UTC",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "rfc1123 with numeric offset",
			input: "Wed, 01 Jan 2025 09:00:00 +0900",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "rfc822 with numeric offset",
			input: "01 Jan 25 09:00 +0900",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "gmt abbreviation",
			input: "Wed, 01 Jan 2025 00:00:00 GMT",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "unix seconds with underscores",
			input: "1_735_689_600",
			expected: expected{
				ok:   true,
				time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "zone abbreviation other than utc or gmt",
			input: "Wed, 01 Jan 2025 00:00:00 EST",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "lowercase z",
			input: "2025-01-01T00:00:00z",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "out of range month",
			input: "2025-13-01",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "clock only",
			input: "12:00:00",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "unix with fraction",
			input: "1735689600.5",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "unix seconds overflow",
			input: "99999999999999999999",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "sign only",
			input: "-",
			expected: expected{
				ok: false,
			},
		},
		{
			name:  "empty",
			input: "",
			expected: expected{
				ok: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := parseTime(test.input)
			if (err == nil) != test.expected.ok {
				t.Errorf(testTemplate, test.input, test.expected.ok, err)
				return
			}
			if test.expected.ok && !actual.Equal(test.expected.time) {
				t.Errorf(testTemplate, test.input, test.expected.time, actual)
			}
			if test.expected.ok && actual.Location() != time.UTC {
				t.Errorf(testTemplate, test.input, time.UTC, actual.Location())
			}
		})
	}
}

func Test_unquote(t *testing.T) {
	tests := []struct {
		name     string
		token    token
		expected string
	}{
		{
			name: "string",
			token: token{
				typ: tokenString,
				v:   `"abc"`,
			},
			expected: "abc",
		},
		{
			name: "raw string",
			token: token{
				typ: tokenRawString,
				v:   "`abc`",
			},
			expected: "abc",
		},
		{
			name: "empty string",
			token: token{
				typ: tokenString,
				v:   `""`,
			},
			expected: "",
		},
		{
			name: "too short",
			token: token{
				typ: tokenString,
				v:   `"`,
			},
			expected: `"`,
		},
		{
			name: "number",
			token: token{
				typ: tokenNumber,
				v:   "42",
			},
			expected: "42",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := unquote(test.token)
			if actual != test.expected {
				t.Errorf(testTemplate, test.token.v, test.expected, actual)
			}
		})
	}
}

// repr converts ast to a string. String and time literals are quoted;
// numbers, durations, and booleans are printed as written.
