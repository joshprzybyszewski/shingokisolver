package model

func GetMinArmsByDir(
	ta []TwoArms,
) (map[Cardinal]int8, bool) {
	if len(ta) == 0 {
		return nil, false
	}

	// I am choosing to write more lines of code in this
	// func so that I avoid a map iteration
	dirA := ta[0].One.Heading
	lenA := ta[0].One.Len

	dirB := ta[0].Two.Heading
	lenB := ta[0].Two.Len

	for i := 1; i < len(ta) && (lenA > 0 || lenB > 0); i++ {
		switch dirA {
		case ta[i].One.Heading:
			if lenA > ta[i].One.Len {
				lenA = ta[i].One.Len
			}
		case ta[i].Two.Heading:
			if lenA > ta[i].Two.Len {
				lenA = ta[i].Two.Len
			}
		default:
			// this TwoArms option doesn't have dirA.
			lenA = 0
		}

		switch dirB {
		case ta[i].One.Heading:
			if lenB > ta[i].One.Len {
				lenB = ta[i].One.Len
			}
		case ta[i].Two.Heading:
			if lenB > ta[i].Two.Len {
				lenB = ta[i].Two.Len
			}
		default:
			// this TwoArms option doesn't have dirB.
			lenB = 0
		}
	}

	return map[Cardinal]int8{
		dirA: lenA,
		dirB: lenB,
	}, len(ta) == 1
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
