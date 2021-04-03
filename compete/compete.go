package compete

import (
	"log"
	"os"
	"sync"

	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

func Run() {

	defer flushLogs()

	wp, err := getPuzzle(5)
	if err != nil {
		panic(err)
	}

	s := solvers.NewSolver(
		wp.pd.NumEdges,
		wp.pd.Nodes,
		solvers.TargetSolverType,
	)

	sr, err := s.Solve()
	if err != nil {
		p := puzzle.NewPuzzle(
			wp.pd.NumEdges,
			wp.pd.Nodes,
		)
		log.Printf("derp. Couldn't solve. %v\n%s\n", err, p)
		return
	}

	submitAnswer(wp, sr)
}

var (
	flushLogsChan = make(chan struct{}, 0)
	logsWritten   sync.WaitGroup
)

func flushLogs() {
	close(flushLogsChan)
	logsWritten.Wait()
}

func writeToFile(
	filename string,
	bytes []byte,
) {
	logsWritten.Add(1)

	go func() {
		defer logsWritten.Done()

		f, err := os.OpenFile(filename,
			os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("error opening file: %v\n", err)
			return
		}
		defer f.Close()

		if _, err := f.Write(bytes); err != nil {
			log.Printf("error WriteString file: %v\n", err)
			return
		}
	}()
}
