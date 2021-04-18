

.PHONY: profile
profile:
	go build -o solver.out .
	./solver.out -includeProfile
	go tool pprof solver.out solverProfile.pprof

.PHONY: debug
debug:
	go run main.go -includeProcessLogs

.PHONY: run
run:
	go run main.go 

.PHONY: compete
compete:
	go run -tags="secretSauce" main.go -competitive

.PHONY: test
test:
	go test -short ./...

.PHONY: longtest
longtest:
	go test ./...

.PHONY: bench
bench:
	go test -benchmem -run="^$$" -bench "^(BenchmarkSolve)$$" ./solvers/...
