package reader

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromWebsiteTask(t *testing.T) {
	testCases := []struct {
		numEdges int
		puzzID   string
		task     string
		expPD    PuzzleDef
	}{{
		numEdges: 5,
		puzzID:   `1,848,688`,
		task:     `bB3B3bW4aB2hW3aB4cB4hB2bB3`,
		expPD: PuzzleDef{
			Description: `PuzzleID: 1,848,688`,
			NumEdges:    5,
			Nodes:       []model.NodeLocation{{
				// TODO figure out what the nodes are
			}},
		},
	}, {
		numEdges: 5,
		puzzID:   `5,258,416`,
		task:     `B5eW2B2bB4cB2fB3eB3h`,
		expPD: PuzzleDef{
			Description: `PuzzleID: 5,258,416`,
			NumEdges:    5,
			Nodes: []model.NodeLocation{{
				Row:     0,
				Col:     0,
				Value:   5,
				IsWhite: false,
			}, {
				Row:     1,
				Col:     0,
				Value:   2,
				IsWhite: true,
			}, {
				Row:     1,
				Col:     1,
				Value:   2,
				IsWhite: false,
			}, {
				Row:     1,
				Col:     4,
				Value:   4,
				IsWhite: false,
			}, {
				Row:     2,
				Col:     2,
				Value:   2,
				IsWhite: false,
			}, {
				Row:     3,
				Col:     3,
				Value:   3,
				IsWhite: false,
			}, {
				Row:     4,
				Col:     3,
				Value:   3,
				IsWhite: false,
			}},
		},
	}, {
		numEdges: 5,
		puzzID:   `20,277`,
		task:     `TODO`, // TODO figure out what this is
		expPD: PuzzleDef{
			Description: `PuzzleID: 20,277`,
			NumEdges:    5,
			Nodes: []model.NodeLocation{{
				Row:     3,
				Col:     2,
				Value:   4,
				IsWhite: false,
			}, {
				Row:     3,
				Col:     5,
				Value:   5,
				IsWhite: true,
			}, {
				Row:     4,
				Col:     0,
				Value:   5,
				IsWhite: true,
			}, {
				Row:     5,
				Col:     1,
				Value:   2,
				IsWhite: false,
			}, {
				Row:     5,
				Col:     3,
				Value:   5,
				IsWhite: false,
			}},
		},
	}}

	for _, tc := range testCases {
		actPD, err := fromWebsiteTask(tc.numEdges, tc.puzzID, tc.task)
		require.NoError(t, err)
		assert.Equal(t, tc.expPD, actPD)
	}
}
