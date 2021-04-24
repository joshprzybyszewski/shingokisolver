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
		nearbyNodes := model.BuildNearbyNodes(n, allTAs, allNodes)
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

type nodeList []model.Node

func (nl nodeList) GetNode(nc model.NodeCoord) (model.Node, bool) {
	// TODO replace with something better?
	for _, n := range nl {
		if n.Coord() == nc {
			return n, true
		}
	}
	return model.Node{}, false
}

func (nl nodeList) toNodeGrid() nodeGrid {
	maxRowIndex := -1
	maxColIndex := -1
	for _, n := range nl {
		if int(n.Coord().Row) > maxRowIndex {
			maxRowIndex = int(n.Coord().Row)
		}
		if int(n.Coord().Col) > maxColIndex {
			maxColIndex = int(n.Coord().Col)
		}
	}

	ng := make(nodeGrid, maxRowIndex+1)
	for r := 0; r < len(ng); r++ {
		row := make([]*model.Node, maxColIndex+1)
		ng[r] = row
	}

	for _, n := range nl {
		n := n
		ng[n.Coord().Row][n.Coord().Col] = &n
	}

	return ng
}

type nodeGrid [][]*model.Node

func (ng nodeGrid) GetNode(nc model.NodeCoord) (model.Node, bool) {
	if ng != nil && int(nc.Row) < len(ng) {
		if nRow := ng[nc.Row]; nRow != nil && int(nc.Col) < len(nRow) {
			if n := nRow[nc.Col]; n != nil {
				return *n, true
			}
		}
	}

	return model.Node{}, false
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

	gnp := nl.toNodeGrid()

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
