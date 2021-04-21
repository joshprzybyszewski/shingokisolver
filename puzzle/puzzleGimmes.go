package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p Puzzle) ClaimGimmes() (Puzzle, model.State) {

	// TODO make these structs if possible
	newState := p.edges.Copy()
	rq := logic.NewQueue(&newState, newState.NumEdges())
	rules := p.rules

	// first we're going to claim any of the gimmes from the "standard"
	// node rules.
	for _, n := range p.nodes {
		nc := n.Coord()
		for _, dir := range model.AllCardinals {
			ep := model.NewEdgePair(nc, dir)

			switch s := updateEdgeFromRules(&newState, ep, rq, rules); s {
			case model.Violation,
				model.Unexpected:
				return Puzzle{}, s
			}
		}
	}

	// now we're going to add all of the extended rules
	for _, n := range p.nodes {
		rules.AddAllTwoArmRules(n, p.getPossibleTwoArms(n))
	}

	// at this point, let's double check the edges surrounding the nodes
	// so that they can catch the extended rules that now apply to them.
	for _, n := range p.nodes {
		nc := n.Coord()
		for _, dir := range model.AllCardinals {
			ep := model.NewEdgePair(nc, dir)

			switch s := updateEdgeFromRules(&newState, ep, rq, rules); s {
			case model.Violation,
				model.Unexpected:
				return Puzzle{}, s
			}
		}
	}

	// run the queue down
	switch s := runQueue(&newState, rq, rules); s {
	case model.Violation, model.Unexpected:
		return Puzzle{}, s
	}

	twoArmOptions := getTwoArmsCache(
		p.nodes,
		p.NumEdges(),
		&newState,
		p,
	)

	return Puzzle{
		edges:         newState,
		rules:         rules,
		twoArmOptions: twoArmOptions,
		nodes:         p.nodes,
	}, model.Incomplete
}
