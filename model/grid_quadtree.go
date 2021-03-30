package model

import (
	"fmt"
	"strconv"
	"strings"
)

type quadTreeNode struct {
	coord NodeCoord
	value OutgoingEdges

	right *quadTreeNode
	above *quadTreeNode
	left  *quadTreeNode
	below *quadTreeNode
}

func (qtn *quadTreeNode) writePrefix(
	sb *strings.Builder,
	depth int,
) {
	for i := 0; i < depth; i++ {
		sb.WriteString("\t")
	}
}

func (qtn *quadTreeNode) writeString(
	sb *strings.Builder,
	depth int,
) {

	if qtn == nil {
		sb.WriteString("(*quadTreeNode)(nil)")
		return
	}
	sb.WriteString("&quadTreeNode{\n")

	qtn.writePrefix(sb, depth)
	sb.WriteString("coord: ")
	sb.WriteString(fmt.Sprintf("%+v", qtn.coord))
	sb.WriteString(",\n")

	qtn.writePrefix(sb, depth)
	sb.WriteString("value: ")
	sb.WriteString(fmt.Sprintf("%+v", qtn.value))
	sb.WriteString(",\n")

	qtn.writePrefix(sb, depth)
	sb.WriteString("right: ")
	qtn.right.writeString(sb, depth+1)
	sb.WriteString(",\n")

	qtn.writePrefix(sb, depth)
	sb.WriteString("above: ")
	qtn.above.writeString(sb, depth+1)
	sb.WriteString(",\n")

	qtn.writePrefix(sb, depth)
	sb.WriteString("left:  ")
	qtn.left.writeString(sb, depth+1)
	sb.WriteString(",\n")

	qtn.writePrefix(sb, depth)
	sb.WriteString("below: ")
	qtn.below.writeString(sb, depth+1)
	sb.WriteString(",\n")

	qtn.writePrefix(sb, depth-1)
	sb.WriteString("}")
}

func (qtn *quadTreeNode) getValue(nc NodeCoord, maxDepth int) (OutgoingEdges, int) {
	if qtn == nil {
		return OutgoingEdges{}, maxDepth
	}

	if qtn.coord == nc {
		return qtn.value, maxDepth
	}

	c := qtn.toHeading(nc)

	targetNode := qtn.below
	switch c {
	case HeadLeft:
		targetNode = qtn.left
	case HeadRight:
		targetNode = qtn.right
	case HeadUp:
		targetNode = qtn.above
	}
	return targetNode.getValue(nc, maxDepth+1)
}

func (qtn *quadTreeNode) toHeading(
	nc NodeCoord,
) Cardinal {
	if nc.Row == qtn.coord.Row {
		if nc.Col < qtn.coord.Col {
			return HeadLeft
		}
		return HeadRight
	}

	if nc.Row < qtn.coord.Row {
		return HeadUp
	}
	return HeadDown
}

func (qtn *quadTreeNode) copy() *quadTreeNode {
	if qtn == nil {
		return nil
	}
	return &quadTreeNode{
		coord: qtn.coord,
		value: qtn.value,
		left:  qtn.left,
		right: qtn.right,
		above: qtn.above,
		below: qtn.below,
	}
}

func (qtn *quadTreeNode) applyUpdate(
	update gridUpdate,
	maxDepth int,
) (_ *quadTreeNode, maxSeenDepth int, addedChild bool) {

	if qtn == nil {
		return &quadTreeNode{
			coord: update.coord,
			value: update.newVal,
		}, maxDepth, true
	} else if qtn.coord == update.coord {
		updatedNode := qtn.copy()
		updatedNode.value = update.newVal
		return updatedNode, maxDepth, false
	}

	c := qtn.toHeading(update.coord)

	var newChild *quadTreeNode
	var child *quadTreeNode

	updatedNode := qtn.copy()

	switch c {
	case HeadLeft:
		child = qtn.left
		defer func() {
			updatedNode.left = newChild
		}()
	case HeadRight:
		child = qtn.right
		defer func() {
			updatedNode.right = newChild
		}()
	case HeadUp:
		child = qtn.above
		defer func() {
			updatedNode.above = newChild
		}()
	case HeadDown:
		child = qtn.below
		defer func() {
			updatedNode.below = newChild
		}()
	}

	newChild, maxSeenDepth, addedChild = child.applyUpdate(
		update,
		maxDepth+1,
	)

	return updatedNode, maxSeenDepth, addedChild
}

type quadTree struct {
	maxEdgeIndex uint8

	root     *quadTreeNode
	numNodes int
}

func newQuadTree(maxEdgeIndex int) Grid {
	return &quadTree{
		maxEdgeIndex: uint8(maxEdgeIndex),
	}
}

var _ Grid = (*quadTree)(nil)

func (t *quadTree) String() string {
	var sb strings.Builder
	sb.WriteString("quadTree{\n\tmaxEdgeIndex: ")
	sb.WriteString(strconv.Itoa(int(t.maxEdgeIndex)))
	sb.WriteString("\n\tnumNodes: ")
	sb.WriteString(strconv.Itoa(t.numNodes))
	sb.WriteString("\n\troot: ")
	t.root.writeString(&sb, 2)
	sb.WriteString("\n}")
	return sb.String()
}

func (t *quadTree) IsInBounds(nc NodeCoord) bool {
	if nc.Row < 0 || nc.Col < 0 {
		return false
	}
	return uint8(nc.Row) <= t.maxEdgeIndex && uint8(nc.Col) <= t.maxEdgeIndex
}

func (t *quadTree) Get(nc NodeCoord) OutgoingEdges {
	oe, _ := t.root.getValue(nc, 0)

	return oe
}

func (t *quadTree) withUpdates(updates []gridUpdate) Grid {
	newRoot := t.root
	numNodes := t.numNodes
	var addedChild bool

	for _, update := range updates {
		newRoot, _, addedChild = newRoot.applyUpdate(update, 0)
		if addedChild {
			numNodes++
		}
	}

	return &quadTree{
		maxEdgeIndex: t.maxEdgeIndex,
		root:         newRoot,
		numNodes:     numNodes,
	}
}

func (t *quadTree) Copy() Grid {
	return &quadTree{
		maxEdgeIndex: t.maxEdgeIndex,
		root:         t.root.copy(),
		numNodes:     t.numNodes,
	}
}
