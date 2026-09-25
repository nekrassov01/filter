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

`nekrassov01/filter` evaluates a focused boolean expression language in Go through a small value-resolver interface.

## Features

- Zero-allocation evaluation in the included benchmark: 1.7x faster than expr and 2.9x faster than CEL ([Benchmarks](#benchmarks))
- Low repeated preparation cost: 9.0x faster than expr and 86x faster than CEL in the same benchmark
- One-method integration with no reflection or map conversion
- Typed comparisons, regular expressions, and logical operators
- Positioned lexing, parsing, and evaluation errors

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

Notes on the API:

- `Parse` returns `*Expr`; `MustParse` panics instead of returning an error, for expressions fixed at build time.
- An `*Expr` is safe to share: `Eval` can run on it from many goroutines at once.
- Build values with `filter.String`, `Int`, `Int64`, `Uint64`, `Float64`, `Duration`, `Time`, and `Bool`, or `filter.ValueOf(any)` when the value is already dynamically typed.
- Integer values retain all 64 bits: `Int` and `Int64` store signed integers, `Uint64` stores unsigned integers, and `Float64` stores floating-point values. `ValueOf` preserves these numeric categories as well.
- `Number` has been removed. Replace `Number(f)` with `Float64(f)` for floating-point values, and replace `Number(float64(i))` with `Int(i)`, `Int64(i)`, or `Uint64(i)` for integers.
- A `Resolve` that returns `false` makes `Eval` fail with `unknown identifier "name"` at the identifier's position.
- Errors from `Parse` and `Eval` are `*filter.Error`; use `errors.As` to read `Kind`, `Line`, and `Col`.

## Examples

Runnable examples are provided for a quick functional check:

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

Results on Apple M2 with Go 1.27.1, median of 5 runs at `-benchtime 5s`:

| Benchmark  | filter                       | expr                            | CEL                               |
| ---------- | ---------------------------- | ------------------------------- | --------------------------------- |
| Prepare    | 2.765 µs, 6.375 KiB, 1 alloc | 24.99 µs, 28.89 KiB, 330 allocs | 236.8 µs, 230.4 KiB, 3,537 allocs |
| Eval Match | 198.5 ns, 0 B, 0 allocs      | 331.6 ns, 146 B, 1 alloc        | 571.5 ns, 147 B, 9 allocs         |
| Eval Miss  | 196.9 ns, 0 B, 0 allocs      | 317.3 ns, 147 B, 1 alloc        | 554.5 ns, 147 B, 9 allocs         |

<details>
<summary>Raw output of five runs</summary>

```powershell
$ make bench
go test -run '^$' -bench . -benchmem -benchtime 5s -count 5 .
goos: darwin
goarch: arm64
pkg: benchmarks
cpu: Apple M2
BenchmarkPrepareFilter-8   	 1562955	      3536 ns/op	    6528 B/op	       1 allocs/op
BenchmarkPrepareFilter-8   	 2128099	      2765 ns/op	    6528 B/op	       1 allocs/op
BenchmarkPrepareFilter-8   	 2184588	      2794 ns/op	    6528 B/op	       1 allocs/op
BenchmarkPrepareFilter-8   	 2216586	      2711 ns/op	    6528 B/op	       1 allocs/op
BenchmarkPrepareFilter-8   	 2214692	      2688 ns/op	    6528 B/op	       1 allocs/op
BenchmarkPrepareExpr-8     	  242307	     25155 ns/op	   29585 B/op	     330 allocs/op
BenchmarkPrepareExpr-8     	  237588	     24747 ns/op	   29585 B/op	     330 allocs/op
BenchmarkPrepareExpr-8     	  243700	     25369 ns/op	   29585 B/op	     330 allocs/op
BenchmarkPrepareExpr-8     	  234158	     24991 ns/op	   29585 B/op	     330 allocs/op
BenchmarkPrepareExpr-8     	  236510	     24588 ns/op	   29585 B/op	     330 allocs/op
BenchmarkPrepareCEL-8      	   26066	    233670 ns/op	  235929 B/op	    3537 allocs/op
BenchmarkPrepareCEL-8      	   25684	    234218 ns/op	  235920 B/op	    3537 allocs/op
BenchmarkPrepareCEL-8      	   23479	    244858 ns/op	  235921 B/op	    3537 allocs/op
BenchmarkPrepareCEL-8      	   25526	    236772 ns/op	  235940 B/op	    3537 allocs/op
BenchmarkPrepareCEL-8      	   25450	    246690 ns/op	  235940 B/op	    3537 allocs/op
BenchmarkEvalFilter/Match-8         	29774217	       199.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkEvalFilter/Match-8         	31408548	       193.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkEvalFilter/Match-8         	30920469	       192.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkEvalFilter/Match-8         	31132731	       203.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkEvalFilter/Match-8         	30747924	       198.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkEvalFilter/Miss-8          	26613837	       196.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkEvalFilter/Miss-8          	30520560	       198.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkEvalFilter/Miss-8          	31283625	       197.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkEvalFilter/Miss-8          	26330175	       193.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkEvalFilter/Miss-8          	31357765	       191.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkEvalExpr/Match-8           	18980235	       387.2 ns/op	     146 B/op	       1 allocs/op
BenchmarkEvalExpr/Match-8           	18640117	       336.4 ns/op	     147 B/op	       1 allocs/op
BenchmarkEvalExpr/Match-8           	18837666	       318.3 ns/op	     147 B/op	       1 allocs/op
BenchmarkEvalExpr/Match-8           	18156148	       329.1 ns/op	     146 B/op	       1 allocs/op
BenchmarkEvalExpr/Match-8           	19022139	       331.6 ns/op	     146 B/op	       1 allocs/op
BenchmarkEvalExpr/Miss-8            	19398187	       323.0 ns/op	     147 B/op	       1 allocs/op
BenchmarkEvalExpr/Miss-8            	18628314	       317.3 ns/op	     147 B/op	       1 allocs/op
BenchmarkEvalExpr/Miss-8            	19545804	       318.6 ns/op	     147 B/op	       1 allocs/op
BenchmarkEvalExpr/Miss-8            	19458451	       314.5 ns/op	     146 B/op	       1 allocs/op
BenchmarkEvalExpr/Miss-8            	19461538	       313.4 ns/op	     147 B/op	       1 allocs/op
BenchmarkEvalCEL/Match-8            	11063966	       528.6 ns/op	     147 B/op	       9 allocs/op
BenchmarkEvalCEL/Match-8            	11558254	       524.5 ns/op	     147 B/op	       9 allocs/op
BenchmarkEvalCEL/Match-8            	11377825	       571.5 ns/op	     147 B/op	       9 allocs/op
BenchmarkEvalCEL/Match-8            	 9214380	       588.3 ns/op	     147 B/op	       9 allocs/op
BenchmarkEvalCEL/Match-8            	 9452922	       604.7 ns/op	     147 B/op	       9 allocs/op
BenchmarkEvalCEL/Miss-8             	10827618	       564.0 ns/op	     147 B/op	       9 allocs/op
BenchmarkEvalCEL/Miss-8             	 8845989	       577.1 ns/op	     147 B/op	       9 allocs/op
BenchmarkEvalCEL/Miss-8             	11497362	       554.5 ns/op	     147 B/op	       9 allocs/op
BenchmarkEvalCEL/Miss-8             	11456950	       538.6 ns/op	     147 B/op	       9 allocs/op
BenchmarkEvalCEL/Miss-8             	11608744	       537.2 ns/op	     147 B/op	       9 allocs/op
PASS
ok  	benchmarks	271.018s
```

</details>

## Syntax

Identifiers are made of Unicode letters, digits, and `_`, with no dots; `true` and `false` in lowercase, title case, or uppercase are literals, not identifiers.

### Literals

| Kind     | Examples                                                                                                                | Notes                                                   |
| -------- | ----------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------- |
| String   | `"Hello"`, `'世界'`, `` `raw\ntext` ``                                                                                  | Double / single / raw (backtick)                        |
| Number   | `42`, `3.14`, `0x1.fp3`                                                                                                 | Subset of Go numeric literals                           |
| Time     | `2023-01-01T00:00:00Z`, `2023-01-01T09:00:00`, `2023-01-01`, `'2023-01-01 09:00:00'`, `'Sun, 01 Jan 2023 09:00:00 GMT'` | Zone-less forms are UTC; quote when it contains a space |
| Duration | `1500ms`, `2s`, `1h30m`, `4000μs`                                                                                       | Go `time.ParseDuration` compatible                      |
| Boolean  | `true`, `false`, `True`, `FALSE`                                                                                        | Lowercase, title case, or uppercase only                |

Time literals accept RFC 3339, `2006-01-02T15:04:05`, `2006-01-02 15:04:05`, `2006-01-02`, RFC 1123 and RFC 822 (each with a named or numeric zone), RFC 850 (named zone), and integer Unix seconds. Rules that follow from Go's `time.Parse`:

- Forms without a zone are read as UTC. A zone abbreviation is accepted only when it is `UTC` or `GMT`; use a numeric offset such as `+0900` for anything else
- Fractional seconds are accepted after any clock time
- Weekday names are not checked against the date
- Two-digit years (RFC 822, RFC 850) map to 1969–2068
- A number compared with a `time.Time` value is read as Unix seconds

Decimal integer literals, including quoted numeric forms, retain their exact value in the range `-9223372036854775808` through `18446744073709551615`. Leading zeros are decimal, and underscores may separate digits. Integers outside this range are rejected when interpreted as numbers. A decimal point or exponent selects floating-point parsing, which can round the literal to `float64` precision.

Comparisons involving an integer are exact, including comparisons with floating-point values: the integer is never rounded to `float64`. For example, integer `9007199254740993` is greater than `9007199254740992.0`, and integer `1` is not equal to `1.0000000001`. When both operands are floating-point values, equality retains the `Epsilon` tolerance.

### Operators

| Category   | Operators                   | Description                                                                                                                                          |
| ---------- | --------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| Comparison | `>` `>=` `<` `<=` `==` `!=` | Ordering for numbers, times, and durations; equality for all types, within `filter.Epsilon` (1e-9) only when both operands are floating-point values |
| Regex      | `=~` `!~`                   | Go regular-expression syntax; the pattern must be a string literal, the value a string                                                               |
| Logical    | `&&` `\|\|` `!`             | Short-circuit                                                                                                                                        |

## Author

[nekrassov01](https://github.com/nekrassov01)

## License

[MIT](https://github.com/nekrassov01/filter/blob/main/LICENSE)
