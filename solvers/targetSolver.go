package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type targetSolver struct {
	puzzle *puzzle.Puzzle

	queue             *partialSolutionQueue
	looseEndConnector *looseEndConnector

	numProcessed int
}

func newTargetSolver(
	size int,
	nl []model.NodeLocation,
) solver {
	if len(nl) == 0 {
		return nil
	}

	return &targetSolver{
		puzzle:            puzzle.NewPuzzle(size, nl),
		queue:             newPartialSolutionQueue(),
		looseEndConnector: &looseEndConnector{},
	}
}

func (d *targetSolver) iterations() int {
	return d.numProcessed
}

func (d *targetSolver) solve() (*puzzle.Puzzle, bool) {
	targets := buildTargets(d.puzzle)

	if len(targets) == 0 {
		return nil, false
	}

	p := d.getSolutionFromDepths(
		&partialSolutionItem{
			puzzle:    d.puzzle.DeepCopy(),
			targeting: targets[0],
		},
	)
	return p, p != nil
}

func (d *targetSolver) getSolutionFromDepths(
	item *partialSolutionItem,
) *puzzle.Puzzle {
	printPartialSolution(`getSolutionFromDepths`, item, d.iterations())

	if item == nil {
		return nil
	}

	if item.targeting == nil {
		item.removeDuplicateLooseEnds()
		printPartialSolution(`getSolutionFromDepths nil target`, item, d.iterations())
		// this means we've iterated through all of the target nodes
		lec := &looseEndConnector{}
		defer func() {
			d.numProcessed += lec.iterations
		}()
		return lec.queueUpLooseEndConnections(item)
	}

	targetCoord := item.targeting.coord

	node, ok := item.puzzle.NodeTargets()[targetCoord]
	if !ok {
		// this should be returning an error, but really it shouldn't be happening
		panic(`what?`)
		// return nil
	}

	// go out in all directions from the target
	// if it's still a valid puzzle, keep going outward
	// until we "complete" the node.
	for i, feeler1 := range model.AllCardinals {
		for _, feeler2 := range model.AllCardinals[i+1:] {
			if node.IsInvalidMotions(feeler1, feeler2) {
				continue
			}

			// then, once we find a completion path, add it to the returned slice
			p := d.buildAllTwoArmsForTraversal(
				item,
				*item.targeting,
				feeler1, feeler2,
			)
			if p != nil {
				return p
			}
		}
	}

	return nil
}

func (d *targetSolver) buildAllTwoArmsForTraversal(
	item *partialSolutionItem,
	t target,
	arm1Heading, arm2Heading model.Cardinal,
) *puzzle.Puzzle {

	if puzzle.IsCompleteNode(item.puzzle, t.coord) {
		// the target node is already complete, perhaps a previous node
		// accidentally completed it. If so, then let's do a sanity check
		// on completion, and then add it as a "partial solution" that
		// has no new loose ends

		oe, inBounds := item.puzzle.GetOutgoingEdgesFrom(t.coord)
		if !inBounds ||
			oe.GetNumInDirection(arm1Heading) == 0 ||
			oe.GetNumInDirection(arm2Heading) == 0 {
			// The node at the coord t.coord is complete, but it doesn't
			// have headings in the directions that we want it to. In order
			// to avoid having duplicate puzzle states, we don't return
			// anything here.
			return nil
		}

		item.targeting = item.targeting.next
		return d.getSolutionFromDepths(item)
	}

	for oneArmLen := int8(1); oneArmLen < t.val; oneArmLen++ {

		printPartialSolution(`buildAllTwoArmsForTraversal`, item, d.iterations())

		twoArmPuzz, arm1End, arm2End := d.sendOutTwoArms(
			item.puzzle,
			t.coord,
			arm1Heading,
			oneArmLen,
			arm2Heading,
			t.val-oneArmLen,
		)

		if !puzzle.IsCompleteNode(twoArmPuzz, t.coord) {
			// we _should_ have added a number of straight edges that will
			// complete our target node.
			// if not, then we don't want to add this to our partials
			continue
		}

		if _, err := twoArmPuzz.IsIncomplete(t.coord); err != nil {
			continue
		}

		psi := &partialSolutionItem{
			puzzle:    twoArmPuzz,
			targeting: item.targeting.next,
			looseEnds: make([]model.NodeCoord, len(item.looseEnds)),
		}
		copy(psi.looseEnds, item.looseEnds)
		psi.looseEnds = append(psi.looseEnds, arm1End, arm2End)
		p := d.getSolutionFromDepths(psi)
		if p != nil {
			return p
		}
	}

	return nil
}

