package model

const (
	Unknown Difficulty = 0
	Easy    Difficulty = 1
	Medium  Difficulty = 2
	Hard    Difficulty = 3
)

var (
	AllDifficulties = []Difficulty{
		Easy,
		Medium,
		Hard,
	}
)

type Difficulty int

func NewDifficulty(s string) Difficulty {
	switch s {
	case `easy`:
		return Easy
	case `medium`:
		return Medium
	case `hard`:
		return Hard
	default:
		return Unknown
	}
}

func (d Difficulty) String() string {
	switch d {
	case Easy:
		return `easy`
	case Medium:
		return `medium`
	case Hard:
		return `hard`
	default:
		return `Unknown Difficulty`
	}
}
