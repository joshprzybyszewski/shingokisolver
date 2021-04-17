package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func (p *Puzzle) ClaimGimmes() model.State {

	defer p.populateTwoArmsCache()

	// first we're going to claim any of the gimmes from the "standard"
	// node rules.
	for nc := range p.nodes {
		for _, dir := range model.AllCardinals {
			ep := model.NewEdgePair(nc, dir)

			switch s := p.updateEdgeFromRules(ep); s {
			case model.Violation,
				model.Unexpected:
				return s
			}
		}
	}

	// now we're going to add all of the extended rules
	for _, n := range p.nodes {
		p.rules.addAllTwoArmRules(n, p.getPossibleTwoArms(n))
	}

	// at this point, let's double check the edges surrounding the nodes
	// so that they can catch the extended rules that now apply to them.
	for nc := range p.nodes {
		for _, dir := range model.AllCardinals {
			ep := model.NewEdgePair(nc, dir)

			switch s := p.updateEdgeFromRules(ep); s {
			case model.Violation,
				model.Unexpected:
				return s
			}
		}
	}

	// run the queue down
	return p.runQueue()
}
