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

func (ta TwoArms) GetAllEdges(start NodeCoord) ([]EdgePair, []EdgePair) {
	existing := make([]EdgePair, 0, ta.One.Len+ta.Two.Len)
	avoided := make([]EdgePair, 0, 2)

	ep := NewEdgePair(start, ta.One.Heading)
	for i := int8(0); i < ta.One.Len; i++ {
		existing = append(existing, ep)
		ep = ep.Next(ta.One.Heading)
	}
	avoided = append(avoided, ep)

	ep = NewEdgePair(start, ta.Two.Heading)
	for i := int8(0); i < ta.Two.Len; i++ {
		existing = append(existing, ep)
		ep = ep.Next(ta.Two.Heading)
	}
	avoided = append(avoided, ep)

	return existing, avoided
}

func (ta TwoArms) equals(other TwoArms) bool {
	if ta == other {
		return true
	}
	return ta.One == other.Two && ta.Two == other.One
}

func ContainsTwoArms(tas []TwoArms, other TwoArms) bool {
	for _, ta := range tas {
		if ta.equals(other) {
			return true
		}
	}
	return false
}
