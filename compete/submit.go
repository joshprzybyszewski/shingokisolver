package compete

import (
	"log"

	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

func submitAnswer(
	res solvers.SolvedResults,
) error {
	// TODO
	log.Printf("TODO need to submit %s\n", res)

	resp, err := post(`https://www.puzzle-shingoki.com/`, nil)
	if err != nil {
		return err
	}
	log.Printf("resp: %s", resp)
	return nil
}
