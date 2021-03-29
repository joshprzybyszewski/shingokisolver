package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeTypeIsInvalidEdges(t *testing.T) {
	testCases := []struct {
		msg                string
		efn                OutgoingEdges
		expNoNodeOutput    bool
		expWhiteNodeOutput bool
		expBlackNodeOutput bool
	}{{
		msg: `up`,
		efn: OutgoingEdges{
			above: 1,
		},
	}, {
		msg: `right`,
		efn: OutgoingEdges{
			right: 1,
		},
	}, {
		msg: `down`,
		efn: OutgoingEdges{
			below: 1,
		},
	}, {
		msg: `left`,
		efn: OutgoingEdges{
			left: 1,
		},
	}, {
		msg: `up and down`,
		efn: OutgoingEdges{
			above: 1,
			below: 1,
		},
		expBlackNodeOutput: true,
	}, {
		msg: `left and right`,
		efn: OutgoingEdges{
			left:  1,
			right: 1,
		},
		expBlackNodeOutput: true,
	}, {
		msg: `up and right`,
		efn: OutgoingEdges{
			above: 1,
			right: 1,
		},
		expWhiteNodeOutput: true,
	}, {
		msg: `right and down`,
		efn: OutgoingEdges{
			below: 1,
			right: 1,
		},
		expWhiteNodeOutput: true,
	}, {
		msg: `down and left`,
		efn: OutgoingEdges{
			below: 1,
			left:  1,
		},
		expWhiteNodeOutput: true,
	}, {
		msg: `left and up`,
		efn: OutgoingEdges{
			above: 1,
			left:  1,
		},
		expWhiteNodeOutput: true,
	}, {
		msg: `up, right, and down`,
		efn: OutgoingEdges{
			above: 1,
			below: 1,
			right: 1,
		},
		expWhiteNodeOutput: true,
		expBlackNodeOutput: true,
	}, {
		msg: `up, right, and left`,
		efn: OutgoingEdges{
			above: 1,
			left:  1,
			right: 1,
		},
		expWhiteNodeOutput: true,
		expBlackNodeOutput: true,
	}, {
		msg: `up, down, and left`,
		efn: OutgoingEdges{
			above: 1,
			below: 1,
			left:  1,
		},
		expWhiteNodeOutput: true,
		expBlackNodeOutput: true,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expNoNodeOutput, noNode.isInvalidEdges(tc.efn))
		assert.Equal(t, tc.expWhiteNodeOutput, WhiteNode.isInvalidEdges(tc.efn))
		assert.Equal(t, tc.expBlackNodeOutput, BlackNode.isInvalidEdges(tc.efn))
	}
}

func TestNodeCopy(t *testing.T) {
	n := Node{
		nType: BlackNode,
		val:   4,
	}

	cpy1 := n.Copy()
	assert.Equal(t, Node{
		nType: BlackNode,
		val:   4,
	}, cpy1)

	assert.Equal(t, uint8(4), n.Value())
	assert.Equal(t, uint8(4), cpy1.Value())
	assert.Equal(t, BlackNode, n.Type())
	assert.Equal(t, BlackNode, cpy1.Type())

	n.val = 5
	n.nType = WhiteNode

	assert.Equal(t, uint8(4), n.Value())
	assert.Equal(t, uint8(4), cpy1.Value())
	assert.Equal(t, WhiteNode, n.Type())
	assert.Equal(t, BlackNode, cpy1.Type())

}
