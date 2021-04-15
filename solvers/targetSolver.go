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

	puzz := d.puzzle.DeepCopy()

	switch s := puzz.ClaimGimmes(); s {
	case model.Incomplete, model.Complete:
		printPuzzleUpdate(`ClaimGimmes`, 0, puzz, targets[0], d.iterations())
	default:
		return nil, false
	}

	if puzz == nil {
		return nil, false
	}

	p := d.getSolutionFromDepths(
		0,
		puzz.DeepCopy(),
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

	switch puzz.GetNodeState(targeting.Coord) {
	case model.Violation:
		return nil

	case model.Complete:
		// the target node is already complete, perhaps a previous node
		// accidentally completed it. If so, then let's do a sanity check
		// on completion, and then add it as a "partial solution" that
		// has no new loose ends

		if targeting.Next == nil {
			switch puzz.GetState(targeting.Coord) {
			case model.Complete:
				return puzz
			default:
				printPuzzleUpdate(`getSolutionFromDepths did not solve!`, depth, puzz, targeting, d.iterations())
			}

			return d.flip(
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
	options := model.BuildTwoArmOptions(node, puzz.NumEdges())
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

func (d *targetSolver) buildAllTwoArmsForTraversal(
	depth int,
	puzz *puzzle.Puzzle,
	curTarget model.Target,
	ta model.TwoArms,
) *puzzle.Puzzle {

	switch puzz.AddEdges(getTwoArmsEdges(
		curTarget.Coord,
		ta,
	)...) {
	case model.Duplicate, model.Incomplete, model.Complete:
	default:
		return nil
	}

	switch puzz.GetState(curTarget.Coord) {
	case model.Complete:
		return puzz
	case model.Incomplete:
		// continue
	default:
		return nil
	}

	if curTarget.Next == nil {
		return d.flip(
			puzz.DeepCopy(),
		)
	}

	return d.getSolutionFromDepths(
		depth+1,
		puzz.DeepCopy(),
		*curTarget.Next,
	)
}

func getTwoArmsEdges(
	start model.NodeCoord,
	ta model.TwoArms,
) []model.EdgePair {

	allEdges := make([]model.EdgePair, 0, ta.One.Len+ta.Two.Len)

	arm1End := start
	for i := int8(0); i < ta.One.Len; i++ {
		allEdges = append(allEdges, model.NewEdgePair(arm1End, ta.One.Heading))
		arm1End = arm1End.Translate(ta.One.Heading)
	}

	arm2End := start
	for i := int8(0); i < ta.Two.Len; i++ {
		allEdges = append(allEdges, model.NewEdgePair(arm2End, ta.Two.Heading))
		arm2End = arm2End.Translate(ta.Two.Heading)
	}

	return allEdges
}
