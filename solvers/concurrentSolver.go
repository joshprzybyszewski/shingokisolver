package solvers

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

const (
	workChanLen = 32

	maxAdditionsAllowed = int32(1 << 20)

	maxAttemptDuration = time.Minute
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

	pendingTargets int32
	pendingFlips   int32
	finished       int32

	started time.Time
}

func (cs *concurrentSolver) solve(
	puzz puzzle.Puzzle,
) (sr SolvedResults, _ error) {
	if numCPU <= 1 {
		log.Printf("Unexpected! Running without concurrency (%d CPUs)", numCPU)
		return solvePuzzleByTargets(puzz)
	}

	cs.started = time.Now()

	defer func() {
		sr.Duration = time.Since(cs.started)
	}()

	puzz, target, state := claimGimmes(puzz)
	switch state {
	case model.Complete:
		return SolvedResults{
			Puzzle:     puzz,
			FinalState: model.Complete,
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

	ctx, cancelFn := context.WithTimeout(context.Background(), maxAttemptDuration)
	defer cancelFn()

	for i := 0; i < numCPU; i++ {
		go cs.startWorker(ctx)
	}

	cs.queuePayload(&unsolved{
		puzz:            puzz,
		target:          target,
		isNodesComplete: state == model.NodesComplete,
		nextUnknown:     model.InvalidEdgePair,
	})

	var sol puzzle.Puzzle
	ok := false
	select {
	case <-ctx.Done():
		//	three minutes is as long as any run these days..
	case sol, ok = <-cs.solution:
	}
	atomic.AddInt32(&cs.finished, 1)
	if !ok {
		log.Printf("solution channel closed!")
		log.Printf("\tpendingTargets:       %d",
			cs.pendingTargets,
		)
		log.Printf("\tpendingFlips:     %d",
			cs.pendingFlips,
		)
		return SolvedResults{}, fmt.Errorf("puzzle unsolvable: %s", puzz.String())
	}
	defer close(cs.solution)

	return SolvedResults{
		Puzzle:     sol,
		FinalState: model.Complete,
	}, nil
}

func (cs *concurrentSolver) isTargetsBackedUp() bool {
	return workChanLen < atomic.LoadInt32(&cs.pendingTargets)
}

func (cs *concurrentSolver) isFlipsBackedUp() bool {
	return workChanLen < atomic.LoadInt32(&cs.pendingFlips)
}

func (cs *concurrentSolver) queuePayload(payload *unsolved) {
	// Note: payload is only a pointer because it's like 300 bytes to copy it.

	if payload.isNodesComplete {
		if !cs.isFlipsBackedUp() {
			atomic.AddInt32(&cs.pendingFlips, 1)
			go cs.queueFlippingPayload(flippingPayload{
				puzz:        payload.puzz,
				nextUnknown: payload.nextUnknown,
			})
			return
		}
	} else {
		if !cs.isTargetsBackedUp() {
			atomic.AddInt32(&cs.pendingTargets, 1)
			go cs.queueTargetPayload(targetPayload{
				puzz:   payload.puzz,
				target: payload.target,
			})
			return
		}
	}
	if atomic.LoadInt32(&cs.finished) > 0 {
		return
	}

	cs.workOnPayloadNow(payload)
}

func (cs *concurrentSolver) queueTargetPayload(
	payload targetPayload,
) {

	// it's alright if the targets channel has been closed
	defer func() {
		_ = recover()
	}()

	cs.targets <- payload
}

func (cs *concurrentSolver) queueFlippingPayload(
	payload flippingPayload,
) {

	// it's alright if the targets channel has been closed
	defer func() {
		_ = recover()
	}()

	cs.flipping <- payload
}

func (cs *concurrentSolver) sendSolution(
	puzz puzzle.Puzzle,
) {

	defer func() {
		// it's alright if the solution channel has been closed
		_ = recover()
	}()

	cs.solution <- puzz
}

func (cs *concurrentSolver) workOnPayloadNow(
	payload *unsolved,
) {

	if payload.isNodesComplete {
		cs.flip(payload.puzz, payload.nextUnknown)
	} else {
		cs.solveAimingAtTarget(payload.puzz, payload.target)
	}
}

func (cs *concurrentSolver) startWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
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
	atomic.AddInt32(&cs.pendingTargets, -1)

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
		up, canProcess := cs.getNextTarget(puzz, targeting)
		if canProcess {
			cs.queuePayload(&up)
		}
		return
	}

	allUnsolved := make([]unsolved, 0, 8)
	// for each of the TwoArm options, we're going to try setting the edges
	// and then descending further into our targets
	for _, option := range targeting.Options {
		allUnsolved = append(allUnsolved,
			cs.addArms(puzz, targeting, option)...,
		)
	}

	for _, up := range allUnsolved {
		cs.queuePayload(&up)
	}
}

