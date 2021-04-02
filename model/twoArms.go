package model

type Arm struct {
	Len     int8
	Heading Cardinal
}

type TwoArms struct {
	One Arm
	Two Arm
}

func BuildTwoArmOptions(n Node) []TwoArms {
	// TODO keep a cache around for this...
	var options []TwoArms

	for i, feeler1 := range AllCardinals {
		for _, feeler2 := range AllCardinals[i+1:] {
			if n.IsInvalidMotions(feeler1, feeler2) {
				continue
			}

			for arm1 := int8(1); arm1 <= n.Value()/2; arm1++ {
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

	return options
}
