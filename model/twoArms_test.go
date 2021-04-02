package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTwoArmOptions(t *testing.T) {

	for val := int8(2); val <= 50; val++ {
		blackNode := NewNode(false, val)
		expTwoArms := longBuildTwoArms(blackNode)
		actTwoArms := BuildTwoArmOptions(blackNode)
		assert.ElementsMatch(t, expTwoArms, actTwoArms, `failed for node %+v`, blackNode)

		// black nodes can "point" in the four cardinal directions.
		// the first two segments are require for this pointing, so
		// the rest of the edges can contribute to this
		expNumOptions := 4 * (1 + int(val-2))
		assert.Equal(t, expNumOptions, len(expTwoArms), `unexpected num options for a black node: %+v`, blackNode)
		assert.Equal(t, expNumOptions, len(actTwoArms), `unexpected num options for a black node: %+v`, blackNode)

		whiteNode := NewNode(true, val)
		expTwoArms = longBuildTwoArms(whiteNode)
		actTwoArms = BuildTwoArmOptions(whiteNode)
		assert.ElementsMatch(t, expTwoArms, actTwoArms, `failed for node %+v`, whiteNode)

		expNumOptions = 2 * (1 + int(val-2))
		assert.Equal(t, expNumOptions, len(expTwoArms), `unexpected num options for a white node: %+v`, whiteNode)
		assert.Equal(t, expNumOptions, len(actTwoArms), `unexpected num options for a white node: %+v`, whiteNode)
	}
}

func longBuildTwoArms(n Node) []TwoArms {
	var options []TwoArms
	var added map[TwoArms]struct{}

	for _, heading1 := range AllCardinals {
		for _, heading2 := range AllCardinals {
			if n.IsInvalidMotions(heading1, heading2) {
				continue
			}
			for len1 := int8(1); len1 < n.Value(); len1++ {
				len2 := n.Value() - len1

				ta := TwoArms{
					One: Arm{
						Len:     len1,
						Heading: heading1,
					},
					Two: Arm{
						Len:     len2,
						Heading: heading2,
					},
				}
				if _, ok := added[ta]; ok {
					continue
				}

				alt := TwoArms{
					One: Arm{
						Len:     len2,
						Heading: heading2,
					},
					Two: Arm{
						Len:     len1,
						Heading: heading1,
					},
				}
				if _, ok := added[alt]; ok {
					continue
				}

				options = append(options, ta)
			}
		}
	}

	return options
}
