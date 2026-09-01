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

`filter` evaluates small boolean filter expressions in Go, faster than expr and CEL on the boolean subset they share.

## Features

- Evaluation allocates nothing: a heavy predicate takes about 200 ns, against 350 ns for expr and 530 ns for CEL, and parsing the same predicate is about 10x and 75–90x faster ([Benchmarks](#benchmarks)).
- Values are supplied through a one-method interface, `Resolve(name string) (filter.Value, bool)`, so nothing is reflected on or copied into maps.
- Comparison, regular-expression, and logical operators over strings, numbers, times, durations, and booleans; the grammar is listed under [Syntax](#syntax).
- Errors are `*filter.Error` values that carry the stage (lex, parse, or eval) and the line and column of the offending token.

## Installation

Install with:

```sh
go get github.com/nekrassov01/filter@latest
```

## Quick start

The target structs must implement a small interface that resolves identifiers to `Value`s. Pass the expression as a string to `Parse`, and use the resulting `Expr` to evaluate the structs one by one.

```go
package main

import (
    "fmt"
    "time"

    "github.com/nekrassov01/filter"
)

// Record is one log line.
type Record struct {
    Time    time.Time
    Level   string
    Status  int
    Latency time.Duration
    Path    string
}

// Resolve maps an identifier to a field of the record.
func (r *Record) Resolve(name string) (filter.Value, bool) {
    switch name {
    case "time":
        return filter.Time(r.Time), true
    case "level":
        return filter.String(r.Level), true
    case "status":
        return filter.Number(float64(r.Status)), true
    case "latency":
        return filter.Duration(r.Latency), true
    case "path":
        return filter.String(r.Path), true
    default:
        return filter.Value{}, false
    }
}

func main() {
    // Parse once; the expression is reused for every record.
    expr, err := filter.Parse(`level == "error" || (status >= 500 && latency > 500ms && path !~ '^/health')`)
    if err != nil {
        panic(err)
    }

    records := []Record{
        {Time: time.Now(), Level: "info", Status: 200, Latency: 12 * time.Millisecond, Path: "/api/users"},
        {Time: time.Now(), Level: "error", Status: 200, Latency: 8 * time.Millisecond, Path: "/api/orders"},
        {Time: time.Now(), Level: "warn", Status: 503, Latency: 900 * time.Millisecond, Path: "/api/search"},
        {Time: time.Now(), Level: "warn", Status: 503, Latency: 900 * time.Millisecond, Path: "/health"},
    }
    for i := range records {
        ok, err := expr.Eval(&records[i])
        if err != nil {
            panic(err)
        }
        if ok {
            fmt.Println(records[i].Level, records[i].Status, records[i].Path)
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

Three choices keep the cost of an evaluation flat as the same expression runs again and again:

- Regex literals are compiled once per distinct pattern in a process-wide cache, so the same pattern in many expressions is compiled once.
- Number, time, and duration literals are validated and converted during parsing, so a malformed literal is a parse error with a position and evaluation compares ready values. Quoted forms such as `"42"`, `"1500ms"`, or `"2023-01-01 09:00:00"` are converted at parse time too when their text reads as a literal, and otherwise at evaluation time against a number, time, or duration value.
- Resolved values are reused within an evaluation: when an identifier appears more than once, its value is cached on first use in a small stack buffer (a heap slice only beyond 16 distinct identifiers), so repeating an identifier does not repeat `Resolve`. Expressions where every identifier appears once skip the cache.

## Benchmarks

The same predicates run through `filter`, [expr](https://github.com/expr-lang/expr), and [CEL](https://github.com/google/cel-go). See [benchmark_test.go](./benchmarks/benchmark_test.go) for the inputs and the environments.

### Setup

> [!NOTE]
> The three libraries differ in scale and purpose: expr and CEL are general expression languages with type checking, functions, and macros; `filter` covers only the boolean subset they are compared on. Each library is given the cheapest equivalent of the same predicate over the same struct fields (expr reads the time and duration bounds from variables because it has no time or duration literals and its `date()` and `duration()` calls run on every evaluation; CEL folds constants and precompiles regular expressions with `OptOptimize`). Treat the numbers as the cost of that subset, not as a ranking of the libraries.

Two inputs are used, each with an ASCII and a Unicode variant:

| Input  | Expression                                                                                                                                                                                                                                                                  |
| ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Simple | <pre>Class == "軍師"</pre>                                                                                                                                                                                                                                                  |
| Heavy  | <pre>Class == "軍師" && Name =~ '^(諸葛亮\|龐統\|法正)' && Name != "" && (<br/>    BirthDate < 0190-01-01T00:00:00Z && ATBGauge >= 20s<br/>) && (<br/>    HitPoint > 50 && MagicPoint > 100 && LifePoint != 0<br/>) && (<br/>    Magic >= 20 \|\| !(Speed < 20)<br/>)</pre> |

Run them from the `benchmarks` module:

```bash
make bench target=filter      # filter only
make bench target=comparison  # filter, expr, and CEL
```

### Results

Results on Apple M2, benchstat center of 5 runs at `-benchtime 100000x` (longer than the Makefile default so that start-up effects average out):

| Benchmark            | filter                     | expr                          | CEL                           |
| -------------------- | -------------------------- | ----------------------------- | ----------------------------- |
| Parse ASCII Simple   | 314 ns, 192 B, 1 alloc     | 6.9 µs, 11.7 KiB, 81 allocs   | 15.3 µs, 16.4 KiB, 331 allocs |
| Eval ASCII Simple    | 9.9 ns, 0 B, 0 allocs      | 70.5 ns, 176 B, 1 alloc       | 61.3 ns, 16 B, 1 alloc        |
| Parse ASCII Heavy    | 3.11 µs, 6.6 KiB, 2 allocs | 31.7 µs, 36.1 KiB, 416 allocs | 240 µs, 215 KiB, 3406 allocs  |
| Eval ASCII Heavy     | 202 ns, 0 B, 0 allocs      | 358 ns, 177 B, 1 alloc        | 523 ns, 147 B, 9 allocs       |
| Parse Unicode Simple | 291 ns, 192 B, 1 alloc     | 6.8 µs, 11.7 KiB, 81 allocs   | 22.4 µs, 25.8 KiB, 452 allocs |
| Eval Unicode Simple  | 10.3 ns, 0 B, 0 allocs     | 72.8 ns, 176 B, 1 alloc       | 64.3 ns, 16 B, 1 alloc        |
| Parse Unicode Heavy  | 3.08 µs, 6.6 KiB, 2 allocs | 30.5 µs, 33.1 KiB, 404 allocs | 270 µs, 257 KiB, 3955 allocs  |
| Eval Unicode Heavy   | 192 ns, 0 B, 0 allocs      | 343 ns, 178 B, 1 alloc        | 531 ns, 147 B, 9 allocs       |

<details>
<summary>Raw output of one run</summary>

```powershell
$ go test -bench '(Simple|Heavy)(Filter|Expr|CEL)$' -benchmem -count 1 -benchtime 100000x .
goos: darwin
goarch: arm64
pkg: benchmarks
cpu: Apple M2
BenchmarkParseASCIISimpleFilter-8         100000               398.6 ns/op           192 B/op          1 allocs/op
BenchmarkEvalASCIISimpleFilter-8          100000                13.21 ns/op            0 B/op          0 allocs/op
BenchmarkParseASCIIHeavyFilter-8          100000              3061 ns/op            6784 B/op          2 allocs/op
BenchmarkEvalASCIIHeavyFilter-8           100000               203.4 ns/op             0 B/op          0 allocs/op
BenchmarkParseUnicodeSimpleFilter-8       100000               293.5 ns/op           192 B/op          1 allocs/op
BenchmarkEvalUnicodeSimpleFilter-8        100000                10.09 ns/op            0 B/op          0 allocs/op
BenchmarkParseUnicodeHeavyFilter-8        100000              3034 ns/op            6784 B/op          2 allocs/op
BenchmarkEvalUnicodeHeavyFilter-8         100000               192.0 ns/op             0 B/op          0 allocs/op
BenchmarkParseASCIISimpleExpr-8           100000              6792 ns/op           11982 B/op         81 allocs/op
BenchmarkEvalASCIISimpleExpr-8            100000                87.51 ns/op          176 B/op          1 allocs/op
BenchmarkParseASCIIHeavyExpr-8            100000             31863 ns/op           36949 B/op        416 allocs/op
BenchmarkEvalASCIIHeavyExpr-8             100000               357.1 ns/op           178 B/op          1 allocs/op
BenchmarkParseUnicodeSimpleExpr-8         100000              6901 ns/op           11982 B/op         81 allocs/op
BenchmarkEvalUnicodeSimpleExpr-8          100000                75.48 ns/op          176 B/op          1 allocs/op
BenchmarkParseUnicodeHeavyExpr-8          100000             30428 ns/op           33875 B/op        404 allocs/op
BenchmarkEvalUnicodeHeavyExpr-8           100000               343.5 ns/op           179 B/op          1 allocs/op
BenchmarkParseASCIISimpleCEL-8            100000             15382 ns/op           16792 B/op        331 allocs/op
BenchmarkEvalASCIISimpleCEL-8             100000                62.36 ns/op           16 B/op          1 allocs/op
BenchmarkParseASCIIHeavyCEL-8             100000            238808 ns/op          219710 B/op       3405 allocs/op
BenchmarkEvalASCIIHeavyCEL-8              100000               514.3 ns/op           147 B/op          9 allocs/op
BenchmarkParseUnicodeSimpleCEL-8          100000             22494 ns/op           26461 B/op        452 allocs/op
BenchmarkEvalUnicodeSimpleCEL-8           100000                62.04 ns/op           16 B/op          1 allocs/op
BenchmarkParseUnicodeHeavyCEL-8           100000            268782 ns/op          262794 B/op       3955 allocs/op
BenchmarkEvalUnicodeHeavyCEL-8            100000               504.3 ns/op           146 B/op          9 allocs/op
PASS
ok  	benchmarks	63.524s
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
