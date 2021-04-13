package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeTypeIsInvalidEdges(t *testing.T) {
	testCases := []struct {
		msg      string
		oe       OutgoingEdges
		expWhite bool
		expBlack bool
	}{{
		msg: `up`,
		oe: OutgoingEdges{
			above: 1,
		},
	}, {
		msg: `right`,
		oe: OutgoingEdges{
			right: 1,
		},
	}, {
		msg: `down`,
		oe: OutgoingEdges{
			below: 1,
		},
	}, {
		msg: `left`,
		oe: OutgoingEdges{
			left: 1,
		},
	}, {
		msg: `up and down`,
		oe: OutgoingEdges{
			above: 1,
			below: 1,
		},
		expBlack: true,
	}, {
		msg: `left and right`,
		oe: OutgoingEdges{
			left:  1,
			right: 1,
		},
		expBlack: true,
	}, {
		msg: `up and right`,
		oe: OutgoingEdges{
			above: 1,
			right: 1,
		},
		expWhite: true,
	}, {
		msg: `right and down`,
		oe: OutgoingEdges{
			below: 1,
			right: 1,
		},
		expWhite: true,
	}, {
		msg: `down and left`,
		oe: OutgoingEdges{
			below: 1,
			left:  1,
		},
		expWhite: true,
	}, {
		msg: `left and up`,
		oe: OutgoingEdges{
			above: 1,
			left:  1,
		},
		expWhite: true,
	}, {
		msg: `up, right, and down`,
		oe: OutgoingEdges{
			above: 1,
			below: 1,
			right: 1,
		},
		expWhite: true,
		expBlack: true,
	}, {
		msg: `up, right, and left`,
		oe: OutgoingEdges{
			above: 1,
			left:  1,
			right: 1,
		},
		expWhite: true,
		expBlack: true,
	}, {
		msg: `up, down, and left`,
		oe: OutgoingEdges{
			above: 1,
			below: 1,
			left:  1,
		},
		expWhite: true,
		expBlack: true,
	}}

	for _, tc := range testCases {
		assert.True(t, noNode.isInvalidEdges(tc.oe))
		assert.Equal(t, tc.expWhite, WhiteNode.isInvalidEdges(tc.oe))
		assert.Equal(t, tc.expBlack, BlackNode.isInvalidEdges(tc.oe))
	}
}

func TestNodeTypeIsInvalidMotions(t *testing.T) {
	testCases := []struct {
		c1       Cardinal
		c2       Cardinal
		expWhite bool
		expBlack bool
	}{{
		c1:       HeadUp,
		c2:       HeadRight,
		expWhite: true,
		expBlack: false,
	}, {
		c1:       HeadUp,
		c2:       HeadDown,
		expWhite: false,
		expBlack: true,
	}, {
		c1:       HeadUp,
		c2:       HeadLeft,
		expWhite: true,
		expBlack: false,
	}, {
		c1:       HeadRight,
		c2:       HeadDown,
		expWhite: true,
		expBlack: false,
	}, {
		c1:       HeadRight,
		c2:       HeadLeft,
		expWhite: false,
		expBlack: true,
	}, {
		c1:       HeadDown,
		c2:       HeadLeft,
		expWhite: true,
		expBlack: false,
	}}

	for _, tc := range testCases {
		assert.True(t, noNode.isInvalidMotions(tc.c1, tc.c2))
		assert.True(t, noNode.isInvalidMotions(tc.c2, tc.c1))
		assert.Equal(t, tc.expWhite, WhiteNode.isInvalidMotions(tc.c1, tc.c2))
		assert.Equal(t, tc.expWhite, WhiteNode.isInvalidMotions(tc.c2, tc.c1))
		assert.Equal(t, tc.expBlack, BlackNode.isInvalidMotions(tc.c1, tc.c2))
		assert.Equal(t, tc.expBlack, BlackNode.isInvalidMotions(tc.c2, tc.c1))
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
	w5 := NewNode(NodeCoord{}, true, 5)

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

	b7 := NewNode(NodeCoord{}, false, 7)
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
