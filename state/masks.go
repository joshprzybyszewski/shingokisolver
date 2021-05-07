package state

import "github.com/joshprzybyszewski/shingokisolver/model"

var (
	masks       = make([]bitData, MaxEdges+1)
	armLenMasks = make([]bitData, len(masks)+1)
)

func init() {
	buildMasks()
	buildArmLenMasks()
}

func buildMasks() {
	for i := 0; i < len(masks); i++ {
		masks[i] = 1 << i
	}
}

func buildArmLenMasks() {
	for i := 0; i < len(armLenMasks); i++ {
		var lenMask bitData
		for mi := 0; mi < i; mi++ {
			lenMask |= masks[mi]
		}
		armLenMasks[i] = lenMask
	}
}

func getMask(
	start model.NodeCoord,
	arm model.Arm,
) bitData {
	switch arm.Heading {
	case model.HeadRight:
		return armLenMasks[arm.Len] << start.Col

	case model.HeadDown:
		return armLenMasks[arm.Len] << start.Row

	case model.HeadLeft:
		lmi := int(arm.Len)
		shift := start.Col - model.ColIndex(arm.Len)
		if shift < 0 {
			lmi += int(shift)
			if lmi < 0 {
				return 0
			}
			shift = 0
		}
		return armLenMasks[lmi] << shift

	case model.HeadUp:
		lmi := int(arm.Len)
		shift := start.Row - model.RowIndex(arm.Len)
		if shift < 0 {
			lmi += int(shift)
			if lmi < 0 {
				return 0
			}
			shift = 0
		}
		return armLenMasks[lmi] << shift

	default:
		return 0
	}
}
