package solvers

import (
	"fmt"
	"log"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

const (
	// TODO determine how big of a channel I want...
	workChanLen = 256
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

	numPayloads  int32
	numProcessed int32
}

func (cs *concurrentSolver) solve(
	puzz puzzle.Puzzle,
) (sr SolvedResults, _ error) {
	if numCPU <= 1 {
		log.Printf("Unexpected! Running without concurrency (%d CPUs)", numCPU)
		return solvePuzzleByTargets(puzz)
	}

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
	cs.work = make(chan unsolvedPayload, workChanLen)
	cs.solution = make(chan puzzle.Puzzle, 1)

	defer close(cs.work)
	defer close(cs.solution)

	for i := 0; i < numCPU; i++ {
		go cs.startWorker()
	}

	cs.queuePayload(unsolvedPayload{
		puzz:            puzz,
		target:          target,
		isNodesComplete: state == model.NodesComplete,
		nextUnknown:     model.InvalidEdgePair,
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

func (cs *concurrentSolver) queuePayload(
	payload unsolvedPayload,
) {

	atomic.AddInt32(&cs.numPayloads, 1)

	defer func() {
		// it's alright if the work channel has been closed
		recover()
	}()

	// TODO consider having a channel for "node complete" and "pre-nodes complete"
	cs.work <- payload
}

func (cs *concurrentSolver) sendSolution(
	puzz puzzle.Puzzle,
) {

	defer func() {
		// it's alright if the solution channel has been closed
		recover()
	}()

	cs.solution <- puzz
}

func (cs *concurrentSolver) startWorker() {
	for payload := range cs.work {
		cs.processPayload(payload)
	}
}

func (cs *concurrentSolver) processPayload(
	payload unsolvedPayload,
) {
	atomic.AddInt32(&cs.numProcessed, 1)

	cs.printPayload(`processPayload`, payload)

	if payload.isNodesComplete {
		cs.flip(payload.puzz, payload.nextUnknown)
		return
	}

	cs.solveAimingAtTarget(payload.puzz, payload.target)
}

func (cs *concurrentSolver) solveAimingAtTarget(
	puzz puzzle.Puzzle,
	targeting model.Target,
) {

	cs.printPuzzleUpdate(`solveAimingAtTarget`, puzz, targeting)

	// Check to see if this node has already been completed.
	switch puzz.GetNodeState(targeting.Node.Coord()) {
	case model.Violation:
		cs.printPuzzleUpdate(`solveAimingAtTarget GetNodeState issue!`, puzz, targeting)
		return

	case model.Complete:
		up, canProcess, isSolved := cs.getNextPayload(puzz, targeting)
		if isSolved {
			cs.sendSolution(puzz)
		} else if canProcess {
			cs.processPayload(up)
		}
		return
	}

	// for each of the TwoArm options, we're going to try setting the edges
	// and then descending further into our targets
	for _, option := range targeting.Options {
		up, canProcess, isSolved := cs.addArms(puzz, targeting, option)
		if isSolved {
			cs.sendSolution(puzz)
		} else if !canProcess {
			continue
		}

		// Here, if we're the only option, then we're just going to descend on this
		// same goroutine. We won't queue up the payload because that would be wasteful
		// and we know we can just keep using this CPU
		if len(targeting.Options) == 1 {
			cs.processPayload(up)
		} else {
			cs.queuePayload(up)
		}
	}
}

func (cs *concurrentSolver) addArms(
	puzz puzzle.Puzzle,
	curTarget model.Target,
	ta model.TwoArms,
) (unsolvedPayload, bool, bool) {
	withArms, ok := addTwoArms(puzz, curTarget.Node.Coord(), ta)
	if !ok {
		return unsolvedPayload{}, false, false
	}
	return cs.getNextPayload(withArms, curTarget)
}

func (cs *concurrentSolver) getNextPayload(
	puzz puzzle.Puzzle,
	curTarget model.Target,
) (unsolvedPayload, bool, bool) {

	cs.printPuzzleUpdate(`getNextPayload`, puzz, curTarget)

	nextTarget, state := puzz.GetNextTarget(curTarget)
	switch state {
	case model.Complete:
		return unsolvedPayload{}, false, true
	case model.Incomplete, model.NodesComplete:
		return unsolvedPayload{
			puzz:            puzz,
			target:          nextTarget,
			isNodesComplete: state == model.NodesComplete,
			nextUnknown:     model.InvalidEdgePair,
		}, true, false
	}
	return unsolvedPayload{}, false, false
}

func (cs *concurrentSolver) flip(
	puzz puzzle.Puzzle,
	ep model.EdgePair,
) {

	cs.printPuzzleUpdate(`flip`, puzz, model.InvalidTarget)
	puzzle.SetNodesComplete(&puzz)

	if ep == model.InvalidEdgePair {
		var ok bool
		ep, ok = puzz.GetUnknownEdge()
		if !ok {
			switch puzz.GetState() {
			case model.Complete:
				cs.sendSolution(puzz)
			}
			return
		}
	}

	cs.justSetEdge(puzz, ep)
	cs.justAvoidEdge(puzz, ep)
}

func (cs *concurrentSolver) justSetEdge(
	base puzzle.Puzzle,
	ep model.EdgePair,
) {
	puzz, state := puzzle.AddEdge(
		base,
		ep,
	)
	switch state {
	case model.Complete, model.Incomplete:
		cs.checkStateAfterFlip(puzz, ep)
	}
}

func (cs *concurrentSolver) justAvoidEdge(
	base puzzle.Puzzle,
	ep model.EdgePair,
) {
	puzz, state := puzzle.AvoidEdge(
		base,
		ep,
	)
	switch state {
	case model.Complete, model.Incomplete:
		cs.checkStateAfterFlip(puzz, ep)
	}
}

func (cs *concurrentSolver) checkStateAfterFlip(
	puzz puzzle.Puzzle,
	ep model.EdgePair,
) {
	nextUnknown, state := puzz.GetStateOfLoop(ep.NodeCoord)
	switch state {
	case model.Complete:
		cs.sendSolution(puzz)
	case model.Incomplete:
		cs.queuePayload(unsolvedPayload{
			puzz:            puzz,
			isNodesComplete: true,
			nextUnknown:     nextUnknown,
		})
	}
}
