

.PHONY:
profile:
	go build -o solver.out .
	./solver.out
	go tool pprof solver.out solverProfile