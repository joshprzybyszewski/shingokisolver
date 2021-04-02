package solvers

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func claimGimmes(
	puzz *puzzle.Puzzle,
) (*puzzle.Puzzle, []model.NodeCoord, int) {
	var state model.State

	var looseEnds []model.NodeCoord
	numProcessed := 0
	// iterate through all four edges. if it's a white node, take either side along the
	// edge. if it's black, claim the inwards facing edge
	for i := 0; i < puzz.NumEdges(); i++ {
		// top (row = 0)
		top := model.NewCoordFromInts(0, i)
		if n, ok := puzz.GetNode(top); ok {
			switch n.Type() {
			case model.BlackNode:
				var end model.NodeCoord

				numProcessed++
				end, puzz, state = puzz.AddEdge(model.HeadDown, top)
				switch state {
				case model.Violation, model.Unexpected:
					panic(`bad dev: ` + state.String())
				}
				looseEnds = append(looseEnds, top, end)
			case model.WhiteNode:
				var left, right model.NodeCoord

				numProcessed++
				left, puzz, state = puzz.AddEdge(model.HeadLeft, top)
				switch state {
				case model.Violation, model.Unexpected:
					panic(`bad dev: ` + state.String())
				}

				numProcessed++
				right, puzz, state = puzz.AddEdge(model.HeadRight, top)
				switch state {
				case model.Violation, model.Unexpected:
					panic(`bad dev: ` + state.String())
				}
				looseEnds = append(looseEnds, left, right)
			}
		}

		// bottom (row = puzz.NumEdges() + 1)
		bottom := model.NewCoordFromInts(puzz.NumEdges()+1, i)
		if n, ok := puzz.GetNode(bottom); ok {
			switch n.Type() {
			case model.BlackNode:
				var end model.NodeCoord

				numProcessed++
				end, puzz, state = puzz.AddEdge(model.HeadUp, bottom)
				switch state {
				case model.Violation, model.Unexpected:
					panic(`bad dev: ` + state.String())
				}
				looseEnds = append(looseEnds, bottom, end)
			case model.WhiteNode:
				var left, right model.NodeCoord

				numProcessed++
				left, puzz, state = puzz.AddEdge(model.HeadLeft, bottom)
				switch state {
				case model.Violation, model.Unexpected:
					panic(`bad dev: ` + state.String())
				}

				numProcessed++
				right, puzz, state = puzz.AddEdge(model.HeadRight, bottom)
				switch state {
				case model.Violation, model.Unexpected:
					panic(`bad dev: ` + state.String())
				}
				looseEnds = append(looseEnds, left, right)
			}
		}

		// left (col = 0)
		left := model.NewCoordFromInts(i, 0)
		if n, ok := puzz.GetNode(left); ok {
			switch n.Type() {
			case model.BlackNode:
				var end model.NodeCoord

				numProcessed++
				end, puzz, state = puzz.AddEdge(model.HeadRight, left)
				switch state {
				case model.Violation, model.Unexpected:
					panic(`bad dev: ` + state.String())
				}
				looseEnds = append(looseEnds, left, end)
			case model.WhiteNode:
				var top, bottom model.NodeCoord

				numProcessed++
				top, puzz, state = puzz.AddEdge(model.HeadUp, left)
				switch state {
				case model.Violation, model.Unexpected:
					panic(`bad dev: ` + state.String())
				}

				numProcessed++
				bottom, puzz, state = puzz.AddEdge(model.HeadDown, left)
				switch state {
				case model.Violation, model.Unexpected:
					panic(`bad dev: ` + state.String())
				}
				looseEnds = append(looseEnds, top, bottom)
			}
		}

		// right (col = puzz.NumEdges() + 1)
		right := model.NewCoordFromInts(i, puzz.NumEdges()+1)
		if n, ok := puzz.GetNode(right); ok {
			switch n.Type() {
			case model.BlackNode:
				var end model.NodeCoord
				end, puzz, state = puzz.AddEdge(model.HeadLeft, right)
				switch state {
				case model.Violation, model.Unexpected:
					panic(`bad dev: ` + state.String())
				}
				looseEnds = append(looseEnds, right, end)
			case model.WhiteNode:
				var top, bottom model.NodeCoord
				top, puzz, state = puzz.AddEdge(model.HeadUp, right)
				switch state {
				case model.Violation, model.Unexpected:
					panic(`bad dev: ` + state.String())
				}
				bottom, puzz, state = puzz.AddEdge(model.HeadDown, right)
				switch state {
				case model.Violation, model.Unexpected:
					panic(`bad dev: ` + state.String())
				}
				looseEnds = append(looseEnds, top, bottom)
			}
		}
	}

	return puzz, looseEnds, numProcessed
}
