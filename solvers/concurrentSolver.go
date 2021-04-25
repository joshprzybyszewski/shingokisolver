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
	workChanLen = 16
)

var (
	numCPU = runtime.NumCPU()
)

type unsolved struct {
	puzz   puzzle.Puzzle
	target model.Target

	isNodesComplete bool
	nextUnknown     model.EdgePair
}

type targetPayload struct {
	puzz   puzzle.Puzzle
	target model.Target
}

type flippingPayload struct {
	puzz        puzzle.Puzzle
	nextUnknown model.EdgePair
}

type concurrentSolver struct {
	targets  chan targetPayload
	flipping chan flippingPayload
	solution chan puzzle.Puzzle

	// TODO I don't want to keep track of these numbers
	// because I don't think they boost my performance
	numTargetsAdded     int32
	numTargetsProcessed int32

	numFlipsAdded     int32
	numFlipsProcessed int32

	numImmediates int32
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
	cs.targets = make(chan targetPayload, workChanLen)
	cs.flipping = make(chan flippingPayload, workChanLen)
	cs.solution = make(chan puzzle.Puzzle, 1)

	defer close(cs.flipping)
	defer close(cs.targets)
	defer close(cs.solution)

	for i := 0; i < numCPU; i++ {
		go cs.startWorker()
	}

	cs.queuePayload(unsolved{
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
	payload unsolved,
) {

	if payload.isNodesComplete {
		cs.queueFlippingPayload(flippingPayload{
			puzz:        payload.puzz,
			nextUnknown: payload.nextUnknown,
		})
	} else {
		atomic.AddInt32(&cs.numTargetsAdded, 1)

		// it's alright if the targets channel has been closed
		defer func() {
			recover()
		}()
		cs.targets <- targetPayload{
			puzz:   payload.puzz,
			target: payload.target,
		}
	}
}

func (cs *concurrentSolver) queueFlippingPayload(
	payload flippingPayload,
) {

	atomic.AddInt32(&cs.numFlipsAdded, 1)

	// it's alright if the targets channel has been closed
	defer func() {
		recover()
	}()

	cs.flipping <- payload
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

func (cs *concurrentSolver) workOnPayloadNow(
	payload unsolved,
) {

	atomic.AddInt32(&cs.numImmediates, 1)
	// cs.printPayload(`workOnPayloadNow`, payload)

	if payload.isNodesComplete {
		cs.flip(payload.puzz, payload.nextUnknown)
	} else {
		cs.solveAimingAtTarget(payload.puzz, payload.target)
	}
}

func (cs *concurrentSolver) startWorker() {
	for {
		select {
		case tp, ok := <-cs.targets:
			if !ok {
				return
			}
			cs.processTargetPayload(tp)
		case fp, ok := <-cs.flipping:
			if !ok {
				return
			}
			cs.processFlipPayload(fp)
		}
	}
}

func (cs *concurrentSolver) processTargetPayload(
	payload targetPayload,
) {
	atomic.AddInt32(&cs.numTargetsProcessed, 1)

	cs.printPayload(`processTargetPayload`, payload)

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
			cs.workOnPayloadNow(up)
		}
		return
	}

	// for each of the TwoArm options, we're going to try setting the edges
	// and then descending further into our targets
	for i, option := range targeting.Options {
		up, canProcess, isSolved := cs.addArms(puzz, targeting, option)
		if isSolved {
			cs.sendSolution(puzz)
			return
		} else if !canProcess {
			continue
		}

		if i == len(targeting.Options)-1 {
			// This is the last option in our targeted options. Instead of
			// sending it for someone else to work on, we're just going to
			// continue using the CPU to see how it do.
			cs.workOnPayloadNow(up)
		} else {
			// we've built the puzzle for another worker to pick up whenever
			// it's free. Send it to them.
			go cs.queuePayload(up)
		}
	}
}

func (cs *concurrentSolver) addArms(
	puzz puzzle.Puzzle,
	curTarget model.Target,
	ta model.TwoArms,
) (unsolved, bool, bool) {
	withArms, ok := addTwoArms(puzz, curTarget.Node.Coord(), ta)
	if !ok {
		return unsolved{}, false, false
	}
	return cs.getNextPayload(withArms, curTarget)
}

func (cs *concurrentSolver) getNextPayload(
	puzz puzzle.Puzzle,
	curTarget model.Target,
) (unsolved, bool, bool) {

	cs.printPuzzleUpdate(`getNextPayload`, puzz, curTarget)

	nextTarget, state := puzz.GetNextTarget(curTarget)
	switch state {
	case model.Complete:
		return unsolved{}, false, true
	case model.Incomplete, model.NodesComplete:
		return unsolved{
			puzz:            puzz,
			target:          nextTarget,
			isNodesComplete: state == model.NodesComplete,
			nextUnknown:     model.InvalidEdgePair,
		}, true, false
	}
	return unsolved{}, false, false
}

func (cs *concurrentSolver) processFlipPayload(
	payload flippingPayload,
) {
	atomic.AddInt32(&cs.numFlipsProcessed, 1)

	cs.printFlippingPayload(`processFlipPayload`, payload)

	cs.flip(payload.puzz, payload.nextUnknown)
}

func (cs *concurrentSolver) flip(
	puzz puzzle.Puzzle,
	ep model.EdgePair,
) {

	puzzle.SetNodesComplete(&puzz)

	cs.printPuzzleUpdate(`flip`, puzz, model.InvalidTarget)

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

	fp1, canContinue, isSolved := cs.justSetEdge(puzz, ep)
	if isSolved {
		return
	}

	fp2, canContinue2, isSolved := cs.justAvoidEdge(puzz, ep)
	if isSolved {
		return
	}
	cs.printFlippingPayload(`finished the justs`, flippingPayload{puzz, ep})

	if canContinue && canContinue2 {
		go cs.queueFlippingPayload(fp1)
		cs.processFlipPayload(fp2)
	} else if canContinue {
		cs.processFlipPayload(fp1)
	} else if canContinue2 {
		cs.processFlipPayload(fp2)
	}
}

func (cs *concurrentSolver) justSetEdge(
	base puzzle.Puzzle,
	ep model.EdgePair,
) (flippingPayload, bool, bool) {
	cs.printFlippingPayload(`justSetEdge`, flippingPayload{base, ep})

	puzz, state := puzzle.AddEdge(
		base,
		ep,
	)
	switch state {
	case model.Complete, model.Incomplete:
		return cs.checkStateAfterFlip(puzz, ep)
	}
	return flippingPayload{}, false, false
}

func (cs *concurrentSolver) justAvoidEdge(
	base puzzle.Puzzle,
	ep model.EdgePair,
) (flippingPayload, bool, bool) {
	cs.printFlippingPayload(`justAvoidEdge`, flippingPayload{base, ep})

	puzz, state := puzzle.AvoidEdge(
		base,
		ep,
	)
	switch state {
	case model.Complete, model.Incomplete:
		return cs.checkStateAfterFlip(puzz, ep)
	}
	return flippingPayload{}, false, false
}

func (cs *concurrentSolver) checkStateAfterFlip(
	puzz puzzle.Puzzle,
	ep model.EdgePair,
) (flippingPayload, bool, bool) {
	nextUnknown, state := puzz.GetStateOfLoop(ep.NodeCoord)
	switch state {
	case model.Complete:
		cs.sendSolution(puzz)
		return flippingPayload{}, false, true
	case model.Incomplete:
		return flippingPayload{
			puzz:        puzz,
			nextUnknown: nextUnknown,
		}, true, false
	}
	return flippingPayload{}, false, false
}
