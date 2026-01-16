<p align="center">
  <h2 align="center">FILTER</h2>
  <p align="center">The minimal filter expressions for Go</p>
  <p align="center">
    <a href="https://github.com/nekrassov01/filter/actions/workflows/ci.yaml"><img src="https://github.com/nekrassov01/filter/actions/workflows/ci.yaml/badge.svg?branch=main" alt="CI" /></a>
    <a href="https://codecov.io/gh/nekrassov01/filter"><img src="https://codecov.io/gh/nekrassov01/filter/graph/badge.svg?token=Z75YW69MQK" alt="Codecov" /></a>
    <a href="https://pkg.go.dev/github.com/nekrassov01/filter"><img src="https://pkg.go.dev/badge/github.com/nekrassov01/filter.svg" alt="Go Reference" /></a>
    <a href="https://goreportcard.com/report/github.com/nekrassov01/filter"><img src="https://goreportcard.com/badge/github.com/nekrassov01/filter" alt="Go Report Card" /></a>
    <img src="https://img.shields.io/github/license/nekrassov01/filter" alt="LICENSE" />
    <a href="https://deepwiki.com/nekrassov01/filter"><img src="https://deepwiki.com/badge.svg" alt="Ask DeepWiki" /></a>
  </p>
</p>

## Overview

`filter` focuses on one task: evaluating small boolean filter expressions in Go without the weight of a general expression engine. The motivation is to avoid large, reflection-heavy or feature-rich DSLs when you only need predictable field filtering. Core traits: minimal syntax (comparisons, basic logical operators, regex, case-insensitive equality), no reflection (caller supplies values via a tiny interface), deterministic errors with positions, and cached regex compilation. This keeps the surface area small while remaining fast and explicit.

## Features

- Comparisons, regex, logical AND / OR / NOT
- Supported types: string, all integer types, float32/64, time.Time, time.Duration, bool
- Case-insensitive equality: `==*` / `!=*`
- Regex: `=~` / `!~`, case-insensitive: `=~*` / `!~*`
- Time literals: [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339) only
- Duration literals: `1500ms`, `2s`, `1h30m`, `4000μs`

## Performance

`filter` intentionally does a small amount of work once, so that evaluating an expression many times stays flat:

- Regex literals: compiled exactly once per distinct pattern (process-wide sync cache). Writing the same "foo.*" pattern many times does not multiply compile cost.
- Numeric & duration RHS literals: parsed eagerly during parsing (including quoted forms like `"42"` or `"1500ms"`); eval just compares pre‑parsed values.
- Field value reuse: per evaluation a tiny map caches each identifier the first time it is requested; referencing the same field dozens of times does not add proportional `GetField` overhead.

## Benchmarks

`filter` is designed to be memory efficient. See [benchmark_test.go](./benchmark_test.go)

### Case 1

Input:

```text
String == "HelloWorld"
```

Result:

```bash
$ go test -bench=Simple$ -benchmem -count 5 -benchtime=10000x
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/filter
cpu: Apple M2
BenchmarkParseSimple-8             10000               992.5 ns/op          4593 B/op          4 allocs/op
BenchmarkParseSimple-8             10000               920.1 ns/op          4592 B/op          4 allocs/op
BenchmarkParseSimple-8             10000               863.9 ns/op          4592 B/op          4 allocs/op
BenchmarkParseSimple-8             10000               802.8 ns/op          4592 B/op          4 allocs/op
BenchmarkParseSimple-8             10000               781.4 ns/op          4592 B/op          4 allocs/op
BenchmarkEvalSimple-8              10000                46.83 ns/op            0 B/op          0 allocs/op
BenchmarkEvalSimple-8              10000                46.67 ns/op            0 B/op          0 allocs/op
BenchmarkEvalSimple-8              10000                47.03 ns/op            0 B/op          0 allocs/op
BenchmarkEvalSimple-8              10000                47.13 ns/op            0 B/op          0 allocs/op
BenchmarkEvalSimple-8              10000                46.83 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/nekrassov01/filter   0.263s
```

### Case 2

Input:

```text
(
	String == "HelloWorld" && StringNumber =~ '^[0-9]+$' && Int > 40
) && (
	Int8 < 10 && Int16 <= 5 && Int32 != 0
) && (
	Float32 >= 2.5 || !(Float64 < 3.0)
) && (
	(Time <= 2023-01-01T00:00:00Z) || (Duration < 2s30ms100μs1000ns) || (Bool == TRUE)
)
```

Result:

