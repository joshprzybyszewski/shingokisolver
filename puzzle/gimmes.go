package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func (p *Puzzle) ClaimGimmes() model.State {
	numEdges := p.NumEdges()
	rq := newRulesQueue()

	for r := 0; r <= numEdges; r++ {
		for c := 0; c <= numEdges; c++ {
			nc := model.NewCoordFromInts(r, c)

			if c < numEdges {
				ep := newEdgePair(nc, model.HeadRight)

				switch s := p.updateEdgeFromRules(rq, ep); s {
				case model.Violation,
					model.Unexpected:
					return s
				}
			}

			if r < numEdges {
				ep := newEdgePair(nc, model.HeadDown)

				switch s := p.updateEdgeFromRules(rq, ep); s {
				case model.Violation,
					model.Unexpected:
					return s
				}
			}
		}
	}

	return p.runQueue(rq)
}
