package logic

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

func (rs *RuleSet) AddAllTwoArmRules(
	nm model.NodeMeta,
) {

	maxArms := model.GetMaxArmsByDir(nm.TwoArmOptions)
	allEPs := make([]model.EdgePair, 0, 16)
	for dir, maxLen := range maxArms {
		start := nm.Coord()
		ep := model.NewEdgePair(start, dir)
		for i := int8(0); i <= maxLen; i++ {
			allEPs = append(allEPs, ep)
			ep = ep.Next(dir)
		}
	}

	for _, ep := range allEPs {
		rs.Get(ep).addInterestedNode(nm.Node)

		perps := ep.Perpendiculars()
		ep2End := ep.NodeCoord.Translate(ep.Cardinal)

		rs.Get(model.NewEdgePair(ep.NodeCoord, perps[0])).addInterestedNode(nm.Node)
		rs.Get(model.NewEdgePair(ep.NodeCoord, perps[1])).addInterestedNode(nm.Node)
		rs.Get(model.NewEdgePair(ep2End, perps[0])).addInterestedNode(nm.Node)
		rs.Get(model.NewEdgePair(ep2End, perps[1])).addInterestedNode(nm.Node)
	}
}
