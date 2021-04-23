package model

type TwoArms struct {
	One Arm
	Two Arm
}

func (ta TwoArms) AfterOne(start NodeCoord) EdgePair {
	return NewEdgePair(
		start.TranslateAlongArm(ta.One),
		ta.One.Heading,
	)
}

func (ta TwoArms) AfterTwo(start NodeCoord) EdgePair {
	return NewEdgePair(
		start.TranslateAlongArm(ta.Two),
		ta.Two.Heading,
	)
}

func (ta TwoArms) GetAllEdges(start NodeCoord) []EdgePair {
	allEdges := make([]EdgePair, 0, ta.One.Len+ta.Two.Len)

	ep := NewEdgePair(start, ta.One.Heading)
	for i := int8(0); i < ta.One.Len; i++ {
		allEdges = append(allEdges, ep)
		ep = ep.Next(ta.One.Heading)
	}

	ep = NewEdgePair(start, ta.Two.Heading)
	for i := int8(0); i < ta.Two.Len; i++ {
		allEdges = append(allEdges, ep)
		ep = ep.Next(ta.Two.Heading)
	}

	return allEdges
}

var (
	// puzzle size to Node to options
	optionsCache map[int]map[Node][]TwoArms = make(map[int]map[Node][]TwoArms, 25)
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

	options := longBuildTwoArms(n, numEdges)

	optionsCache[numEdges][n] = options

	return options
}

func longBuildTwoArms(
	n Node,
	numEdges int,
) []TwoArms {

	var options []TwoArms
	contains := func(ta TwoArms) bool {
		for _, o := range options {
			if ta == o {
				return true
			}
		}
		return false
	}

	isOutOfBounds := func(n Node, a Arm) bool {
		endCoord := a.EndFrom(n.Coord())
		return endCoord.Row < 0 ||
			endCoord.Col < 0 ||
			endCoord.Row > RowIndex(numEdges) ||
			endCoord.Col > ColIndex(numEdges)
	}

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
				if isOutOfBounds(n, arm1) {
					continue
				}

				arm2 := Arm{
					Len:     n.Value() - len1,
					Heading: heading2,
				}
				if isOutOfBounds(n, arm2) {
					continue
				}

				if contains(TwoArms{
					One: arm1,
					Two: arm2,
				}) || contains(TwoArms{
					One: arm2,
					Two: arm1,
				}) {
					continue
				}

				options = append(options, TwoArms{
					One: arm1,
					Two: arm2,
				})
			}
		}
	}

	return options
}
