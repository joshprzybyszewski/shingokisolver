package model

import "fmt"

type NodeType uint8

const (
	// no constraints
	noNode NodeType = 0
	// must be passed through in a straight line
	WhiteNode NodeType = 1
	// must be turned upon
	BlackNode NodeType = 2
)

func (nt NodeType) String() string {
	switch nt {
	case WhiteNode:
		return `w`
	case BlackNode:
		return `b`
	default:
		return `unknown NodeType`
	}
}

func (nt NodeType) isInvalidMotions(c1, c2 Cardinal) bool {
	if c1 == c2 {
		return true
	}

	switch nt {
	case WhiteNode:
		// white nodes need to be straight. therefore, they're
		// invalid if they have opposing directions set
		return (c1+c2)%2 == 1
	case BlackNode:
		// black nodes need to be bent. therefore, they're
		// invalid if they have a straight line in them
		// This is an optimization because of the defined values
		// of each of the cardinals
		return (c1+c2)%2 == 0
	default:
		return true
	}
}

type Node struct {
	coord NodeCoord
	nType NodeType
	val   int8
}

func NewNode(coord NodeCoord, isWhite bool, v int8) Node {
	nt := BlackNode
	if isWhite {
		nt = WhiteNode
	}
	return Node{
		coord: coord,
		nType: nt,
		val:   v,
	}
}

func (n Node) String() string {
	return fmt.Sprintf("%s%2d @ %+v", n.nType, n.val, n.coord)
}

func (n Node) Coord() NodeCoord {
	return n.coord
}

func (n Node) Type() NodeType {
	return n.nType
}

func (n Node) Value() int8 {
	return n.val
}

func (n Node) Copy() Node {
	return n
}

func (n Node) IsInvalidMotions(c1, c2 Cardinal) bool {
	return n.Type().isInvalidMotions(c1, c2)
}
