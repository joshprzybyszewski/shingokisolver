package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type targetSolver struct {
	puzzle *puzzle.Puzzle

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

func getLooseEnds(
	prev []model.NodeCoord,
	start, arm1End, arm2End model.NodeCoord,
) []model.NodeCoord {
	ends := make([]model.NodeCoord, len(prev), len(prev)+2)
	copy(ends, prev)

	ends = append(ends, arm1End, arm2End)

	isBetween := func(val, inclusive, exclusive int) bool {
		if inclusive < exclusive {
			return val >= inclusive && val < exclusive
		}
		return val <= inclusive && val > exclusive
	}

	endPoints := []model.NodeCoord{arm1End, arm2End}
	shouldRemove := func(nc model.NodeCoord) bool {
		sameRow := nc.Row == start.Row
		sameCol := nc.Col == start.Col
		if !sameRow && !sameCol {
			return false
		}

		if sameRow {
			if sameCol {
				// this looseEnd matches our start node
				return true
			}

			for _, end := range endPoints {
				if end.Row == start.Row && isBetween(int(nc.Col), int(start.Col), int(end.Col)) {
					return true
				}
			}
			return false
		}

		for _, end := range endPoints {
			if end.Col == start.Col && isBetween(int(nc.Row), int(start.Row), int(end.Row)) {
				return true
			}
		}

		return false
	}

	for i := 0; i < len(ends); i++ {
		if shouldRemove(ends[i]) {
			ends = append(ends[:i], ends[i+1:]...)
			i--
		}
	}

	return ends
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
			arm1ext, twoArmPuzzWithOneBranch, err := twoArmPuzz.AddEdge(fromArm1, arm1Coord)
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
			arm2ext, twoArmPuzzWithTwoBranches, err := twoArmPuzzWithOneBranch.AddEdge(fromArm2, arm2Coord)
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

type looseEndConnector struct {
	iterations int
}

func (lec *looseEndConnector) connect(
	partial *partialSolutionItem,
) *puzzle.Puzzle {
	printPartialSolution(`connect`, partial, lec.iterations)

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

		// we only need to look at the first loose end in the
		// puzzle, so we return the following.
		return lec.iterateMorePartials(partial, start, morePartials)
	}

	return nil
}

func (lec *looseEndConnector) iterateMorePartials(
	partial *partialSolutionItem,
	start model.NodeCoord,
	morePartials map[model.NodeCoord][]*puzzle.Puzzle,
) *puzzle.Puzzle {
	looseEndsWithoutStart := remove(partial.looseEnds, start)

	for hitGoal, slice := range morePartials {
		for _, nextPuzzle := range slice {
			puzz := lec.connect(&partialSolutionItem{
				puzzle:    nextPuzzle,
				looseEnds: remove(looseEndsWithoutStart, hitGoal),
			})
			if puzz != nil {
				return puzz
			}
		}
	}
	return nil
}

func remove(orig []model.NodeCoord, exclude model.NodeCoord) []model.NodeCoord {
	cpy := make([]model.NodeCoord, 0, len(orig)-1)
	for _, le := range orig {
		if le != exclude {
			cpy = append(cpy, le)
		}
	}
	return cpy
}
