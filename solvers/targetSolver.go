package solvers

import (
	"fmt"
	"log"
	"sort"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

type target struct {
	coord model.NodeCoord
	val   int

	next *target
}

func buildTargets(p *puzzle.Puzzle) []*target {
	targets := make([]*target, 0, len(p.NodeTargets()))

	for nc, n := range p.NodeTargets() {
		targets = append(targets, &target{
			coord: nc,
			val:   int(n.Value()),
		})
	}

	maxRowColVal := p.NumEdges()
	isOnTheSide := func(coord model.NodeCoord) bool {
		return coord.Row == 0 ||
			coord.Row == model.RowIndex(maxRowColVal) ||
			coord.Col == 0 ||
			coord.Col == model.ColIndex(maxRowColVal)
	}
	sort.Slice(targets, func(i, j int) bool {
		// rank higher valued nodes at the start of the target list
		if targets[i].val != targets[j].val {
			return targets[i].val > targets[j].val
		}

		// put nodes with more limitations (i.e. on the sides of
		// of the graph) higher up on the list
		iIsEdge := isOnTheSide(targets[i].coord)
		jIsEdge := isOnTheSide(targets[j].coord)
		if iIsEdge && !jIsEdge {
			return true
		} else if jIsEdge && !iIsEdge {
			return false
		}

		// at this point, we just want a consistent ordering.
		// let's put nodes closer to (0,0) higher up in the list
		if targets[i].coord.Row != targets[j].coord.Row {
			return targets[i].coord.Row < targets[j].coord.Row
		}
		return targets[i].coord.Col < targets[j].coord.Col
	})

	for i := 1; i < len(targets); i++ {
		targets[i-1].next = targets[i]
	}

	return targets
}

type partialSolutionQueue struct {
	items []*partialSolutionItem
}

func newPartialSolutionQueue() *partialSolutionQueue {
	return &partialSolutionQueue{
		items: make([]*partialSolutionItem, 0),
	}
}

func (q *partialSolutionQueue) isEmpty() bool {
	return len(q.items) == 0
}

func (q *partialSolutionQueue) pop() *partialSolutionItem {
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *partialSolutionQueue) push(items ...*partialSolutionItem) {
	// TODO for each item, if it already exists in the queue, don't
	// add it a second time? But how is it non-deterministic? How
	// can we generate the same partial solution?
	q.items = append(q.items, items...)
}

type partialSolutionItem struct {
	puzzle *puzzle.Puzzle

	targeting *target

	looseEnds []model.NodeCoord
}

// eliminates loose ends that don't actually exist
func (partial *partialSolutionItem) removeDuplicateLooseEnds() {
	sort.Slice(partial.looseEnds, func(i, j int) bool {
		if partial.looseEnds[i].Row != partial.looseEnds[j].Row {
			return partial.looseEnds[i].Row < partial.looseEnds[j].Row
		}
		return partial.looseEnds[i].Col < partial.looseEnds[j].Col
	})

	for i := 0; i < len(partial.looseEnds)-1; i++ {
		if partial.looseEnds[i] == partial.looseEnds[i+1] {
			partial.looseEnds = append(
				partial.looseEnds[:i],
				partial.looseEnds[i+2:]...)
			i--
		}
	}
}

func printPartialSolution(
	partial *partialSolutionItem,
	iterations int,
) {
	if !includeProgressLogs {
		return
	}
	if partial.puzzle == nil ||
		partial.targeting == nil {
		return
	}

	log.Printf("printPartialSolution (%d iterations): (targeting %v) %v",
		iterations,
		partial.targeting,
		partial.looseEnds,
	)
	log.Printf(":\n%s\n", partial.puzzle.String())
	fmt.Scanf("hello there")
}

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
		d.queue.push(d.getPartialSolutionsForNode(
			&partialSolutionItem{
				puzzle:    d.puzzle.DeepCopy(),
				targeting: target,
			},
		)...)
	}

	for !d.queue.isEmpty() {
		item := d.queue.pop()
		if item.targeting == nil {
			d.looseEndConnector.addPartialSolutions(item)
			continue
		}

		partials := d.getPartialSolutionsForNode(item)

		d.queue.push(partials...)
	}
}

