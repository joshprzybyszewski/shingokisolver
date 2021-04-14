package puzzle

import (
	"fmt"
	"strings"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p *Puzzle) String() string {
	return p.string(true)
}

func (p *Puzzle) Solution() string {
	return p.string(false)
}

func (p *Puzzle) string(
	includeXs bool,
) string {
	if p == nil {
		return `(*Puzzle)<nil>`
	}
	var sb strings.Builder
	sb.WriteString("\n")
	for r := 0; r < p.numNodes(); r++ {
		var below strings.Builder
		for c := 0; c < p.numNodes(); c++ {
			nc := model.NewCoordFromInts(r, c)
			// write a node
			sb.WriteString(`(`)
			if n, ok := p.nodes[nc]; ok {
				if n.Type() == model.WhiteNode {
					sb.WriteString(`w`)
				} else {
					sb.WriteString(`b`)
				}
				sb.WriteString(fmt.Sprintf("%2d", n.Value()))
			} else {
				sb.WriteString(`(-)`)
			}
			sb.WriteString(`)`)

			// now draw an edge
			ep := NewEdgePair(nc, model.HeadRight)
			if p.edges.isInBounds(ep) {
				switch p.edges.GetEdge(ep) {
				case model.EdgeExists:
					sb.WriteString(`---`)
				case model.EdgeAvoided:
					if includeXs {
						sb.WriteString(` X `)
					} else {
						sb.WriteString(`   `)
					}
				default:
					sb.WriteString(`   `)
				}
			}

			// now draw any edges that are below
			below.WriteString(`  `)
			ep = NewEdgePair(nc, model.HeadDown)
			if p.edges.isInBounds(ep) {
				switch p.edges.GetEdge(ep) {
				case model.EdgeExists:
					below.WriteString(`|`)
				case model.EdgeAvoided:
					if includeXs {
						below.WriteString(`X`)
					} else {
						below.WriteString(` `)
					}
				default:
					below.WriteString(` `)
				}
			}
			below.WriteString(`     `)
		}
		sb.WriteString("\n")
		sb.WriteString(below.String())
		sb.WriteString("\n")
	}
	return sb.String()
}
