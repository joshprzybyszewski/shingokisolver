package puzzle

import (
	"fmt"
	"strings"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p *Puzzle) String() string {
	// return p.DebugString()
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
			ep, err := standardizeInput(nc, model.HeadRight)
			if err == nil && p.edges.isInBounds(ep) {
				switch p.edges.GetEdge(ep) {
				case model.EdgeExists:
					sb.WriteString(`---`)
				case model.EdgeAvoided:
					sb.WriteString(`XXX`)
				default:
					sb.WriteString(`   `)
				}
			}

			// now draw any edges that are below
			below.WriteString(`  `)
			ep, err = standardizeInput(nc, model.HeadDown)
			if err == nil && p.edges.isInBounds(ep) {
				switch p.edges.GetEdge(ep) {
				case model.EdgeExists:
					below.WriteString(`|`)
				case model.EdgeAvoided:
					below.WriteString(`X`)
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

func (p *Puzzle) DebugString() string {
	if p == nil {
		return `(*Puzzle)<nil>`
	}
	var sb strings.Builder
	sb.WriteString("\n")
	for r := 0; r < p.numNodes(); r++ {
		var above strings.Builder
		var aboveNums strings.Builder
		var row strings.Builder
		var belowNums strings.Builder
		var below strings.Builder
		for c := 0; c < p.numNodes(); c++ {
			nc := model.NewCoordFromInts(r, c)
			oe, _ := p.GetOutgoingEdgesFrom(nc)

			// draw the left side
			if oe.IsLeft() {
				row.WriteString(`<-`)
				row.WriteString(fmt.Sprintf("%2d", oe.Left()))
			} else {
				row.WriteString(`    `)
			}
			above.WriteString(`    `)
			aboveNums.WriteString(`    `)
			belowNums.WriteString(`    `)
			below.WriteString(`    `)

			// write a node
			row.WriteString(`(`)
			if n, ok := p.nodes[nc]; ok {
				if n.Type() == model.WhiteNode {
					row.WriteString(`w`)
				} else {
					row.WriteString(`b`)
				}
				row.WriteString(fmt.Sprintf("%2d", n.Value()))
			} else {
				row.WriteString(`XXX`)
			}
			row.WriteString(`)`)

			if oe.IsAbove() {
				above.WriteString(` ^^^ `)
				aboveNums.WriteString(fmt.Sprintf(" %2d ", oe.Above()))
			} else {
				above.WriteString(`     `)
				aboveNums.WriteString(`     `)
			}

			if oe.IsBelow() {
				belowNums.WriteString(fmt.Sprintf(" %2d ", oe.Below()))
				below.WriteString(` vvv `)
			} else {
				belowNums.WriteString(`     `)
				below.WriteString(`     `)
			}

			// now draw the right side
			if oe.IsRight() {
				row.WriteString(fmt.Sprintf("%2d", oe.Right()))
				row.WriteString(`->`)
			} else {
				row.WriteString(`    `)
			}
			row.WriteString(` `)

			above.WriteString(`     `)
			aboveNums.WriteString(`     `)
			belowNums.WriteString(`     `)
			below.WriteString(`     `)
		}
		sb.WriteString(above.String())
		sb.WriteString("\n")
		sb.WriteString(aboveNums.String())
		sb.WriteString("\n")
		sb.WriteString(row.String())
		sb.WriteString("\n")
		sb.WriteString(belowNums.String())
		sb.WriteString("\n")
		sb.WriteString(below.String())
		sb.WriteString("\n")
		sb.WriteString("\n")
	}
	return sb.String()
}
