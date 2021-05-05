package puzzle

import (
	"github.com/joshprzybyszewski/shingokisolver/model"
)

// GetNodeState returns the state of the Node at the given Coord.
// Warning: this method could be expensive!
// If the node is satisfied, then we will walk the outgoing
// edges from it to see if there exists a loop that shouldn't exist.
func (p Puzzle) GetNodeState(
	nc model.NodeCoord,
) model.State {

	var nm *model.NodeMeta
	for _, m := range p.metas {
		if m.Coord() == nc {
			nm = m
			break
		}
	}
	if nm == nil {
		// why are you asking about this?
		return model.Unexpected
	}

	// TODO convert this to a check on CheckComplete of the NodeMeta
	ns := nm.Filter(&p.edges)
	if ns != model.Complete {
		return ns
	}

	// We know the node is complete at this point.
	// Let's pre-emptively check for a loop that shouldn't exist
	_, s := p.getStateOfLoop(nm.Coord())
	if s != model.Complete && s != model.Incomplete {
		return model.Violation
	}

	return model.Complete
}

// getStateOfNodes checks all of the nodes' state quickly
// Returns `model.NodesComplete` if all nodes are satisfied.
// Does not check for loops.
func (p Puzzle) getStateOfNodes() model.State {
	if p.areNodesComplete() {
		return model.NodesComplete
	}

	for _, nm := range p.metas {
		if nm.IsComplete {
			continue
		}

		// TODO this Filter might need to be CheckComplete instead. This might
		// destroy larger puzzles...
		switch s := nm.Filter(&p.edges); s {
		case model.Complete:
		default:
			return s
		}
	}

	return model.NodesComplete
}
