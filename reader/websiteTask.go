package reader

import (
	"errors"
	"strconv"
	"strings"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

const (
	websiteWhiteNode byte = 'W'
	websiteBlackNode byte = 'B'

	websiteDecodeCharMagicNum byte = 96
)

func FromWebsiteTask(
	numEdges int,
	input string,
) ([]model.NodeLocation, error) {

	// This func is based on the following source, grabbed from the js of the website:
	/*
		for (var t = [], e = 0, r = 0; r < i.task.length; r++)
		                    "W" == i.task[r] ? (t[e] = parseInt(i.task.substring(r + 1)),
		                    e++) : "B" == i.task[r] ? (t[e] = -parseInt(i.task.substring(r + 1)),
		                    e++) : i.task[r] >= "0" && i.task[r] <= "9" || (e += this.decodeChar(i.task[r]));
		                for (var r = 0; r < i.puzzleHeight + 1; r++) {
		                    this.task[r] = [],
		                    this.currentState.taskStatus[r] = [],
		                    this.currentState.cellHatch[r] = [];
		                    for (var e = 0; e < i.puzzleWidth + 1; e++) {
		                        var s = r * (i.puzzleWidth + 1) + e;
		                        this.currentState.taskStatus[r][e] = !1,
		                        "undefined" == typeof t[s] ? this.task[r][e] = 0 : this.task[r][e] = t[s]
		                    }
		                }
	*/

	pd := PuzzleDef{}
	pd.NumEdges = numEdges
	maxNodexIndex := numEdges * numEdges

	reader := strings.NewReader(input)
	for nodeIndex := 0; nodeIndex <= maxNodexIndex; {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		if b != websiteWhiteNode && b != websiteBlackNode {
			// increase
			nodeIndex += int(b - websiteDecodeCharMagicNum)
			continue
		}

		isWhite := b == websiteWhiteNode

		var value []byte
		b, err = reader.ReadByte()
		for err == nil && b >= '0' && b <= '9' {
			value = append(value, b)
			b, err = reader.ReadByte()
		}
		if err != nil {
			return nil, errors.New(`problem on ReadByte`)
		}
		err = reader.UnreadByte()
		if err != nil {
			return nil, errors.New(`problem on UnreadByte`)
		}
		val, err := strconv.Atoi(string(value))
		if err != nil {
			return nil, errors.New(`expected value from bytes`)
		}

		// read the next char(s) as an int
		pd.Nodes = append(pd.Nodes, model.NodeLocation{
			Row:     nodeIndex / numEdges,
			Col:     nodeIndex % numEdges,
			IsWhite: isWhite,
			Value:   int8(val),
		})
	}

	return pd.Nodes, nil
}
