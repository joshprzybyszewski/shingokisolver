package puzzle

type getEdgesFromNoder interface {
	getOutgoingEdgesFrom(nodeCoord) (edgesFromNode, bool)
}

type walker interface {
	walk() (map[nodeCoord]struct{}, bool)
}

type simpleWalker struct {
	provider   getEdgesFromNoder
	start      nodeCoord
	cur        nodeCoord
	seen       map[nodeCoord]struct{}
	foundStart bool
}

func newWalker(
	getEFNer getEdgesFromNoder,
	start nodeCoord,
) walker {
	return &simpleWalker{
		provider: getEFNer,
		start:    start,
		cur:      start,
		seen:     map[nodeCoord]struct{}{},
	}
}

func (sw *simpleWalker) walk() (map[nodeCoord]struct{}, bool) {
	move := sw.walkToNextPoint(headNowhere)
	if move == headNowhere {
		// our path all the way around was incomplete
		return nil, false
	}

	for sw.cur.row != sw.start.row || sw.cur.col != sw.start.col {
		move = sw.walkToNextPoint(move)
		if move == headNowhere || sw.foundStart {
			// if we can't go anywhere, then we'll break out of the loop
			// because this means we don't have a loop.
			break
		}
	}

	if sw.foundStart || sw.cur.row == sw.start.row && sw.cur.col == sw.start.col {
		return sw.seen, true
	}

	return nil, false
}

func (sw *simpleWalker) walkToNextPoint(
	avoid cardinal,
) cardinal {

	efn, ok := sw.provider.getOutgoingEdgesFrom(sw.cur)
	if !ok {
		return headNowhere
	}

	if efn.isabove() && avoid != headUp {
		nextRow := sw.cur.row - rowIndex(efn.above)
		sw.markNodesAsSeen(nextRow, sw.cur.row, sw.cur.col, sw.cur.col)
		sw.cur.row = nextRow
		return headDown
	}

	if efn.isleft() && avoid != headLeft {
		nextCol := sw.cur.col - colIndex(efn.left)
		sw.markNodesAsSeen(sw.cur.row, sw.cur.row, nextCol, sw.cur.col)
		sw.cur.col = nextCol
		return headRight
	}

	if efn.isbelow() && avoid != headDown {
		nextRow := sw.cur.row + rowIndex(efn.below)
		sw.markNodesAsSeen(sw.cur.row, nextRow, sw.cur.col, sw.cur.col)
		sw.cur.row = nextRow
		return headUp
	}

	if efn.isright() && avoid != headRight {
		nextCol := sw.cur.col + colIndex(efn.right)
		sw.markNodesAsSeen(sw.cur.row, sw.cur.row, sw.cur.col, nextCol)
		sw.cur.col = nextCol
		return headLeft
	}

	return headNowhere
}

func (sw *simpleWalker) markNodesAsSeen(
	minR, maxR rowIndex,
	minC, maxC colIndex,
) {
	for r := minR; r <= maxR; r++ {
		for c := minC; c <= maxC; c++ {
			nc := nodeCoord{
				row: r,
				col: c,
			}
			if nc == sw.start {
				if _, ok := sw.seen[nc]; ok {
					// this means we've seen the starting node before
					sw.foundStart = true
				}
			}
			sw.seen[nc] = struct{}{}
		}
	}
}
