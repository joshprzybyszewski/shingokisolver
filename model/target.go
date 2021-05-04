package model

var (
	InvalidTarget = Target{
		Node: InvalidNode,
	}
)

type Target struct {
	Parent  *Target
	Options []TwoArms
	Node    Node
}
