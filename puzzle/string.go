// nolint:gocyclo
package puzzle

import (
	"fmt"
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

type colorWriter interface {
	addNormal(*strings.Builder, string)
	addGreen(*strings.Builder, string)
	addRed(*strings.Builder, string)
	addYellow(*strings.Builder, string)
	addDim(*strings.Builder, string)
}

type writer int

const (
	noColor  writer = 0
	useColor writer = 1
)

func (w writer) addNormal(sb *strings.Builder, t string) {
	addNormal(sb, t)
}
func (w writer) addGreen(sb *strings.Builder, t string) {
	if w == noColor {
		addNormal(sb, t)
		return
	}
	addGreen(sb, t)
}
func (w writer) addRed(sb *strings.Builder, t string) {
	if w == noColor {
		addNormal(sb, t)
		return
	}
	addRed(sb, t)
}
func (w writer) addYellow(sb *strings.Builder, t string) {
	if w == noColor {
		addNormal(sb, t)
		return
	}
	addYellow(sb, t)
}
func (w writer) addDim(sb *strings.Builder, t string) {
	if w == noColor {
		addNormal(sb, t)
		return
	}
	addDim(sb, t)
}

func (p Puzzle) numNodes() int {
	return p.numEdges() + 1
}

func (p Puzzle) String() string {
	return p.string(true, useColor)
}

func (p Puzzle) Solution() string {
	return p.string(false, useColor)
}

func (p Puzzle) BlandString() string {
	return p.string(true, noColor)
}

func (p Puzzle) BlandSolution() string {
	return p.string(false, noColor)
}

func (p Puzzle) string(
	includeXs bool,
	cw colorWriter,
) string {

	isBland := cw == noColor

	var sb strings.Builder

	sb.WriteString("\n")

	if !isBland {
		cw.addNormal(&sb, `    `)
		for c := 0; c < p.numNodes(); c++ {
			cw.addNormal(&sb, fmt.Sprintf(" c%2d ", c))
			cw.addNormal(&sb, `   `)
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
	}

	for r := 0; r < p.numNodes(); r++ {
		var below strings.Builder
		if !isBland {
			cw.addNormal(&sb, fmt.Sprintf("r%2d ", r))
			cw.addNormal(&below, `    `)
		}

		for c := 0; c < p.numNodes(); c++ {
			nc := model.NewCoordFromInts(r, c)
			// write a node
			if n, ok := p.gn.GetNode(nc); ok {
				nOut, isMax := getSumOutgoingStraightLines(nc, &p.edges)
				if nOut == n.Value() {
					cw.addGreen(&sb, n.PrettyString())
				} else if nOut > n.Value() {
					cw.addRed(&sb, n.PrettyString())
				} else if isMax {
					cw.addRed(&sb, n.PrettyString())
				} else if nOut, nAvoid := getNumOutEdges(nc, &p.edges); nOut > 2 {
					cw.addRed(&sb, n.PrettyString())
				} else if nOut == 2 {
					cw.addYellow(&sb, n.PrettyString())
				} else if nAvoid+nOut == 4 || nAvoid == 3 {
					cw.addRed(&sb, n.PrettyString())
				} else {
					cw.addNormal(&sb, n.PrettyString())
				}
			} else {
				if nOut, nAvoid := getNumOutEdges(nc, &p.edges); nOut > 2 {
					cw.addRed(&sb, emptyNode)
				} else if nOut == 2 {
					cw.addGreen(&sb, emptyNode)
				} else if (nOut != 0 && nAvoid+nOut == 4) || nAvoid == 3 {
					cw.addRed(&sb, emptyNode)
				} else if nOut == 1 {
					cw.addYellow(&sb, emptyNode)
				} else {
					cw.addNormal(&sb, emptyNode)
				}
			}

			// now draw an edge
			ep := model.NewEdgePair(nc, model.HeadRight)
			if p.edges.IsInBounds(ep) {
				switch p.edges.GetEdge(ep) {
				case model.EdgeExists:
					cw.addGreen(&sb, `---`)
				case model.EdgeAvoided:
					if includeXs {
						cw.addDim(&sb, ` X `)
					} else {
						sb.WriteString(`   `)
					}
				case model.EdgeUnknown:
					sb.WriteString(`   `)
				default:
					cw.addRed(&sb, `???`)
				}
			}

			// now draw any edges that are below
			below.WriteString(` `)
			ep = model.NewEdgePair(nc, model.HeadDown)
			if p.edges.IsInBounds(ep) {
				switch p.edges.GetEdge(ep) {
				case model.EdgeExists:
					cw.addGreen(&below, ` | `)
				case model.EdgeAvoided:
					if includeXs {
						cw.addDim(&below, ` X `)
					} else {
						below.WriteString(`   `)
					}
				case model.EdgeUnknown:
					below.WriteString(`   `)
				default:
					cw.addRed(&below, `???`)
				}
			}
			below.WriteString(`    `)
		}
		if !isBland {
			cw.addNormal(&sb, fmt.Sprintf(" r%2d", r))
			cw.addNormal(&below, `    `)
		}

		sb.WriteString("\n")
		sb.WriteString(below.String())
		sb.WriteString("\n")
	}

	if !isBland {
		cw.addNormal(&sb, `    `)
		for c := 0; c < p.numNodes(); c++ {
			cw.addNormal(&sb, fmt.Sprintf(" c%2d ", c))
			cw.addNormal(&sb, `   `)
		}
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
