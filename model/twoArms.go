package model

type Arm struct {
	Len     int8
	Heading Cardinal
}

type TwoArms struct {
	One Arm
	Two Arm
}

var (
	optionsCache map[NodeType]map[int8][]TwoArms
)

func init() {
	optionsCache = map[NodeType]map[int8][]TwoArms{
		WhiteNode: make(map[int8][]TwoArms, 25),
		BlackNode: make(map[int8][]TwoArms, 25),
	}
}

func BuildTwoArmOptions(n Node) []TwoArms {
	if answer, ok := optionsCache[n.Type()][n.Value()]; ok {
		return answer
	}

	options := longBuildTwoArms(n)

	optionsCache[n.Type()][n.Value()] = options

	return options
}

func longBuildTwoArms(n Node) []TwoArms {
	var options []TwoArms
	contains := func(ta TwoArms) bool {
		for _, o := range options {
			if ta == o {
				return true
			}
		}
		return false
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
				arm2 := Arm{
					Len:     n.Value() - len1,
					Heading: heading2,
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
