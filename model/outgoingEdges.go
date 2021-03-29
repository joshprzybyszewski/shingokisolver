package model

type OutgoingEdges struct {
	above int8
	below int8
	left  int8
	right int8
}

func (oe OutgoingEdges) TotalEdges() int8 {
	return oe.above + oe.below + oe.left + oe.right
}

func (oe OutgoingEdges) Above() int8 {
	return oe.above
}

func (oe OutgoingEdges) Below() int8 {
	return oe.below
}

func (oe OutgoingEdges) Left() int8 {
	return oe.left
}

func (oe OutgoingEdges) Right() int8 {
	return oe.right
}

func (oe OutgoingEdges) IsAbove() bool {
	return oe.above != 0
}

func (oe OutgoingEdges) IsBelow() bool {
	return oe.below != 0
}

func (oe OutgoingEdges) IsLeft() bool {
	return oe.left != 0
}

func (oe OutgoingEdges) IsRight() bool {
	return oe.right != 0
}

func (oe OutgoingEdges) GetNumOutgoingDirections() int8 {
	var numBranches int8

	if oe.IsAbove() {
		numBranches++
	}
	if oe.IsBelow() {
		numBranches++
	}
	if oe.IsLeft() {
		numBranches++
	}
	if oe.IsRight() {
		numBranches++
	}

	return numBranches
}
