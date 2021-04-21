package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p Puzzle) ClaimGimmes() (Puzzle, model.State) {

	rq := logic.NewQueue(p.edges, p.NumEdges())

	// first we're going to claim any of the gimmes from the "standard"
	// node rules.
	for _, n := range p.nodes {
		nc := n.Coord()
		for _, dir := range model.AllCardinals {
			ep := model.NewEdgePair(nc, dir)

			switch s := p.updateEdgeFromRules(ep, rq); s {
			case model.Violation,
				model.Unexpected:
				return Puzzle{}, s
			}
		}
	}

	// now we're going to add all of the extended rules
	for _, n := range p.nodes {
		p.rules.AddAllTwoArmRules(n, p.getPossibleTwoArms(n))
	}

	// at this point, let's double check the edges surrounding the nodes
	// so that they can catch the extended rules that now apply to them.
	for _, n := range p.nodes {
		nc := n.Coord()
		for _, dir := range model.AllCardinals {
			ep := model.NewEdgePair(nc, dir)

			switch s := p.updateEdgeFromRules(ep, rq); s {
			case model.Violation,
				model.Unexpected:
				return Puzzle{}, s
			}
		}
	}

	// run the queue down
	switch s := p.runQueue(rq); s {
	case model.Violation, model.Unexpected:
		return Puzzle{}, s
	}

	p.twoArmOptions = getTwoArmsCache(p.nodes, p.NumEdges(), p.edges, p)
	return p, model.Incomplete
}
