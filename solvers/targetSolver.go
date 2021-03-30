package solvers

import (
	"sort"

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

	d.buildAllPartials(targets)

	p := d.getFullSolutionFromLooseEndConnector()

	return p, p != nil
}

func (d *targetSolver) buildAllPartials(
	targets []*target,
) {

	for i, target := range targets {
		if i > 0 {
			break
		}
		d.queue.push(d.getAllPartialSolutionsForItem(
			&partialSolutionItem{
				puzzle:    d.puzzle.DeepCopy(),
				targeting: target,
			},
		)...)
	}

	for !d.queue.isEmpty() {
		item := d.queue.pop()
		if item.targeting == nil {
			// Since we're not targeting a particular node, it means
			// we've completed them all. Let's add this as a "completed"
			// partial solution to our loose end connector.
			d.looseEndConnector.addPartialSolutions(item)
			continue
		}

		partials := d.getAllPartialSolutionsForItem(item)

		d.queue.push(partials...)
	}
}

func (d *targetSolver) getAllPartialSolutionsForItem(
	item *partialSolutionItem,
) (partials []*partialSolutionItem) {

	if item == nil || item.targeting == nil {
		return nil
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
			partials = append(partials, d.getPartialSolutionsForTwoArmedTarget(
				item,
				*item.targeting,
				feeler1, feeler2,
			)...)
		}
	}

	return partials
}

func (d *targetSolver) getPartialSolutionsForTwoArmedTarget(
	item *partialSolutionItem,
	t target,
	arm1Heading, arm2Heading model.Cardinal,
) (partials []*partialSolutionItem) {

	if puzzle.IsCompleteNode(item.puzzle, t.coord) {
		// the target node is already complete, perhaps a previous node
		// accidentally completed it. If so, then let's do a sanity check
		// on completion, and then add it as a "partial solution" that
		// has no new loose ends
		if _, err := item.puzzle.IsIncomplete(t.coord); err != nil {
			return nil
		}

		d.numProcessed++
		psi := &partialSolutionItem{
			puzzle:    item.puzzle.DeepCopy(),
			targeting: item.targeting.next,
			looseEnds: make([]model.NodeCoord, len(item.looseEnds)),
		}
		copy(psi.looseEnds, item.looseEnds)
		// once we find a completion path, add it to the returned slice
		partials = append(partials, psi)
		printPartialSolution(`IsCompleteNode`, psi, d.iterations())
		return partials
	}

	arm1End := t.coord
	arm2End := t.coord

	var err error
	oneArmPuzz := item.puzzle

	for numArms1 := 1; numArms1 < t.val; numArms1++ {
		d.numProcessed++
		arm1End, oneArmPuzz, err = oneArmPuzz.AddEdge(arm1Heading, arm1End)
		if err != nil {
			break
		}

		// reset the "end" of the second edge to be at the target
		arm2End = t.coord
		twoArmPuzz := oneArmPuzz
		for numArms2 := 1; numArms1+numArms2 <= t.val; numArms2++ {
			d.numProcessed++
			arm2End, twoArmPuzz, err = twoArmPuzz.AddEdge(arm2Heading, arm2End)
			if err != nil || twoArmPuzz == nil {
				break
			}
		}

		if twoArmPuzz == nil {
			continue
		}

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

		// once we find a completion path, add it to the returned slice
		partials = append(partials, psi)
		printPartialSolution(`end of arm building`, psi, d.iterations())
	}

	return partials
}

func (d *targetSolver) getFullSolutionFromLooseEndConnector() *puzzle.Puzzle {

	defer func() {
		d.numProcessed += d.looseEndConnector.iterations
	}()

	return d.looseEndConnector.solve()
}

type looseEndConnector struct {
	partials   []*partialSolutionItem
	iterations int
}

func (lec *looseEndConnector) addPartialSolutions(
	partials ...*partialSolutionItem,
) {
	lec.partials = append(lec.partials, partials...)
}

func (lec *looseEndConnector) solve() *puzzle.Puzzle {
	lec.prepPartials()

	return lec.attemptAllPartials()
}

func (lec *looseEndConnector) prepPartials() {
	// remove duplicate loose ends from all of our solutions
	for _, partial := range lec.partials {
		partial.removeDuplicateLooseEnds()
	}

	// sort the partial solutions so that the solutions with the most
	// connections (least number of loose ends) are at the front
	sort.Slice(lec.partials, func(i, j int) bool {
		return len(lec.partials[i].looseEnds) < len(lec.partials[j].looseEnds)
	})
}

func (lec *looseEndConnector) attemptAllPartials() *puzzle.Puzzle {
	// iterate through the partial solutions, trying to connect all of their
	// loose ends.
	attemptedCache := newPuzzleCache()
	for _, partial := range lec.partials {
		if attemptedCache.contains(partial.puzzle) {
			continue
		}

		p := lec.queueUpLooseEndConnections(partial)
		if p != nil {
			return p
		}

		attemptedCache.add(partial.puzzle)
	}
	return nil
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
	printPartialSolution(`connectLooseEnds`, partial, 7)

	for i, start := range partial.looseEnds {
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
