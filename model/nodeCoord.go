package model

type RowIndex int8
type ColIndex int8

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
