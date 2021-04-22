package logic

import (
	"fmt"
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
	"github.com/stretchr/testify/assert"
)

func getTestStates(t *testing.T) (
	a2ThroughA8, a3ThroughA8, a2ThroughA5, a3ThroughA6 model.GetEdger,
) {
	a2ThroughA8 = state.BuildGetEdger(
		t,
		7,
		model.NewCoord(1, 0),
		model.HeadDown,
		model.HeadDown,
		model.HeadDown,
		model.HeadDown,
		model.HeadDown,
		model.HeadDown,
	)
	a3ThroughA8 = state.BuildGetEdger(
		t,
		7,
		model.NewCoord(2, 0),
		model.HeadDown,
		model.HeadDown,
		model.HeadDown,
		model.HeadDown,
		model.HeadDown,
	)
	a2ThroughA5 = state.BuildGetEdger(
		t,
		7,
		model.NewCoord(1, 0),
		model.HeadDown,
		model.HeadDown,
		model.HeadDown,
	)
	a3ThroughA6 = state.BuildGetEdger(
		t,
		7,
		model.NewCoord(2, 0),
		model.HeadDown,
		model.HeadDown,
		model.HeadDown,
	)

	return a2ThroughA8, a3ThroughA8, a2ThroughA5, a3ThroughA6
}

func TestGetNumInDirection(t *testing.T) {
	a2ThroughA8, a3ThroughA8, a2ThroughA5, a3ThroughA6 := getTestStates(t)

	testCases := []struct {
		ge            model.GetEdger
		nc            model.NodeCoord
		nValue        int
		dir           model.Cardinal
		expPossible   int
		expStartIndex int
		expNumInChain int
	}{{
		ge:            a2ThroughA8,
		nc:            model.NewCoord(0, 0),
		nValue:        6,
		dir:           model.HeadDown,
		expPossible:   7,
		expStartIndex: 1,
		expNumInChain: 6,
	}, {
		ge:            a3ThroughA8,
		nc:            model.NewCoord(0, 0),
		nValue:        6,
		dir:           model.HeadDown,
		expPossible:   7,
		expStartIndex: 2,
		expNumInChain: 5,
	}, {
		ge:            a2ThroughA5,
		nc:            model.NewCoord(0, 0),
		nValue:        3,
		dir:           model.HeadDown,
		expPossible:   4,
		expStartIndex: 1,
		expNumInChain: 3,
	}, {
		ge:            a2ThroughA5,
		nc:            model.NewCoord(3, 0),
		nValue:        5,
		dir:           model.HeadDown,
		expPossible:   4,
		expStartIndex: -1,
		expNumInChain: 0,
	}, {
		ge:            a3ThroughA6,
		nc:            model.NewCoord(0, 0),
		nValue:        3,
		dir:           model.HeadDown,
		expPossible:   2,
		expStartIndex: -1,
		expNumInChain: 0,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Heading %s from %s (value: %d)", tc.dir, tc.nc, tc.nValue), func(t *testing.T) {
			ge := tc.ge
			t.Logf("Running on puzzle:\n%s", ge)

			numPossibleEdges, existingChainStartIndex, numInLastExistingChain := getNumInDirection(
				ge, tc.nValue, tc.nc, tc.dir,
			)

			assert.Equal(t, tc.expPossible, numPossibleEdges)
			assert.Equal(t, tc.expStartIndex, existingChainStartIndex)
			assert.Equal(t, tc.expNumInChain, numInLastExistingChain)
		})
	}
}

