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

`filter` focuses on one task: evaluating small boolean filter expressions in Go without the weight of a general expression engine. The motivation is to avoid large, reflection-heavy or feature-rich DSLs when you only need predictable value filtering. Core traits: minimal syntax (comparisons, basic logical operators, regex), no reflection (caller supplies values via a tiny interface), deterministic errors with positions, and cached regex compilation. This keeps the surface area small while remaining fast and explicit.

## Features

- Comparisons, regex, logical AND / OR / NOT
- Values via a one-method `Resolver` interface: `Resolve(name string) (filter.Value, bool)`
- Errors are `*filter.Error` with `Kind`, `Line`, and `Col`
- Supported types: string, all integer types, float32/64, time.Time, time.Duration, bool
- Regex: `=~` / `!~`
- Time literals: RFC 3339, RFC 1123, RFC 850, RFC 822, `2006-01-02T15:04:05`, `2006-01-02 15:04:05`, `2006-01-02`, Unix seconds; zone-less forms are UTC
- Duration literals: `1500ms`, `2s`, `1h30m`, `4000μs`

## Performance

`filter` intentionally does a small amount of work once, so that evaluating an expression many times stays flat:

- Regex literals: compiled exactly once per distinct pattern (process-wide sync cache). Writing the same "foo.*" pattern many times does not multiply compile cost.
- Number, time, and duration RHS literals: validated and converted once during parsing, so malformed literals are reported as parse errors with their position; eval just compares pre‑parsed values. Quoted forms like `"42"`, `"1500ms"`, or `"2023-01-01 09:00:00"` are converted during parsing too when their text reads as a literal, and otherwise at evaluation time when compared against a numeric, time, or duration value.
- Resolved value reuse: when an identifier appears more than once, each evaluation caches its value on first use in a small stack buffer (a heap slice only beyond 16 distinct identifiers); referencing the same identifier dozens of times does not add proportional `Resolve` overhead. Expressions where every identifier appears once skip the cache entirely.

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
BenchmarkParseASCIISimple-8                10000               491.0 ns/op           240 B/op          2 allocs/op
BenchmarkParseASCIISimple-8                10000               325.8 ns/op           240 B/op          2 allocs/op
BenchmarkParseASCIISimple-8                10000               321.5 ns/op           240 B/op          2 allocs/op
BenchmarkParseASCIISimple-8                10000               339.7 ns/op           240 B/op          2 allocs/op
BenchmarkParseASCIISimple-8                10000               318.7 ns/op           240 B/op          2 allocs/op
BenchmarkEvalASCIISimple-8                 10000                18.78 ns/op           16 B/op          1 allocs/op
BenchmarkEvalASCIISimple-8                 10000                19.23 ns/op           16 B/op          1 allocs/op
BenchmarkEvalASCIISimple-8                 10000                18.33 ns/op           16 B/op          1 allocs/op
BenchmarkEvalASCIISimple-8                 10000                17.57 ns/op           16 B/op          1 allocs/op
BenchmarkEvalASCIISimple-8                 10000                18.72 ns/op           16 B/op          1 allocs/op
BenchmarkParseUnicodeSimple-8              10000               375.7 ns/op           241 B/op          2 allocs/op
BenchmarkParseUnicodeSimple-8              10000               350.9 ns/op           240 B/op          2 allocs/op
BenchmarkParseUnicodeSimple-8              10000               328.4 ns/op           240 B/op          2 allocs/op
BenchmarkParseUnicodeSimple-8              10000               300.8 ns/op           240 B/op          2 allocs/op
BenchmarkParseUnicodeSimple-8              10000               302.2 ns/op           240 B/op          2 allocs/op
BenchmarkEvalUnicodeSimple-8               10000                17.00 ns/op           16 B/op          1 allocs/op
BenchmarkEvalUnicodeSimple-8               10000                17.55 ns/op           16 B/op          1 allocs/op
BenchmarkEvalUnicodeSimple-8               10000                17.86 ns/op           16 B/op          1 allocs/op
BenchmarkEvalUnicodeSimple-8               10000                17.26 ns/op           16 B/op          1 allocs/op
BenchmarkEvalUnicodeSimple-8               10000                17.54 ns/op           16 B/op          1 allocs/op
PASS
ok      github.com/nekrassov01/filter/benchmarks        0.397s
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
BenchmarkParseASCIIHeavy-8         10000              3737 ns/op            6833 B/op          3 allocs/op
BenchmarkParseASCIIHeavy-8         10000              3642 ns/op            6832 B/op          3 allocs/op
BenchmarkParseASCIIHeavy-8         10000              3274 ns/op            6832 B/op          3 allocs/op
BenchmarkParseASCIIHeavy-8         10000              3126 ns/op            6832 B/op          3 allocs/op
BenchmarkParseASCIIHeavy-8         10000              3104 ns/op            6832 B/op          3 allocs/op
BenchmarkEvalASCIIHeavy-8          10000               233.6 ns/op           307 B/op          7 allocs/op
BenchmarkEvalASCIIHeavy-8          10000               237.1 ns/op           307 B/op          7 allocs/op
BenchmarkEvalASCIIHeavy-8          10000               238.4 ns/op           307 B/op          7 allocs/op
BenchmarkEvalASCIIHeavy-8          10000               241.7 ns/op           307 B/op          7 allocs/op
BenchmarkEvalASCIIHeavy-8          10000               230.4 ns/op           307 B/op          7 allocs/op
BenchmarkParseUnicodeHeavy-8       10000              3189 ns/op            6833 B/op          3 allocs/op
BenchmarkParseUnicodeHeavy-8       10000              3103 ns/op            6832 B/op          3 allocs/op
BenchmarkParseUnicodeHeavy-8       10000              3115 ns/op            6832 B/op          3 allocs/op
BenchmarkParseUnicodeHeavy-8       10000              7892 ns/op            6832 B/op          3 allocs/op
BenchmarkParseUnicodeHeavy-8       10000              3150 ns/op            6832 B/op          3 allocs/op
BenchmarkEvalUnicodeHeavy-8        10000               241.4 ns/op           307 B/op          7 allocs/op
BenchmarkEvalUnicodeHeavy-8        10000               235.1 ns/op           307 B/op          7 allocs/op
BenchmarkEvalUnicodeHeavy-8        10000               242.2 ns/op           307 B/op          7 allocs/op
BenchmarkEvalUnicodeHeavy-8        10000               231.4 ns/op           307 B/op          7 allocs/op
BenchmarkEvalUnicodeHeavy-8        10000               237.6 ns/op           307 B/op          7 allocs/op
PASS
ok      github.com/nekrassov01/filter/benchmarks        0.719s
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
BenchmarkParseASCIIRepeated-8              10000             89408 ns/op          180529 B/op          3 allocs/op
BenchmarkParseASCIIRepeated-8              10000             86024 ns/op          180529 B/op          3 allocs/op
BenchmarkParseASCIIRepeated-8              10000             87320 ns/op          180528 B/op          3 allocs/op
BenchmarkParseASCIIRepeated-8              10000             86209 ns/op          180528 B/op          3 allocs/op
BenchmarkParseASCIIRepeated-8              10000             87820 ns/op          180528 B/op          3 allocs/op
BenchmarkEvalASCIIRepeated-8               10000              5614 ns/op             304 B/op          7 allocs/op
BenchmarkEvalASCIIRepeated-8               10000              6135 ns/op             307 B/op          7 allocs/op
BenchmarkEvalASCIIRepeated-8               10000              5861 ns/op             307 B/op          7 allocs/op
BenchmarkEvalASCIIRepeated-8               10000              6843 ns/op             307 B/op          7 allocs/op
BenchmarkEvalASCIIRepeated-8               10000              7280 ns/op             307 B/op          7 allocs/op
BenchmarkParseUnicodeRepeated-8            10000            110986 ns/op          180531 B/op          3 allocs/op
BenchmarkParseUnicodeRepeated-8            10000             99490 ns/op          180528 B/op          3 allocs/op
BenchmarkParseUnicodeRepeated-8            10000             86554 ns/op          180528 B/op          3 allocs/op
BenchmarkParseUnicodeRepeated-8            10000             93112 ns/op          180528 B/op          3 allocs/op
BenchmarkParseUnicodeRepeated-8            10000             87269 ns/op          180528 B/op          3 allocs/op
BenchmarkEvalUnicodeRepeated-8             10000              5316 ns/op             304 B/op          7 allocs/op
BenchmarkEvalUnicodeRepeated-8             10000              5187 ns/op             307 B/op          7 allocs/op
BenchmarkEvalUnicodeRepeated-8             10000              5321 ns/op             307 B/op          7 allocs/op
BenchmarkEvalUnicodeRepeated-8             10000              5626 ns/op             307 B/op          7 allocs/op
BenchmarkEvalUnicodeRepeated-8             10000              5219 ns/op             307 B/op          7 allocs/op
PASS
ok      github.com/nekrassov01/filter/benchmarks        10.079s
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
func (r *Record) Resolve(name string) (filter.Value, bool) {
    switch name {
    case "Name":
        return filter.String(r.Name), true
    case "Latency":
        return filter.Duration(r.Latency), true
    case "Retries", "RetryCount":
        return filter.Number(float64(r.Retries)), true
    case "Enabled":
        return filter.Bool(r.Enabled), true
    default:
        return filter.Value{}, false
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

| Category   | Operators                   | Description                             |
| ---------- | --------------------------- | --------------------------------------- |
| Comparison | `>` `>=` `<` `<=` `==` `!=` | Strings, integers, times, and durations |
| Regex      | `=~` `!~`                   | Cached per pattern string               |
| Logical    | `&&` `\|\|` `!`             | Short-circuit                           |

## Author

[nekrassov01](https://github.com/nekrassov01)

## License

[MIT](https://github.com/nekrassov01/filter/blob/main/LICENSE)
