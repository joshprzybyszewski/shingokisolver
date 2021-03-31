

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
	go run main.go -competitive
