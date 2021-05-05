package model

type NodeMeta struct {
	Node

	Nearby        NearbyNodes
	TwoArmOptions []TwoArms

	IsComplete bool
}

func (nm *NodeMeta) Copy() *NodeMeta {
	if nm == nil {
		return nil
	}
	if nm.IsComplete {
		// instead of allocing another meta to heap,
		// we can just return this pointer because
		// we only care about Node and IsComplete
		return nm
	}

	tao := make([]TwoArms, len(nm.TwoArmOptions))
	copy(tao, nm.TwoArmOptions)

	return &NodeMeta{
		Node:          nm.Node,
		Nearby:        nm.Nearby,
		TwoArmOptions: tao,
	}
}

func (nm *NodeMeta) Filter(ge GetEdger) State {
	if nm == nil {
		return Incomplete
	}
	if nm.IsComplete {
		return Complete
	}

	nm.TwoArmOptions = nm.GetFilteredOptions(
		nm.TwoArmOptions,
		ge,
		nm.Nearby,
	)

	return nm.CheckComplete(ge)
}

func (nm *NodeMeta) CheckComplete(ge GetEdger) State {
	if nm.IsComplete {
		return Complete
	}

	switch len(nm.TwoArmOptions) {
	case 0:
		return Violation
	case 1:
	default:
		return Incomplete
	}

	if ge.AllExist(nm.Node.Coord(), nm.TwoArmOptions[0].One) &&
		ge.IsAvoided(nm.TwoArmOptions[0].AfterOne(nm.Node.Coord())) &&
		ge.AllExist(nm.Node.Coord(), nm.TwoArmOptions[0].Two) &&
		ge.IsAvoided(nm.TwoArmOptions[0].AfterTwo(nm.Node.Coord())) {

		nm.IsComplete = true
		nm.TwoArmOptions = nil

		return Complete
	}

	// there's only one option, but not all of the necessary edges exist
	return Incomplete
}
