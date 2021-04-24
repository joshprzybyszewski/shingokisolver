package model

var (
	InvalidTarget = Target{
		Node: Node{
			coord: InvalidNodeCoord,
		},
	}
)

type Target struct {
	Parent  *Target
	Options []TwoArms
	Node    Node
}
