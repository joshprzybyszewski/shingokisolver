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
	diff difficulty,
	secretMeta map[string]string,
) (string, http.Header) {
	return ``, nil
}

func getPuzzleInfo(
	doc *goquery.Document,
	secretMeta map[string]string,
) (string, string, error) {
	return ``, ``, ``, errors.New(`getPuzzleInfo public sauce`)
}

func getPostSolutionData(
	wp websitePuzzle,
	res solvers.SolvedResults,
) (http.Header, io.Reader) {
	panic(`public sauce: getPostSolutionData`)
	return ``, nil, nil
}

func getHallOfFameSubmission(
	wp websitePuzzle,
	doc *goquery.Document,
) (string, http.Header, io.Reader, error) {
	return ``, nil, nil, errors.New(`public sauce: getHallOfFameSubmission`)
}
