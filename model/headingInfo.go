package model

type headingInfo struct {
	armLensDoubleAvoids        []bool
	maxArmLenUntilIncomingEdge int8
}

func buildHeadingInfos(
	node Node,
	ge GetEdger,
) []headingInfo {
	// TODO this func is expensive!
	// Is there a way to cache it for re-use, or to
	// harvest the data in a more performant fashion?
	res := make([]headingInfo, len(AllCardinals)+1)

	nc := node.Coord()
	maxLen := node.Value() - 1

	for _, heading := range AllCardinals {
		hInfo := headingInfo{
			armLensDoubleAvoids: make([]bool, maxLen+1),
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
				hInfo.armLensDoubleAvoids[i] = true
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

	if hi.armLensDoubleAvoids[arm.Len] {
		return false
	}

	return true
}
