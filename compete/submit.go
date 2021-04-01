package compete

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

func submitAnswer(
	wp websitePuzzle,
	res solvers.SolvedResults,
) error {

	defer func() {
		writeToFile(`./temp/answer.txt`, []byte(res.Puzzle.String()))
	}()

	header, data := getPostSolutionData(wp, res)

	resp, err := post(`https://www.puzzle-shingoki.com/`, header, data)
	if err != nil {
		return err
	}
	writeToFile(`./temp/postAnswer.html`, resp)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return err
	}

	hallUrl, header, data := getHallOfFameSubmission(doc)
	resp, err = post(hallUrl, header, data)
	if err != nil {
		return err
	}
	writeToFile(`./temp/postHallOfFame.html`, resp)

	return nil
}
