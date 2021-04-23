package puzzle

import (
	"log"
	"sort"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func ClaimGimmes(p Puzzle) (Puzzle, model.State) {
	outPuzz, ms := claimGimmes(p)
	if ms != model.Incomplete {
		log.Printf("ClaimGimmes() second performUpdates got unexpected state: %s", ms)
		return Puzzle{}, ms
	}

	outPuzz.twoArmOptions = buildTwoArmsCache(
		outPuzz.nodes,
		outPuzz.numEdges(),
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
	allNodeEdgesToCheck = make(map[model.Node]map[model.Cardinal]int8, len(obviousFilled.nodes))
	taRules := make([]twoArmRuleContext, 0, len(obviousFilled.nodes))
	for _, n := range obviousFilled.nodes {
		allTAs := model.BuildTwoArmOptions(n, obviousFilled.numEdges())
		nearbyNodes := model.BuildNearbyNodes(n, allTAs, obviousFilled)
		possibleTAs := n.GetFilteredOptions(allTAs, &obviousFilled.edges, nearbyNodes)
		taRules = append(taRules, twoArmRuleContext{
			node:    n,
			options: possibleTAs,
		})
		allNodeEdgesToCheck[n] = model.GetMaxArmsByDir(possibleTAs)
	}

	sort.Slice(taRules, func(i, j int) bool {
		return len(taRules[i].options) < len(taRules[j].options)
	})

	for _, taRule := range taRules {
		obviousFilled.rules.AddAllTwoArmRules(
			taRule.node,
			obviousFilled,
			taRule.options,
		)
	}

	return performUpdates(obviousFilled, updates{
		nodes: allNodeEdgesToCheck,
	})

}

type twoArmRuleContext struct {
	options []model.TwoArms
	node    model.Node
}
