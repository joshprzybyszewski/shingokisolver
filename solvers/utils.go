package solvers

import "github.com/joshprzybyszewski/shingokisolver/model"

func getLooseEndsNotOnArms(
	prev []model.NodeCoord,
	start, arm1End, arm2End model.NodeCoord,
) []model.NodeCoord {
	newLooseEnds := make([]model.NodeCoord, 0, len(prev)+2)

	endPoints := []model.NodeCoord{arm1End, arm2End}

	isBetween := func(val, inclusive, exclusive int) bool {
		if inclusive < exclusive {
			return val >= inclusive && val < exclusive
		}
		return val <= inclusive && val > exclusive
	}

	shouldInclude := func(pnc model.NodeCoord) bool {
		sameRow := pnc.Row == start.Row
		sameCol := pnc.Col == start.Col

		if !sameRow && !sameCol {
			return true
		}
		if sameRow {
			if sameCol {
				// this looseEnd matches our start node. don't add it
				return false
			}

			for _, end := range endPoints {
				if end.Row == start.Row &&
					isBetween(int(pnc.Col), int(start.Col), int(end.Col)) {
					return false
				}
			}
			return true
		}

		for _, end := range endPoints {
			if end.Col == start.Col &&
				isBetween(int(pnc.Row), int(start.Row), int(end.Row)) {
				return false
			}
		}
		return true
	}

	for _, pnc := range prev {
		if shouldInclude(pnc) {
			newLooseEnds = append(newLooseEnds, pnc)
		}
	}

	return append(newLooseEnds, endPoints...)
}
