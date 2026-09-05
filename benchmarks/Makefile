benchtime ?= 5s
count ?= 5

.PHONY: bench clean
bench:
	go test -run '^$$' -bench . -benchmem -benchtime $(benchtime) -count $(count) .

clean:
	go clean
