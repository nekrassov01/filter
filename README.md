<p align="center">
  <img src="./assets/logo.png" alt="filter logo" width="120">
</p>
<h1 align="center">FILTER</h1>

<p align="center">The minimal filter expressions for Go</p>
<p align="center">
    <a href="https://github.com/nekrassov01/filter/actions/workflows/ci.yml"><img src="https://github.com/nekrassov01/filter/actions/workflows/ci.yml/badge.svg?branch=main" alt="CI" /></a>
    <a href="https://pkg.go.dev/github.com/nekrassov01/filter"><img src="https://pkg.go.dev/badge/github.com/nekrassov01/filter.svg" alt="Go Reference" /></a>
    <img src="https://img.shields.io/github/license/nekrassov01/filter" alt="LICENSE" />
    <a href="https://deepwiki.com/nekrassov01/filter"><img src="https://deepwiki.com/badge.svg" alt="Ask DeepWiki" /></a>
</p>

## Table of contents

- [Table of contents](#table-of-contents)
- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Examples](#examples)
- [Performance](#performance)
- [Benchmarks](#benchmarks)
  - [Setup](#setup)
  - [Results](#results)
- [Syntax](#syntax)
  - [Literals](#literals)
  - [Operators](#operators)
- [Author](#author)
- [License](#license)

## Overview

`filter` evaluates a focused boolean expression language in Go through a small value-resolver interface.

## Features

- Zero-allocation evaluation in the included benchmark: 1.6x faster than expr and 2.8x faster than CEL ([Benchmarks](#benchmarks))
- Low repeated preparation cost: 9.3x faster than expr and 89x faster than CEL in the same benchmark
- One-method integration with no reflection or map conversion
- Typed comparisons, regular expressions, and logical operators
- Positioned parse and evaluation errors

## Installation

Install with:

```sh
go get github.com/nekrassov01/filter@latest
```

## Quick start

The target structs must implement a small interface that returns a `filter.Value` for each identifier. Pass the expression as a string to `Parse`, and use the resulting `Expr` to evaluate the structs one by one.

```go
package main

import (
    "fmt"
    "time"

    "github.com/nekrassov01/filter"
)

// LogLine represents one application log entry.
type LogLine struct {
    Time    time.Time
    Level   string
    Status  int
    Latency time.Duration
    Path    string
}

// Resolve maps an identifier to a field of the log line.
func (o *LogLine) Resolve(name string) (filter.Value, bool) {
    switch name {
    case "time":
        return filter.Time(o.Time), true
    case "level":
        return filter.String(o.Level), true
    case "status":
        return filter.Number(float64(o.Status)), true
    case "latency":
        return filter.Duration(o.Latency), true
    case "path":
        return filter.String(o.Path), true
    default:
        return filter.Value{}, false
    }
}

func main() {
    lines := []LogLine{
        {Time: time.Now(), Level: "info", Status: 200, Latency: 12 * time.Millisecond, Path: "/api/users"},
        {Time: time.Now(), Level: "error", Status: 200, Latency: 8 * time.Millisecond, Path: "/api/orders"},
        {Time: time.Now(), Level: "warn", Status: 503, Latency: 900 * time.Millisecond, Path: "/api/search"},
        {Time: time.Now(), Level: "warn", Status: 503, Latency: 900 * time.Millisecond, Path: "/health"},
    }
    condition := `level == "error" || (status >= 500 && latency > 500ms && path !~ '^/health')`
    expr, err := filter.Parse(condition)
    if err != nil {
        panic(err)
    }
    for _, line := range lines {
        ok, err := expr.Eval(&line)
        if err != nil {
            panic(err)
        }
        if ok {
            fmt.Println(line.Level, line.Status, line.Path)
        }
    }
    // Output:
    // error 200 /api/orders
    // warn 503 /api/search
}
```

Notes on the API:

- `Parse` returns `*Expr`; `MustParse` panics instead of returning an error, for expressions fixed at build time.
- An `*Expr` is safe to share: `Eval` can run on it from many goroutines at once.
- Build values with `filter.String`, `Number`, `Duration`, `Time`, and `Bool`, or `filter.ValueOf(any)` when the value is already dynamically typed.
- A `Resolve` that returns `false` makes `Eval` fail with `unknown identifier "name"` at the identifier's position.
- Errors from `Parse` and `Eval` are `*filter.Error`; use `errors.As` to read `Kind`, `Line`, and `Col`.

## Examples

A sample test is provided for a quick functional check:

```sh
go test ./examples/
```

## Performance

`Parse` performs work that can be reused by every subsequent evaluation:

- A successfully compiled regex literal is stored in a process-wide cache, and later parses of the same pattern reuse it.
- Number, time, and duration literals are validated and converted during parsing, so a malformed literal is a parse error with a position and evaluation compares ready values. Quoted forms such as `"42"`, `"1500ms"`, or `"2023-01-01 09:00:00"` are converted at parse time too when their text reads as a literal, and otherwise at evaluation time against a number, time, or duration value.
- Resolved values are reused within an evaluation: when an identifier appears more than once, its value is cached on first use in a small stack buffer (a heap slice only beyond 16 distinct identifiers), so repeating an identifier does not repeat `Resolve`. Expressions where every identifier appears once skip the cache.

## Benchmarks

The same expression runs through `filter`, [expr](https://github.com/expr-lang/expr), and [CEL](https://github.com/google/cel-go). See [benchmark_test.go](./benchmarks/benchmark_test.go) for the inputs and the environments.

### Setup

> [!NOTE]
> These numbers compare only the shared boolean subset, not the libraries as a whole:
>
> - Scope: expr and CEL are general expression languages with type checking, functions, and macros; `filter` covers only the boolean subset used here.
> - Equivalent setup: each library receives the cheapest equivalent expression over the same struct fields. expr reads the time and duration bounds from variables; CEL folds constants and precompiles regular expressions with `OptOptimize`.
> - Prepare: measures `filter.Parse`, `expr.Compile`, or CEL compilation, constant folding, and program construction. Reusable environment and option setup is excluded. After the first parse, `filter` reuses its process-wide regular-expression cache.
> - Eval: prepares each expression once. filter and expr receive the same `examples.Stats` value; the CEL map is built from that value before measurement.
>
> Treat the results as the cost of this subset, not as a ranking of the libraries.

One input is used:

```text
Class == "軍師" &&
Name =~ '^(諸葛亮|龐統|法正)' &&
Name != "" &&
BirthDate < 0190-01-01T00:00:00Z &&
ATBGauge >= 20s &&
HitPoint > 50 &&
MagicPoint > 100 &&
LifePoint != 0 &&
Speed >= 20
```

Run them from the `benchmarks` module:

```bash
make bench
```

### Results

Results on Apple M2 with Go 1.27.0, median of 5 runs at `-benchtime 5s`:

| Benchmark  | filter                       | expr                            | CEL                               |
| ---------- | ---------------------------- | ------------------------------- | --------------------------------- |
| Prepare    | 2.645 µs, 6.375 KiB, 1 alloc | 24.65 µs, 28.89 KiB, 330 allocs | 235.9 µs, 230.4 KiB, 3,537 allocs |
| Eval Match | 197.5 ns, 0 B, 0 allocs      | 322.6 ns, 146 B, 1 alloc        | 549.7 ns, 147 B, 9 allocs         |
| Eval Miss  | 194.6 ns, 0 B, 0 allocs      | 319.7 ns, 146 B, 1 alloc        | 543.0 ns, 147 B, 9 allocs         |

<details>
<summary>Raw output of five runs</summary>

```powershell
$ make bench
go test -run '^$' -bench . -benchmem -benchtime 5s -count 5 .
goos: darwin
goarch: arm64
pkg: benchmarks
cpu: Apple M2
BenchmarkPrepareFilter-8            2268176        2663 ns/op      6528 B/op        1 allocs/op
BenchmarkPrepareFilter-8            2246515        2645 ns/op      6528 B/op        1 allocs/op
BenchmarkPrepareFilter-8            2289106        2624 ns/op      6528 B/op        1 allocs/op
BenchmarkPrepareFilter-8            2269144        2653 ns/op      6528 B/op        1 allocs/op
BenchmarkPrepareFilter-8            2270137        2629 ns/op      6528 B/op        1 allocs/op
BenchmarkPrepareExpr-8               213518       24599 ns/op     29585 B/op      330 allocs/op
BenchmarkPrepareExpr-8               247897       24536 ns/op     29585 B/op      330 allocs/op
BenchmarkPrepareExpr-8               233426       24854 ns/op     29585 B/op      330 allocs/op
BenchmarkPrepareExpr-8               252259       24823 ns/op     29585 B/op      330 allocs/op
BenchmarkPrepareExpr-8               253591       24646 ns/op     29585 B/op      330 allocs/op
BenchmarkPrepareCEL-8                 26018      234934 ns/op    235925 B/op     3537 allocs/op
BenchmarkPrepareCEL-8                 25333      242909 ns/op    235913 B/op     3537 allocs/op
BenchmarkPrepareCEL-8                 25419      235659 ns/op    235938 B/op     3537 allocs/op
BenchmarkPrepareCEL-8                 25440      235885 ns/op    235943 B/op     3537 allocs/op
BenchmarkPrepareCEL-8                 25336      237116 ns/op    235952 B/op     3537 allocs/op
BenchmarkEvalFilter/Match-8        29084169       198.0 ns/op         0 B/op        0 allocs/op
BenchmarkEvalFilter/Match-8        31754208       192.2 ns/op         0 B/op        0 allocs/op
BenchmarkEvalFilter/Match-8        29332669       197.5 ns/op         0 B/op        0 allocs/op
BenchmarkEvalFilter/Match-8        31774808       197.2 ns/op         0 B/op        0 allocs/op
BenchmarkEvalFilter/Match-8        31623531       201.6 ns/op         0 B/op        0 allocs/op
BenchmarkEvalFilter/Miss-8         30802072       200.9 ns/op         0 B/op        0 allocs/op
BenchmarkEvalFilter/Miss-8         31632187       194.6 ns/op         0 B/op        0 allocs/op
BenchmarkEvalFilter/Miss-8         31694331       193.8 ns/op         0 B/op        0 allocs/op
BenchmarkEvalFilter/Miss-8         30615394       192.6 ns/op         0 B/op        0 allocs/op
BenchmarkEvalFilter/Miss-8         31775536       195.8 ns/op         0 B/op        0 allocs/op
BenchmarkEvalExpr/Match-8          18688982       332.0 ns/op       146 B/op        1 allocs/op
BenchmarkEvalExpr/Match-8          19136402       319.1 ns/op       146 B/op        1 allocs/op
BenchmarkEvalExpr/Match-8          19175406       317.3 ns/op       146 B/op        1 allocs/op
BenchmarkEvalExpr/Match-8          19129795       326.4 ns/op       147 B/op        1 allocs/op
BenchmarkEvalExpr/Match-8          19149080       322.6 ns/op       146 B/op        1 allocs/op
BenchmarkEvalExpr/Miss-8           19069545       318.2 ns/op       146 B/op        1 allocs/op
BenchmarkEvalExpr/Miss-8           18322120       319.8 ns/op       146 B/op        1 allocs/op
BenchmarkEvalExpr/Miss-8           18774595       319.6 ns/op       146 B/op        1 allocs/op
BenchmarkEvalExpr/Miss-8           18598035       319.7 ns/op       146 B/op        1 allocs/op
BenchmarkEvalExpr/Miss-8           18969115       331.7 ns/op       147 B/op        1 allocs/op
BenchmarkEvalCEL/Match-8           10345711       549.7 ns/op       147 B/op        9 allocs/op
BenchmarkEvalCEL/Match-8           11293833       536.9 ns/op       147 B/op        9 allocs/op
BenchmarkEvalCEL/Match-8           10855304       547.6 ns/op       147 B/op        9 allocs/op
BenchmarkEvalCEL/Match-8           11381338       553.3 ns/op       147 B/op        9 allocs/op
BenchmarkEvalCEL/Match-8           10970720       559.9 ns/op       147 B/op        9 allocs/op
BenchmarkEvalCEL/Miss-8            10850890       550.3 ns/op       147 B/op        9 allocs/op
BenchmarkEvalCEL/Miss-8            10868491       542.0 ns/op       147 B/op        9 allocs/op
BenchmarkEvalCEL/Miss-8            11033768       543.0 ns/op       147 B/op        9 allocs/op
BenchmarkEvalCEL/Miss-8            10614069       544.9 ns/op       147 B/op        9 allocs/op
BenchmarkEvalCEL/Miss-8            11225605       542.8 ns/op       147 B/op        9 allocs/op
PASS
ok      benchmarks      272.084s
```

</details>

## Syntax

Identifiers are made of Unicode letters, digits, and `_`, with no dots; `true` and `false` in any letter case are literals, not identifiers.

### Literals

| Kind     | Examples                                                                                                                | Notes                                                   |
| -------- | ----------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------- |
| String   | `"Hello"`, `'世界'`, `` `raw\ntext` ``                                                                                  | Double / single / raw (backtick)                        |
| Number   | `42`, `3.14`, `0x1.fp3`                                                                                                 | Subset of Go numeric literals                           |
| Time     | `2023-01-01T00:00:00Z`, `2023-01-01T09:00:00`, `2023-01-01`, `'2023-01-01 09:00:00'`, `'Sun, 01 Jan 2023 09:00:00 GMT'` | Zone-less forms are UTC; quote when it contains a space |
| Duration | `1500ms`, `2s`, `1h30m`, `4000μs`                                                                                       | Go `time.ParseDuration` compatible                      |
| Boolean  | `true`, `false`, `True`, `FALSE`                                                                                        | Any letter case; compared as `true` / `false`           |

Time literals accept RFC 3339, `2006-01-02T15:04:05`, `2006-01-02 15:04:05`, `2006-01-02`, RFC 1123 and RFC 822 (each with a named or numeric zone), RFC 850 (named zone), and integer Unix seconds. Rules that follow from Go's `time.Parse`:

- Forms without a zone are read as UTC. A zone abbreviation is accepted only when it is `UTC` or `GMT`; use a numeric offset such as `+0900` for anything else
- Fractional seconds are accepted after any clock time
- Weekday names are not checked against the date
- Two-digit years (RFC 822, RFC 850) map to 1969–2068
- A number compared with a `time.Time` value is read as Unix seconds

### Operators

| Category   | Operators                   | Description                                                                                                    |
| ---------- | --------------------------- | -------------------------------------------------------------------------------------------------------------- |
| Comparison | `>` `>=` `<` `<=` `==` `!=` | Ordering for numbers, times, and durations; equality for all types, within `filter.Epsilon` (1e-9) for numbers |
| Regex      | `=~` `!~`                   | Go regular-expression syntax; the pattern must be a string literal, the value a string                         |
| Logical    | `&&` `\|\|` `!`             | Short-circuit                                                                                                  |

## Author

[nekrassov01](https://github.com/nekrassov01)

## License

[MIT](https://github.com/nekrassov01/filter/blob/main/LICENSE)