func (d *targetSolver) sendOutTwoArms(
	puzz *puzzle.Puzzle,
	start model.NodeCoord,
	arm1Heading model.Cardinal,
	arm1Length int8,
	arm2Heading model.Cardinal,
	arm2Length int8,
) (*puzzle.Puzzle, model.NodeCoord, model.NodeCoord) {

	var err error

	arm1End := start
	for {
		if oe, inBounds := puzz.GetOutgoingEdgesFrom(start); !inBounds || oe.GetNumInDirection(arm1Heading) >= arm1Length {
			break
		}

		d.numProcessed++
		arm1End, puzz, err = puzz.AddEdge(arm1Heading, arm1End)
		if err != nil {
			return nil, model.NodeCoord{}, model.NodeCoord{}
		}
	}

	arm2End := start
	for {
		if oe, inBounds := puzz.GetOutgoingEdgesFrom(start); !inBounds || oe.GetNumInDirection(arm2Heading) >= arm2Length {
			break
		}

		d.numProcessed++
		arm2End, puzz, err = puzz.AddEdge(arm2Heading, arm2End)
		if err != nil {
			return nil, model.NodeCoord{}, model.NodeCoord{}
		}
	}

	return puzz, arm1End, arm2End
}

type looseEndConnector struct {
	iterations int
}

func (lec *looseEndConnector) queueUpLooseEndConnections(
	psi *partialSolutionItem,
) *puzzle.Puzzle {
	printPartialSolution(`queueUpLooseEndConnections`, psi, lec.iterations)

	psq := newPartialSolutionQueue()
	psq.push(psi)

	for !psq.isEmpty() {
		partial := psq.pop()

		puzz := lec.connectLooseEnds(psq, partial)
		if puzz != nil {
			return puzz
		}
	}

	return nil
}

func (lec *looseEndConnector) connectLooseEnds(
	psq *partialSolutionQueue,
	partial *partialSolutionItem,
) *puzzle.Puzzle {
	printPartialSolution(`connectLooseEnds`, partial, lec.iterations)

	numLooseEndsByNode := make(map[model.NodeCoord]int, len(partial.looseEnds))
	for _, g := range partial.looseEnds {
		numLooseEndsByNode[g] += 1
	}

	for i, start := range partial.looseEnds {
		if numLooseEndsByNode[start] > 1 {
			// This means that two nodes have a "loose end" that meets
			// up at the same point. Let's just skip trying to find this
			// "loose" end a buddy since it already has one.
			continue
		}

		dfs := newDFSSolverForPartialSolution()
		p, morePartials, sol := dfs.solveForGoals(
			partial.puzzle,
			start,
			partial.looseEnds[i+1:],
		)
		lec.iterations += dfs.iterations()

		switch sol {
		case solvedPuzzle:
			return p
		}

		lec.pushMorePartialSolutions(psq, partial, start, morePartials)
	}

	return nil
}

func (lec *looseEndConnector) pushMorePartialSolutions(
	psq *partialSolutionQueue,
	partial *partialSolutionItem,
	start model.NodeCoord,
	morePartials map[model.NodeCoord][]*puzzle.Puzzle,
) {
	for hitGoal, slice := range morePartials {
		for _, nextPuzzle := range slice {
			newLooseEnds, foundBoth := copyWithout(
				partial.looseEnds,
				start, hitGoal,
			)
			if !foundBoth {
				panic(`bad time`)
			}
			psq.push(&partialSolutionItem{
				puzzle:    nextPuzzle,
				looseEnds: newLooseEnds,
			})
		}
	}
}

func copyWithout(orig []model.NodeCoord, exclude1, exclude2 model.NodeCoord) ([]model.NodeCoord, bool) {
	newLooseEnds := make([]model.NodeCoord, 0, len(orig))
	found1, found2 := false, false
	for _, le := range orig {
		switch le {
		case exclude1:
			if found1 {
				return nil, false // bad state
			}
			found1 = true
		case exclude2:
			if found2 {
				return nil, false // bad state
			}
			found2 = true
		default:
			newLooseEnds = append(newLooseEnds, le)
		}
	}
	return newLooseEnds, found1 && found2
}
