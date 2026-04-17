export GO111MODULE=on

.PHONY: bench
bench:
	go test -bench . -benchmem -count 5 -benchtime 10000x -cpuprofile cpu.prof -memprofile mem.prof .

.PHONY: clean
clean:
	go clean
	rm -f cpu.prof mem.prof benchmarks.test
