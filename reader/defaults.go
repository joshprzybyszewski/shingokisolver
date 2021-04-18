package reader

import "github.com/joshprzybyszewski/shingokisolver/model"

func DefaultPuzzles() []PuzzleDef {
	pds := []PuzzleDef{{
		Description: `Manual 2x2 (contrived example)`,
		NumEdges:    2,
		Nodes: []model.NodeLocation{{
			Row:     1,
			Col:     1,
			IsWhite: false,
			Value:   2,
		}},
	}, {
		Description: `Manual Easy`,
		NumEdges:    5,
		Nodes: []model.NodeLocation{{
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

	puzzles := []struct {
		desc string
		puzz string
	}{{
		desc: `Manual Easy 7x7`,
		puzz: `.......b11
w2..b4w3...
..w3.b2...
...w4....
.b3...b4..
........
...w4.b5..
b10.w6.....`,
	}, {
		desc: `Manual normal 7x7`,
		puzz: `........
.b2.w3....
.b6....w2.
w4w2...w4.b2
..w2.....
..b3.....
........
..b4....b4`,
	}, {
		desc: `Manual second normal 7x7`,
		puzz: `.w3......
...b2..w2w3
....b3...
...b4.b5b5.
b6.b3w2....
..b2.....
....b3...
..b3..b5..`,
	}, {
		desc: `Manual hard 7x7`,
		puzz: `b8.....b2.
.w6b6..w2..
...b3.b5..
........
.......b2
......w2.
........
.......b3`,
	}, {
		desc: `Manual normal 10x10`,
		puzz: `.b2..b2.....b4
b2..b2......b2
.b2.b2b2b3....b2
b2.w2...b4.w4..
.......w2..b4
b2.b3b2.....b2.
.........b3.
....b4.b2..w2.
...w3..w2.b3..
.......b3b2..
.b3...b4.b2..b3`,
	}, {
		desc: `Manual easy 15x15`,
		puzz: `.......b3..w2b3b3..b6
..b4b9..w3.......w4.
....b2..b4........
w2....b4..........
............w4...
b3.b3w7..w2.w2....b7..
........b4.......
......w5...b3.....
.b2....b2....b4..w4.
.....b2b3..b3.....b7
.b2b3b2w3.w2.....w5.b6.
...b4...b2w2..b3....
....b4.......b4...
.w2..w2......b4...w6
.....w5.......w4w5.
b12..........w3....`,
	}}
	for _, s := range puzzles {
		pd, err := FromString(s.puzz)
		if err == nil {
			pd.Description = s.desc
			pds = append(pds, pd)
		}
	}
	return pds
}
