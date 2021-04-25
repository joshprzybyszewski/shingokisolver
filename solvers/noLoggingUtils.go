// +build !debug

package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func (cs *concurrentSolver) printPuzzleUpdate(string, puzzle.Puzzle, model.Target) {
}

func (cs *concurrentSolver) printPayload(string, targetPayload) {
}

func (cs *concurrentSolver) printFlippingPayload(string, flippingPayload) {
}

func (cs *concurrentSolver) logMeta() {
}

func printPuzzleUpdate(string, puzzle.Puzzle, model.Target) {
}
