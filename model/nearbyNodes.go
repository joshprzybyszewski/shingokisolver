package model

// The first slice is indexed as a Cardinal.
// The second slice is indexed as distance.
// Example: for "white 5" in the bottom left
// (b 3)
//   |
// (   )
//   |
// (w 5)---(w 5)---(   )---(   )   (   )   (b 7)
// Results:
// [][]*Node{
//   nil, // for HeadNowhere = 0
//   []*Node{nil, "w 5", nil, nil, nil, "b 7"}, // for HeadRight = 1
//   []*Node{nil, nil, "b 3"}, // for HeadUp = 2
//   nil, // for HeadLeft = 3
//   nil, // for HeadDown = 4
// }
type NearbyNodes [][]*Node

func (nn NearbyNodes) Get(c Cardinal) []*Node {
	return nn[c]
}
