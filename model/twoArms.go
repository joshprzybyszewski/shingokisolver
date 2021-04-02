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
	optionsCache = map[Node][]TwoArms{}
)

func BuildTwoArmOptions(n Node) []TwoArms {
	if answer, ok := optionsCache[n]; ok {
		return answer
	}

	var options []TwoArms

	for _, feeler1 := range AllCardinals {
		for _, feeler2 := range AllCardinals {
			if n.IsInvalidMotions(feeler1, feeler2) {
				continue
			}

			for arm1 := int8(1); arm1 < n.Value(); arm1++ {
				options = append(options, TwoArms{
					One: Arm{
						Len:     arm1,
						Heading: feeler1,
					},
					Two: Arm{
						Len:     n.Value() - arm1,
						Heading: feeler2,
					},
				})
			}
		}
	}

	optionsCache[n] = options

	return options
}
