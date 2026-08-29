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

`filter` focuses on one task: evaluating small boolean filter expressions in Go without the weight of a general expression engine. The motivation is to avoid large, reflection-heavy or feature-rich DSLs when you only need predictable value filtering. Core traits: minimal syntax (comparisons, basic logical operators, regex, case-insensitive equality), no reflection (caller supplies values via a tiny interface), deterministic errors with positions, and cached regex compilation. This keeps the surface area small while remaining fast and explicit.

## Features

- Comparisons, regex, logical AND / OR / NOT
- Values via a one-method `Resolver` interface: `Resolve(name string) (any, bool)`
- Errors are `*filter.Error` with `Kind`, `Line`, and `Col`
- Supported types: string, all integer types, float32/64, time.Time, time.Duration, bool
- Case-insensitive equality: `==*` / `!=*`
- Regex: `=~` / `!~`, case-insensitive: `=~*` / `!~*`
- Time literals: RFC 3339, RFC 1123, RFC 850, RFC 822, `2006-01-02T15:04:05`, `2006-01-02 15:04:05`, `2006-01-02`, Unix seconds; zone-less forms are UTC
- Duration literals: `1500ms`, `2s`, `1h30m`, `4000μs`

## Performance

`filter` intentionally does a small amount of work once, so that evaluating an expression many times stays flat:

- Regex literals: compiled exactly once per distinct pattern (process-wide sync cache). Writing the same "foo.*" pattern many times does not multiply compile cost.
- Number, time, and duration RHS literals: validated and converted once during parsing, so malformed literals are reported as parse errors with their position; eval just compares pre‑parsed values. Quoted forms like `"42"`, `"1500ms"`, or `"2023-01-01 09:00:00"` are converted during parsing too when their text reads as a literal, and otherwise at evaluation time when compared against a numeric, time, or duration value.
- Resolved value reuse: when an identifier appears more than once, each evaluation caches its value on first use in a small stack buffer (a heap slice only beyond 8 distinct identifiers); referencing the same identifier dozens of times does not add proportional `Resolve` overhead. Expressions where every identifier appears once skip the cache entirely.

## Benchmarks

`filter` is designed to be memory efficient. See [benchmark_test.go](./benchmarks/benchmark_test.go)

### Case 1

Input:

```text
Class == "軍師"
```

Result:

