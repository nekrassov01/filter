# filter

[![CI](https://github.com/nekrassov01/filter/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/nekrassov01/filter/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/nekrassov01/filter/graph/badge.svg?token=Z75YW69MQK)](https://codecov.io/gh/nekrassov01/filter)
[![Go Reference](https://pkg.go.dev/badge/github.com/nekrassov01/filter.svg)](https://pkg.go.dev/github.com/nekrassov01/filter)
[![Go Report Card](https://goreportcard.com/badge/github.com/nekrassov01/filter)](https://goreportcard.com/report/github.com/nekrassov01/filter)

`filter` is a simple library for parsing and evaluating boolean filter expressions in Go. It parses short expressions and evaluates them against field values you provide.

## Purpose

It is not a general-purpose expression engine; it provides only the minimal syntax needed for filtering. It stays lightweight and avoids reflection.

## Features

- Comparisons, regex, logical AND / OR / NOT
- Supported types: string, all integer types, float32/64, time.Duration, bool
- Case‑insensitive equality: `==*` / `!=*`
- Regex: `=~` / `!~` (right-hand side must be a string literal)
- Duration literals: `1500ms`, `2s`, `1h30m`, `4000μs`
- Floating equality: absolute epsilon `1e-9` for `==` / `!=`

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

## Literals

| Kind     | Examples                               | Notes                              |
| -------- | -------------------------------------- | ---------------------------------- |
| String   | `"Hello"`, `'世界'`, `` `raw\ntext` `` | Double / single / raw (backtick)   |
| Number   | `42`, `3.14`, `0x1.fp3`                | Subset of Go numeric literals      |
| Duration | `1500ms`, `2s`, `1h30m`, `4000μs`      | Go `time.ParseDuration` compatible |
| Boolean  | `true`, `false`, `True`, `FALSE`       | Case-insensitive variants accepted |

## Operators

| Category                  | Operators         | Description                                        |
| ------------------------- | ----------------- | -------------------------------------------------- |
| Comparison                | `> >= < <= == !=` | Numbers / durations (`==` / `!=` also for strings) |
| Case-insensitive (string) | `==* !=*`         | Unicode case folding                               |
| Regex                     | `=~ !~`           | RE2 (Go regex), cached per pattern                 |
| Logical                   | `&&` `\|\|` `!`   | Short-circuit                                      |

## Benchmarks

See `benchmark_test.go`.

```bash
go test -run=^$ -bench=. -benchmem -count 5 -benchtime=10000x
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
