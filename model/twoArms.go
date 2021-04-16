package model

import "fmt"

type Arm struct {
	Len     int8
	Heading Cardinal
}

func (a Arm) String() string {
	return fmt.Sprintf("Arm{Len: %d, Heading: %s}", a.Len, a.Heading)
}

func (a Arm) EndFrom(
	start NodeCoord,
) NodeCoord {
	for i := int8(0); i < a.Len; i++ {
		start = start.Translate(a.Heading)
	}
	return start
}

type TwoArms struct {
	One Arm
	Two Arm
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
