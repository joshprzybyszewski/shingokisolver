package model

type EdgeState uint8

const (
	EdgeErrored     EdgeState = 0
	EdgeOutOfBounds EdgeState = 1
	EdgeUnknown     EdgeState = 2
	EdgeExists      EdgeState = 3
	EdgeAvoided     EdgeState = 4
)

func (es EdgeState) String() string {
	switch es {
	case EdgeErrored:
		return `EdgeErrored`
	case EdgeOutOfBounds:
		return `EdgeOutOfBounds`
	case EdgeUnknown:
		return `EdgeUnknown`
	case EdgeExists:
		return `EdgeExists`
	case EdgeAvoided:
		return `EdgeAvoided`
	default:
		return `Unknown EdgeState`
	}
}
