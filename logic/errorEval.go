package logic

import (
	"log"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

var _ evaluator = errEval(``)

type errEval string

func (ee errEval) evaluate(model.GetEdger) model.EdgeState {
	log.Printf("errEval: %s", ee)
	return model.EdgeErrored
}
