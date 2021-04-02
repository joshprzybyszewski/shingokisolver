package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

func (p *Puzzle) Targets() []model.Target {
	return model.BuildTargets(p.nodes, p.NumEdges())
}
