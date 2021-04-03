package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

type paths struct {
	slice []model.NodeCoord
}

func newPaths(maxLooseEnds int) *paths {
	return &paths{
		slice: make([]model.NodeCoord, 0, maxLooseEnds),
	}
}

func (p *paths) getLooseEnd() (model.NodeCoord, bool) {
	var le model.NodeCoord
	if p.numLooseEnds() == 0 {
		return le, false
	}

	return p.slice[0], true
}

func (p *paths) numLooseEnds() int {
	return len(p.slice)
}

func (p *paths) copy(maxSize int) *paths {
	sliceCpy := make([]model.NodeCoord, len(p.slice), cap(p.slice))
	copy(sliceCpy, p.slice)
	return &paths{
		slice: sliceCpy,
	}
}

func (p *paths) add(start, end model.NodeCoord) {
	hasStart, hasEnd := false, false
	for i := 0; i < len(p.slice); i++ {
		if p.slice[i] != start && p.slice[i] != end {
			continue
		}
		if p.slice[i] == start {
			hasStart = true
		} else {
			hasEnd = true
		}
		p.slice[i] = p.slice[len(p.slice)-1]
		p.slice = p.slice[:len(p.slice)-1]
		i--
	}
	if !hasStart {
		p.slice = append(p.slice, start)
	}
	if !hasEnd {
		p.slice = append(p.slice, end)
	}
}
