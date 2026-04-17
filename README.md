<p align="center">
  <h2 align="center">FILTER</h2>
  <p align="center">The minimal filter expressions for Go</p>
  <p align="center">
    <a href="https://github.com/nekrassov01/filter/actions/workflows/ci.yml"><img src="https://github.com/nekrassov01/filter/actions/workflows/ci.yml/badge.svg?branch=main" alt="CI" /></a>
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
$ go test -bench Simple$ -benchmem -count 5 -benchtime 10000x ./benchmarks/
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/filter/benchmarks
cpu: Apple M2
BenchmarkParseSimple-8             10000              1086 ns/op            4832 B/op          5 allocs/op
BenchmarkParseSimple-8             10000               801.3 ns/op          4832 B/op          5 allocs/op
BenchmarkParseSimple-8             10000               874.7 ns/op          4832 B/op          5 allocs/op
BenchmarkParseSimple-8             10000               747.2 ns/op          4832 B/op          5 allocs/op
BenchmarkParseSimple-8             10000               769.5 ns/op          4832 B/op          5 allocs/op
BenchmarkEvalSimple-8              10000                58.45 ns/op           16 B/op          1 allocs/op
BenchmarkEvalSimple-8              10000                62.60 ns/op           16 B/op          1 allocs/op
BenchmarkEvalSimple-8              10000                60.58 ns/op           16 B/op          1 allocs/op
BenchmarkEvalSimple-8              10000                59.67 ns/op           16 B/op          1 allocs/op
BenchmarkEvalSimple-8              10000                60.77 ns/op           16 B/op          1 allocs/op
PASS
ok      github.com/nekrassov01/filter/benchmarks        0.386s
```

### Case 2

Input:

```text
Class == "軍師" && Name =~ '^(諸葛亮|龐統|法正)' && Name != "" && (
    BirthDate < '0190-01-01T00:00:00Z' && ActiveTimeBattleGauge >= '20s'
) && (
    HitPoint > "50" && MagicPoint > 100 && LifePoint != 0
) && (
    Magic >= 20 || !(Speed < 20)
)
```

Result:

```bash
$ go test -bench Heavy$ -benchmem -count 5 -benchtime 10000x ./benchmarks/
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/filter/benchmarks
cpu: Apple M2
BenchmarkParseHeavy-8              10000              8176 ns/op           13481 B/op          9 allocs/op
BenchmarkParseHeavy-8              10000              5449 ns/op           13480 B/op          9 allocs/op
BenchmarkParseHeavy-8              10000              5424 ns/op           13480 B/op          9 allocs/op
BenchmarkParseHeavy-8              10000              5394 ns/op           13480 B/op          9 allocs/op
BenchmarkParseHeavy-8              10000              5403 ns/op           13480 B/op          9 allocs/op
BenchmarkEvalHeavy-8               10000               632.0 ns/op           703 B/op          9 allocs/op
BenchmarkEvalHeavy-8               10000               660.9 ns/op           699 B/op          9 allocs/op
BenchmarkEvalHeavy-8               10000               648.1 ns/op           699 B/op          9 allocs/op
BenchmarkEvalHeavy-8               10000               666.6 ns/op           699 B/op          9 allocs/op
BenchmarkEvalHeavy-8               10000               640.2 ns/op           703 B/op          9 allocs/op
PASS
ok      github.com/nekrassov01/filter/benchmarks        0.667s
```

### Case 3

Input:

Concatenate Case 2 with `&&` 30 times

Result:

```bash
$ go test -bench Repeated$ -benchmem -count 5 -benchtime 10000x ./benchmarks/
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/filter/benchmarks
cpu: Apple M2
BenchmarkParseRepeated-8           10000            161514 ns/op          472234 B/op         14 allocs/op
BenchmarkParseRepeated-8           10000            159965 ns/op          472233 B/op         14 allocs/op
BenchmarkParseRepeated-8           10000            163515 ns/op          472233 B/op         14 allocs/op
BenchmarkParseRepeated-8           10000            161376 ns/op          472233 B/op         14 allocs/op
BenchmarkParseRepeated-8           10000            160894 ns/op          472234 B/op         14 allocs/op
BenchmarkEvalRepeated-8            10000             14868 ns/op             703 B/op          9 allocs/op
BenchmarkEvalRepeated-8            10000             15078 ns/op             703 B/op          9 allocs/op
BenchmarkEvalRepeated-8            10000             15318 ns/op             707 B/op          9 allocs/op
BenchmarkEvalRepeated-8            10000             15306 ns/op             700 B/op          9 allocs/op
BenchmarkEvalRepeated-8            10000             14730 ns/op             703 B/op          9 allocs/op
PASS
ok      github.com/nekrassov01/filter/benchmarks        9.211s
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

// MyTarget represents the example filter target.
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
