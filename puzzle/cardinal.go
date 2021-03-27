package puzzle

type cardinal uint8

const (
	headNowhere cardinal = 0
	headRight   cardinal = 1
	headUp      cardinal = 2
	headLeft    cardinal = 3
	headDown    cardinal = 4
)

func (c cardinal) String() string {
	switch c {
	case headNowhere:
		return `headNowhere`
	case headRight:
		return `headRight`
	case headUp:
		return `headUp`
	case headLeft:
		return `headLeft`
	case headDown:
		return `headDown`
	default:
		return `unknown cardinal`
	}
}
