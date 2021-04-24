package model

func GetMinArmsByDir(
	ta []TwoArms,
) (map[Cardinal]int8, bool) {
	if len(ta) == 0 {
		return nil, false
	}

	res := make(map[Cardinal]int8, 2)
	res[ta[0].One.Heading] = ta[0].One.Len
	res[ta[0].Two.Heading] = ta[0].Two.Len

	if len(ta) == 1 {
		return res, true
	}

	for i := 1; i < len(ta) && len(res) > 0; i++ {
		for k, v := range res {
			switch k {
			case ta[i].One.Heading:
				if v > ta[i].One.Len {
					res[k] = ta[i].One.Len
				}
			case ta[i].Two.Heading:
				if v > ta[i].Two.Len {
					res[k] = ta[i].Two.Len
				}
			default:
				delete(res, k)
			}
		}
	}

	return res, false
}

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
