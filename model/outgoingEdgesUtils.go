package model

const (
	rightShift int = 0
	aboveShift int = 8
	leftShift  int = 16
	belowShift int = 24
)

const (
	defaultExpectedOneDirection = 4
	defaultTotalExpectedUpdates = 2 + defaultExpectedOneDirection + defaultExpectedOneDirection
)

type gridUpdate struct {
	coord  NodeCoord
	newVal OutgoingEdges
}

func UpdateGridConnections(
	grid Grid,
	motion Cardinal,
	startCoord, endCoord NodeCoord,
) Grid {
	start := grid.Get(startCoord)
	end := grid.Get(endCoord)

	updates := make([]gridUpdate, 0, defaultTotalExpectedUpdates)

	switch motion {
	case HeadLeft:
		start.left = end.left + 1
		end.right = start.right + 1
	case HeadRight:
		start.right = end.right + 1
		end.left = start.left + 1
	case HeadUp:
		start.above = end.above + 1
		end.below = start.below + 1
	case HeadDown:
		start.below = end.below + 1
		end.above = start.above + 1
	}

	updates = append(updates,
		gridUpdate{
			coord:  startCoord,
			newVal: start,
		},
		gridUpdate{
			coord:  endCoord,
			newVal: end,
		},
	)

	switch motion {
	case HeadLeft:
		updates = append(updates,
			getRowConnectionUpdates(
				grid,
				endCoord, startCoord,
				end, start,
			)...,
		)
	case HeadRight:
		updates = append(updates,
			getRowConnectionUpdates(
				grid,
				startCoord, endCoord,
				start, end,
			)...,
		)
	case HeadUp:
		updates = append(updates,
			getColConnectionUpdates(
				grid,
				endCoord, startCoord,
				end, start,
			)...,
		)
	case HeadDown:
		updates = append(updates,
			getColConnectionUpdates(
				grid,
				startCoord, endCoord,
				start, end,
			)...,
		)
	}

	return grid.withUpdates(updates)
}

func getRowConnectionUpdates(
	grid Grid,
	leftNode, rightNode NodeCoord,
	initialLeftVal, initialRightVal OutgoingEdges,
) []gridUpdate {
	updates := make([]gridUpdate, 0, defaultExpectedOneDirection)

	curCoord := rightNode
	cur := initialRightVal
	for cur.IsRight() {
		nextCoord := curCoord.Translate(HeadRight)
		if !grid.IsInBounds(nextCoord) {
			break
		}
		next := grid.Get(nextCoord)
		next.left = cur.left + 1
		updates = append(updates, gridUpdate{
			coord:  nextCoord,
			newVal: next,
		})

		cur = next
		curCoord = nextCoord
	}

	curCoord = leftNode
	cur = initialLeftVal
	for cur.IsLeft() {
		nextCoord := curCoord.Translate(HeadLeft)
		if !grid.IsInBounds(nextCoord) {
			break
		}
		next := grid.Get(nextCoord)
		next.right = cur.right + 1
		updates = append(updates, gridUpdate{
			coord:  nextCoord,
			newVal: next,
		})

		cur = next
		curCoord = nextCoord
	}

	return updates
}

func getColConnectionUpdates(
	grid Grid,
	topNode, bottomNode NodeCoord,
	initialTopVal, initialBottomVal OutgoingEdges,
) []gridUpdate {
	updates := make([]gridUpdate, 0, defaultExpectedOneDirection)

	curCoord := topNode
	cur := initialTopVal
	for cur.IsAbove() {
		nextCoord := curCoord.Translate(HeadUp)
		if !grid.IsInBounds(nextCoord) {
			break
		}
		next := grid.Get(nextCoord)
		next.below = cur.below + 1
		updates = append(updates, gridUpdate{
			coord:  nextCoord,
			newVal: next,
		})

		cur = next
		curCoord = nextCoord
	}

	curCoord = bottomNode
	cur = initialBottomVal
	for cur.IsBelow() {
		nextCoord := curCoord.Translate(HeadDown)
		if !grid.IsInBounds(nextCoord) {
			break
		}
		next := grid.Get(nextCoord)
		next.above = cur.above + 1
		updates = append(updates, gridUpdate{
			coord:  nextCoord,
			newVal: next,
		})

		cur = next
		curCoord = nextCoord
	}

	return updates
}

func newOutgoingEdgesFromInt32(input int32) OutgoingEdges {
	return OutgoingEdges{
		right: int8(input >> rightShift),
		above: int8(input >> aboveShift),
		left:  int8(input >> leftShift),
		below: int8(input >> belowShift),
	}
}

func outgoingEdgesToInt32(oe OutgoingEdges) int32 {
	return int32(oe.right)<<rightShift |
		int32(oe.above)<<aboveShift |
		int32(oe.left)<<leftShift |
		int32(oe.below)<<belowShift
}
