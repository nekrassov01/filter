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

## Features

- Comparisons, regex, logical AND / OR / NOT
- Supported types: string, all integer types, float32/64, time.Duration, bool
- Case‑insensitive equality: `==*` / `!=*`
- Regex: `=~` / `!~` (right-hand side must be a string literal)
- Duration literals: `1500ms`, `2s`, `1h30m`, `4000μs`

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
	// matched
}
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

## Benchmarks

`filter` is designed to be memory efficient.

```bash
$ go test -run=^$ -bench=. -benchmem -count 5 -benchtime=10000x
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/filter
cpu: Apple M2
BenchmarkParse-8           10000              5942 ns/op            2945 B/op          4 allocs/op
BenchmarkParse-8           10000              3744 ns/op            2944 B/op          4 allocs/op
BenchmarkParse-8           10000              3615 ns/op            2944 B/op          4 allocs/op
BenchmarkParse-8           10000              3582 ns/op            2944 B/op          4 allocs/op
BenchmarkParse-8           10000              3438 ns/op            2944 B/op          4 allocs/op
BenchmarkEval-8            10000               208.5 ns/op            67 B/op          5 allocs/op
BenchmarkEval-8            10000               205.8 ns/op            67 B/op          5 allocs/op
BenchmarkEval-8            10000               221.9 ns/op            67 B/op          5 allocs/op
BenchmarkEval-8            10000               233.2 ns/op            67 B/op          5 allocs/op
BenchmarkEval-8            10000               216.5 ns/op            64 B/op          5 allocs/op
PASS
ok      github.com/nekrassov01/filter   0.428s
```

## Author

[nekrassov01](https://github.com/nekrassov01)

## License

[MIT](https://github.com/nekrassov01/filter/blob/main/LICENSE)