func (d *targetSolver) getPartialSolutionsForNode(
	prevPartial *partialSolutionItem,
) (partials []*partialSolutionItem) {

	if prevPartial == nil || prevPartial.targeting == nil {
		return nil
	}
	targetCoord := prevPartial.targeting.coord

	node, ok := prevPartial.puzzle.NodeTargets()[targetCoord]
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
			partials = append(partials, d.getSolutionsForNodeInDirections(
				prevPartial,
				targetCoord,
				node,
				feeler1, feeler2,
			)...)
		}
	}

	return partials
}

func (d *targetSolver) getSolutionsForNodeInDirections(
	prevItem *partialSolutionItem,
	targetCoord model.NodeCoord,
	node model.Node,
	arm1Heading, arm2Heading model.Cardinal,
) (partials []*partialSolutionItem) {

	if puzzle.IsCompleteNode(prevItem.puzzle, targetCoord) {
		// the target node is already complete, perhaps a previous node
		// accidentally completed it. If so, then let's do a sanity check
		// on completion, and then add it as a "partial solution" that
		// has no new loose ends
		if _, err := prevItem.puzzle.IsIncomplete(targetCoord); err != nil {
			return nil
		}

		d.numProcessed++
		item := &partialSolutionItem{
			puzzle:    prevItem.puzzle.DeepCopy(),
			targeting: prevItem.targeting.next,
			looseEnds: make([]model.NodeCoord, len(prevItem.looseEnds)),
		}
		copy(item.looseEnds, prevItem.looseEnds)
		// once we find a completion path, add it to the returned slice
		partials = append(partials, item)
		printPartialSolution(item, d.iterations())
		return partials
	}

	arm1End := targetCoord
	arm2End := targetCoord

	var err error
	oneArmPuzz := prevItem.puzzle

	for i := 0; i <= int(node.Value()); i++ {
		if i > 0 {
			d.numProcessed++
			arm1End, oneArmPuzz, err = oneArmPuzz.AddEdge(arm1Heading, arm1End)
			if err != nil {
				break
			}
		}

		// reset the "end" of the second edge to be at the target
		arm2End = targetCoord
		twoArmPuzz := oneArmPuzz
		for j := 1; i+j <= int(node.Value()); j++ {
			d.numProcessed++
			arm2End, twoArmPuzz, err = twoArmPuzz.AddEdge(arm2Heading, arm2End)
			if err != nil || twoArmPuzz == nil {
				break
			}
		}

		if twoArmPuzz == nil {
			continue
		}

		if !puzzle.IsCompleteNode(twoArmPuzz, targetCoord) {
			// we _should_ have added a number of straight edges that will
			// complete our target node.
			// if not, then we don't want to add this to our partials
			continue
		}

		if _, err := twoArmPuzz.IsIncomplete(targetCoord); err != nil {
			continue
		}

		item := &partialSolutionItem{
			puzzle:    twoArmPuzz,
			targeting: prevItem.targeting.next,
			looseEnds: make([]model.NodeCoord, len(prevItem.looseEnds)),
		}
		copy(item.looseEnds, prevItem.looseEnds)
		item.looseEnds = append(item.looseEnds, arm1End, arm2End)

		// once we find a completion path, add it to the returned slice
		partials = append(partials, item)
		printPartialSolution(item, d.iterations())
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
	printPartialSolution(psi, lec.iterations)

	psq := newPartialSolutionQueue()
	psq.push(psi)

	minLooseEnds := 1000
	var bestPartial *partialSolutionItem
	startingPuzzle := psi.puzzle.DeepCopy()

	for !psq.isEmpty() {
		partial := psq.pop()

		if len(partial.looseEnds) < minLooseEnds {
			minLooseEnds = len(partial.looseEnds)
			bestPartial = partial
		}

		puzz := lec.connectLooseEnds(psq, partial)
		if puzz != nil {
			return puzz
		}
	}

	log.Printf("started from: \n%s\n", startingPuzzle.String())
	if bestPartial != nil {
		log.Printf("bestPartial: targeting %v, %+v\n%s\n",
			bestPartial.targeting,
			bestPartial.looseEnds,
			"todo", //bestPartial.puzzle.String(),
		)
	} else {
		log.Printf("best partial was nil\n")
	}
	return nil
}

func (lec *looseEndConnector) connectLooseEnds(
	psq *partialSolutionQueue,
	partial *partialSolutionItem,
) *puzzle.Puzzle {
	printPartialSolution(partial, 7)

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
