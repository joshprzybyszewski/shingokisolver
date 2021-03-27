package puzzle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEdges(t *testing.T) {
	e := newEdges()
	for i := 0; i < MAX_EDGES; i++ {
		assert.False(t, e.isEdge(i))
	}
	t.Logf("e: %b\n", e)

	v2 := 5
	e2 := e.addEdge(v2)
	t.Logf("e2: %b\n", e2)
	assert.True(t, e2.isEdge(v2))
	assert.False(t, e.isEdge(v2))

	v3 := 13
	e3 := e2.addEdge(v3)
	t.Logf("e3: %b\n", e3)
	assert.True(t, e3.isEdge(v3))
	assert.False(t, e2.isEdge(v3))
	assert.False(t, e.isEdge(v3))

	v4 := 0
	e4 := e3.addEdge(v4)
	t.Logf("e4: %b\n", e4)
	assert.True(t, e4.isEdge(v4))
	assert.False(t, e3.isEdge(v4))
	assert.False(t, e2.isEdge(v4))
	assert.False(t, e.isEdge(v4))

	v5 := -3
	e5 := e4.addEdge(v5)
	t.Logf("e5: %b\n", e5)
	assert.False(t, e5.isEdge(v5))
	assert.False(t, e4.isEdge(v5))
	assert.False(t, e3.isEdge(v5))
	assert.False(t, e2.isEdge(v5))
	assert.False(t, e.isEdge(v5))

	v6 := MAX_EDGES + 1
	e6 := e5.addEdge(v6)
	t.Logf("e6: %b\n", e6)
	assert.False(t, e6.isEdge(v6))
	assert.False(t, e5.isEdge(v6))
	assert.False(t, e4.isEdge(v6))
	assert.False(t, e3.isEdge(v6))
	assert.False(t, e2.isEdge(v6))
	assert.False(t, e.isEdge(v6))
}
