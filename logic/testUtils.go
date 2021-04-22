package logic

import "github.com/joshprzybyszewski/shingokisolver/model"

func GetMinArmCheckEvaluation(
	node model.Node,
	myDir model.Cardinal,
	myIndex int,
	ge model.GetEdger,
) model.EdgeState {

	return newMinArmCheckEvaluator(node, myDir, myIndex).evaluate(ge)
}

func GetNumInDirectionExperimental(
	ge model.GetEdger,
	nValue int,
	start model.NodeCoord,
	dir model.Cardinal,
) (numPossibleEdges, existingChainStartIndex, numInLastExistingChain int) {
	return getNumInDirection(ge, nValue, start, dir)
}
