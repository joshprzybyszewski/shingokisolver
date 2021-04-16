package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func (p *Puzzle) ClaimGimmes() model.State {
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

	return p.runQueue()
}
