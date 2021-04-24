

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

.PHONY: crun
crun:
	go run . -concurrency

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

.PHONY: lint
lint: ## Runs linters (via golangci-lint) on golang code
	golangci-lint run -v ./...

.PHONY: build
build: ## builds the go binary with gcflags to see what the compiler's doing
	go build -gcflags='-m -m' ./...