func (cs *concurrentSolver) addArms(
	puzz puzzle.Puzzle,
	curTarget model.Target,
	ta model.TwoArms,
) []unsolved {
	withArms, ok := addTwoArms(puzz, curTarget.Node.Coord(), ta)
	if !ok {
		return nil
	}

	return cs.getNextPayload(withArms, curTarget, ta)
}

func (cs *concurrentSolver) getNextPayload(
	puzz puzzle.Puzzle,
	curTarget model.Target,
	ta model.TwoArms,
) []unsolved {

	cs.printPuzzleUpdate(`getNextPayload`, puzz, curTarget)

	// arm 1
	endOf1 := curTarget.Node.Coord().TranslateAlongArm(ta.One)
	endOf2 := curTarget.Node.Coord().TranslateAlongArm(ta.Two)

	withArmAndPerps := make([]unsolved, 0, 4)

	for _, perpDir1 := range ta.One.Heading.Perpendiculars() {
		for _, perpDir2 := range ta.Two.Heading.Perpendiculars() {
			withPerps, ms := puzzle.AddEdges(puzz, []model.EdgePair{
				model.NewEdgePair(endOf1, perpDir1),
				model.NewEdgePair(endOf2, perpDir2),
			})

			switch ms {
			case model.Complete:
				cs.sendSolution(withPerps)
				return nil
			case model.NodesComplete:
				withArmAndPerps = append(withArmAndPerps,
					unsolved{
						puzz:            puzz,
						isNodesComplete: true,
						nextUnknown:     model.InvalidEdgePair,
					},
				)
			case model.Incomplete:
				u, ok := cs.getNextTarget(withPerps, curTarget)
				if ok {
					withArmAndPerps = append(withArmAndPerps, u)
				}
			}

		}
	}
	return withArmAndPerps
}

func (cs *concurrentSolver) getNextTarget(
	puzz puzzle.Puzzle,
	curTarget model.Target,
) (unsolved, bool) {

	nextTarget, state := puzz.GetNextTarget(curTarget)
	switch state {
	case model.Complete:
		cs.sendSolution(puzz)
		return unsolved{}, false
	case model.Incomplete, model.NodesComplete:
		return unsolved{
			puzz:            puzz,
			target:          nextTarget,
			isNodesComplete: state == model.NodesComplete,
			nextUnknown:     model.InvalidEdgePair,
		}, true
	}
	return unsolved{}, false
}

func (cs *concurrentSolver) processFlipPayload(
	payload flippingPayload,
) {
	atomic.AddInt32(&cs.pendingFlips, -1)

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
		return cs.checkStateAfterFlip(puzz)
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
		return cs.checkStateAfterFlip(puzz)
	}
	return flippingPayload{}, false, false
}

func (cs *concurrentSolver) checkStateAfterFlip(
	puzz puzzle.Puzzle,
) (flippingPayload, bool, bool) {
	nextUnknown, state := puzz.GetStateOfLoop()
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
