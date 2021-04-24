package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEdgePairStandardize(t *testing.T) {
	testCases := []struct {
		input EdgePair
		exp   EdgePair
	}{{
		input: EdgePair{
			NodeCoord: NewCoord(0, 0),
			Cardinal:  HeadRight,
		},
		exp: EdgePair{
			NodeCoord: NewCoord(0, 0),
			Cardinal:  HeadRight,
		},
	}, {
		input: EdgePair{
			NodeCoord: NewCoord(3, 3),
			Cardinal:  HeadDown,
		},
		exp: EdgePair{
			NodeCoord: NewCoord(3, 3),
			Cardinal:  HeadDown,
		},
	}, {
		input: EdgePair{
			NodeCoord: NewCoord(2, 2),
			Cardinal:  HeadLeft,
		},
		exp: EdgePair{
			NodeCoord: NewCoord(2, 1),
			Cardinal:  HeadRight,
		},
	}, {
		input: EdgePair{
			NodeCoord: NewCoord(17, 17),
			Cardinal:  HeadUp,
		},
		exp: EdgePair{
			NodeCoord: NewCoord(16, 17),
			Cardinal:  HeadDown,
		},
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.exp, tc.input.Standardize())
	}
}

func TestEdgePairNext(t *testing.T) {
	testCases := []struct {
		input EdgePair
		exp   EdgePair
		dir   Cardinal
	}{{
		input: NewEdgePair(NewCoord(0, 0), HeadRight),
		exp:   NewEdgePair(NewCoord(0, 1), HeadRight),
		dir:   HeadRight,
	}, {
		input: NewEdgePair(NewCoord(3, 3), HeadDown),
		exp:   NewEdgePair(NewCoord(4, 3), HeadDown),
		dir:   HeadDown,
	}, {
		input: NewEdgePair(NewCoord(2, 2), HeadLeft),
		exp:   NewEdgePair(NewCoord(2, 0), HeadRight),
		dir:   HeadLeft,
	}, {
		input: NewEdgePair(NewCoord(17, 17), HeadUp),
		exp:   NewEdgePair(NewCoord(15, 17), HeadDown),
		dir:   HeadUp,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.exp, tc.input.Next(tc.dir))
	}
}
