package model

import "fmt"

type Arm struct {
	Len     int8
	Heading Cardinal
}

func (a Arm) String() string {
	return fmt.Sprintf("Arm{Len: %d, Heading: %s}", a.Len, a.Heading)
}

func (a Arm) EndFrom(
	start NodeCoord,
) NodeCoord {
	return start.TranslateAlongArm(a)
}
