package reader

import (
	"errors"
	"strconv"
	"strings"

	"github.com/joshprzybyszewski/shingokisolver/puzzle"
)

const (
	emptyNode byte = '.'
	blackNode byte = 'b'
	whiteNode byte = 'w'
)

type PuzzleDef struct {
	NumEdges int
	Nodes    []puzzle.NodeLocation
}

func FromString(input string) (PuzzleDef, error) {
	pd := PuzzleDef{}
	lines := strings.Split(input, "\n")
	pd.NumEdges = len(lines) - 1

	for lineIndex, l := range lines {
		r := strings.NewReader(l)
		colIndex := 0
		for ; ; colIndex++ {
			b, err := r.ReadByte()
			if err != nil {
				break
			}
			if b == emptyNode {
				continue
			}
			if (b != blackNode) && (b != whiteNode) {
				return PuzzleDef{}, errors.New(`unexpected byte: `)
			}
			isWhite := b == whiteNode
			var value []byte
			b, err = r.ReadByte()
			if err != nil {
				return PuzzleDef{}, errors.New(`expected bytes for value: `)
			}
			value = append(value, b)
			b, err = r.ReadByte()
			if err == nil {
				if b == emptyNode || b == blackNode || b == whiteNode {
					err = r.UnreadByte()
					if err != nil {
						return PuzzleDef{}, errors.New(`problem on UnreadByte: `)
					}
				} else {
					value = append(value, b)
				}
			}
			val, err := strconv.Atoi(string(value))
			if err != nil {
				return PuzzleDef{}, errors.New(`expected value from bytes: `)
			}

			pd.Nodes = append(pd.Nodes, puzzle.NodeLocation{
				Row:     lineIndex,
				Col:     colIndex,
				IsWhite: isWhite,
				Value:   int8(val),
			})
		}
		if colIndex != len(lines) {
			return PuzzleDef{}, errors.New(`a line had a wrong number of nodes in it`)
		}

	}
	return pd, nil
}
