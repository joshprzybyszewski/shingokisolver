package puzzle

import (
	"log"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p Puzzle) ClaimGimmes() (Puzzle, model.State) {
	outPuzz, ms := claimGimmes(p)
	if ms != model.Incomplete {
		log.Printf("ClaimGimmes() second performUpdates got unexpected state: %s", ms)
		return Puzzle{}, ms
	}

	outPuzz.twoArmOptions = buildTwoArmsCache(
		outPuzz.nodes,
		outPuzz.NumEdges(),
		&outPuzz.edges,
		outPuzz,
	)
	return outPuzz, outPuzz.GetState()
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
		nearbyNodes := model.BuildNearbyNodes(n, allTAs, gn)
		res = append(res, nodeWithOptions{
			Node:    n,
			Options: n.GetFilteredOptions(allTAs, ge, nearbyNodes),
		})
	}

	return res
}
