package model

type headingInfo struct {
	armLensDoubleAvoids        map[int8]struct{}
	maxArmLenUntilIncomingEdge int8
}

func buildHeadingInfos(
	node Node,
	ge GetEdger,
) map[Cardinal]headingInfo {
	res := make(map[Cardinal]headingInfo, 4)

	nc := node.Coord()
	maxLen := node.Value() - 1

	for _, heading := range AllCardinals {
		hInfo := headingInfo{
			armLensDoubleAvoids: make(map[int8]struct{}, 2),
		}
		cur := nc.Translate(heading)
		for i := int8(1); i <= maxLen; i++ {
			bothAvoided := true
			for _, dir := range heading.Perpendiculars() {
				if ge.IsEdge(NewEdgePair(cur, dir)) {
					hInfo.maxArmLenUntilIncomingEdge = i
					break
				}
				if !ge.IsAvoided(NewEdgePair(cur, dir)) {
					bothAvoided = false
					break
				}
			}

			if hInfo.maxArmLenUntilIncomingEdge > 0 {
				break
			}
			if bothAvoided {
				hInfo.armLensDoubleAvoids[i] = struct{}{}
			}

			cur = cur.Translate(heading)
		}

		res[heading] = hInfo
	}

	return res
}

func (hi headingInfo) isValidArm(arm Arm) bool {
	if hi.maxArmLenUntilIncomingEdge > 0 &&
		hi.maxArmLenUntilIncomingEdge < arm.Len {
		return false
	}

	if _, ok := hi.armLensDoubleAvoids[arm.Len]; ok {
		return false
	}

	return true
}
