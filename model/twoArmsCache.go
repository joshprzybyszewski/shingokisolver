package model

var (
	// puzzle size to Node to options
	// should allow for state.MaxEdges
	optionsCache map[int]map[Node][]TwoArms = make(map[int]map[Node][]TwoArms, 50)
)

func BuildTwoArmOptions(
	n Node,
	numEdges int,
) []TwoArms {

	if _, ok := optionsCache[numEdges]; !ok {
		optionsCache[numEdges] = make(map[Node][]TwoArms, numEdges*numEdges/4)
	}

	if answer, ok := optionsCache[numEdges][n]; ok {
		return answer
	}

	options := longBuildTwoArms(n, int8(numEdges))

	optionsCache[numEdges][n] = options

	return options
}

func longBuildTwoArms(
	n Node,
	numEdges int8,
) []TwoArms {

	var options []TwoArms

	for _, heading1 := range AllCardinals {
		for _, heading2 := range AllCardinals {
			if n.IsInvalidMotions(heading1, heading2) {
				continue
			}

			for len1 := int8(1); len1 < n.Value(); len1++ {

				arm1 := Arm{
					Len:     len1,
					Heading: heading1,
				}
				if isOutOfBounds(n, numEdges, arm1) {
					continue
				}

				arm2 := Arm{
					Len:     n.Value() - len1,
					Heading: heading2,
				}
				if isOutOfBounds(n, numEdges, arm2) {
					continue
				}

				ta := TwoArms{
					One: arm1,
					Two: arm2,
				}
				if ContainsTwoArms(options, ta) {
					continue
				}

				options = append(options, ta)
			}
		}
	}

	return options
}

func isOutOfBounds(n Node, numEdges int8, a Arm) bool {
	endCoord := a.EndFrom(n.Coord())
	return endCoord.Row < 0 ||
		endCoord.Col < 0 ||
		endCoord.Row > RowIndex(numEdges) ||
		endCoord.Col > ColIndex(numEdges)
}
