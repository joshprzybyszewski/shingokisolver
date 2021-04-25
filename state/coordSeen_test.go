package state

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/stretchr/testify/assert"
)

func TestCoordSeen(t *testing.T) {
	cs := NewCoordSeen(5)

	cs.Mark(model.NewCoord(3, 3))
	assert.True(t, cs.IsCoordSeen(model.NewCoord(3, 3)))
	assert.False(t, cs.IsCoordSeen(model.NewCoord(3, 2)))
	assert.False(t, cs.IsCoordSeen(model.NewCoord(3, 4)))
	assert.False(t, cs.IsCoordSeen(model.NewCoord(4, 3)))
	assert.False(t, cs.IsCoordSeen(model.NewCoord(2, 3)))

	cs = NewCoordSeen(MaxEdges)

	cs.Mark(model.NewCoord(3, 3))
	assert.True(t, cs.IsCoordSeen(model.NewCoord(3, 3)))
	assert.False(t, cs.IsCoordSeen(model.NewCoord(3, 2)))
	assert.False(t, cs.IsCoordSeen(model.NewCoord(3, 4)))
	assert.False(t, cs.IsCoordSeen(model.NewCoord(4, 3)))
	assert.False(t, cs.IsCoordSeen(model.NewCoord(2, 3)))

	cs.Mark(model.NewCoord(MaxEdges, MaxEdges))
	assert.True(t, cs.IsCoordSeen(model.NewCoord(MaxEdges, MaxEdges)))
	assert.False(t, cs.IsCoordSeen(model.NewCoord(MaxEdges-1, MaxEdges)))
	assert.False(t, cs.IsCoordSeen(model.NewCoord(MaxEdges, MaxEdges-1)))
}
