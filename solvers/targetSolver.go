package solvers

import (
	"fmt"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type targetSolver struct {
	puzzle *puzzle.Puzzle

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
		puzzle: puzzle.NewPuzzle(size, nl),
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

	puzz, looseEnds, numProcessed := claimGimmes(d.puzzle)
	d.numProcessed += numProcessed

	printPuzzleUpdate(`claimGimmes`, 0, puzz, targets[0], looseEnds, d.iterations())

	p := d.getSolutionFromDepths(
		0,
		puzz,
		targets[0],
		looseEnds,
	)
	return p, p != nil
}

func (d *targetSolver) getSolutionFromDepths(
	depth int,
	puzz *puzzle.Puzzle,
	targeting target,
	looseEnds []model.NodeCoord,
) *puzzle.Puzzle {

	printPuzzleUpdate(`getSolutionFromDepths`, depth, puzz, targeting, looseEnds, d.iterations())

	targetCoord := targeting.coord

	node, ok := puzz.NodeTargets()[targetCoord]
	if !ok {
		// this should be returning an error, but really it shouldn't be happening
		panic(`what?`)
		// return nil
	}

	if tCoord := targeting.coord; puzzle.IsCompleteNode(puzz, tCoord) {
		// the target node is already complete, perhaps a previous node
		// accidentally completed it. If so, then let's do a sanity check
		// on completion, and then add it as a "partial solution" that
		// has no new loose ends

		leCpy := make([]model.NodeCoord, len(looseEnds))
		copy(leCpy, looseEnds)

		if targeting.next == nil {
			return d.getSolutionByConnectingLooseEnds(
				puzz.DeepCopy(),
				leCpy,
			)
		}

		return d.getSolutionFromDepths(
			depth+1,
			puzz.DeepCopy(),
			*targeting.next,
			leCpy,
		)
	}

	// go out in all directions from the target
	// if it's still a valid puzzle, keep going outward
	// until we "complete" the node.
	options := buildTwoArmOptions(node)
	for _, option := range options {
		leCpy := make([]model.NodeCoord, len(looseEnds))
		copy(leCpy, looseEnds)

		// then, once we find a completion path, add it to the returned slice
		p := d.buildAllTwoArmsForTraversal(
			depth,
			puzz.DeepCopy(),
			targeting,
			leCpy,
			option,
		)
		if p != nil {
			return p
		}
	}

	return nil
}

func (d *targetSolver) getSolutionByConnectingLooseEnds(
	puzz *puzzle.Puzzle,
	looseEnds []model.NodeCoord,
) *puzzle.Puzzle {

	printAllTargetsHit(`getSolutionByConnectingLooseEnds`, puzz, looseEnds, d.iterations())

	// this means we've iterated through all of the target nodes
	return d.connect(puzz, looseEnds)
}

func (d *targetSolver) buildAllTwoArmsForTraversal(
	depth int,
	puzz *puzzle.Puzzle,
	curTarget target,
	looseEnds []model.NodeCoord,
	ta twoArms,
) *puzzle.Puzzle {

	printPuzzleUpdate(fmt.Sprintf(`buildAllTwoArmsForTraversal(%+v)`, ta), depth, puzz, curTarget, looseEnds, d.iterations())

	twoArmPuzz, arm1End, arm2End := d.sendOutTwoArms(
		puzz.DeepCopy(),
		curTarget.coord,
		ta,
	)

	if !puzzle.IsCompleteNode(twoArmPuzz, curTarget.coord) {
		printPuzzleUpdate(fmt.Sprintf(`!puzzle.IsCompleteNode(twoArmPuzz, %+v)`, curTarget.coord), depth, twoArmPuzz, curTarget, looseEnds, d.iterations())
		// we _should_ have added a number of straight edges that will
		// complete our target node.
		// if not, then we don't want to add this to our partials
		return nil
	}

	if isIncomplete, err := twoArmPuzz.IsIncomplete(arm1End); err != nil {
		printPuzzleUpdate(fmt.Sprintf(`twoArmPuzz.IsIncomplete(%+v) = %v`, curTarget.coord, err), depth, twoArmPuzz, curTarget, looseEnds, d.iterations())
		return nil
	} else if !isIncomplete {
		printPuzzleUpdate(`!isIncomplete`, depth, twoArmPuzz, curTarget, looseEnds, d.iterations())
		return twoArmPuzz
	}

	newLooseEnds := getNewLooseEndsForBranches(looseEnds, curTarget.coord, arm1End, arm2End)

	if curTarget.next == nil {
		return d.getSolutionByConnectingLooseEnds(
			twoArmPuzz.DeepCopy(),
			newLooseEnds,
		)
	}

	return d.getSolutionFromDepths(
		depth+1,
		twoArmPuzz,
		*curTarget.next,
		newLooseEnds,
	)
}

func (d *targetSolver) sendOutTwoArms(
	puzz *puzzle.Puzzle,
	start model.NodeCoord,
	ta twoArms,
) (*puzzle.Puzzle, model.NodeCoord, model.NodeCoord) {

	var err error

	arm1End := start
	for {
		if oe, inBounds := puzz.GetOutgoingEdgesFrom(start); !inBounds || oe.GetNumInDirection(ta.arm1Heading) >= ta.arm1Len {
			break
		}

		prevPuzz := puzz
		prevEnd := arm1End

		d.numProcessed++
		arm1End, puzz, err = puzz.AddEdge(ta.arm1Heading, arm1End)
		if err != nil {
			if err != puzzle.ErrEdgeAlreadyExists {
				return nil, model.NodeCoord{}, model.NodeCoord{}
			}
			// if the edge already exists, let's allow it
			// and see if the puzzle will be valid
			puzz = prevPuzz
			arm1End = prevEnd.Translate(ta.arm1Heading)
		}
	}

	arm2End := start
	for {
		if oe, inBounds := puzz.GetOutgoingEdgesFrom(start); !inBounds || oe.GetNumInDirection(ta.arm2Heading) >= ta.arm2Len {
			break
		}

		prevPuzz := puzz
		prevEnd := arm2End

		d.numProcessed++
		arm2End, puzz, err = puzz.AddEdge(ta.arm2Heading, arm2End)
		if err != nil {
			if err != puzzle.ErrEdgeAlreadyExists {
				return nil, model.NodeCoord{}, model.NodeCoord{}
			}
			// if the edge already exists, let's allow it
			// and see if the puzzle will be valid
			puzz = prevPuzz
			arm2End = prevEnd.Translate(ta.arm2Heading)
		}
	}

	return puzz, arm1End, arm2End
}
