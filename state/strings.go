package state

import (
	"strings"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (edges TriEdges) String() string {
	return edges.string(true)
}

func (edges TriEdges) Solution() string {
	return edges.string(false)
}

func (edges TriEdges) string(
	includeXs bool,
) string {
	var sb strings.Builder
	sb.WriteString("\n")
	for r := 0; r < edges.NumEdges()+1; r++ {
		var below strings.Builder
		for c := 0; c < edges.NumEdges()+1; c++ {
			nc := model.NewCoordFromInts(r, c)
			// write a node
			sb.WriteString(`(   )`)

			// now draw an edge
			ep := model.NewEdgePair(nc, model.HeadRight)
			if edges.IsInBounds(ep) {
				switch edges.GetEdge(ep) {
				case model.EdgeExists:
					sb.WriteString(`---`)
				case model.EdgeAvoided:
					if includeXs {
						sb.WriteString(` X `)
					} else {
						sb.WriteString(`   `)
					}
				case model.EdgeUnknown:
					sb.WriteString(`   `)
				default:
					sb.WriteString(`???`)
				}
			}

			// now draw any edges that are below
			below.WriteString(` `)
			ep = model.NewEdgePair(nc, model.HeadDown)
			if edges.IsInBounds(ep) {
				switch edges.GetEdge(ep) {
				case model.EdgeExists:
					below.WriteString(` | `)
				case model.EdgeAvoided:
					if includeXs {
						below.WriteString(` X `)
					} else {
						below.WriteString(`   `)
					}
				case model.EdgeUnknown:
					below.WriteString(`   `)
				default:
					below.WriteString(`???`)
				}
			}
			below.WriteString(`    `)
		}
		sb.WriteString("\n")
		sb.WriteString(below.String())
		sb.WriteString("\n")
	}
	return sb.String()
}
