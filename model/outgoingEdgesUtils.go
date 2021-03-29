package model

func UpdateGridConnections(
	grid GridSetterAndGetter,
	startCoord, endCoord NodeCoord,
	motion Cardinal,
) {
	start := grid.Get(startCoord)
	end := grid.Get(endCoord)

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

	grid.Set(startCoord, start)
	grid.Set(endCoord, end)

	switch motion {
	case HeadLeft:
		updateRowConnections(grid, endCoord, startCoord)
	case HeadRight:
		updateRowConnections(grid, startCoord, endCoord)
	case HeadUp:
		updateColConnections(grid, endCoord, startCoord)
	case HeadDown:
		updateColConnections(grid, startCoord, endCoord)
	}
}

func updateRowConnections(
	grid GridSetterAndGetter,
	leftNode, rightNode NodeCoord,
) {
	curCoord := rightNode
	cur := grid.Get(curCoord)
	for cur.IsRight() {
		nextCoord := curCoord.Translate(HeadRight)
		if !grid.IsInBounds(nextCoord) {
			break
		}
		next := grid.Get(nextCoord)
		next.left = cur.left + 1
		grid.Set(nextCoord, next)

		cur = next
		curCoord = nextCoord
	}

	curCoord = leftNode
	cur = grid.Get(curCoord)
	for cur.IsLeft() {
		nextCoord := curCoord.Translate(HeadLeft)
		if !grid.IsInBounds(nextCoord) {
			break
		}
		next := grid.Get(nextCoord)
		next.right = cur.right + 1
		grid.Set(nextCoord, next)

		cur = next
		curCoord = nextCoord
	}
}

func updateColConnections(
	grid GridSetterAndGetter,
	topNode, bottomNode NodeCoord,
) {
	curCoord := topNode
	cur := grid.Get(curCoord)
	for cur.IsAbove() {
		nextCoord := curCoord.Translate(HeadUp)
		if !grid.IsInBounds(nextCoord) {
			break
		}
		next := grid.Get(nextCoord)
		next.below = cur.below + 1
		grid.Set(nextCoord, next)

		cur = next
		curCoord = nextCoord
	}

	curCoord = bottomNode
	cur = grid.Get(curCoord)
	for cur.IsBelow() {
		nextCoord := curCoord.Translate(HeadDown)
		if !grid.IsInBounds(nextCoord) {
			break
		}
		next := grid.Get(nextCoord)
		next.above = cur.above + 1
		grid.Set(nextCoord, next)

		cur = next
		curCoord = nextCoord
	}
}
