package reader

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

const (
	websiteWhiteNode byte = 'W'
	websiteBlackNode byte = 'B'

	websiteDecodeCharMagicNum byte = 96
)

const (
	websiteCachePuzzlesFilename = `./reader/websitecache/puzzles.txt`
)

func FromWebsiteTask(
	numEdges int,
	puzzID string,
	input string,
) (PuzzleDef, error) {

	go cacheTaskToFile(numEdges, puzzID, input)

	return fromWebsiteTask(numEdges, puzzID, input)
}

func cacheTaskToFile(numEdges int, puzzID, input string) {
	f, err := os.OpenFile(websiteCachePuzzlesFilename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error opening file: %v\n", err)
		return
	}
	defer f.Close()

	line := fmt.Sprintf("%d:%s:%s\n", numEdges, puzzID, input)

	if _, err := f.WriteString(line); err != nil {
		log.Printf("error WriteString file: %v\n", err)
		return
	}
}

func fromWebsiteTask(
	numEdges int,
	puzzID string,
	input string,
) (PuzzleDef, error) {

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

	pd := PuzzleDef{
		NumEdges:    numEdges,
		Description: `PuzzleID: ` + puzzID,
	}
	numNodes := numEdges + 1
	maxNodexIndex := numNodes * numNodes

	reader := strings.NewReader(input)
	for nodeIndex := 0; nodeIndex <= maxNodexIndex; {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return PuzzleDef{}, fmt.Errorf(`problem on ReadByte: %+v`, err)
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
			if err == io.EOF {
				break
			}
			return PuzzleDef{}, fmt.Errorf(`problem on ReadByte: %+v`, err)
		}
		err = reader.UnreadByte()
		if err != nil {
			return PuzzleDef{}, fmt.Errorf(`problem on UnreadByte: %+v`, err)
		}
		val, err := strconv.Atoi(string(value))
		if err != nil {
			return PuzzleDef{}, fmt.Errorf(`expected value from bytes: %+v`, err)
		}

		// read the next char(s) as an int
		pd.Nodes = append(pd.Nodes, model.NodeLocation{
			Row:     nodeIndex / numNodes,
			Col:     nodeIndex % numNodes,
			IsWhite: isWhite,
			Value:   int8(val),
		})
		nodeIndex++
	}

	return pd, nil
}
