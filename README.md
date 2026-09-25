<p align="center">
  <img src="./assets/logo.png" alt="filter logo" width="120">
</p>
<h1 align="center">FILTER</h1>

<p align="center">A minimal filter expression language for Go</p>
<p align="center">
    <a href="https://github.com/nekrassov01/filter/actions/workflows/ci.yml"><img src="https://github.com/nekrassov01/filter/actions/workflows/ci.yml/badge.svg?branch=main" alt="CI" /></a>
    <a href="https://pkg.go.dev/github.com/nekrassov01/filter"><img src="https://pkg.go.dev/badge/github.com/nekrassov01/filter.svg" alt="Go Reference" /></a>
    <img src="https://img.shields.io/github/license/nekrassov01/filter" alt="LICENSE" />
    <a href="https://deepwiki.com/nekrassov01/filter"><img src="https://deepwiki.com/badge.svg" alt="Ask DeepWiki" /></a>
</p>

## Table of contents

Start with [Quick start](#quick-start) for usage or [Syntax](#syntax) for expression rules.

- [Table of contents](#table-of-contents)
- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Building values](#building-values)
- [Error handling](#error-handling)
- [Examples](#examples)
- [Performance](#performance)
  - [Parsing](#parsing)
  - [Evaluation](#evaluation)
- [Benchmarks](#benchmarks)
  - [Setup](#setup)
  - [Results](#results)
- [Syntax](#syntax)
  - [Identifiers](#identifiers)
  - [Strings](#strings)
  - [Numbers](#numbers)
  - [Times](#times)
  - [Durations](#durations)
  - [Booleans](#booleans)
  - [Operators](#operators)
- [Author](#author)
- [License](#license)

## Overview

`nekrassov01/filter` filters Go values with boolean expressions and a one-method resolver.

## Features

Designed for filtering existing Go values:

- Zero-allocation evaluation in the included benchmark: 1.7x faster than expr and 2.9x faster than CEL ([Benchmarks](#benchmarks))
- Low repeated preparation cost: 9.0x faster than expr and 86x faster than CEL in the same benchmark
- One-method integration with no reflection or map conversion
- Typed comparisons, regular expressions, and logical operators
- Lexing, parsing, and evaluation errors with source positions

## Installation

Add the latest release to your Go module:

```sh
go get github.com/nekrassov01/filter@latest
```

## Quick start

This example filters application logs by level, status, latency, and path.

1. Implement `filter.Resolver` to map identifiers to values.
2. Parse the expression with `filter.Parse`.
3. Call `Eval` on the returned `*Expr` for each input.

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
        return filter.Int(o.Status), true
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

Reuse the parsed `*Expr` across inputs and goroutines. For fixed expressions, `filter.MustParse` returns `*Expr` or panics on a parse error.

## Building values

Match the constructor to the Go value's type:

| Input type      | Constructor          |
| --------------- | -------------------- |
| `string`        | `filter.String(v)`   |
| `int`           | `filter.Int(v)`      |
| `int64`         | `filter.Int64(v)`    |
| `uint64`        | `filter.Uint64(v)`   |
| `float64`       | `filter.Float64(v)`  |
| `time.Duration` | `filter.Duration(v)` |
| `time.Time`     | `filter.Time(v)`     |
| `bool`          | `filter.Bool(v)`     |
| `any`           | `filter.ValueOf(v)`  |

For built-in numeric types, `ValueOf` distinguishes signed integers, unsigned integers, and floating-point values. See [Numbers](#numbers) for precision rules.

## Error handling

Use `errors.As` to inspect `*filter.Error` values returned by `Parse` or `Eval`:

| Field         | Meaning                                      |
| ------------- | -------------------------------------------- |
| `Kind`        | Lexing, parsing, or evaluation stage         |
| `Line`, `Col` | 1-based position; both zero when unavailable |

If `Resolve` returns `false`, `Eval` reports `unknown identifier "name"` at the identifier's position.

## Examples

Run the example tests to check their output:

```sh
go test ./examples/
```

## Performance

Parsed data is reused across evaluations; lookup results are cached only within each evaluation.

### Parsing

`Parse` prepares values for later evaluations:

| Input                                         | Preparation                                                       |
| --------------------------------------------- | ----------------------------------------------------------------- |
| Regular expressions                           | Compile once; reuse through a process-wide cache                  |
| Number, time, and duration literals           | Validate and convert; report errors with positions                |
| Quoted literals such as `"42"` and `"1500ms"` | Convert when recognized; otherwise defer conversion to evaluation |

Deferred conversion applies when comparing with a number, time, or duration.

### Evaluation

Each repeated identifier is resolved once per evaluation:

- Cache up to 16 distinct identifiers on the stack; use a heap slice beyond that.
- Skip the cache when no identifiers repeat.

## Benchmarks

The benchmarks compare `filter` with [expr](https://github.com/expr-lang/expr) and [CEL](https://github.com/google/cel-go).
See [benchmark_test.go](./benchmarks/benchmark_test.go) for the implementation.

### Setup

The libraries evaluate equivalent expressions over the same input values.

> [!NOTE]
> These results cover the shared boolean subset, not the libraries as a whole.
>
> | Library | Preparation                                                                         | Evaluation input                 |
> | ------- | ----------------------------------------------------------------------------------- | -------------------------------- |
> | filter  | Parse; reuse cached regular expressions after the first parse                       | `examples.Stats`                 |
> | expr    | Compile; read time and duration bounds from variables                               | The same `examples.Stats`        |
> | CEL     | Compile, fold constants, and precompile regex with `OptOptimize`; build the program | A map built from the same values |
>
> - Prepare: excludes reusable environment and option setup.
> - Eval: excludes expression preparation and CEL map construction.

Expression:

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

Run from the repository root or the `benchmarks` directory:

```sh
make bench
```

### Results

Apple M2, Go 1.27.1. Medians of 5 runs at `-benchtime 5s`.

| Benchmark  | filter                       | expr                            | CEL                               |
| ---------- | ---------------------------- | ------------------------------- | --------------------------------- |
| Prepare    | 2.765 µs, 6.375 KiB, 1 alloc | 24.99 µs, 28.89 KiB, 330 allocs | 236.8 µs, 230.4 KiB, 3,537 allocs |
| Eval Match | 198.5 ns, 0 B, 0 allocs      | 331.6 ns, 146 B, 1 alloc        | 571.5 ns, 147 B, 9 allocs         |
| Eval Miss  | 196.9 ns, 0 B, 0 allocs      | 317.3 ns, 147 B, 1 alloc        | 554.5 ns, 147 B, 9 allocs         |

Raw output of five runs:

```powershell
$ make bench
go test -run '^$' -bench . -benchmem -benchtime 5s -count 5 .
goos: darwin
goarch: arm64
pkg: benchmarks
cpu: Apple M2
BenchmarkPrepareFilter-8      1562955    3536 ns/op    6528 B/op     1 allocs/op
BenchmarkPrepareFilter-8      2128099    2765 ns/op    6528 B/op     1 allocs/op
BenchmarkPrepareFilter-8      2184588    2794 ns/op    6528 B/op     1 allocs/op
BenchmarkPrepareFilter-8      2216586    2711 ns/op    6528 B/op     1 allocs/op
BenchmarkPrepareFilter-8      2214692    2688 ns/op    6528 B/op     1 allocs/op
BenchmarkPrepareExpr-8         242307   25155 ns/op   29585 B/op   330 allocs/op
BenchmarkPrepareExpr-8         237588   24747 ns/op   29585 B/op   330 allocs/op
BenchmarkPrepareExpr-8         243700   25369 ns/op   29585 B/op   330 allocs/op
BenchmarkPrepareExpr-8         234158   24991 ns/op   29585 B/op   330 allocs/op
BenchmarkPrepareExpr-8         236510   24588 ns/op   29585 B/op   330 allocs/op
BenchmarkPrepareCEL-8           26066  233670 ns/op  235929 B/op  3537 allocs/op
BenchmarkPrepareCEL-8           25684  234218 ns/op  235920 B/op  3537 allocs/op
BenchmarkPrepareCEL-8           23479  244858 ns/op  235921 B/op  3537 allocs/op
BenchmarkPrepareCEL-8           25526  236772 ns/op  235940 B/op  3537 allocs/op
BenchmarkPrepareCEL-8           25450  246690 ns/op  235940 B/op  3537 allocs/op
BenchmarkEvalFilter/Match-8  29774217   199.7 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Match-8  31408548   193.3 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Match-8  30920469   192.5 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Match-8  31132731   203.6 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Match-8  30747924   198.5 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Miss-8   26613837   196.9 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Miss-8   30520560   198.0 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Miss-8   31283625   197.5 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Miss-8   26330175   193.1 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Miss-8   31357765   191.7 ns/op       0 B/op     0 allocs/op
BenchmarkEvalExpr/Match-8    18980235   387.2 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Match-8    18640117   336.4 ns/op     147 B/op     1 allocs/op
BenchmarkEvalExpr/Match-8    18837666   318.3 ns/op     147 B/op     1 allocs/op
BenchmarkEvalExpr/Match-8    18156148   329.1 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Match-8    19022139   331.6 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Miss-8     19398187   323.0 ns/op     147 B/op     1 allocs/op
BenchmarkEvalExpr/Miss-8     18628314   317.3 ns/op     147 B/op     1 allocs/op
BenchmarkEvalExpr/Miss-8     19545804   318.6 ns/op     147 B/op     1 allocs/op
BenchmarkEvalExpr/Miss-8     19458451   314.5 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Miss-8     19461538   313.4 ns/op     147 B/op     1 allocs/op
BenchmarkEvalCEL/Match-8     11063966   528.6 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Match-8     11558254   524.5 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Match-8     11377825   571.5 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Match-8      9214380   588.3 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Match-8      9452922   604.7 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Miss-8      10827618   564.0 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Miss-8       8845989   577.1 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Miss-8      11497362   554.5 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Miss-8      11456950   538.6 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Miss-8      11608744   537.2 ns/op     147 B/op     9 allocs/op
PASS
ok      benchmarks  271.018s
```

## Syntax

Expressions combine identifiers, literals, and operators.

### Identifiers

Identifiers name the values supplied by your resolver:

- Use Unicode letters, digits, and `_`; no dots.
- Boolean literals are reserved; see [Booleans](#booleans).

### Strings

Enclose strings in double quotes, single quotes, or backticks:

| Form            | Example           |
| --------------- | ----------------- |
| Double quotes   | `"Hello"`         |
| Single quotes   | `'世界'`          |
| Backticks (raw) | `` `raw\ntext` `` |

### Numbers

Numeric literals follow a subset of Go syntax, such as `42`, `3.14`, and `0x1.fp3`. Integer and floating-point forms use different rules.

| Literal form                            | Parsing                                                          |
| --------------------------------------- | ---------------------------------------------------------------- |
| Decimal integer, including quoted forms | Exact from `-9223372036854775808` through `18446744073709551615` |
| Integer outside that range              | Rejected when interpreted as a number                            |
| Leading zeros                           | Decimal                                                          |
| Underscores                             | Allowed between digits                                           |
| Decimal point or exponent               | Floating-point parsing; may round to `float64` precision         |

Comparisons involving an integer are exact, including comparisons with floats. Equality uses `filter.Epsilon` (`1e-9`) only when both operands are floating-point values.

For example:

- Integer `9007199254740993` is greater than `9007199254740992.0`.
- Integer `1` is not equal to `1.0000000001`.

### Times

Quote time literals containing spaces, such as `'2023-01-01 09:00:00'`. Accepted formats:

- RFC 3339
- `2006-01-02T15:04:05`, `2006-01-02 15:04:05`, or `2006-01-02`
- RFC 1123 and RFC 822, each with a named or numeric zone
- RFC 850 with a named zone
- Integer Unix seconds when compared with a `time.Time` value

Interpretation rules:

- Missing zone: UTC
- Named zones: `UTC` or `GMT` only; otherwise use a numeric offset such as `+0900`
- Fractional seconds are accepted after any clock time
- Weekday names are not checked against the date
- Two-digit years (RFC 822, RFC 850) map to 1969–2068

### Durations

Duration literals use `time.ParseDuration` syntax. Examples: `1500ms`, `2s`, `1h30m`, or `4000μs`.

### Booleans

Boolean literals accept these three case forms:

| Case       | True   | False   |
| ---------- | ------ | ------- |
| Lowercase  | `true` | `false` |
| Title case | `True` | `False` |
| Uppercase  | `TRUE` | `FALSE` |

### Operators

Combine comparisons into boolean conditions with these operators:

| Category | Operators         | Supported values                                                |
| -------- | ----------------- | --------------------------------------------------------------- |
| Ordering | `>` `>=` `<` `<=` | Numbers, times, and durations                                   |
| Equality | `==` `!=`         | All types; see [Numbers](#numbers) for floating-point tolerance |
| Regex    | `=~` `!~`         | Strings; the pattern must be a string literal                   |
| Logical  | `&&` `\|\|` `!`   | Boolean expressions; short-circuit evaluation                   |

Regular expressions use Go syntax.

## Author

[nekrassov01](https://github.com/nekrassov01)

## License

[MIT](https://github.com/nekrassov01/filter/blob/main/LICENSE)
