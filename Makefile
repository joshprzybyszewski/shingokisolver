

.PHONY:
profile:
	go build -o solver.out .
	./solver.out -includeProfile
	go tool pprof solver.out solverProfile