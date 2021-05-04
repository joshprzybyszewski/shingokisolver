package model

func GetMaxArmsByDir(
	tas []TwoArms,
) map[Cardinal]int8 {
	if len(tas) == 0 {
		return nil
	}

	res := make(map[Cardinal]int8, 4)
	for _, ta := range tas {
		if ta.One.Len > res[ta.One.Heading] {
			res[ta.One.Heading] = ta.One.Len
		}
		if ta.Two.Len > res[ta.Two.Heading] {
			res[ta.Two.Heading] = ta.Two.Len
		}
	}
	return res
}