```bash
$ go test -bench Simple$ -benchmem -count 5 -benchtime 10000x ./benchmarks/
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/filter/benchmarks
cpu: Apple M2
BenchmarkParseSimple-8             10000               504.9 ns/op           241 B/op          2 allocs/op
BenchmarkParseSimple-8             10000               433.2 ns/op           240 B/op          2 allocs/op
BenchmarkParseSimple-8             10000               393.4 ns/op           240 B/op          2 allocs/op
BenchmarkParseSimple-8             10000               379.4 ns/op           240 B/op          2 allocs/op
BenchmarkParseSimple-8             10000               367.2 ns/op           240 B/op          2 allocs/op
BenchmarkEvalSimple-8              10000                18.46 ns/op           16 B/op          1 allocs/op
BenchmarkEvalSimple-8              10000                18.31 ns/op           16 B/op          1 allocs/op
BenchmarkEvalSimple-8              10000                18.38 ns/op           16 B/op          1 allocs/op
BenchmarkEvalSimple-8              10000                18.78 ns/op           16 B/op          1 allocs/op
BenchmarkEvalSimple-8              10000                18.24 ns/op           16 B/op          1 allocs/op
PASS
ok      github.com/nekrassov01/filter/benchmarks    0.434s
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
BenchmarkParseHeavy-8              10000              4905 ns/op            6832 B/op          3 allocs/op
BenchmarkParseHeavy-8              10000              4177 ns/op            6833 B/op          3 allocs/op
BenchmarkParseHeavy-8              10000              4125 ns/op            6832 B/op          3 allocs/op
BenchmarkParseHeavy-8              10000              4093 ns/op            6832 B/op          3 allocs/op
BenchmarkParseHeavy-8              10000              4155 ns/op            6832 B/op          3 allocs/op
BenchmarkEvalHeavy-8               10000               253.5 ns/op           311 B/op          7 allocs/op
BenchmarkEvalHeavy-8               10000               239.0 ns/op           311 B/op          7 allocs/op
BenchmarkEvalHeavy-8               10000               242.3 ns/op           307 B/op          7 allocs/op
BenchmarkEvalHeavy-8               10000               238.4 ns/op           311 B/op          7 allocs/op
BenchmarkEvalHeavy-8               10000               244.2 ns/op           307 B/op          7 allocs/op
PASS
ok      github.com/nekrassov01/filter/benchmarks    0.610s
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
BenchmarkParseRepeated-8           10000            125342 ns/op          196912 B/op          3 allocs/op
BenchmarkParseRepeated-8           10000            123763 ns/op          196912 B/op          3 allocs/op
BenchmarkParseRepeated-8           10000            122811 ns/op          196912 B/op          3 allocs/op
BenchmarkParseRepeated-8           10000            124326 ns/op          196912 B/op          3 allocs/op
BenchmarkParseRepeated-8           10000            123979 ns/op          196912 B/op          3 allocs/op
BenchmarkEvalRepeated-8            10000              5186 ns/op             307 B/op          7 allocs/op
BenchmarkEvalRepeated-8            10000              5180 ns/op             307 B/op          7 allocs/op
BenchmarkEvalRepeated-8            10000              5207 ns/op             307 B/op          7 allocs/op
BenchmarkEvalRepeated-8            10000              5223 ns/op             307 B/op          7 allocs/op
BenchmarkEvalRepeated-8            10000              5168 ns/op             311 B/op          7 allocs/op
PASS
ok      github.com/nekrassov01/filter/benchmarks    6.837s
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

// Record is the example value to filter.
type Record struct {
    Name    string
    Latency time.Duration
    Retries int
    Enabled bool
}

// Resolve maps an identifier to its value.
func (r *Record) Resolve(name string) (any, bool) {
    switch name {
    case "Name":
        return r.Name, true
    case "Latency":
        return r.Latency, true
    case "Retries", "RetryCount":
        return r.Retries, true
    case "Enabled":
        return r.Enabled, true
    default:
        return nil, false
    }
}

func main() {
    input := `Name =~ '^foo' && (Latency < 1500ms || Retries != 0) && Enabled == true`

    expr, err := filter.Parse(input)
    if err != nil {
        panic(err)
    }

    record := &Record{
        Name:    "foobar",
        Latency: 100 * time.Millisecond,
        Retries: 3,
        Enabled: true,
    }

    ok, err := expr.Eval(record)
    if err != nil {
        panic(err)
    }
    fmt.Println("matched:", ok)
}
```

## Syntax

### Literals

| Kind     | Examples                                                                                                                | Notes                                                   |
| -------- | ----------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------- |
| String   | `"Hello"`, `'世界'`, `` `raw\ntext` ``                                                                                  | Double / single / raw (backtick)                        |
| Number   | `42`, `3.14`, `0x1.fp3`                                                                                                 | Subset of Go numeric literals                           |
| Time     | `2023-01-01T00:00:00Z`, `2023-01-01T09:00:00`, `2023-01-01`, `'2023-01-01 09:00:00'`, `'Sun, 01 Jan 2023 09:00:00 GMT'` | Zone-less forms are UTC; quote when it contains a space |
| Duration | `1500ms`, `2s`, `1h30m`, `4000μs`                                                                                       | Go `time.ParseDuration` compatible                      |
| Boolean  | `true`, `false`, `True`, `FALSE`                                                                                        | Case-insensitive variants accepted                      |

Time literals accept RFC 3339, `2006-01-02T15:04:05`, `2006-01-02 15:04:05`, `2006-01-02`, RFC 1123, RFC 850, RFC 822 (each with a named or numeric zone), and integer Unix seconds. Rules that follow from Go's `time.Parse`:

- Forms without a zone are read as UTC. A zone abbreviation is accepted only when it is `UTC` or `GMT`; use a numeric offset such as `+0900` for anything else
- Fractional seconds are accepted after any clock time
- Weekday names are not checked against the date
- Two-digit years (RFC 822, RFC 850) map to 1969–2068
- A number compared with a `time.Time` value is read as Unix seconds

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
