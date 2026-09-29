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
  - [IP addresses](#ip-addresses)
  - [Booleans](#booleans)
  - [Operators](#operators)
- [Author](#author)
- [License](#license)

## Overview

`nekrassov01/filter` filters Go values with boolean expressions and a one-method resolver.

## Features

Designed for filtering existing Go values:

- Zero-allocation evaluation in the included benchmark: 1.7x faster than expr and 2.8x faster than CEL ([Benchmarks](#benchmarks))
- Low repeated preparation cost: 9.7x faster than expr and 92x faster than CEL in the same benchmark
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
| `time.Time`     | `filter.Time(v)`     |
| `time.Duration` | `filter.Duration(v)` |
| `netip.Addr`    | `filter.Addr(v)`     |
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
| Number, time, duration, and IP literals       | Validate and convert; report errors with positions                |
| Quoted literals such as `"42"` and `"1500ms"` | Convert when recognized; otherwise defer conversion to evaluation |

Deferred conversion applies when comparing with a number, time, duration, or IP address.

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

| Benchmark  | filter                   | expr                            | CEL                               |
| ---------- | ------------------------ | ------------------------------- | --------------------------------- |
| Prepare    | 2.519 µs, 4 KiB, 1 alloc | 24.45 µs, 28.89 KiB, 330 allocs | 232.8 µs, 230.4 KiB, 3,537 allocs |
| Eval Match | 185.9 ns, 0 B, 0 allocs  | 309.1 ns, 146 B, 1 alloc        | 523.7 ns, 147 B, 9 allocs         |
| Eval Miss  | 185.5 ns, 0 B, 0 allocs  | 308.9 ns, 146 B, 1 alloc        | 528.9 ns, 147 B, 9 allocs         |

Raw output of five runs:

```powershell
$ make bench
go test -run '^$' -bench . -benchmem -benchtime 5s -count 5 .
goos: darwin
goarch: arm64
pkg: benchmarks
cpu: Apple M2
BenchmarkPrepareFilter-8      2366385    2519 ns/op    4096 B/op     1 allocs/op
BenchmarkPrepareFilter-8      2296510    2605 ns/op    4096 B/op     1 allocs/op
BenchmarkPrepareFilter-8      2284668    2563 ns/op    4096 B/op     1 allocs/op
BenchmarkPrepareFilter-8      2404842    2507 ns/op    4096 B/op     1 allocs/op
BenchmarkPrepareFilter-8      2405040    2515 ns/op    4096 B/op     1 allocs/op
BenchmarkPrepareExpr-8         249426   24449 ns/op   29585 B/op   330 allocs/op
BenchmarkPrepareExpr-8         252070   24342 ns/op   29585 B/op   330 allocs/op
BenchmarkPrepareExpr-8         224080   24607 ns/op   29585 B/op   330 allocs/op
BenchmarkPrepareExpr-8         252404   24398 ns/op   29585 B/op   330 allocs/op
BenchmarkPrepareExpr-8         237146   24892 ns/op   29585 B/op   330 allocs/op
BenchmarkPrepareCEL-8           26157  230388 ns/op  235926 B/op  3537 allocs/op
BenchmarkPrepareCEL-8           26084  231415 ns/op  235915 B/op  3537 allocs/op
BenchmarkPrepareCEL-8           25908  233023 ns/op  235929 B/op  3537 allocs/op
BenchmarkPrepareCEL-8           25828  232761 ns/op  235933 B/op  3537 allocs/op
BenchmarkPrepareCEL-8           25388  233987 ns/op  235941 B/op  3537 allocs/op
BenchmarkEvalFilter/Match-8  32449746   185.9 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Match-8  32548308   185.6 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Match-8  32811486   185.6 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Match-8  30899841   186.1 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Match-8  32820894   189.3 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Miss-8   32101962   184.9 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Miss-8   32391811   187.0 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Miss-8   32950225   185.5 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Miss-8   32680621   185.1 ns/op       0 B/op     0 allocs/op
BenchmarkEvalFilter/Miss-8   32978432   185.8 ns/op       0 B/op     0 allocs/op
BenchmarkEvalExpr/Match-8    19468948   308.2 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Match-8    19108596   309.1 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Match-8    19723540   307.9 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Match-8    19765957   310.2 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Match-8    19504303   315.9 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Miss-8     18889473   309.5 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Miss-8     19664346   308.9 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Miss-8     19140988   308.7 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Miss-8     19578049   308.4 ns/op     146 B/op     1 allocs/op
BenchmarkEvalExpr/Miss-8     19420458   309.8 ns/op     146 B/op     1 allocs/op
BenchmarkEvalCEL/Match-8     11329592   523.7 ns/op     146 B/op     9 allocs/op
BenchmarkEvalCEL/Match-8     11227610   522.7 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Match-8     11680453   521.9 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Match-8     11341896   531.6 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Match-8     11207260   539.0 ns/op     146 B/op     9 allocs/op
BenchmarkEvalCEL/Miss-8      11450215   531.5 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Miss-8      11527864   526.7 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Miss-8      11424969   526.3 ns/op     147 B/op     9 allocs/op
BenchmarkEvalCEL/Miss-8      11298447   528.9 ns/op     146 B/op     9 allocs/op
BenchmarkEvalCEL/Miss-8      11491068   529.0 ns/op     147 B/op     9 allocs/op
PASS
ok      benchmarks  271.035s
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

### IP addresses

Use `filter.Addr` or `filter.ValueOf` for `netip.Addr` values. IP literals accept IPv4, IPv6, and IPv6 zones:

```text
IP == 192.0.2.1
IP >= 2001:db8::1 && IP < 2001:db8::100
IP == fe80::1%eth0
```

Quotes are optional unless the zone contains whitespace, operators, parentheses, or quotes.

Comparisons follow `netip.Addr.Compare`:

- IPv4 sorts before IPv6; IPv4-mapped IPv6 remains distinct from IPv4.
- IPv6 zones participate in equality and ordering.
- An invalid Go address sorts before valid addresses.
- Invalid unquoted literals fail during parsing; quoted literals are validated when compared as addresses.

A Go `string` containing an IP address still uses string comparison. Convert `net.IP` with `netip.AddrFromSlice` before building a value.

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
| Ordering | `>` `>=` `<` `<=` | Numbers, times, durations, and IP addresses                     |
| Equality | `==` `!=`         | All types; see [Numbers](#numbers) for floating-point tolerance |
| Regex    | `=~` `!~`         | Strings; the pattern must be a string literal                   |
| Logical  | `&&` `\|\|` `!`   | Boolean expressions; short-circuit evaluation                   |

Regular expressions use Go syntax.

## Author

[nekrassov01](https://github.com/nekrassov01)

## License

[MIT](https://github.com/nekrassov01/filter/blob/main/LICENSE)
