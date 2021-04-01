// +build !secretSauce

package compete

import (
	"io"
	"net/http"

	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

func getPostSolutionData(
	wp websitePuzzle,
	res solvers.SolvedResults,
) (http.Header, io.Reader) {
	return nil, nil
}
