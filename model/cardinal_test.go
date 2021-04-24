package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllCards(t *testing.T) {
	assert.Len(t, AllCardinals, 4)
	assert.Len(t, AllCardinalsMap, 4)
}

func TestPerpendiculars(t *testing.T) {
	tests := []struct {
		name string
		want []Cardinal
		args Cardinal
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
			if got := tt.args.Perpendiculars(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Perpendiculars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpposite(t *testing.T) {
	tests := []struct {
		name string
		args Cardinal
		want Cardinal
	}{{
		args: HeadUp,
		want: HeadDown,
	}, {
		args: HeadDown,
		want: HeadUp,
	}, {
		args: HeadRight,
		want: HeadLeft,
	}, {
		args: HeadLeft,
		want: HeadRight,
	}, {
		args: HeadNowhere,
		want: HeadNowhere,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.Opposite(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Opposite() = %v, want %v", got, tt.want)
			}
		})
	}
}
