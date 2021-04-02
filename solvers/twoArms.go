package solvers

import "github.com/joshprzybyszewski/shingokisolver/model"

type twoArms struct {
	arm1Len int8
	arm2Len int8

	arm1Heading model.Cardinal
	arm2Heading model.Cardinal
}

func buildTwoArmOptions(n model.Node) []twoArms {
	// TODO keep a cache around for this...
	var options []twoArms

	for i, feeler1 := range model.AllCardinals {
		for _, feeler2 := range model.AllCardinals[i+1:] {
			if n.IsInvalidMotions(feeler1, feeler2) {
				continue
			}

			for arm1 := int8(1); arm1 <= n.Value()/2; arm1++ {
				options = append(options, twoArms{
					arm1Heading: feeler1,
					arm2Heading: feeler2,
					arm1Len:     arm1,
					arm2Len:     n.Value() - arm1,
				})
			}
		}
	}

	return options
}
