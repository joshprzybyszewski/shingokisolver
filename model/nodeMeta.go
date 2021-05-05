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
		return &NodeMeta{
			Node:       nm.Node,
			IsComplete: true,
		}
	}

	tao := make([]TwoArms, len(nm.TwoArmOptions))
	copy(tao, nm.TwoArmOptions)

	return &NodeMeta{
		Node:          nm.Node,
		Nearby:        nm.Nearby, // TODO copy nearby?
		TwoArmOptions: tao,
	}
}

func (nm *NodeMeta) Filter(ge GetEdger) State {
	if nm == nil || nm.IsComplete {
		return Incomplete
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
		return Complete
	}

	// there's only one option, but not all of the necessary edges exist
	return Incomplete
}
