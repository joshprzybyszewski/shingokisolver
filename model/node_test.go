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

	assert.Equal(t, int8(4), n.Value())
	assert.Equal(t, int8(4), cpy1.Value())
	assert.Equal(t, BlackNode, n.Type())
	assert.Equal(t, BlackNode, cpy1.Type())

	n.val = 5
	n.nType = WhiteNode

	assert.Equal(t, int8(5), n.Value())
	assert.Equal(t, int8(4), cpy1.Value())
	assert.Equal(t, WhiteNode, n.Type())
	assert.Equal(t, BlackNode, cpy1.Type())

}

func TestIsInvalidMotions(t *testing.T) {
	w5 := NewNode(true, 5)

	assert.True(t, w5.IsInvalidMotions(HeadUp, HeadRight))
	assert.True(t, w5.IsInvalidMotions(HeadUp, HeadLeft))
	assert.True(t, w5.IsInvalidMotions(HeadDown, HeadRight))
	assert.True(t, w5.IsInvalidMotions(HeadDown, HeadLeft))
	assert.True(t, w5.IsInvalidMotions(HeadRight, HeadUp))
	assert.True(t, w5.IsInvalidMotions(HeadRight, HeadDown))
	assert.True(t, w5.IsInvalidMotions(HeadLeft, HeadUp))
	assert.True(t, w5.IsInvalidMotions(HeadLeft, HeadDown))

	assert.False(t, w5.IsInvalidMotions(HeadUp, HeadDown))
	assert.False(t, w5.IsInvalidMotions(HeadDown, HeadUp))
	assert.False(t, w5.IsInvalidMotions(HeadLeft, HeadRight))
	assert.False(t, w5.IsInvalidMotions(HeadRight, HeadLeft))

	b7 := NewNode(false, 7)
	assert.False(t, b7.IsInvalidMotions(HeadUp, HeadRight))
	assert.False(t, b7.IsInvalidMotions(HeadUp, HeadLeft))
	assert.False(t, b7.IsInvalidMotions(HeadDown, HeadRight))
	assert.False(t, b7.IsInvalidMotions(HeadDown, HeadLeft))
	assert.False(t, b7.IsInvalidMotions(HeadRight, HeadUp))
	assert.False(t, b7.IsInvalidMotions(HeadRight, HeadDown))
	assert.False(t, b7.IsInvalidMotions(HeadLeft, HeadUp))
	assert.False(t, b7.IsInvalidMotions(HeadLeft, HeadDown))

	assert.True(t, b7.IsInvalidMotions(HeadUp, HeadDown))
	assert.True(t, b7.IsInvalidMotions(HeadDown, HeadUp))
	assert.True(t, b7.IsInvalidMotions(HeadLeft, HeadRight))
	assert.True(t, b7.IsInvalidMotions(HeadRight, HeadLeft))
}
