<p align="center">
  <h2 align="center">FILTER</h2>
  <p align="center">The lightweight filter expression for Go</p>
  <p align="center">
    <a href="https://github.com/nekrassov01/filter/actions/workflows/test.yml"><img src="https://github.com/nekrassov01/filter/actions/workflows/test.yml/badge.svg?branch=main" alt="CI" /></a>
    <a href="https://codecov.io/gh/nekrassov01/filter"><img src="https://codecov.io/gh/nekrassov01/filter/graph/badge.svg?token=Z75YW69MQK" alt="Go Report Card" /></a>
    <a href="https://pkg.go.dev/github.com/nekrassov01/filter"><img src="https://pkg.go.dev/badge/github.com/nekrassov01/filter.svg" alt="Go Reference" /></a>
    <a href="https://goreportcard.com/report/github.com/nekrassov01/filter"><img src="https://goreportcard.com/badge/github.com/nekrassov01/filter" alt="Go Report Card" /></a>
    <img src="https://img.shields.io/github/license/nekrassov01/filter" alt="LICENSE" />
    <a href="https://deepwiki.com/nekrassov01/filter"><img src="https://deepwiki.com/badge.svg" alt="Ask DeepWiki" /></a>
  </p>
</p>

## Overview

`filter` focuses on one task: evaluating small boolean filter expressions in Go without the weight of a general expression engine. The motivation is to avoid large, reflection‑heavy or feature‑rich DSLs when you only need predictable field filtering. Core traits: minimal syntax (comparisons, basic logical operators, regex, case‑insensitive equality), no reflection (caller supplies values via a tiny interface), deterministic errors with positions, and cached regex compilation. This keeps the surface area small while remaining fast and explicit.

## Installation

```sh
go get github.com/nekrassov01/filter@latest
```

## Example

```go
// Target interface
type Target interface {
	GetField(key string) (any, error)
}

expr, err := filter.Parse(`Name =~ '^foo' && (Latency < 1500ms || Retries != 0) && Enabled == true`)
if err != nil {
	log.Fatal(err)
}
ok, err := expr.Eval(myTarget)
if err != nil {
	log.Fatal(err)
}
if ok {
	fmt.Println("matched")
}
```

## Features

- Comparisons, regex, logical AND / OR / NOT
- Supported types: string, all integer types, float32/64, time.Duration, bool
- Case‑insensitive equality: `==*` / `!=*`
- Regex: `=~` / `!~` (right-hand side must be a string literal)
- Duration literals: `1500ms`, `2s`, `1h30m`, `4000μs`

## Performance

`filter` intentionally does a small amount of work once, so that evaluating an expression many times stays flat:

- Regex literals: compiled exactly once per distinct pattern (process‑wide sync cache). Writing the same `/foo.*/` style pattern many times does not multiply compile cost.
- Numeric & duration RHS literals: parsed eagerly during parsing (including quoted forms like `"42"` or `"1500ms"`); eval just compares pre‑parsed values.
- Field value reuse: per evaluation a tiny map caches each identifier the first time it is requested; referencing the same field dozens of times (common in generated filters) does not add proportional `GetField` overhead.

Net effect: expressions with high token repetition scale sub‑linearly in both time and allocations compared to naïve re‑parsing / re‑compiling approaches.

## Benchmarks

`filter` is designed to be memory efficient. See [benchmark_test.go](./benchmark_test.go)

```bash
$ go test -run=^$ -bench=. -benchmem -count 5 -benchtime=10000x
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/filter
cpu: Apple M2
BenchmarkParse-8           10000              4876 ns/op            3728 B/op          6 allocs/op
BenchmarkParse-8           10000              4403 ns/op            3728 B/op          6 allocs/op
BenchmarkParse-8           10000              3895 ns/op            3728 B/op          6 allocs/op
BenchmarkParse-8           10000              3816 ns/op            3728 B/op          6 allocs/op
BenchmarkParse-8           10000              3789 ns/op            3728 B/op          6 allocs/op
BenchmarkEval-8            10000               285.5 ns/op            48 B/op          4 allocs/op
BenchmarkEval-8            10000               294.9 ns/op            48 B/op          4 allocs/op
BenchmarkEval-8            10000               314.6 ns/op            51 B/op          4 allocs/op
BenchmarkEval-8            10000               299.5 ns/op            51 B/op          4 allocs/op
BenchmarkEval-8            10000               315.0 ns/op            51 B/op          4 allocs/op
PASS
ok      github.com/nekrassov01/filter   0.449s
```

Stable even when Input is concatenated 50 times with `&&`.

```bash
...
BenchmarkParse-8           10000            198312 ns/op          198547 B/op         11 allocs/op
BenchmarkParse-8           10000            186675 ns/op          198546 B/op         11 allocs/op
BenchmarkParse-8           10000            187490 ns/op          198545 B/op         11 allocs/op
BenchmarkParse-8           10000            189315 ns/op          198545 B/op         11 allocs/op
BenchmarkParse-8           10000            191222 ns/op          198545 B/op         11 allocs/op
BenchmarkEval-8            10000             12252 ns/op              48 B/op          4 allocs/op
BenchmarkEval-8            10000             12142 ns/op              55 B/op          4 allocs/op
BenchmarkEval-8            10000             12581 ns/op              48 B/op          4 allocs/op
BenchmarkEval-8            10000             12167 ns/op              51 B/op          4 allocs/op
BenchmarkEval-8            10000             12156 ns/op              48 B/op          4 allocs/op
PASS
ok      github.com/nekrassov01/filter   10.349s
```

## Syntax

### Literals

| Kind     | Examples                               | Notes                              |
| -------- | -------------------------------------- | ---------------------------------- |
| String   | `"Hello"`, `'世界'`, `` `raw\ntext` `` | Double / single / raw (backtick)   |
| Number   | `42`, `3.14`, `0x1.fp3`                | Subset of Go numeric literals      |
| Duration | `1500ms`, `2s`, `1h30m`, `4000μs`      | Go `time.ParseDuration` compatible |
| Boolean  | `true`, `false`, `True`, `FALSE`       | Case-insensitive variants accepted |

### Operators

| Category                  | Operators         | Description                                        |
| ------------------------- | ----------------- | -------------------------------------------------- |
| Comparison                | `> >= < <= == !=` | Numbers / durations (`==` / `!=` also for strings) |
| Case-insensitive (string) | `==* !=*`         | Unicode case folding                               |
| Regex                     | `=~ !~`           | RE2 (Go regex), cached per pattern string          |
| Logical                   | `&&` `\|\|` `!`   | Short-circuit                                      |

## Author

[nekrassov01](https://github.com/nekrassov01)

## License

[MIT](https://github.com/nekrassov01/filter/blob/main/LICENSE)
