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

	updateCache(&outPuzz)

	return outPuzz, outPuzz.GetState()
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
	updateCache(&obviousFilled)

	allNodeEdgesToCheck = make(map[model.Node]map[model.Cardinal]int8, len(obviousFilled.nodes))

	nodesAndOptions := make([]nodeAndOptions, len(obviousFilled.nodes))

	for i := range obviousFilled.nodes {
		n := obviousFilled.nodes[i]

		optionsCopy := make([]model.TwoArms, len(obviousFilled.twoArmOptions[i]))
		copy(optionsCopy, obviousFilled.twoArmOptions[i])

		allNodeEdgesToCheck[n] = model.GetMaxArmsByDir(optionsCopy)
		nodesAndOptions[i] = nodeAndOptions{
			Node:    n,
			Options: optionsCopy,
		}
	}

	// Now we're going to add rules for the nodes that have the fewest options first.
	// This helps with evaluation because it means we're going to need to do less
	// filtering on the first evaluators.
	sort.Slice(nodesAndOptions, func(i, j int) bool {
		return len(nodesAndOptions[i].Options) < len(nodesAndOptions[j].Options)
	})

	for _, nwo := range nodesAndOptions {
		obviousFilled.rules.AddAllTwoArmRules(
			nwo.Node,
			obviousFilled.gn,
			nwo.Options,
		)
	}

	return performUpdates(obviousFilled, updates{
		nodes: allNodeEdgesToCheck,
	})

}

type nodeAndOptions struct {
	Options []model.TwoArms
	model.Node
}
