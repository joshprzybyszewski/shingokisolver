package model

type GetEdger interface {
	GetEdge(EdgePair) EdgeState
	IsEdge(EdgePair) bool
	IsAvoided(EdgePair) bool
	IsInBounds(EdgePair) bool

	AllExist(NodeCoord, Arm) bool
	Any(NodeCoord, Arm) (bool, bool)
	AnyAvoided(NodeCoord, Arm) bool

	NumEdges() int
}

type GetNoder interface {
	GetNode(NodeCoord) (Node, bool)
}
