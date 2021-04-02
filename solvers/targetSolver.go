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
	targets := buildTargets(d.puzzle)

	if len(targets) == 0 {
		return nil, false
	}

	puzz, looseEnds, numProcessed := claimGimmes(d.puzzle)
	d.numProcessed += numProcessed

	p := d.getSolutionFromDepths(
		&partialSolutionItem{
			puzzle:    puzz,
			targeting: targets[0],
			looseEnds: looseEnds,
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
		return lec.connect(item)
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
				item.targeting,
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
	t *target,
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

		leCpy := make([]model.NodeCoord, len(item.looseEnds))
		copy(leCpy, item.looseEnds)
		return d.getSolutionFromDepths(&partialSolutionItem{
			puzzle:    item.puzzle.DeepCopy(),
			targeting: item.targeting.next,
			looseEnds: leCpy,
		})
	}

	for oneArmLen := int8(1); oneArmLen < t.node.Value(); oneArmLen++ {

		printPartialSolution(`buildAllTwoArmsForTraversal`, item, d.iterations())

		twoArmPuzz, arm1End, arm2End := d.sendOutTwoArms(
			item.puzzle.DeepCopy(),
			t.coord,
			arm1Heading,
			oneArmLen,
			arm2Heading,
			t.node.Value()-oneArmLen,
		)

		if !puzzle.IsCompleteNode(twoArmPuzz, t.coord) {
			// we _should_ have added a number of straight edges that will
			// complete our target node.
			// if not, then we don't want to add this to our partials
			continue
		}

		// if isComplete, err := twoArmPuzz.IsComplete(t.coord); err != nil {
		// 	continue
		// } else if isComplete {
		// 	return twoArmPuzz
		// }

		if isIncomplete, err := twoArmPuzz.IsIncomplete(t.coord); err != nil {
			continue
		} else if !isIncomplete {
			return twoArmPuzz
		}

		// TODO figure out what's busted with adding perps
		// p := d.buildPerpendicularsOnArms(
		// 	twoArmPuzz,
		// 	arm1End, arm2End,
		// 	arm1Heading, arm2Heading,
		// 	item.targeting.next,
		// 	item.looseEnds,
		// )

		p := d.getSolutionFromDepths(&partialSolutionItem{
			puzzle:    twoArmPuzz,
			targeting: item.targeting.next,
			looseEnds: getLooseEnds(item.looseEnds, t.coord, arm1End, arm2End),
		})
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

		prevPuzz := puzz
		prevEnd := arm1End

		d.numProcessed++
		arm1End, puzz, err = puzz.AddEdge(arm1Heading, arm1End)
		if err != nil {
			if err != puzzle.ErrEdgeAlreadyExists {
				return nil, model.NodeCoord{}, model.NodeCoord{}
			}
			// if the edge already exists, let's allow it
			// and see if the puzzle will be valid
			puzz = prevPuzz
			arm1End = prevEnd.Translate(arm1Heading)
		}
	}

	arm2End := start
	for {
		if oe, inBounds := puzz.GetOutgoingEdgesFrom(start); !inBounds || oe.GetNumInDirection(arm2Heading) >= arm2Length {
			break
		}

		prevPuzz := puzz
		prevEnd := arm2End

		d.numProcessed++
		arm2End, puzz, err = puzz.AddEdge(arm2Heading, arm2End)
		if err != nil {
			if err != puzzle.ErrEdgeAlreadyExists {
				return nil, model.NodeCoord{}, model.NodeCoord{}
			}
			// if the edge already exists, let's allow it
			// and see if the puzzle will be valid
			puzz = prevPuzz
			arm2End = prevEnd.Translate(arm2Heading)
		}
	}

	return puzz, arm1End, arm2End
}

func (d *targetSolver) buildPerpendicularsOnArms(
	twoArmPuzz *puzzle.Puzzle,
	arm1Coord, arm2Coord model.NodeCoord,
	arm1Heading, arm2Heading model.Cardinal,
	nextTarget *target,
	prevLooseEnds []model.NodeCoord,
) *puzzle.Puzzle {
	// TODO if our new arm hits a node, then we should just add edges so that node is satisfied...

	// Each of our arms needs a branch that's perpendicular to where it ended.
	for _, fromArm1 := range model.Perpendiculars(arm1Heading) {
		for _, fromArm2 := range model.Perpendiculars(arm2Heading) {

			d.numProcessed++
			arm1ext, twoArmPuzzWithOneBranch, err := twoArmPuzz.DeepCopy().AddEdge(fromArm1, arm1Coord)
			if err != nil {
				if err != puzzle.ErrEdgeAlreadyExists {
					continue
				}
				arm1ext = arm1Coord
				twoArmPuzzWithOneBranch = twoArmPuzz
			}

			if isIncomplete, err := twoArmPuzzWithOneBranch.IsIncomplete(arm1ext); err != nil {
				continue
			} else if !isIncomplete {
				return twoArmPuzzWithOneBranch
			}

			d.numProcessed++
			arm2ext, twoArmPuzzWithTwoBranches, err := twoArmPuzzWithOneBranch.DeepCopy().AddEdge(fromArm2, arm2Coord)
			if err != nil {
				if err != puzzle.ErrEdgeAlreadyExists {
					continue
				}
				arm2ext = arm2Coord
				twoArmPuzzWithTwoBranches = twoArmPuzzWithOneBranch
			}

			if isIncomplete, err := twoArmPuzzWithTwoBranches.IsIncomplete(arm2ext); err != nil {
				continue
			} else if !isIncomplete {
				return twoArmPuzzWithTwoBranches
			}

			looseEnds := make([]model.NodeCoord, len(prevLooseEnds))
			copy(looseEnds, prevLooseEnds)
			p := d.getSolutionFromDepths(&partialSolutionItem{
				puzzle:    twoArmPuzzWithTwoBranches,
				targeting: nextTarget,
				looseEnds: append(looseEnds, arm1ext, arm2ext),
			})
			if p != nil {
				return p
			}
		}
	}
	return nil
}
