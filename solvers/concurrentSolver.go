package solvers

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

var (
	numCPU = runtime.NumCPU()
)

type unsolvedPayload struct {
	puzz   puzzle.Puzzle
	target model.Target

	isNodesComplete bool
	nextUnknown     model.EdgePair
}

type concurrentSolver struct {
	work     chan unsolvedPayload
	solution chan puzzle.Puzzle
}

func (cs concurrentSolver) solve(
	puzz puzzle.Puzzle,
) (sr SolvedResults, _ error) {
	defer func(t0 time.Time) {
		sr.Duration = time.Since(t0)
	}(time.Now())

	puzz, target, state := claimGimmes(puzz)
	switch state {
	case model.Complete:
		return SolvedResults{
			Puzzle: puzz,
		}, nil
	case model.Incomplete, model.NodesComplete:
		// move down into concurrency
	default:
		return SolvedResults{}, fmt.Errorf("puzzle unsolvable: %s", puzz.String())
	}

	// Start concurrency to find our solution
	cs.work = make(chan unsolvedPayload, numCPU)
	cs.solution = make(chan puzzle.Puzzle, 1)

	defer close(cs.work)
	defer close(cs.solution)

	for i := 0; i < numCPU; i++ {
		go cs.startWorker()
	}

	cs.processPayload(unsolvedPayload{
		puzz:            puzz,
		target:          target,
		isNodesComplete: state == model.NodesComplete,
	})

	sol, ok := <-cs.solution
	if !ok {
		log.Printf("solution channel closed!")
		return SolvedResults{}, fmt.Errorf("puzzle unsolvable: %s", puzz.String())
	}

	return SolvedResults{
		Puzzle: sol,
	}, nil
}

func (cs concurrentSolver) startWorker() {
	for payload := range cs.work {
		cs.processPayload(payload)
	}
}

func (cs concurrentSolver) processPayload(
	payload unsolvedPayload,
) {
	if payload.isNodesComplete {
		// TODO flip
		return
	}

	cs.solveAimingAtTarget(payload.puzz, payload.target)
}

func (cs concurrentSolver) solveAimingAtTarget(
	puzz puzzle.Puzzle,
	targeting model.Target,
) {

	printPuzzleUpdate(`solveAimingAtTarget`, puzz, targeting)

	// Check to see if this node has already been completed.
	switch puzz.GetNodeState(targeting.Node.Coord()) {
	case model.Violation:
		printPuzzleUpdate(`solveAimingAtTarget GetNodeState issue!`, puzz, targeting)
		return

	case model.Complete:
		cs.stepAheadToNextTarget(puzz, targeting)
		return
	}

	// for each of the TwoArm options, we're going to try setting the edges
	// and then descending further into our targets
	for _, option := range targeting.Options {
		// then, once we find a completion path, add it to the returned slice
		withArms, ok := addTwoArms(puzz, targeting.Node.Coord(), option)
		if ok {
			cs.stepAheadToNextTarget(withArms, targeting)
		}
	}
}

func (cs concurrentSolver) stepAheadToNextTarget(
	puzz puzzle.Puzzle,
	curTarget model.Target,
) {

	printPuzzleUpdate(`stepAheadToNextTarget`, puzz, curTarget)

	nextTarget, state := puzz.GetNextTarget(curTarget)
	switch state {
	case model.Complete:
		cs.solution <- puzz
	case model.Incomplete, model.NodesComplete:
		cs.work <- unsolvedPayload{
			puzz:            puzz,
			target:          nextTarget,
			isNodesComplete: state == model.NodesComplete,
		}
	}
}

func (cs concurrentSolver) firstFlip(
	puzz puzzle.Puzzle,
) {
	printPuzzleUpdate(`firstFlip`, puzz, model.InvalidTarget)

	ep, ok := puzz.GetUnknownEdge()
	if !ok {
		switch puzz.GetState() {
		case model.Complete:
			cs.solution <- puzz
		}
		return
	}

	cs.flip(puzz, ep)
}

func (cs concurrentSolver) flip(
	puzz puzzle.Puzzle,
	ep model.EdgePair,
) {

	printPuzzleUpdate(`flip`, puzz, model.InvalidTarget)

	var nextUnknown model.EdgePair

	puzzWithEdge, state := puzzle.AddEdge(
		puzz,
		ep,
	)
	switch state {
	case model.Complete, model.Incomplete:
		nextUnknown, state = puzzWithEdge.GetStateOfLoop(ep.NodeCoord)
		switch state {
		case model.Complete:
			cs.solution <- puzzWithEdge
		case model.Incomplete:
			cs.work <- unsolvedPayload{
				puzz:            puzz,
				isNodesComplete: true,
				nextUnknown:     nextUnknown,
			}
		}
	}

	puzzWithoutEdge, state := puzzle.AvoidEdge(
		puzz,
		ep,
	)
	switch state {
	case model.Complete, model.Incomplete:
		nextUnknown, state = puzzWithoutEdge.GetStateOfLoop(ep.NodeCoord)
		switch state {
		case model.Complete:
			cs.solution <- puzzWithEdge
		case model.Incomplete:
			cs.work <- unsolvedPayload{
				puzz:            puzz,
				isNodesComplete: true,
				nextUnknown:     nextUnknown,
			}
		}
	}
}
