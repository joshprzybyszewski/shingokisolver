package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeCoordDistanceTo(t *testing.T) {
	testCases := []struct {
		start   NodeCoord
		other   NodeCoord
		expDist int
	}{{
		start:   NewCoord(0, 0),
		other:   NewCoord(0, 4),
		expDist: 4,
	}, {
		start:   NewCoord(2, 2),
		other:   NewCoord(6, 2),
		expDist: 4,
	}, {
		start:   NewCoord(17, 17),
		other:   NewCoord(21, 21),
		expDist: 8,
	}, {
		start:   NewCoord(6, 2),
		other:   NewCoord(2, 2),
		expDist: 4,
	}, {
		start:   NewCoord(2, 8),
		other:   NewCoord(2, 2),
		expDist: 6,
	}, {
		start:   NewCoord(6, 8),
		other:   NewCoord(2, 2),
		expDist: 10,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expDist, tc.start.DistanceTo(tc.other))
	}
}

func TestNodeCoordTranslate(t *testing.T) {
	testCases := []struct {
		start NodeCoord
		exp   NodeCoord
		dir   Cardinal
	}{{
		start: NewCoord(0, 0),
		dir:   HeadDown,
		exp:   NewCoord(1, 0),
	}, {
		start: NewCoord(2, 2),
		dir:   HeadRight,
		exp:   NewCoord(2, 3),
	}, {
		start: NewCoord(10, 10),
		dir:   HeadUp,
		exp:   NewCoord(9, 10),
	}, {
		start: NewCoord(2, 2),
		dir:   HeadLeft,
		exp:   NewCoord(2, 1),
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.exp, tc.start.Translate(tc.dir))
	}
}

func TestNodeCoordTranslateAlongArm(t *testing.T) {
	testCases := []struct {
		start NodeCoord
		exp   NodeCoord
		arm   Arm
	}{{
		start: NewCoord(0, 0),
		arm: Arm{
			Heading: HeadDown,
			Len:     6,
		},
		exp: NewCoord(6, 0),
	}, {
		start: NewCoord(2, 2),
		arm: Arm{
			Heading: HeadRight,
			Len:     3,
		},
		exp: NewCoord(2, 5),
	}, {
		start: NewCoord(10, 10),
		arm: Arm{
			Heading: HeadUp,
			Len:     3,
		},
		exp: NewCoord(7, 10),
	}, {
		start: NewCoord(2, 2),
		arm: Arm{
			Heading: HeadLeft,
			Len:     2,
		},
		exp: NewCoord(2, 0),
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.exp, tc.start.TranslateAlongArm(tc.arm))
	}
}
