package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGridGetSetCopy(t *testing.T) {
	for i := 1; i < 3; i++ {
		g := NewGrid(i)

		nc00 := NewCoordFromInts(0, 0)
		oe00 := OutgoingEdges{
			below: 1,
			right: 1,
		}
		g.Set(nc00, oe00)
		assert.Equal(t, oe00, g.Get(nc00))

		gCpy := g.Copy()
		assert.Equal(t, g, gCpy)
	}

	for i := 3; i < 7; i++ {
		g := NewGrid(i)

		nc00 := NewCoordFromInts(0, 0)
		oe00 := OutgoingEdges{
			below: 1,
			right: 1,
		}
		g.Set(nc00, oe00)
		assert.Equal(t, oe00, g.Get(nc00))

		nc11 := NewCoordFromInts(1, 1)
		oe11 := OutgoingEdges{
			below: 2,
			right: 2,
		}
		g.Set(nc11, oe11)
		assert.Equal(t, oe11, g.Get(nc11))

		gCpy := g.Copy()
		assert.Equal(t, g, gCpy)
	}

	for i := 7; i < MAX_EDGES; i++ {
		g := NewGrid(i)

		nc00 := NewCoordFromInts(0, 0)
		oe00 := OutgoingEdges{
			below: 1,
			right: 1,
		}
		g.Set(nc00, oe00)
		assert.Equal(t, oe00, g.Get(nc00))

		nc11 := NewCoordFromInts(1, 1)
		oe11 := OutgoingEdges{
			below: 2,
			right: 2,
		}
		g.Set(nc11, oe11)
		assert.Equal(t, oe11, g.Get(nc11))

		nc77 := NewCoordFromInts(7, 7)
		oe77 := OutgoingEdges{
			below: 7,
			right: 7,
		}
		g.Set(nc77, oe77)
		assert.Equal(t, oe77, g.Get(nc77))

		gCpy := g.Copy()
		assert.Equal(t, g, gCpy)
	}
}
