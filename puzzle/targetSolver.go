package puzzle

import (
	"fmt"
	"log"
	"sort"
)

type target struct {
	coord nodeCoord
	val   int
}

func buildTargets(p *puzzle) []target {
	targets := make([]target, 0, len(p.nodes))

	for nc, n := range p.nodes {
		targets = append(targets, target{
			coord: nc,
			val:   int(n.val),
		})
	}

	maxRowColVal := p.numEdges
	isOnTheSide := func(coord nodeCoord) bool {
		return coord.row == 0 ||
			coord.row == rowIndex(maxRowColVal) ||
			coord.col == 0 ||
			coord.col == colIndex(maxRowColVal)
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
		if targets[i].coord.row != targets[j].coord.row {
			return targets[i].coord.row < targets[j].coord.row
		}
		return targets[i].coord.col < targets[j].coord.col
	})

	return targets
}

type partialSolutionQueue struct {
	items []*partialSolutionItem
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
	puzzle         *puzzle
	numSolvedNodes int

	looseEnds []nodeCoord
}

// eliminates loose ends that don't actually exist
func (partial *partialSolutionItem) removeDuplicateLooseEnds() {
	sort.Slice(partial.looseEnds, func(i, j int) bool {
		if partial.looseEnds[i].row != partial.looseEnds[j].row {
			return partial.looseEnds[i].row < partial.looseEnds[j].row
		}
		return partial.looseEnds[i].col < partial.looseEnds[j].col
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
	if !IncludeProgressLogs {
		return
	}
	// if !(iterations < 10 ||
	// iterations%100 == 0 ||
	if partial.puzzle == nil ||
		partial.numSolvedNodes < len(partial.puzzle.nodes) {
		return
	}

	log.Printf("printPartialSolution (%d iterations): (%d nodes solved) %v",
		iterations,
		partial.numSolvedNodes,
		partial.looseEnds,
	)
	log.Printf(":\n%s\n", partial.puzzle.String())
	fmt.Scanf("hello there")
}

type targetSolver struct {
	puzzle *puzzle

	queue             *partialSolutionQueue
	looseEndConnector *looseEndConnector

	numProcessed int
}

func newTargetSolver(
	size int,
	nl []NodeLocation,
) solver {
	if len(nl) == 0 {
		return nil
	}

	return &targetSolver{
		puzzle:            newPuzzle(size, nl),
		queue:             &partialSolutionQueue{},
		looseEndConnector: &looseEndConnector{},
	}
}

func (d *targetSolver) iterations() int {
	return d.numProcessed
}

func (d *targetSolver) solve() (*puzzle, bool) {
	targets := buildTargets(d.puzzle)

	p := d.getFullSolution(targets)
	return p, p != nil
}

func (d *targetSolver) getFullSolution(
	targets []target,
) *puzzle {

	d.queue.push(d.getPartialSolutionsForNode(
		targets,
		&partialSolutionItem{
			puzzle:         d.puzzle,
			numSolvedNodes: 0,
		},
	)...)

	for !d.queue.isEmpty() {
		prev := d.queue.pop()

		partials := d.getPartialSolutionsForNode(targets, prev)
		if prev.numSolvedNodes < len(targets)-1 {
			d.queue.push(partials...)
			continue
		}

		d.looseEndConnector.addPartialSolutions(partials...)
	}

	p := d.getFullSolutionFromLooseEndConnector()
	if p != nil {
		return p
	}

	return nil
}

func (d *targetSolver) getPartialSolutionsForNode(
	targets []target,
	prevPartial *partialSolutionItem,
) (partials []*partialSolutionItem) {
	targetCoord := targets[prevPartial.numSolvedNodes].coord

	node, ok := prevPartial.puzzle.nodes[targetCoord]
	if !ok {
		// this should be returning an error, but really it shouldn't be happening
		panic(`what?`)
		// return nil
	}

	// go out in all directions from the target
	// if it's still a valid puzzle, keep going outward
	// until we "complete" the node.
	for i, feeler1 := range allDirections {
		for _, feeler2 := range allDirections[i+1:] {
			if node.nType.isInvalidMotions(feeler1, feeler2) {
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
	prevPartial *partialSolutionItem,
	targetCoord nodeCoord,
	node node,
	c1, c2 cardinal,
) (partials []*partialSolutionItem) {

	if isCompleteNode(prevPartial.puzzle, targetCoord) {
		// the target node is already complete, perhaps a previous node
		// accidentally completed it. If so, then let's do a sanity check
		// on completion, and then add it as a "partial solution" that
		// has no new loose ends
		if _, err := prevPartial.puzzle.IsIncomplete(targetCoord); err != nil {
			return nil
		}

		d.numProcessed++
		item := &partialSolutionItem{
			puzzle:         prevPartial.puzzle.deepCopy(),
			numSolvedNodes: prevPartial.numSolvedNodes + 1,
			looseEnds:      make([]nodeCoord, len(prevPartial.looseEnds)),
		}
		copy(item.looseEnds, prevPartial.looseEnds)
		// once we find a completion path, add it to the returned slice
		partials = append(partials, item)
		printPartialSolution(item, d.iterations())
		return partials
	}

	endOfEdge1 := targetCoord
	endOfEdge2 := targetCoord

	var err error
	pWithEdge1 := prevPartial.puzzle

	for i := 1; i <= int(node.val); i++ {
		d.numProcessed++
		endOfEdge1, pWithEdge1, err = pWithEdge1.AddEdge(c1, endOfEdge1)
		if err != nil {
			break
		}

		// reset the "end" of the second edge to be at the target
		endOfEdge2 = targetCoord
		pWithEdge2 := pWithEdge1
		for j := 1; i+j <= int(node.val); j++ {
			d.numProcessed++
			endOfEdge2, pWithEdge2, err = pWithEdge2.AddEdge(c2, endOfEdge2)
			if err != nil || pWithEdge2 == nil {
				break
			}
		}

		if pWithEdge2 == nil {
			continue
		}

		if !isCompleteNode(pWithEdge2, targetCoord) {
			// we _should_ have added a number of straight edges that will
			// complete our target node.
			// if not, then we don't want to add this to our partials
			continue
		}

		if _, err := pWithEdge2.IsIncomplete(targetCoord); err != nil {
			continue
		}

		item := &partialSolutionItem{
			puzzle:         pWithEdge2,
			numSolvedNodes: prevPartial.numSolvedNodes + 1,
			looseEnds:      make([]nodeCoord, len(prevPartial.looseEnds)),
		}
		copy(item.looseEnds, prevPartial.looseEnds)
		item.looseEnds = append(item.looseEnds, endOfEdge1, endOfEdge2)

		// once we find a completion path, add it to the returned slice
		partials = append(partials, item)
		printPartialSolution(item, d.iterations())
	}

	return partials
}

func (d *targetSolver) getFullSolutionFromLooseEndConnector() *puzzle {

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

func (lec *looseEndConnector) solve() *puzzle {
	// remove duplicate loose ends from all of our solutions
	for _, partial := range lec.partials {
		partial.removeDuplicateLooseEnds()
	}

	// sort the partial solutions so that the solutions with the most
	// connections (least number of loose ends) are at the front
	sort.Slice(lec.partials, func(i, j int) bool {
		return len(lec.partials[i].looseEnds) < len(lec.partials[j].looseEnds)
	})

	// iterate through the partial solutions, trying to connect all of their
	// loose ends.
	attemptedCache := newPuzzleCache()
	for _, partial := range lec.partials {
		if attemptedCache.contains(partial.puzzle) {
			continue
		}
		p := lec.connectLooseEnds(partial)
		if p != nil {
			return p
		}
		attemptedCache.add(partial.puzzle)
	}

	return nil
}

func (lec *looseEndConnector) connectLooseEnds(
	partial *partialSolutionItem,
) *puzzle {
	printPartialSolution(partial, lec.iterations)

	dfs := newDFSSolverForPartialSolution(partial.puzzle)
	defer func() {
		lec.iterations += dfs.iterations()
	}()
	for i, start := range partial.looseEnds {
		for j := i; j < len(partial.looseEnds); j++ {
			g, sol := dfs.solveForGoal(start, partial.looseEnds[j])
			switch sol {
			case solvedPuzzle:
				return g
			case foundGoal:
				// TODO take note
			case badState:
				// just continue on with life...
			}
		}
	}

	return nil
}
