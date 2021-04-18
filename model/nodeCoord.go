package model

import "fmt"

var (
	InvalidNodeCoord NodeCoord = NodeCoord{
		Row: -1,
		Col: -1,
	}
)

type RowIndex int16
type ColIndex int16

type NodeCoord struct {
	Row RowIndex
	Col ColIndex
}

func NewCoord(r RowIndex, c ColIndex) NodeCoord {
	return NodeCoord{
		Row: r,
		Col: c,
	}
}

func (nc NodeCoord) String() string {
	return fmt.Sprintf("NodeCoord{Row: %d, Col: %d}", nc.Row, nc.Col)
}

func (nc NodeCoord) DistanceTo(
	other NodeCoord,
) int {
	dist := int(nc.Row - other.Row)
	if dist < 0 {
		dist *= -1
	}
	if nc.Col < other.Col {
		return dist + int(other.Col-nc.Col)
	} else {
		return dist + int(nc.Col-other.Col)
	}
}

func NewCoordFromInts(r, c int) NodeCoord {
	return NodeCoord{
		Row: RowIndex(r),
		Col: ColIndex(c),
	}
}

func (nc NodeCoord) Translate(
	move Cardinal,
) NodeCoord {
	switch move {
	case HeadUp:
		nc.Row--
	case HeadDown:
		nc.Row++
	case HeadLeft:
		nc.Col--
	case HeadRight:
		nc.Col++
	}
	return nc
}

func (nc NodeCoord) TranslateAlongArm(
	arm Arm,
) NodeCoord {
	switch arm.Heading {
	case HeadUp:
		nc.Row -= RowIndex(arm.Len)
	case HeadDown:
		nc.Row += RowIndex(arm.Len)
	case HeadLeft:
		nc.Col -= ColIndex(arm.Len)
	case HeadRight:
		nc.Col += ColIndex(arm.Len)
	}
	return nc
}
