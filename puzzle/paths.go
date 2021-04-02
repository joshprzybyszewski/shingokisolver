package puzzle

import "github.com/joshprzybyszewski/shingokisolver/model"

type paths struct {
	connections map[model.NodeCoord]model.NodeCoord
}

func (p *paths) copy() *paths {
	connCpy := make(map[model.NodeCoord]model.NodeCoord)
	for k, v := range p.connections {
		connCpy[k] = v
	}
	return &paths{
		connections: connCpy,
	}
}

func (p *paths) looseEnds() []model.NodeCoord {
	ends := make([]model.NodeCoord, 0, len(p.connections))
	for k := range p.connections {
		ends = append(ends, k)
	}
	return ends
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
