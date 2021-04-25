.PHONY: help
help: ## Prints out the options available in this makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: profile
profile: ## Run the solver and grab a CPU profile using pprof
	go build -o solver.out .
	./solver.out -includeProfile
	go tool pprof solver.out solverProfile.pprof

.PHONY: serialdebug
serialdebug: ## Run the serial solver and include progress logs
	go run -tags="debug" .

.PHONY: serialrun
serialrun: ## Run the solver using the serial solver
	go run . 

.PHONY: run
run: ## Run the solver
	go run . -concurrency

.PHONY: debug
debug: ## Run the solver and include progress logs
	go run -tags="debug" . -concurrency

.PHONY: results
results: ## Run the solver and write the captured results to the READMEs
	go run . -concurrency -results

.PHONY: serialresults
serialresults: ## Run the serial solver and write the captured results to the READMEs
	go run . -results

.PHONY: compete
compete: ## Fetch puzzles from the internet and submit the answers for the hall of fame (requires secret sauce).
	rm -rf temp/*
	go run -tags="secretSauce" . -competitive

.PHONY: test
test: ## Run unit tests
	go test -short ./...

.PHONY: longtest
longtest: ## Runs unit tests, including potentially long-running ones
	go test ./...

.PHONY: bench
bench: ## Runs benchmarks (replaced by 'make results')
	go test -benchmem -run="^$$" -bench "^(BenchmarkSolve)$$" ./solvers/...

.PHONY: lint
lint: ## Runs linters (via golangci-lint) on golang code
	golangci-lint run -v ./...

.PHONY: build
build: ## builds the go binary with gcflags to see what the compiler's doing
	go build -gcflags='-m -m' ./...