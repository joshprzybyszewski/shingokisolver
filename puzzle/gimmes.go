package puzzle

import (
	"log"
	"sort"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func ClaimGimmes(p Puzzle) (Puzzle, model.State) {
	outPuzz, ms := claimGimmes(p)
	if ms != model.Incomplete {
		log.Printf("ClaimGimmes() claimGimmes got unexpected state: %s", ms)
		return Puzzle{}, ms
	}

	outPuzz.twoArmOptions = buildTwoArmsCache(
		outPuzz.nodes,
		outPuzz.numEdges(),
		&outPuzz.edges,
	)
	return outPuzz, outPuzz.GetState()
}

func buildTwoArmsCache(
	allNodes nodeList,
	numEdges int,
	ge model.GetEdger,
) []nodeWithOptions {
	res := make([]nodeWithOptions, 0, len(allNodes))

	for _, n := range allNodes {
		allTAs := model.BuildTwoArmOptions(n, numEdges)
		maxLensByDir := model.GetMaxArmsByDir(allTAs)
		nearbyNodes := model.BuildNearbyNodes(n, allNodes, maxLensByDir)
		options := n.GetFilteredOptions(allTAs, ge, nearbyNodes)
		res = append(res, nodeWithOptions{
			Node:    n,
			Options: options,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return len(res[i].Options) < len(res[j].Options)
	})

	return res
}

func claimGimmes(
	p Puzzle,
) (Puzzle, model.State) {
	allNodeEdgesToCheck := make(map[model.Node]map[model.Cardinal]int8, len(p.nodes))
	for _, n := range p.nodes {
		allNodeEdgesToCheck[n] = model.GetMaxArmsByDir(
			model.BuildTwoArmOptions(n, p.numEdges()),
		)
	}
	obviousFilled, ms := performUpdates(p, updates{
		nodes: allNodeEdgesToCheck,
	})

	if ms != model.Incomplete {
		log.Printf("ClaimGimmes() first performUpdates got unexpected state: %s", ms)
		return Puzzle{}, ms
	}

	// now we're going to add all of the extended rules
	var nl nodeList = obviousFilled.nodes
	twoArmOptions := buildTwoArmsCache(
		nl,
		obviousFilled.numEdges(),
		&obviousFilled.edges,
	)

	// TODO figure out the best way to give a GetNoder to the rules
	var gnp model.GetNoder = nl

	allNodeEdgesToCheck = make(map[model.Node]map[model.Cardinal]int8, len(obviousFilled.nodes))
	for _, tao := range twoArmOptions {
		allNodeEdgesToCheck[tao.Node] = model.GetMaxArmsByDir(tao.Options)
		obviousFilled.rules.AddAllTwoArmRules(
			tao.Node,
			gnp,
			tao.Options,
		)
	}

	return performUpdates(obviousFilled, updates{
		nodes: allNodeEdgesToCheck,
	})

}
