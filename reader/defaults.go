package reader

import "github.com/joshprzybyszewski/shingokisolver/puzzlegrid"

func DefaultPuzzles() []PuzzleDef {
	pds := []PuzzleDef{{
		NumEdges: 2,
		Nodes: []puzzlegrid.NodeLocation{{
			Row:     1,
			Col:     1,
			IsWhite: false,
			Value:   2,
		}},
	}, {
		NumEdges: 5,
		Nodes: []puzzlegrid.NodeLocation{{
			Row:     3,
			Col:     2,
			IsWhite: false,
			Value:   4,
		}, {
			Row:     3,
			Col:     5,
			IsWhite: true,
			Value:   5,
		}, {
			Row:     4,
			Col:     0,
			IsWhite: true,
			Value:   5,
		}, {
			Row:     5,
			Col:     1,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     5,
			Col:     3,
			IsWhite: false,
			Value:   5,
		}},
	}}

	puzzles := []string{
		`.......b11
w2..b4w3...
..w3.b2...
...w4....
.b3...b4..
........
...w4.b5..
b10.w6.....`, // easy 7x7
		`...b2.b3..
w4...w2w2..
...b2..b2b2
.......b5
.w4....b2.
..w3.....
w2b2.b2.b2..
........`, // normal 7x7
		`.b2..b2.....b4
b2..b2......b2
.b2.b2b2b3....b2
b2.w2...b4.w4..
.......w2..b4
b2.b3b2.....b2.
.........b3.
....b4.b2..w2.
...w3..w2.b3..
.......b3b2..
.b3...b4.b2..b3`, // easy 10x10
	}
	for _, s := range puzzles {
		pd, err := FromString(s)
		if err == nil {
			pds = append(pds, pd)
		}
	}
	return pds
}
