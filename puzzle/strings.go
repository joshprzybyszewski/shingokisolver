package puzzle

import (
	"fmt"
	"strings"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (p *Puzzle) String() string {
	return p.string(true, true)
}

func (p *Puzzle) Solution() string {
	return p.string(false, false)
}

func (p *Puzzle) string(
	includeXs bool,
	includeQueue bool,
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
				sb.WriteString(`   `)
			}
			sb.WriteString(`)`)

			// now draw an edge
			ep := model.NewEdgePair(nc, model.HeadRight)
			if p.edges.isInBounds(ep) {
				// enQueued := includeQueue && ep.IsIn(p.rq.toCheck...)
				_, ok := p.rq.toCheck[ep]
				enQueued := includeQueue && ok
				if enQueued {
					sb.WriteString(`qq`)
					// sb.WriteString(fmt.Sprintf("%2d", ep.IndexOf(p.rq.toCheck...)))
				}
				switch p.edges.GetEdge(ep) {
				case model.EdgeExists:
					if !enQueued {
						sb.WriteString(`--`)
					}
					sb.WriteString(`-`)
				case model.EdgeAvoided:
					if enQueued {
						if includeXs {
							sb.WriteString(`X`)
						} else {
							sb.WriteString(` `)
						}
					} else if includeXs {
						sb.WriteString(` X `)
					} else {
						sb.WriteString(`   `)
					}
				case model.EdgeUnknown:
					if !enQueued {
						sb.WriteString(`  `)
					}
					sb.WriteString(` `)
				default:
					if !enQueued {
						sb.WriteString(`??`)
					}
					sb.WriteString(`?`)
				}
			}

			// now draw any edges that are below
			below.WriteString(` `)
			ep = model.NewEdgePair(nc, model.HeadDown)
			if p.edges.isInBounds(ep) {
				// enQueued := includeQueue && ep.IsIn(p.rq.toCheck...)
				_, ok := p.rq.toCheck[ep]
				enQueued := includeQueue && ok
				if enQueued {
					below.WriteString(`qq`)
					// below.WriteString(fmt.Sprintf("%2d", ep.IndexOf(p.rq.toCheck...)))
				}

				switch p.edges.GetEdge(ep) {
				case model.EdgeExists:
					if !enQueued {
						below.WriteString(` `)
					}
					below.WriteString(`|`)
					if !enQueued {
						below.WriteString(` `)
					}
				case model.EdgeAvoided:
					if enQueued {
						if includeXs {
							below.WriteString(`X`)
						} else {
							below.WriteString(` `)
						}
					} else if includeXs {
						below.WriteString(` X `)
					} else {
						below.WriteString(`   `)
					}
				case model.EdgeUnknown:
					if !enQueued {
						below.WriteString(`  `)
					}
					below.WriteString(` `)
				default:
					if !enQueued {
						below.WriteString(`??`)
					}
					below.WriteString(`?`)
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
