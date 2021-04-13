package model

type Cardinal uint8

const (
	HeadNowhere Cardinal = 0
	HeadRight   Cardinal = 1
	HeadUp      Cardinal = 2
	HeadLeft    Cardinal = 3
	HeadDown    Cardinal = 4
)

func (c Cardinal) String() string {
	switch c {
	case HeadNowhere:
		return `HeadNowhere`
	case HeadRight:
		return `HeadRight`
	case HeadUp:
		return `HeadUp`
	case HeadLeft:
		return `HeadLeft`
	case HeadDown:
		return `HeadDown`
	default:
		return `unknown Cardinal`
	}
}

func (c Cardinal) Opposite() Cardinal {
	switch c {
	case HeadRight:
		return HeadLeft
	case HeadLeft:
		return HeadRight
	case HeadUp:
		return HeadDown
	case HeadDown:
		return HeadUp
	}
	return HeadNowhere
}

func (c Cardinal) Perpendiculars() []Cardinal {
	switch c {
	case HeadRight, HeadLeft:
		return []Cardinal{HeadUp, HeadDown}
	case HeadUp, HeadDown:
		return []Cardinal{HeadRight, HeadLeft}
	}
	return nil
}

var (
	AllCardinals = []Cardinal{
		HeadRight,
		HeadUp,
		HeadLeft,
		HeadDown,
	}
	AllCardinalsMap = map[Cardinal]struct{}{
		HeadRight: {},
		HeadUp:    {},
		HeadLeft:  {},
		HeadDown:  {},
	}
)