```bash
$ go test -bench=Heavy$ -benchmem -count 5 -benchtime=10000x
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/filter
cpu: Apple M2
BenchmarkParseHeavy-8              10000              8030 ns/op           13240 B/op          8 allocs/op
BenchmarkParseHeavy-8              10000              6834 ns/op           13240 B/op          8 allocs/op
BenchmarkParseHeavy-8              10000              7145 ns/op           13240 B/op          8 allocs/op
BenchmarkParseHeavy-8              10000              6989 ns/op           13240 B/op          8 allocs/op
BenchmarkParseHeavy-8              10000              6896 ns/op           13240 B/op          8 allocs/op
BenchmarkEvalHeavy-8               10000               572.8 ns/op           616 B/op          3 allocs/op
BenchmarkEvalHeavy-8               10000               621.7 ns/op           616 B/op          3 allocs/op
BenchmarkEvalHeavy-8               10000               645.1 ns/op           616 B/op          3 allocs/op
BenchmarkEvalHeavy-8               10000               574.2 ns/op           616 B/op          3 allocs/op
BenchmarkEvalHeavy-8               10000               604.0 ns/op           616 B/op          3 allocs/op
PASS
ok      github.com/nekrassov01/filter   0.629s
```

### Case 3

Input:

Concatenate Case 2 with `&&` 30 times

Result:

```bash
$ go test -bench=Repeated$ -benchmem -count 5 -benchtime=10000x
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/filter
cpu: Apple M2
BenchmarkParseRepeated-8           10000            209364 ns/op          471998 B/op         13 allocs/op
BenchmarkParseRepeated-8           10000            211333 ns/op          471997 B/op         13 allocs/op
BenchmarkParseRepeated-8           10000            209689 ns/op          471996 B/op         13 allocs/op
BenchmarkParseRepeated-8           10000            209426 ns/op          471996 B/op         13 allocs/op
BenchmarkParseRepeated-8           10000            210669 ns/op          471996 B/op         13 allocs/op
BenchmarkEvalRepeated-8            10000             13089 ns/op             616 B/op          3 allocs/op
BenchmarkEvalRepeated-8            10000             13002 ns/op             616 B/op          3 allocs/op
BenchmarkEvalRepeated-8            10000             12970 ns/op             616 B/op          3 allocs/op
BenchmarkEvalRepeated-8            10000             13090 ns/op             616 B/op          3 allocs/op
BenchmarkEvalRepeated-8            10000             13008 ns/op             616 B/op          3 allocs/op
PASS
ok      github.com/nekrassov01/filter   11.370s
```

## Installation

```sh
go get github.com/nekrassov01/filter@latest
```

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/nekrassov01/filter"
)

// Example target type.
type MyTarget struct {
	Name    string
	Latency time.Duration
	Retries int
	Enabled bool
}

// GetField maps a field name to its value.
func (t *MyTarget) GetField(key string) (any, error) {
	switch key {
	case "Name":
		return t.Name, nil
	case "Latency":
		return t.Latency, nil
	case "Retries", "RetryCount":
		return t.Retries, nil
	case "Enabled":
		return t.Enabled, nil
	default:
		return nil, fmt.Errorf("field not found: %q", key)
	}
}

func main() {
	input := `Name =~ '^foo' && (Latency < 1500ms || Retries != 0) && Enabled == true`

	expr, err := filter.Parse(input)
	if err != nil {
		panic(err)
	}

	target := &MyTarget{
		Name:    "foobar",
		Latency: 100 * time.Millisecond,
		Retries: 3,
		Enabled: true,
	}

	ok, err := expr.Eval(target)
	if err != nil {
		panic(err)
	}
	fmt.Println("matched:", ok)
}
```

## Syntax

### Literals

| Kind     | Examples                               | Notes                              |
| -------- | -------------------------------------- | ---------------------------------- |
| String   | `"Hello"`, `'世界'`, `` `raw\ntext` `` | Double / single / raw (backtick)   |
| Number   | `42`, `3.14`, `0x1.fp3`                | Subset of Go numeric literals      |
| Time     | `2023-01-01T00:00:00Z`                 | Go `time.RFC3339` compatible       |
| Duration | `1500ms`, `2s`, `1h30m`, `4000μs`      | Go `time.ParseDuration` compatible |
| Boolean  | `true`, `false`, `True`, `FALSE`       | Case-insensitive variants accepted |

### Operators

| Category                  | Operators                   | Description                                          |
| ------------------------- | --------------------------- | ---------------------------------------------------- |
| Comparison                | `>` `>=` `<` `<=` `==` `!=` | Strings, integers, times, and durations              |
| Case-insensitive (string) | `==*` `!=*`                 | Unicode case folding                                 |
| Regex                     | `=~` `!~` `=~*` `!~*`       | Cached per pattern string; `*` adds case-insensitive |
| Logical                   | `&&` `\|\|` `!`             | Short-circuit                                        |

## Author

[nekrassov01](https://github.com/nekrassov01)

## License

[MIT](https://github.com/nekrassov01/filter/blob/main/LICENSE)
