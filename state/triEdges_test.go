package state

import (
	"testing"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/stretchr/testify/assert"
)

func TestAllExist(t *testing.T) {

	ets := New(5)

	nc00 := model.NewCoord(0, 0)
	assert.Equal(t, model.Incomplete, ets.setEdge(model.NewEdgePair(nc00, model.HeadRight)))

	arm := model.Arm{
		Heading: model.HeadRight,
		Len:     1,
	}
	assert.True(t, ets.AllExist(nc00, arm))
	arm.Len = 2
	assert.False(t, ets.AllExist(nc00, arm))
	arm.Len = 1
	arm.Heading = model.HeadLeft
	assert.False(t, ets.AllExist(nc00, arm))
	arm.Heading = model.HeadDown
	assert.False(t, ets.AllExist(nc00, arm))
	arm.Heading = model.HeadUp
	assert.False(t, ets.AllExist(nc00, arm))

	nc22 := model.NewCoord(2, 2)
	assert.Equal(t, model.Incomplete, ets.setEdge(
		model.NewEdgePair(nc22, model.HeadRight),
	))
	nc23 := model.NewCoord(2, 3)
	assert.Equal(t, model.Incomplete, ets.setEdge(
		model.NewEdgePair(nc23, model.HeadRight),
	))

	arm = model.Arm{
		Heading: model.HeadRight,
		Len:     1,
	}
	assert.True(t, ets.AllExist(nc22, arm))
	assert.True(t, ets.AllExist(nc23, arm))
	arm.Len = 2
	assert.True(t, ets.AllExist(nc22, arm))
	assert.False(t, ets.AllExist(nc23, arm))
	arm.Len = 3
	assert.False(t, ets.AllExist(nc22, arm))
	arm.Len = 1
	arm.Heading = model.HeadLeft
	assert.False(t, ets.AllExist(nc22, arm))
	assert.True(t, ets.AllExist(nc23, arm))
	arm.Heading = model.HeadDown
	assert.False(t, ets.AllExist(nc22, arm))
	arm.Heading = model.HeadUp
	assert.False(t, ets.AllExist(nc22, arm))

	nc20 := model.NewCoord(2, 0)
	assert.Equal(t, model.Incomplete, ets.setEdge(
		model.NewEdgePair(nc20, model.HeadRight),
	))

	arm = model.Arm{
		Heading: model.HeadLeft,
		Len:     2,
	}
	assert.False(t, ets.AllExist(nc22, arm))
}

func TestAny(t *testing.T) {
	ets := New(5)

	nc00 := model.NewCoord(0, 0)
	nc01 := model.NewCoord(0, 1)
	assert.Equal(t, model.Incomplete, ets.setEdge(model.NewEdgePair(nc00, model.HeadRight)))

	arm := model.Arm{
		Heading: model.HeadRight,
		Len:     1,
	}
	var anyExist, anyAvoided bool
	anyExist, anyAvoided = ets.Any(nc00, arm)
	assert.True(t, anyExist)
	assert.False(t, anyAvoided)
	arm.Len = 2
	anyExist, anyAvoided = ets.Any(nc00, arm)
	assert.True(t, anyExist)
	assert.False(t, anyAvoided)
	arm.Len = 1
	arm.Heading = model.HeadLeft
	anyExist, anyAvoided = ets.Any(nc00, arm)
	assert.False(t, anyExist)
	assert.True(t, anyAvoided)
	arm.Heading = model.HeadDown
	anyExist, anyAvoided = ets.Any(nc00, arm)
	assert.False(t, anyExist)
	assert.False(t, anyAvoided)
	arm.Heading = model.HeadUp
	anyExist, anyAvoided = ets.Any(nc00, arm)
	assert.False(t, anyExist)
	assert.True(t, anyAvoided)
	arm.Len = 2
	arm.Heading = model.HeadLeft
	anyExist, anyAvoided = ets.Any(nc01, arm)
	assert.True(t, anyExist)
	assert.True(t, anyAvoided)

	nc22 := model.NewCoord(2, 2)
	assert.Equal(t, model.Incomplete, ets.setEdge(
		model.NewEdgePair(nc22, model.HeadRight),
	))
	nc23 := model.NewCoord(2, 3)
	assert.Equal(t, model.Incomplete, ets.setEdge(
		model.NewEdgePair(nc23, model.HeadRight),
	))

	arm = model.Arm{
		Heading: model.HeadRight,
		Len:     1,
	}
	anyExist, anyAvoided = ets.Any(nc22, arm)
	assert.True(t, anyExist)
	assert.False(t, anyAvoided)
	anyExist, anyAvoided = ets.Any(nc23, arm)
	assert.True(t, anyExist)
	assert.False(t, anyAvoided)
	arm.Len = 2
	anyExist, anyAvoided = ets.Any(nc22, arm)
	assert.True(t, anyExist)
	assert.False(t, anyAvoided)
	anyExist, anyAvoided = ets.Any(nc23, arm)
	assert.True(t, anyExist)
	assert.False(t, anyAvoided)
	arm.Len = 3
	anyExist, anyAvoided = ets.Any(nc22, arm)
	assert.True(t, anyExist)
	assert.False(t, anyAvoided)
	arm.Len = 1
	arm.Heading = model.HeadLeft
	anyExist, anyAvoided = ets.Any(nc22, arm)
	assert.False(t, anyExist)
	assert.False(t, anyAvoided)
	anyExist, anyAvoided = ets.Any(nc23, arm)
	assert.True(t, anyExist)
	assert.False(t, anyAvoided)
	arm.Heading = model.HeadDown
	anyExist, anyAvoided = ets.Any(nc22, arm)
	assert.False(t, anyExist)
	assert.False(t, anyAvoided)
	arm.Heading = model.HeadUp
	anyExist, anyAvoided = ets.Any(nc22, arm)
	assert.False(t, anyExist)
	assert.False(t, anyAvoided)

	nc20 := model.NewCoord(2, 0)
	assert.Equal(t, model.Incomplete, ets.setEdge(
		model.NewEdgePair(nc20, model.HeadRight),
	))

	arm = model.Arm{
		Heading: model.HeadLeft,
		Len:     2,
	}
	anyExist, anyAvoided = ets.Any(nc22, arm)
	assert.True(t, anyExist)
	assert.False(t, anyAvoided)

	nc24 := model.NewCoord(2, 4)
	assert.Equal(t, model.Incomplete, ets.avoidEdge(
		model.NewEdgePair(nc24, model.HeadRight),
	))

	arm = model.Arm{
		Heading: model.HeadRight,
		Len:     3,
	}
	anyExist, anyAvoided = ets.Any(nc22, arm)
	assert.True(t, anyExist)
	assert.True(t, anyAvoided)
	anyExist, anyAvoided = ets.Any(nc23, arm)
	assert.True(t, anyExist)
	assert.True(t, anyAvoided)
	anyExist, anyAvoided = ets.Any(nc24, arm)
	assert.False(t, anyExist)
	assert.True(t, anyAvoided)
}

