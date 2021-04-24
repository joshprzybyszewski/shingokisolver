// nolint:gocyclo
package reader

import (
	"errors"
	"strconv"
	"strings"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

const (
	emptyNode byte = '.'
	BlackNode byte = 'b'
	WhiteNode byte = 'w'
)

func FromString(input string) (model.Definition, error) {
	pd := model.Definition{}
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
			if (b != BlackNode) && (b != WhiteNode) {
				return model.Definition{}, errors.New(`unexpected byte: `)
			}
			isWhite := b == WhiteNode
			var value []byte
			b, err = r.ReadByte()
			if err != nil {
				return model.Definition{}, errors.New(`expected bytes for value: `)
			}
			value = append(value, b)
			b, err = r.ReadByte()
			if err == nil {
				if b == emptyNode || b == BlackNode || b == WhiteNode {
					err = r.UnreadByte()
					if err != nil {
						return model.Definition{}, errors.New(`problem on UnreadByte: `)
					}
				} else {
					value = append(value, b)
				}
			}
			val, err := strconv.Atoi(string(value))
			if err != nil {
				return model.Definition{}, errors.New(`expected value from bytes: `)
			}

			pd.Nodes = append(pd.Nodes, model.NodeLocation{
				Row:     lineIndex,
				Col:     colIndex,
				IsWhite: isWhite,
				Value:   int8(val),
			})
		}
		if colIndex != len(lines) {
			return model.Definition{}, errors.New(`a line had a wrong number of nodes in it`)
		}

	}
	return pd, nil
}
