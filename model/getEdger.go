package model

type GetEdger interface {
	// TODO add nuanced methods instead
	// like IsEdge, IsAvoided, IsDefined
	GetEdge(EdgePair) EdgeState

	AllExist(NodeCoord, Arm) bool
	Any(NodeCoord, Arm) (bool, bool)
	AnyAvoided(NodeCoord, Arm) bool
}
