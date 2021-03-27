package puzzle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEdgesFromNodeGetNumCardinals(t *testing.T) {
	efn := edgesFromNode{}
	assert.Zero(t, efn.getNumOutgoingDirections())
	efn = edgesFromNode{
		above: 1,
		below: 1,
		left:  1,
		right: 1,
	}
	assert.Equal(t, int8(4), efn.getNumOutgoingDirections())
}

func TestNodeTypeIsInvalidEdges(t *testing.T) {
	testCases := []struct {
		msg                string
		efn                edgesFromNode
		expNoNodeOutput    bool
		expWhiteNodeOutput bool
		expBlackNodeOutput bool
	}{{
		msg: `up`,
		efn: edgesFromNode{
			above: 1,
		},
	}, {
		msg: `right`,
		efn: edgesFromNode{
			right: 1,
		},
	}, {
		msg: `down`,
		efn: edgesFromNode{
			below: 1,
		},
	}, {
		msg: `left`,
		efn: edgesFromNode{
			left: 1,
		},
	}, {
		msg: `up and down`,
		efn: edgesFromNode{
			above: 1,
			below: 1,
		},
		expBlackNodeOutput: true,
	}, {
		msg: `left and right`,
		efn: edgesFromNode{
			left:  1,
			right: 1,
		},
		expBlackNodeOutput: true,
	}, {
		msg: `up and right`,
		efn: edgesFromNode{
			above: 1,
			right: 1,
		},
		expWhiteNodeOutput: true,
	}, {
		msg: `right and down`,
		efn: edgesFromNode{
			below: 1,
			right: 1,
		},
		expWhiteNodeOutput: true,
	}, {
		msg: `down and left`,
		efn: edgesFromNode{
			below: 1,
			left:  1,
		},
		expWhiteNodeOutput: true,
	}, {
		msg: `left and up`,
		efn: edgesFromNode{
			above: 1,
			left:  1,
		},
		expWhiteNodeOutput: true,
	}, {
		msg: `up, right, and down`,
		efn: edgesFromNode{
			above: 1,
			below: 1,
			right: 1,
		},
		expWhiteNodeOutput: true,
		expBlackNodeOutput: true,
	}, {
		msg: `up, right, and left`,
		efn: edgesFromNode{
			above: 1,
			left:  1,
			right: 1,
		},
		expWhiteNodeOutput: true,
		expBlackNodeOutput: true,
	}, {
		msg: `up, down, and left`,
		efn: edgesFromNode{
			above: 1,
			below: 1,
			left:  1,
		},
		expWhiteNodeOutput: true,
		expBlackNodeOutput: true,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expNoNodeOutput, noNode.isInvalidEdges(tc.efn))
		assert.Equal(t, tc.expWhiteNodeOutput, whiteNode.isInvalidEdges(tc.efn))
		assert.Equal(t, tc.expBlackNodeOutput, blackNode.isInvalidEdges(tc.efn))
	}
}

func TestNodeCopy(t *testing.T) {
	n := node{
		nType: blackNode,
		val:   4,
	}
	cpy1 := n.copy()
	assert.Equal(t, node{
		nType: blackNode,
		val:   4,
	}, cpy1)
	n.val = 5

	assert.Equal(t, node{
		nType: blackNode,
		val:   4,
	}, cpy1)
}