func TestMinArmCheckEvaluator(t *testing.T) {
	_, _, a2ThroughA5, _ := getTestStates(t)

	custom1 := state.BuildGetEdgerWithInput(t, 7, state.BuildInput{
		Existing: []model.EdgePair{
			model.NewEdgePair(model.NewCoord(2, 0), model.HeadDown),
			model.NewEdgePair(model.NewCoord(3, 0), model.HeadDown),
		},
		Avoided: []model.EdgePair{
			model.NewEdgePair(model.NewCoord(2, 1), model.HeadRight),
		},
	})

	custom2 := state.BuildGetEdgerWithInput(t, 5, state.BuildInput{
		Existing: []model.EdgePair{
			model.NewEdgePair(model.NewCoord(3, 5), model.HeadDown),
			model.NewEdgePair(model.NewCoord(4, 5), model.HeadDown),
			model.NewEdgePair(model.NewCoord(4, 4), model.HeadDown),
			model.NewEdgePair(model.NewCoord(5, 4), model.HeadRight),
		},
		Avoided: []model.EdgePair{
			model.NewEdgePair(model.NewCoord(3, 4), model.HeadLeft),
		},
	})

	custom3 := state.BuildGetEdgerWithInput(t, 7, state.BuildInput{
		Existing: []model.EdgePair{
			model.NewEdgePair(model.NewCoord(2, 0), model.HeadDown),
			model.NewEdgePair(model.NewCoord(3, 0), model.HeadDown),
			model.NewEdgePair(model.NewCoord(5, 0), model.HeadDown),
			model.NewEdgePair(model.NewCoord(6, 0), model.HeadDown),
		},
	})

	custom4 := state.BuildGetEdgerWithInput(t, 5, state.BuildInput{
		Existing: []model.EdgePair{
			model.NewEdgePair(model.NewCoord(1, 1), model.HeadRight),
			model.NewEdgePair(model.NewCoord(1, 2), model.HeadRight),
		},
		Avoided: []model.EdgePair{
			model.NewEdgePair(model.NewCoord(1, 0), model.HeadRight),
		},
	})

	testCases := []struct {
		ge           model.GetEdger
		node         model.Node
		myDir        model.Cardinal
		myIndex      int
		expEdgeState model.EdgeState
	}{{
		ge:           custom4,
		node:         model.NewNode(model.NewCoord(1, 1), false, 5),
		myDir:        model.HeadRight,
		myIndex:      2,
		expEdgeState: model.EdgeUnknown,
	}, {
		ge:           custom3,
		node:         model.NewNode(model.NewCoord(3, 0), true, 3),
		myDir:        model.HeadUp,
		myIndex:      0,
		expEdgeState: model.EdgeExists,
	}, {
		ge:           custom3,
		node:         model.NewNode(model.NewCoord(3, 0), true, 3),
		myDir:        model.HeadUp,
		myIndex:      1,
		expEdgeState: model.EdgeExists,
	}, {
		ge:           custom3,
		node:         model.NewNode(model.NewCoord(3, 0), true, 3),
		myDir:        model.HeadDown,
		myIndex:      0,
		expEdgeState: model.EdgeUnknown,
	}, {
		ge:           custom3,
		node:         model.NewNode(model.NewCoord(3, 0), true, 3),
		myDir:        model.HeadDown,
		myIndex:      1,
		expEdgeState: model.EdgeUnknown,
	}, {
		ge:           custom2,
		node:         model.NewNode(model.NewCoord(3, 4), false, 3),
		myDir:        model.HeadUp,
		myIndex:      0,
		expEdgeState: model.EdgeUnknown,
	}, {
		ge:           custom2,
		node:         model.NewNode(model.NewCoord(3, 4), false, 3),
		myDir:        model.HeadUp,
		myIndex:      1,
		expEdgeState: model.EdgeUnknown,
	}, {
		ge:           custom2,
		node:         model.NewNode(model.NewCoord(3, 4), false, 3),
		myDir:        model.HeadDown,
		myIndex:      0,
		expEdgeState: model.EdgeUnknown,
	}, {
		ge:           custom2,
		node:         model.NewNode(model.NewCoord(3, 4), false, 3),
		myDir:        model.HeadDown,
		myIndex:      1,
		expEdgeState: model.EdgeUnknown,
	}, {
		ge:           custom1,
		node:         model.NewNode(model.NewCoord(2, 0), false, 5),
		myDir:        model.HeadUp,
		myIndex:      0,
		expEdgeState: model.EdgeAvoided,
	}, {
		ge:           custom1,
		node:         model.NewNode(model.NewCoord(2, 0), false, 5),
		myDir:        model.HeadUp,
		myIndex:      1,
		expEdgeState: model.EdgeUnknown,
	}, {
		ge:           custom1,
		node:         model.NewNode(model.NewCoord(2, 0), false, 5),
		myDir:        model.HeadDown,
		myIndex:      0,
		expEdgeState: model.EdgeExists,
	}, {
		ge:           custom1,
		node:         model.NewNode(model.NewCoord(2, 0), false, 5),
		myDir:        model.HeadDown,
		myIndex:      1,
		expEdgeState: model.EdgeExists,
	}, {
		ge:           custom1,
		node:         model.NewNode(model.NewCoord(2, 0), false, 5),
		myDir:        model.HeadDown,
		myIndex:      2,
		expEdgeState: model.EdgeExists,
	}, {
		ge:           custom1,
		node:         model.NewNode(model.NewCoord(2, 0), false, 5),
		myDir:        model.HeadDown,
		myIndex:      3,
		expEdgeState: model.EdgeExists,
	}, {
		ge:           a2ThroughA5,
		node:         model.NewNode(model.NewCoord(3, 0), true, 5),
		myDir:        model.HeadDown,
		myIndex:      0,
		expEdgeState: model.EdgeExists,
	}, {
		ge:           a2ThroughA5,
		node:         model.NewNode(model.NewCoord(3, 0), true, 5),
		myDir:        model.HeadDown,
		myIndex:      1,
		expEdgeState: model.EdgeExists,
	}, {
		ge:           a2ThroughA5,
		node:         model.NewNode(model.NewCoord(3, 0), true, 5),
		myDir:        model.HeadDown,
		myIndex:      2,
		expEdgeState: model.EdgeUnknown,
	}, {
		ge:           a2ThroughA5,
		node:         model.NewNode(model.NewCoord(3, 0), true, 7),
		myDir:        model.HeadUp,
		myIndex:      2,
		expEdgeState: model.EdgeExists,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Node: %s, dir: %s, index: %d", tc.node, tc.myDir, tc.myIndex), func(t *testing.T) {
			ge := tc.ge
			t.Logf("Running on state:\n%s", ge)

			e := newMinArmCheckEvaluator(tc.node, tc.myDir, tc.myIndex)

			assert.Equal(t, tc.expEdgeState, e.evaluate(ge))
		})
	}

}
