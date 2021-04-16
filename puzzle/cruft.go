package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func (p *Puzzle) GetUnknownEdge() (model.EdgePair, bool) {
	// TODO find a loose end instead
	for r := 0; r <= p.NumEdges(); r++ {
		for c := 0; c <= p.NumEdges(); c++ {
			nc := model.NewCoordFromInts(r, c)

			ep := model.NewEdgePair(nc, model.HeadRight)
			if p.GetEdgeState(ep) == model.EdgeUnknown {
				return ep, true
			}

			ep = model.NewEdgePair(nc, model.HeadDown)
			if p.GetEdgeState(ep) == model.EdgeUnknown {
				return ep, true
			}
		}
	}
	return model.EdgePair{}, false
}
