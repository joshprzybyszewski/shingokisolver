package puzzle

func updateConnections(
	grid gridSetterAndGetter,
	startCoord, endCoord nodeCoord,
	motion cardinal,
) {
	start := grid.get(startCoord)
	end := grid.get(endCoord)

	switch motion {
	case headLeft:
		start.left = end.left + 1
		end.right = start.right + 1
	case headRight:
		start.right = end.right + 1
		end.left = start.left + 1
	case headUp:
		start.above = end.above + 1
		end.below = start.below + 1
	case headDown:
		start.below = end.below + 1
		end.above = start.above + 1
	}

	grid.set(startCoord, start)
	grid.set(endCoord, end)

	switch motion {
	case headLeft:
		updateRowConnections(grid, endCoord, startCoord)
	case headRight:
		updateRowConnections(grid, startCoord, endCoord)
	case headUp:
		updateColConnections(grid, endCoord, startCoord)
	case headDown:
		updateColConnections(grid, startCoord, endCoord)
	}
}

func updateRowConnections(
	grid gridSetterAndGetter,
	leftNode, rightNode nodeCoord,
) {
	curCoord := rightNode
	cur := grid.get(curCoord)
	for cur.isright() {
		nextCoord := curCoord.translate(headRight)
		if !grid.isInBounds(nextCoord) {
			break
		}
		next := grid.get(nextCoord)
		next.left = cur.left + 1
		grid.set(nextCoord, next)

		cur = next
		curCoord = nextCoord
	}

	curCoord = leftNode
	cur = grid.get(curCoord)
	for cur.isleft() {
		nextCoord := curCoord.translate(headLeft)
		if !grid.isInBounds(nextCoord) {
			break
		}
		next := grid.get(nextCoord)
		next.right = cur.right + 1
		grid.set(nextCoord, next)

		cur = next
		curCoord = nextCoord
	}
}

func updateColConnections(
	grid gridSetterAndGetter,
	topNode, bottomNode nodeCoord,
) {
	curCoord := topNode
	cur := grid.get(curCoord)
	for cur.isabove() {
		nextCoord := curCoord.translate(headUp)
		if !grid.isInBounds(nextCoord) {
			break
		}
		next := grid.get(nextCoord)
		next.below = cur.below + 1
		grid.set(nextCoord, next)

		cur = next
		curCoord = nextCoord
	}

	curCoord = bottomNode
	cur = grid.get(curCoord)
	for cur.isbelow() {
		nextCoord := curCoord.translate(headDown)
		if !grid.isInBounds(nextCoord) {
			break
		}
		next := grid.get(nextCoord)
		next.above = cur.above + 1
		grid.set(nextCoord, next)

		cur = next
		curCoord = nextCoord
	}
}
