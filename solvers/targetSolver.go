package solvers

import (
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
	targets := d.puzzle.Targets()

	if len(targets) == 0 {
		return nil, false
	}

	puzz, numProcessed := claimGimmes(d.puzzle.DeepCopy())
	d.numProcessed += numProcessed

	printPuzzleUpdate(`claimGimmes`, 0, puzz, targets[0], d.iterations())

	p := d.getSolutionFromDepths(
		0,
		puzz,
		targets[0],
	)
	return p, p != nil
}

func (d *targetSolver) getSolutionFromDepths(
	depth int,
	puzz *puzzle.Puzzle,
	targeting model.Target,
) *puzzle.Puzzle {

	printPuzzleUpdate(`getSolutionFromDepths`, depth, puzz, targeting, d.iterations())

	targetCoord := targeting.Coord

	node, ok := puzz.GetNode(targetCoord)
	if !ok {
		// this should be returning an error, but really it shouldn't be happening
		panic(`what?`)
		// return nil
	}

	if tCoord := targeting.Coord; puzzle.IsCompleteNode(puzz, tCoord) {
		// the target node is already complete, perhaps a previous node
		// accidentally completed it. If so, then let's do a sanity check
		// on completion, and then add it as a "partial solution" that
		// has no new loose ends

		if targeting.Next == nil {
			return d.getSolutionByConnectingLooseEnds(
				puzz.DeepCopy(),
			)
		}

		return d.getSolutionFromDepths(
			depth+1,
			puzz.DeepCopy(),
			*targeting.Next,
		)
	}

	// go out in all directions from the target
	// if it's still a valid puzzle, keep going outward
	// until we "complete" the node.
	options := model.BuildTwoArmOptions(node)
	for _, option := range options {
		// then, once we find a completion path, add it to the returned slice
		p := d.buildAllTwoArmsForTraversal(
			depth,
			puzz.DeepCopy(),
			targeting,
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
) *puzzle.Puzzle {

	printAllTargetsHit(`getSolutionByConnectingLooseEnds`, puzz, d.iterations())

	// this means we've iterated through all of the target nodes
	return d.connect(puzz)
}

func (d *targetSolver) buildAllTwoArmsForTraversal(
	depth int,
	puzz *puzzle.Puzzle,
	curTarget model.Target,
	ta model.TwoArms,
) *puzzle.Puzzle {

	twoArmPuzz := d.sendOutTwoArms(
		puzz.DeepCopy(),
		curTarget.Coord,
		ta,
	)

	if !puzzle.IsCompleteNode(twoArmPuzz, curTarget.Coord) {
		// we _should_ have added a number of straight edges that will
		// complete our target node.
		// if not, then we don't want to add this to our partials
		return nil
	}

	switch puzz.GetState() {
	case model.Complete:
		return twoArmPuzz
	case model.Incomplete:
		// continue
	default:
		return nil
	}

	if curTarget.Next == nil {
		return d.getSolutionByConnectingLooseEnds(
			twoArmPuzz.DeepCopy(),
		)
	}

	return d.getSolutionFromDepths(
		depth+1,
		twoArmPuzz,
		*curTarget.Next,
	)
}

func (d *targetSolver) sendOutTwoArms(
	puzz *puzzle.Puzzle,
	start model.NodeCoord,
	ta model.TwoArms,
) *puzzle.Puzzle {

	var state model.State

	arm1End := start
	for {
		if oe, inBounds := puzz.GetOutgoingEdgesFrom(start); !inBounds || oe.GetNumInDirection(ta.One.Heading) >= ta.One.Len {
			break
		}

		prevPuzz := puzz
		prevEnd := arm1End

		d.numProcessed++
		arm1End, state = puzz.AddEdge(ta.One.Heading, arm1End)
		switch state {
		case model.Duplicate:
			// if the edge already exists, let's allow it
			// and see if the puzzle will be valid
			puzz = prevPuzz
			arm1End = prevEnd.Translate(ta.One.Heading)
		case model.Incomplete, model.Complete:
			// these cases are "ok"
		default:
			// there was a problem. Early return
			return nil
		}
	}

	arm2End := start
	for {
		if oe, inBounds := puzz.GetOutgoingEdgesFrom(start); !inBounds || oe.GetNumInDirection(ta.Two.Heading) >= ta.Two.Len {
			break
		}

		prevPuzz := puzz
		prevEnd := arm2End

		d.numProcessed++
		arm2End, state = puzz.AddEdge(ta.Two.Heading, arm2End)
		switch state {
		case model.Duplicate:
			// if the edge already exists, let's allow it
			// and see if the puzzle will be valid
			puzz = prevPuzz
			arm2End = prevEnd.Translate(ta.Two.Heading)
		case model.Incomplete, model.Complete:
			// these cases are "ok"
		default:
			// there was a problem. Early return
			return nil
		}
	}

	return puzz
}
