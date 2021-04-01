package compete

import (
	"log"

	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

func submitAnswer(
	wp websitePuzzle,
	res solvers.SolvedResults,
) error {

	header, data := getPostSolutionData(wp, res)

	resp, err := post(`https://www.puzzle-shingoki.com/`, header, data)
	if err != nil {
		return err
	}
	log.Printf("resp: %s", resp)

	return nil
}
