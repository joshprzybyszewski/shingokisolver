package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func ClaimGimmes(p Puzzle) (Puzzle, model.State) {
	outPuzz, ms := claimGimmes(p)
	if ms != model.Incomplete {
		return Puzzle{}, ms
	}

	return outPuzz, outPuzz.GetState()
}

func claimGimmes(
	p Puzzle,
) (Puzzle, model.State) {

	return performUpdates(p, updates{
		metas: p.getMetasCopy(),
	})
}
