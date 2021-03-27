package main

import (
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/puzzlegrid"
	"github.com/joshprzybyszewski/shingokisolver/reader"
)

func main() {
	f, err := os.Create(`solverProfile.pprof`)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	t0 := time.Now()

	for _, pd := range reader.DefaultPuzzles() {
		err := puzzlegrid.NewDFSSolver(pd.NumEdges, pd.Nodes).Solve()
		// puzzlegrid.NewBFSSolver(numEdges, nodes).Solve()
		if err != nil {
			panic(err)
		}
		if time.Since(t0) > 4*time.Second {
			break
		}
	}
}
