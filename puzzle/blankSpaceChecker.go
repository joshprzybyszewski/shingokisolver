package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

// TODO I don't think I should use this code because my
// loose end detection isn't advanced enough
func CheckBlankSpace(
	puzz Puzzle,
	ep model.EdgePair,
) bool {

	b := bfs{
		ge:        &puzz.edges,
		toLookAt:  make([]model.NodeCoord, 0, 16),
		seen:      make(map[model.NodeCoord]struct{}, 16),
		looseEnds: make(map[model.NodeCoord]struct{}, 4),
	}

	start := ep.NodeCoord
	if !b.isLooseEnd(start) {
		start = ep.NodeCoord.Translate(ep.Cardinal)
		if !b.isLooseEnd(start) {
			return true
		}
	}

	b.start = start
	b.toLookAt = append(b.toLookAt, start)

	for len(b.toLookAt) > 0 {
		nc := b.toLookAt[0]
		b.lookAt(nc)
		b.toLookAt = b.toLookAt[1:]
	}

	return len(b.looseEnds)%2 == 0
}

type bfs struct {
	seen      map[model.NodeCoord]struct{}
	looseEnds map[model.NodeCoord]struct{}

	ge model.GetEdger

	toLookAt []model.NodeCoord
	start    model.NodeCoord
}

func (b *bfs) lookAt(nc model.NodeCoord) {
	if _, ok := b.seen[nc]; ok {
		return
	}

	b.seen[nc] = struct{}{}
	if b.isLooseEnd(nc) {
		b.looseEnds[nc] = struct{}{}
		if nc != b.start {
			// we only care about reaching other loose ends. Once we've found
			// a loose end, we don't need to iterate out from there
			return
		}
	}

	for _, dir := range model.AllCardinals {
		if b.ge.GetEdge(model.NewEdgePair(nc, dir)) == model.EdgeUnknown {
			b.toLookAt = append(b.toLookAt, nc.Translate(dir))
		}
	}
}

func (b *bfs) isLooseEnd(nc model.NodeCoord) bool {
	// TODO this could be way faster
	numEdges := 0
	for _, dir := range model.AllCardinals {
		if b.ge.IsEdge(model.NewEdgePair(nc, dir)) {
			numEdges++
		}
	}
	return numEdges == 1
}
