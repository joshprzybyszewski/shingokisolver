

.PHONY: profile
profile:
	go build -o solver.out .
	./solver.out -includeProfile
	go tool pprof solver.out solverProfile.pprof

.PHONY: debug
debug:
	go run . -includeProcessLogs

.PHONY: run
run:
	go run . 

.PHONY: results
results:
	go run . -results

.PHONY: compete
compete:
	rm -rf temp/*
	go run -tags="secretSauce" . -competitive

.PHONY: test
test:
	go test -short ./...

.PHONY: longtest
longtest:
	go test ./...

.PHONY: bench
bench:
	go test -benchmem -run="^$$" -bench "^(BenchmarkSolve)$$" ./solvers/...
