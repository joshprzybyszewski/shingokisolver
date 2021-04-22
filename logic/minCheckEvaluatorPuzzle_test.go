package logic_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshprzybyszewski/shingokisolver/logic"
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

func TestMinArmCheckBlackNode(t *testing.T) {
	puzz := puzzle.BuildTestPuzzleWithNoRules(
		t,
		5,
		[]model.NodeLocation{{
			Row:     3,
			Col:     4,
			IsWhite: false,
			Value:   3,
		}},
		model.NewCoord(4, 4),
		model.HeadDown,
		model.HeadRight,
		model.HeadUp,
		model.HeadUp,
		model.HeadLeft,
	)

	target, state := puzz.GetFirstTarget()
	assert.Equal(t, model.Incomplete, state)

	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		0,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		1,
		puzz.Alpha(),
	))

	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadUp,
		0,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadUp,
		1,
		puzz.Alpha(),
	))
}

func TestMinArmCheckEmptySpaces(t *testing.T) {
	puzz := puzzle.NewPuzzle(7, []model.NodeLocation{{
		Row:     3,
		Col:     3,
		IsWhite: false,
		Value:   5,
	}})
	puzz, state := puzzle.AvoidEdge(puzz, model.NewEdgePair(
		model.NewCoord(0, 3),
		model.HeadDown,
	))
	assert.Equal(t, model.Incomplete, state)
	puzz, state = puzzle.AvoidEdge(puzz, model.NewEdgePair(
		model.NewCoord(3, 0),
		model.HeadRight,
	))
	assert.Equal(t, model.Incomplete, state)
	puzz, state = puzzle.AvoidEdge(puzz, model.NewEdgePair(
		model.NewCoord(3, 5),
		model.HeadRight,
	))
	assert.Equal(t, model.Incomplete, state)
	t.Logf("TestMinArmCheckEmptySpaces puzzle: \n%s\n", puzz)

	target, state := puzz.GetFirstTarget()
	assert.Equal(t, model.Incomplete, state)

	assert.Equal(t, model.EdgeExists, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		0,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeExists, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		1,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeExists, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		2,
		puzz.Alpha(),
	))

	assert.Equal(t, model.EdgeAvoided, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadUp,
		0,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadUp,
		1,
		puzz.Alpha(),
	))

	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadLeft,
		0,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadLeft,
		1,
		puzz.Alpha(),
	))

	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadRight,
		0,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadRight,
		1,
		puzz.Alpha(),
	))
}

func TestMinArmCheckOverExtended(t *testing.T) {
	puzz := puzzle.NewPuzzle(7, []model.NodeLocation{{
		Row:     3,
		Col:     0,
		IsWhite: false,
		Value:   3,
	}})
	puzz, state := puzzle.AddEdge(puzz, model.NewEdgePair(
		model.NewCoord(4, 0),
		model.HeadDown,
	))
	assert.Equal(t, model.Incomplete, state)
	puzz, state = puzzle.AddEdge(puzz, model.NewEdgePair(
		model.NewCoord(5, 0),
		model.HeadDown,
	))
	assert.Equal(t, model.Incomplete, state)
	puzz, state = puzzle.AddEdge(puzz, model.NewEdgePair(
		model.NewCoord(3, 0),
		model.HeadRight,
	))
	assert.Equal(t, model.Incomplete, state)
	puzz, state = puzzle.AvoidEdge(puzz, model.NewEdgePair(
		model.NewCoord(3, 1),
		model.HeadRight,
	))
	assert.Equal(t, model.Incomplete, state)
	t.Logf("puzz: \n%s\n", puzz)

	target, state := puzz.GetFirstTarget()
	assert.Equal(t, model.Incomplete, state)

	assert.Equal(t, model.EdgeAvoided, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		0,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		1,
		puzz.Alpha(),
	))

	assert.Equal(t, model.EdgeExists, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadUp,
		0,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeExists, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadUp,
		1,
		puzz.Alpha(),
	))

	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadRight,
		0,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadRight,
		1,
		puzz.Alpha(),
	))
}

func TestMinArmCheckBlackNodeOutThere(t *testing.T) {
	puzz := puzzle.BuildTestPuzzleWithNoRules(
		t,
		7,
		[]model.NodeLocation{{
			Row:     0,
			Col:     0,
			IsWhite: false,
			Value:   6,
		}},
		model.NewCoord(1, 0),
		model.HeadUp,
		model.HeadRight,
		model.HeadRight,
		model.HeadRight,
	)

	target, state := puzz.GetFirstTarget()
	assert.Equal(t, model.Incomplete, state)

	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		0,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		1,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		2,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		3,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		4,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadDown,
		5,
		puzz.Alpha(),
	))

	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadRight,
		0,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadRight,
		1,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadRight,
		2,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadRight,
		3,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadRight,
		4,
		puzz.Alpha(),
	))
	assert.Equal(t, model.EdgeUnknown, logic.GetMinArmCheckEvaluation(
		target.Node,
		model.HeadRight,
		5,
		puzz.Alpha(),
	))

}
