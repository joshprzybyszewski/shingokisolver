// nolint:gocyclo
package puzzle

import (
	"strings"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

const (
	esc    = "\033"
	red    = esc + `[31m`
	green  = esc + `[32m`
	yellow = esc + `[33m`
	dim    = esc + `[2m`
	end    = esc + `[0m`

	emptyNode = `(   )`
)

func addNormal(
	sb *strings.Builder,
	text string,
) {
	sb.WriteString(text)
}

func addGreen(
	sb *strings.Builder,
	text string,
) {
	sb.WriteString(green)
	sb.WriteString(text)
	sb.WriteString(end)
}

func addRed(
	sb *strings.Builder,
	text string,
) {
	sb.WriteString(red)
	sb.WriteString(text)
	sb.WriteString(end)
}

func addYellow(
	sb *strings.Builder,
	text string,
) {
	sb.WriteString(yellow)
	sb.WriteString(text)
	sb.WriteString(end)
}

func addDim(
	sb *strings.Builder,
	text string,
) {
	sb.WriteString(dim)
	sb.WriteString(text)
	sb.WriteString(end)
}

func (p Puzzle) numNodes() int {
	return p.numEdges() + 1
}

func (p Puzzle) String() string {
	return p.string(true)
}

func (p Puzzle) Solution() string {
	return p.string(false)
}

func (p Puzzle) string(
	includeXs bool,
) string {
	var sb strings.Builder
	sb.WriteString("\n")
	for r := 0; r < p.numNodes(); r++ {
		var below strings.Builder
		for c := 0; c < p.numNodes(); c++ {
			nc := model.NewCoordFromInts(r, c)
			// write a node
			if n, ok := p.GetNode(nc); ok {
				nOut, isMax := getSumOutgoingStraightLines(nc, &p.edges)
				if nOut == n.Value() {
					addGreen(&sb, n.PrettyString())
				} else if isMax {
					addRed(&sb, n.PrettyString())
				} else if nOut, nAvoid := getNumOutEdges(nc, &p.edges); nOut > 2 {
					addRed(&sb, n.PrettyString())
				} else if nOut == 2 {
					addYellow(&sb, n.PrettyString())
				} else if nAvoid+nOut == 4 || nAvoid == 3 {
					addRed(&sb, n.PrettyString())
				} else {
					addNormal(&sb, n.PrettyString())
				}
			} else {
				if nOut, nAvoid := getNumOutEdges(nc, &p.edges); nOut > 2 {
					addRed(&sb, emptyNode)
				} else if nOut == 2 {
					addGreen(&sb, emptyNode)
				} else if (nOut != 0 && nAvoid+nOut == 4) || nAvoid == 3 {
					addRed(&sb, emptyNode)
				} else if nOut == 1 {
					addYellow(&sb, emptyNode)
				} else {
					addNormal(&sb, emptyNode)
				}
			}

			// now draw an edge
			ep := model.NewEdgePair(nc, model.HeadRight)
			if p.edges.IsInBounds(ep) {
				switch p.edges.GetEdge(ep) {
				case model.EdgeExists:
					addGreen(&sb, `---`)
				case model.EdgeAvoided:
					if includeXs {
						addDim(&sb, ` X `)
					} else {
						sb.WriteString(`   `)
					}
				case model.EdgeUnknown:
					sb.WriteString(`   `)
				default:
					addRed(&sb, `???`)
				}
			}

			// now draw any edges that are below
			below.WriteString(` `)
			ep = model.NewEdgePair(nc, model.HeadDown)
			if p.edges.IsInBounds(ep) {
				switch p.edges.GetEdge(ep) {
				case model.EdgeExists:
					addGreen(&below, ` | `)
				case model.EdgeAvoided:
					if includeXs {
						addDim(&below, ` X `)
					} else {
						below.WriteString(`   `)
					}
				case model.EdgeUnknown:
					below.WriteString(`   `)
				default:
					addRed(&below, `???`)
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

func getNumOutEdges(
	nc model.NodeCoord,
	ge model.GetEdger,
) (int, int) {
	nOut := 0
	nAvoid := 0
	for _, dir := range model.AllCardinals {
		if ge.IsEdge(model.NewEdgePair(nc, dir)) {
			nOut++
		} else if ge.IsAvoided(model.NewEdgePair(nc, dir)) {
			nAvoid++
		}
	}
	return nOut, nAvoid
}
