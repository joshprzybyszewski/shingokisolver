package solvers

import "github.com/joshprzybyszewski/shingokisolver/model"

func getLooseEnds(
	prev []model.NodeCoord,
	start, arm1End, arm2End model.NodeCoord,
) []model.NodeCoord {
	ends := make([]model.NodeCoord, len(prev), len(prev)+2)
	copy(ends, prev)

	ends = append(ends, arm1End, arm2End)

	isBetween := func(val, inclusive, exclusive int) bool {
		if inclusive < exclusive {
			return val >= inclusive && val < exclusive
		}
		return val <= inclusive && val > exclusive
	}

	endPoints := []model.NodeCoord{arm1End, arm2End}
	shouldRemove := func(nc model.NodeCoord) bool {
		sameRow := nc.Row == start.Row
		sameCol := nc.Col == start.Col
		if !sameRow && !sameCol {
			return false
		}

		if sameRow {
			if sameCol {
				// this looseEnd matches our start node
				return true
			}

			for _, end := range endPoints {
				if end.Row == start.Row && isBetween(int(nc.Col), int(start.Col), int(end.Col)) {
					return true
				}
			}
			return false
		}

		for _, end := range endPoints {
			if end.Col == start.Col && isBetween(int(nc.Row), int(start.Row), int(end.Row)) {
				return true
			}
		}

		return false
	}

	for i := 0; i < len(ends); i++ {
		if shouldRemove(ends[i]) {
			ends = append(ends[:i], ends[i+1:]...)
			i--
		}
	}

	return ends
}
