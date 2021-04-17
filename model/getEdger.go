package model

type GetEdger interface {
	GetEdge(EdgePair) EdgeState
	IsEdge(EdgePair) bool
	IsAvoided(EdgePair) bool

	AllExist(NodeCoord, Arm) bool
	Any(NodeCoord, Arm) (bool, bool)
	AnyAvoided(NodeCoord, Arm) bool
}
