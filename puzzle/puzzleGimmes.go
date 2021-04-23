package puzzle

import (
	"log"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p Puzzle) ClaimGimmes() (Puzzle, model.State) {
	outPuzz, ms := claimGimmes(p, p.nodes)
	switch ms {
	case model.Violation, model.Unexpected:
		log.Printf("ClaimGimmes() second performUpdates got unexpected state: %s", ms)
		return Puzzle{}, ms
	}

	outPuzz.twoArmOptions = buildTwoArmsCache(
		outPuzz.nodes,
		outPuzz.NumEdges(),
		&outPuzz.edges,
		outPuzz,
	)
	return outPuzz, outPuzz.GetState(outPuzz.getRandomCoord())
}

func buildTwoArmsCache(
	allNodes []model.Node,
	numEdges int,
	ge model.GetEdger,
	gn model.GetNoder,
) []nodeWithOptions {
	res := make([]nodeWithOptions, 0, len(allNodes))

	for _, n := range allNodes {
		allTAs := model.BuildTwoArmOptions(n, numEdges)
		res = append(res, nodeWithOptions{
			Node:    n,
			Options: n.GetFilteredOptions(allTAs, ge, gn),
		})
	}

	return res
}
