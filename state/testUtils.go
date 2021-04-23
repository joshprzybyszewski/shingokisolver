package state

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/stretchr/testify/require"
)

func BuildGetEdger(
	t *testing.T,
	numEdges int,
	startCoord model.NodeCoord,
	steps ...model.Cardinal,
) model.GetEdger {
	if numEdges > MaxEdges {
		t.Error(`bad input numEdges`)
	}

	edges := New(numEdges)

	c := startCoord
	var ms model.State
	for _, dir := range steps {
		ms = edges.setEdge(model.NewEdgePair(c, dir))
		switch ms {
		case model.Unexpected, model.Violation, model.Duplicate:
			require.Fail(t, "failure building puzzle", "unexpected state (%s) after adding edge: %+v, %+v\n", ms, dir, c)
		}
		c = c.Translate(dir)
	}
	t.Logf("BuildGetEdger produced: \n%s\n", edges)
	return &edges
}

type BuildInput struct {
	Existing []model.EdgePair
	Avoided  []model.EdgePair
}

func BuildGetEdgerWithInput(
	t *testing.T,
	numEdges int,
	input BuildInput,
) model.GetEdger {
	if numEdges > MaxEdges {
		t.Error(`bad input numEdges`)
	}

	edges := New(numEdges)

	var ms model.State
	for _, e := range input.Existing {
		ms = edges.setEdge(e)
		switch ms {
		case model.Unexpected, model.Violation, model.Duplicate:
			require.Fail(t, "failure building puzzle", "unexpected state (%s) after adding edge: %s\n", ms, e)
		}
	}
	for _, e := range input.Avoided {
		ms = edges.avoidEdge(e)
		switch ms {
		case model.Unexpected, model.Violation, model.Duplicate:
			require.Fail(t, "failure building puzzle", "unexpected state (%s) after avoiding edge: %s\n", ms, e)
		}
	}

	t.Logf("BuildGetEdger produced: \n%s\n", edges)
	return &edges
}
