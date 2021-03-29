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

type partialSolutionItem struct {
	puzzle         *puzzle
	numSolvedNodes int

	looseEnds []nodeCoord
}

type targetSolver struct {
	puzzle *puzzle

	queue []*partialSolutionItem

	numProcessed int
	targets      []target
}

func newTargetSolver(
	size int,
	nl []NodeLocation,
) solver {
	if len(nl) == 0 {
		return nil
	}

	return &targetSolver{
		puzzle: newPuzzle(size, nl),
	}
}

func (d *targetSolver) iterations() int {
	return d.numProcessed
}

func (d *targetSolver) solve() (*puzzle, bool) {
	d.buildTargets()

	p := d.solveAllTargets()
	return p, p != nil
}

func (d *targetSolver) solveAllTargets() *puzzle {
	d.queue = d.getSolutionsForNode(&partialSolutionItem{
		puzzle:         d.puzzle,
		numSolvedNodes: 0,
	})

	for len(d.queue) > 0 {
		prev := d.queue[0]
		d.queue = d.queue[1:]

		snis := d.getSolutionsForNode(prev)
		if prev.numSolvedNodes < len(d.targets)-1 {
			d.queue = append(d.queue, snis...)
			continue
		}

		p := d.connectTheDots(snis)
		if p != nil {
			return p
		}
	}

	return nil
}

func (d *targetSolver) getSolutionsForNode(
	prevPartial *partialSolutionItem,
) (partials []*partialSolutionItem) {
	targetCoord := d.targets[prevPartial.numSolvedNodes].coord

	node, ok := prevPartial.puzzle.nodes[targetCoord]
	if !ok {
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
		}
		copy(item.looseEnds, prevPartial.looseEnds)
		// once we find a completion path, add it to the returned slice
		partials = append(partials, item)
		d.showPuzzle(item)
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
			looseEnds:      append(prevPartial.looseEnds, endOfEdge1, endOfEdge2),
		}
		// once we find a completion path, add it to the returned slice
		partials = append(partials, item)
		d.showPuzzle(item)
	}

	return partials
}

func (d *targetSolver) connectTheDots(snis []*partialSolutionItem) *puzzle {
	// TODO
	d.printSNIs(snis)
	return nil
}

func (d *targetSolver) printSNIs(partials []*partialSolutionItem) {
	for _, p := range partials {
		d.showPuzzle(p)
	}
}

func (d *targetSolver) showPuzzle(partial *partialSolutionItem) {
	if d.iterations() > 10 &&
		d.iterations()%100 != 0 &&
		partial.numSolvedNodes < len(d.targets) {
		return
	}
	log.Printf("showPuzzle (%d iterations): (%d nodes) %v", d.iterations(), partial.numSolvedNodes, partial.looseEnds)
	p := partial.puzzle
	if p == nil {
		log.Printf(": (nil)\n")
		return
	}
	log.Printf(":\n%s\n", p.String())
	fmt.Scanf("hello there")
}

func (d *targetSolver) buildTargets() {
	for nc, n := range d.puzzle.nodes {
		d.targets = append(d.targets, target{
			coord: nc,
			val:   int(n.val),
		})
	}
	maxRowColVal := d.puzzle.numEdges
	isOnTheSide := func(coord nodeCoord) bool {
		return coord.row == 0 || coord.col == 0 ||
			coord.row == rowIndex(maxRowColVal) || coord.col == colIndex(maxRowColVal)
	}
	sort.Slice(d.targets, func(i, j int) bool {
		// rank higher valued nodes at the start of the target list
		if d.targets[i].val != d.targets[j].val {
			return d.targets[i].val > d.targets[j].val
		}

		// put nodes with more limitations (i.e. on the sides of
		// of the graph) higher up on the list
		iIsEdge := isOnTheSide(d.targets[i].coord)
		jIsEdge := isOnTheSide(d.targets[j].coord)
		if iIsEdge && !jIsEdge {
			return true
		} else if jIsEdge && !iIsEdge {
			return false
		}

		// at this point, we just want a consistent ordering.
		// let's put nodes closer to (0,0) higher up in the list
		if d.targets[i].coord.row != d.targets[j].coord.row {
			return d.targets[i].coord.row < d.targets[j].coord.row
		}
		return d.targets[i].coord.col < d.targets[j].coord.col
	})

	fmt.Printf("buildTargets complete!\n%+v\n", d.targets)
}
