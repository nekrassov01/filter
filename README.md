<p align="center">
  <h2 align="center">FILTER</h2>
  <p align="center">The minimum filter expression for Go</p>
  <p align="center">
    <a href="https://github.com/nekrassov01/filter/actions/workflows/test.yml"><img src="https://github.com/nekrassov01/filter/actions/workflows/test.yml/badge.svg?branch=main" alt="CI" /></a>
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

* Comparisons, regex, logical AND / OR / NOT
* Supported types: string, all integer types, float32/64, time.Duration, bool
* Case-insensitive equality: `==*` / `!=*`
* Regex: `=~` / `!~`, case-insensitive: `=~*` / `!~*`
* Duration literals: `1500ms`, `2s`, `1h30m`, `4000μs`

## Performance

`filter` intentionally does a small amount of work once, so that evaluating an expression many times stays flat:

* Regex literals: compiled exactly once per distinct pattern (process-wide sync cache). Writing the same "foo.*" pattern many times does not multiply compile cost.
* Numeric & duration RHS literals: parsed eagerly during parsing (including quoted forms like `"42"` or `"1500ms"`); eval just compares pre‑parsed values.
* Field value reuse: per evaluation a tiny map caches each identifier the first time it is requested; referencing the same field dozens of times (common in generated filters) does not add proportional `GetField` overhead.

Net effect: expressions with high token repetition scale sub-linearly in both time and allocations compared to naïve re-parsing / re-compiling approaches.

## Benchmarks

`filter` is designed to be memory efficient. See [benchmark_test.go](./benchmark_test.go)

Simple input: `String == "HelloWorld"`

```bash
$ go test -bench=Simple$ -benchmem -count 5 -benchtime=10000x
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/filter
cpu: Apple M2
BenchmarkParseSimple-8             10000              1147 ns/op            4625 B/op          6 allocs/op
BenchmarkParseSimple-8             10000              1016 ns/op            4624 B/op          6 allocs/op
BenchmarkParseSimple-8             10000               918.7 ns/op          4624 B/op          6 allocs/op
BenchmarkParseSimple-8             10000               814.8 ns/op          4624 B/op          6 allocs/op
BenchmarkParseSimple-8             10000               761.7 ns/op          4624 B/op          6 allocs/op
BenchmarkEvalSimple-8              10000                45.50 ns/op            0 B/op          0 allocs/op
BenchmarkEvalSimple-8              10000                44.95 ns/op            0 B/op          0 allocs/op
BenchmarkEvalSimple-8              10000                45.02 ns/op            0 B/op          0 allocs/op
BenchmarkEvalSimple-8              10000                45.15 ns/op            0 B/op          0 allocs/op
BenchmarkEvalSimple-8              10000                45.10 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/nekrassov01/filter   0.277s
```

Even when given complex input, performance does not drop drastically.

```bash
$ go test -bench=Heavy$ -benchmem -count 5 -benchtime=10000x
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/filter
cpu: Apple M2
BenchmarkParseHeavy-8              10000              8116 ns/op           13272 B/op         10 allocs/op
BenchmarkParseHeavy-8              10000              7374 ns/op           13272 B/op         10 allocs/op
BenchmarkParseHeavy-8              10000              7905 ns/op           13272 B/op         10 allocs/op
BenchmarkParseHeavy-8              10000              7130 ns/op           13272 B/op         10 allocs/op
BenchmarkParseHeavy-8              10000              7134 ns/op           13272 B/op         10 allocs/op
BenchmarkEvalHeavy-8               10000               528.6 ns/op           616 B/op          3 allocs/op
BenchmarkEvalHeavy-8               10000               535.5 ns/op           616 B/op          3 allocs/op
BenchmarkEvalHeavy-8               10000               567.8 ns/op           616 B/op          3 allocs/op
BenchmarkEvalHeavy-8               10000               533.7 ns/op           616 B/op          3 allocs/op
BenchmarkEvalHeavy-8               10000               519.4 ns/op           616 B/op          3 allocs/op
PASS
ok      github.com/nekrassov01/filter   0.634s
```

Stable even when heavy input is concatenated 30 times with `&&`.

```bash
$ go test -bench=Repeated$ -benchmem -count 5 -benchtime=10000x
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/filter
cpu: Apple M2
BenchmarkParseRepeated-8           10000            220931 ns/op          472029 B/op         15 allocs/op
BenchmarkParseRepeated-8           10000            230419 ns/op          472028 B/op         15 allocs/op
BenchmarkParseRepeated-8           10000            237931 ns/op          472028 B/op         15 allocs/op
BenchmarkParseRepeated-8           10000            219241 ns/op          472028 B/op         15 allocs/op
BenchmarkParseRepeated-8           10000            219584 ns/op          472028 B/op         15 allocs/op
BenchmarkEvalRepeated-8            10000             12536 ns/op             616 B/op          3 allocs/op
BenchmarkEvalRepeated-8            10000             12356 ns/op             616 B/op          3 allocs/op
BenchmarkEvalRepeated-8            10000             12398 ns/op             616 B/op          3 allocs/op
BenchmarkEvalRepeated-8            10000             12357 ns/op             616 B/op          3 allocs/op
BenchmarkEvalRepeated-8            10000             12334 ns/op             616 B/op          3 allocs/op
PASS
ok      github.com/nekrassov01/filter   12.127s
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
	case "Retries":
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
| Duration | `1500ms`, `2s`, `1h30m`, `4000μs`      | Go `time.ParseDuration` compatible |
| Time     | `2023-01-01T00:00:00Z`                 | Go `time.RFC3339` compatible       |
| Boolean  | `true`, `false`, `True`, `FALSE`       | Case-insensitive variants accepted |

### Operators

| Category                  | Operators         | Description                                                          |
| ------------------------- | ----------------- | -------------------------------------------------------------------- |
| Comparison                | `> >= < <= == !=` | Numbers / durations (`==` / `!=` also for strings)                   |
| Case-insensitive (string) | `==* !=*`         | Unicode case folding                                                 |
| Regex                     | `=~ !~ =~* !~*`   | RE2 (Go regex), cached per pattern string; `*` adds case-insensitive |
| Logical                   | `&&` `\|\|` `!`   | Short-circuit                                                        |

## Author

[nekrassov01](https://github.com/nekrassov01)

## License

[MIT](https://github.com/nekrassov01/filter/blob/main/LICENSE)
