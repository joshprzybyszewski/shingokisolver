package reader

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/puzzlegrid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromString(t *testing.T) {
	input := `...b2.b3..
w4...w2w2..
...b2..b2b2
.......b5
.w4....b2.
..w3.....
w2b2.b2.b2..
........`

	actPD, err := FromString(input)
	require.NoError(t, err)

	expPD := PuzzleDef{
		NumEdges: 7,
		Nodes: []puzzlegrid.NodeLocation{{
			Row:     0,
			Col:     3,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     0,
			Col:     5,
			IsWhite: false,
			Value:   3,
		}, {
			Row:     1,
			Col:     0,
			IsWhite: true,
			Value:   4,
		}, {
			Row:     1,
			Col:     4,
			IsWhite: true,
			Value:   2,
		}, {
			Row:     1,
			Col:     5,
			IsWhite: true,
			Value:   2,
		}, {
			Row:     2,
			Col:     3,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     2,
			Col:     6,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     2,
			Col:     7,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     3,
			Col:     7,
			IsWhite: false,
			Value:   5,
		}, {
			Row:     4,
			Col:     1,
			IsWhite: true,
			Value:   4,
		}, {
			Row:     4,
			Col:     6,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     5,
			Col:     2,
			IsWhite: true,
			Value:   3,
		}, {
			Row:     6,
			Col:     0,
			IsWhite: true,
			Value:   2,
		}, {
			Row:     6,
			Col:     1,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     6,
			Col:     3,
			IsWhite: false,
			Value:   2,
		}, {
			Row:     6,
			Col:     5,
			IsWhite: false,
			Value:   2,
		}},
	}

	assert.Equal(t, expPD, actPD)
}
