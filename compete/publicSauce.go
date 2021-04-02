// +build publicSauce

package compete

import (
	"errors"
	"io"
	"net/http"

	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

func requestNewPuzzle(
	edges int,
) (string, http.Header) {
	return ``, nil
}

func getPuzzleInfo(
	doc *goquery.Document,
) (string, string, string, error) {
	return ``, ``, ``, errors.New(`getPuzzleInfo public sauce`)
}

func getPostSolutionData(
	wp websitePuzzle,
	res solvers.SolvedResults,
) (http.Header, io.Reader) {
	panic(`public sauce: getPostSolutionData`)
	return nil, nil
}

func getHallOfFameSubmission(
	doc *goquery.Document,
) (string, http.Header, io.Reader) {
	panic(`public sauce: getHallOfFameSubmission`)
	return nil, nil
}
