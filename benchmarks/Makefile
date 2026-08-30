export GO111MODULE=on

# target selects the benchmarks to run: filter (default) or comparison.
target ?= filter
ifeq ($(target),filter)
BENCH := (Simple|Heavy)Filter$$
else ifeq ($(target),comparison)
BENCH := (Simple|Heavy)(Filter|Expr|CEL)$$
else
$(error unknown target "$(target)": use filter or comparison)
endif

.PHONY: bench
bench:
	go test -bench '$(BENCH)' -benchmem -count 5 -benchtime 10000x -cpuprofile cpu.prof -memprofile mem.prof .

.PHONY: clean
clean:
	go clean
	rm -f cpu.prof mem.prof benchmarks.test
