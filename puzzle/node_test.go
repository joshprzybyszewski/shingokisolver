package puzzle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEdgesFromNodeGetNumCardinals(t *testing.T) {
	efn := edgesFromNode{}
	assert.Zero(t, efn.getNumCardinals())
	efn = edgesFromNode{
		isabove: true,
		isbelow: true,
		isleft:  true,
		isright: true,
	}
	assert.Equal(t, int8(4), efn.getNumCardinals())
}

func TestNodeTypeIsInvalidEdges(t *testing.T) {
	testCases := []struct {
		msg                string
		efn                *edgesFromNode
		expNoNodeOutput    bool
		expWhiteNodeOutput bool
		expBlackNodeOutput bool
	}{{
		msg: `up`,
		efn: &edgesFromNode{
			isabove: true,
			isbelow: false,
			isleft:  false,
			isright: false,
		},
	}, {
		msg: `right`,
		efn: &edgesFromNode{
			isabove: false,
			isbelow: false,
			isleft:  false,
			isright: true,
		},
	}, {
		msg: `down`,
		efn: &edgesFromNode{
			isabove: false,
			isbelow: true,
			isleft:  false,
			isright: false,
		},
	}, {
		msg: `left`,
		efn: &edgesFromNode{
			isabove: false,
			isbelow: false,
			isleft:  true,
			isright: false,
		},
	}, {
		msg: `up and down`,
		efn: &edgesFromNode{
			isabove: true,
			isbelow: true,
			isleft:  false,
			isright: false,
		},
		expBlackNodeOutput: true,
	}, {
		msg: `left and right`,
		efn: &edgesFromNode{
			isabove: false,
			isbelow: false,
			isleft:  true,
			isright: true,
		},
		expBlackNodeOutput: true,
	}, {
		msg: `up and right`,
		efn: &edgesFromNode{
			isabove: true,
			isbelow: false,
			isleft:  false,
			isright: true,
		},
		expWhiteNodeOutput: true,
	}, {
		msg: `right and down`,
		efn: &edgesFromNode{
			isabove: false,
			isbelow: true,
			isleft:  false,
			isright: true,
		},
		expWhiteNodeOutput: true,
	}, {
		msg: `down and left`,
		efn: &edgesFromNode{
			isabove: false,
			isbelow: true,
			isleft:  true,
			isright: false,
		},
		expWhiteNodeOutput: true,
	}, {
		msg: `left and up`,
		efn: &edgesFromNode{
			isabove: true,
			isbelow: false,
			isleft:  true,
			isright: false,
		},
		expWhiteNodeOutput: true,
	}, {
		msg: `up, right, and down`,
		efn: &edgesFromNode{
			isabove: true,
			isbelow: true,
			isleft:  false,
			isright: true,
		},
		expWhiteNodeOutput: true,
		expBlackNodeOutput: true,
	}, {
		msg: `up, right, and left`,
		efn: &edgesFromNode{
			isabove: true,
			isbelow: false,
			isleft:  true,
			isright: true,
		},
		expWhiteNodeOutput: true,
		expBlackNodeOutput: true,
	}, {
		msg: `up, down, and left`,
		efn: &edgesFromNode{
			isabove: true,
			isbelow: true,
			isleft:  true,
			isright: false,
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
	var n *node
	assert.Nil(t, n.copy())

	n = &node{
		nType: blackNode,
		val:   4,
	}
	assert.Equal(t, &node{
		nType: blackNode,
		val:   4,
	}, n.copy())

	n.seen = true
	assert.Equal(t, &node{
		nType: blackNode,
		val:   4,
	}, n.copy())
}
