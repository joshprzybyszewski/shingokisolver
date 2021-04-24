package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArmEndFrom(t *testing.T) {
	r4 := Arm{
		Len:     4,
		Heading: HeadRight,
	}
	d4 := Arm{
		Len:     4,
		Heading: HeadDown,
	}

	testCases := []struct {
		arm    Arm
		start  NodeCoord
		expEnd NodeCoord
	}{{
		arm:    r4,
		start:  NewCoord(0, 0),
		expEnd: NewCoord(0, 4),
	}, {
		arm:    r4,
		start:  NewCoord(3, 3),
		expEnd: NewCoord(3, 7),
	}, {
		arm:    d4,
		start:  NewCoord(2, 2),
		expEnd: NewCoord(6, 2),
	}, {
		arm:    d4,
		start:  NewCoord(17, 17),
		expEnd: NewCoord(21, 17),
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expEnd, tc.arm.EndFrom(tc.start))
	}
}