func TestGetMask(t *testing.T) {
	testCases := []struct {
		start       model.NodeCoord
		arm         model.Arm
		expBitaData bitData
	}{{
		start: model.NewCoord(0, 0),
		arm: model.Arm{
			Heading: model.HeadRight,
			Len:     1,
		},
		expBitaData: 1 << 0,
	}, {
		start: model.NewCoord(0, 0),
		arm: model.Arm{
			Heading: model.HeadDown,
			Len:     1,
		},
		expBitaData: 1 << 0,
	}, {
		start: model.NewCoord(0, 1),
		arm: model.Arm{
			Heading: model.HeadRight,
			Len:     1,
		},
		expBitaData: 1 << 1,
	}, {
		start: model.NewCoord(1, 0),
		arm: model.Arm{
			Heading: model.HeadDown,
			Len:     1,
		},
		expBitaData: 1 << 1,
	}, {
		start: model.NewCoord(0, 1),
		arm: model.Arm{
			Heading: model.HeadLeft,
			Len:     1,
		},
		expBitaData: 1 << 0,
	}, {
		start: model.NewCoord(1, 0),
		arm: model.Arm{
			Heading: model.HeadUp,
			Len:     1,
		},
		expBitaData: 1 << 0,
	}, {
		start: model.NewCoord(0, 0),
		arm: model.Arm{
			Heading: model.HeadRight,
			Len:     2,
		},
		expBitaData: 1<<1 | 1<<0,
	}, {
		start: model.NewCoord(0, 0),
		arm: model.Arm{
			Heading: model.HeadDown,
			Len:     2,
		},
		expBitaData: 1<<1 | 1<<0,
	}, {
		start: model.NewCoord(0, 2),
		arm: model.Arm{
			Heading: model.HeadLeft,
			Len:     2,
		},
		expBitaData: 1<<1 | 1<<0,
	}, {
		start: model.NewCoord(2, 0),
		arm: model.Arm{
			Heading: model.HeadUp,
			Len:     2,
		},
		expBitaData: 1<<1 | 1<<0,
	}, {
		start: model.NewCoord(0, 1),
		arm: model.Arm{
			Heading: model.HeadLeft,
			Len:     2,
		},
		expBitaData: 1 << 0,
	}, {
		start: model.NewCoord(0, 1),
		arm: model.Arm{
			Heading: model.HeadLeft,
			Len:     3,
		},
		expBitaData: 1 << 0,
	}, {
		start: model.NewCoord(0, 2),
		arm: model.Arm{
			Heading: model.HeadLeft,
			Len:     3,
		},
		expBitaData: 1<<1 | 1<<0,
	}, {
		start: model.NewCoord(1, 0),
		arm: model.Arm{
			Heading: model.HeadUp,
			Len:     2,
		},
		expBitaData: 1 << 0,
	}, {
		start: model.NewCoord(0, 0),
		arm: model.Arm{
			Heading: model.HeadDown,
			Len:     MaxEdges,
		},
		expBitaData: armLenMasks[len(armLenMasks)-1],
	}, {
		start: model.NewCoord(1, 0),
		arm: model.Arm{
			Heading: model.HeadDown,
			Len:     MaxEdges,
		},
		expBitaData: armLenMasks[len(armLenMasks)-1] << 1,
	}, {
		start: model.NewCoord(1, 0),
		arm: model.Arm{
			Heading: model.HeadUp,
			Len:     MaxEdges,
		},
		expBitaData: 1 << 0,
	}}

	for _, tc := range testCases {
		actMask := getMask(tc.start, tc.arm)
		assert.Equal(t, tc.expBitaData, actMask, "getMask(%s, %s) expected: 0x%08b but was 0x%08b", tc.start, tc.arm, tc.expBitaData, actMask)
	}
}
