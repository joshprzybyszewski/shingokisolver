package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEdgesFromNodeGetNumCardinals(t *testing.T) {
	oe := OutgoingEdges{}
	assert.Zero(t, oe.GetNumOutgoingDirections())
	oe = OutgoingEdges{
		above: 1,
		below: 1,
		left:  1,
		right: 1,
	}
	assert.Equal(t, int8(4), oe.GetNumOutgoingDirections())
}
