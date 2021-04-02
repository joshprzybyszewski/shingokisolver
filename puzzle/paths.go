package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

type paths struct {
	connections map[model.NodeCoord]model.NodeCoord
}

func (p *paths) getLooseEnd() (model.NodeCoord, bool) {
	var le model.NodeCoord
	if p.numLooseEnds() == 0 {
		return le, false
	}

	for k := range p.connections {
		le = k
		break
	}
	return le, true
}

func (p *paths) numLooseEnds() int {
	return len(p.connections)
}

func (p *paths) copy(maxSize int) *paths {
	connCpy := make(map[model.NodeCoord]model.NodeCoord, maxSize)
	for k, v := range p.connections {
		connCpy[k] = v
	}
	return &paths{
		connections: connCpy,
	}
}

func (p *paths) add(start, end model.NodeCoord) {
	startBuddy, hasStart := p.connections[start]

	endBuddy, hasEnd := p.connections[end]

	if hasStart && hasEnd {
		delete(p.connections, start)
		delete(p.connections, end)
		if startBuddy == endBuddy {
			// we've connected a loop. delete both of the buddies.
			delete(p.connections, startBuddy)
			delete(p.connections, endBuddy)
			return
		}

		// we've connected two disparate paths.
		// notify the respective ends.
		p.connections[startBuddy] = endBuddy
		p.connections[endBuddy] = startBuddy
		return
	}

	if hasStart {
		delete(p.connections, start)
		if startBuddy == end {
			// we've connected a loop. delete both of the entries.
			delete(p.connections, startBuddy)
			delete(p.connections, end)
			return
		}

		p.connections[startBuddy] = end
		p.connections[end] = startBuddy
		return
	}

	if hasEnd {
		delete(p.connections, end)
		if start == endBuddy {
			// we've connected a loop. delete both of the entries.
			delete(p.connections, start)
			delete(p.connections, endBuddy)
			return
		}

		p.connections[start] = endBuddy
		p.connections[endBuddy] = start
		return
	}

	p.connections[start] = end
	p.connections[end] = start
}
