package model

import (
	"reflect"
	"testing"
)

func TestPerpendiculars(t *testing.T) {
	tests := []struct {
		name string
		args Cardinal
		want []Cardinal
	}{{
		args: HeadUp,
		want: []Cardinal{HeadRight, HeadLeft},
	}, {
		args: HeadDown,
		want: []Cardinal{HeadRight, HeadLeft},
	}, {
		args: HeadRight,
		want: []Cardinal{HeadUp, HeadDown},
	}, {
		args: HeadLeft,
		want: []Cardinal{HeadUp, HeadDown},
	}, {
		args: HeadNowhere,
		want: nil,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Perpendiculars(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Perpendiculars() = %v, want %v", got, tt.want)
			}
		})
	}
}
