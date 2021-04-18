package compete

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

func submitAnswer(
	wp websitePuzzle,
	res solvers.SolvedResults,
) error {

	defer func() {
		writeToFile(`answer.txt`, []byte(res.Puzzle.Solution()))
	}()

	postURL, header, data := getPostSolutionData(wp, res)

	resp, respHeaders, err := post(postURL, header, data)
	if err != nil {
		return err
	}
	writeToFile(`postAnswer.html`, resp)
	writeToFile(`postRequestHeaders.txt`, []byte(fmt.Sprintf("%+v", header)))
	writeToFile(`postAnswerHeaders.txt`, []byte(fmt.Sprintf("%+v", respHeaders)))

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return err
	}

	hallUrl, header, data, err := getHallOfFameSubmission(wp, doc)
	if err != nil {
		return err
	}
	resp, respHeaders, err = post(hallUrl, header, data)
	if err != nil {
		return err
	}
	writeToFile(`postHallOfFame.html`, resp)
	writeToFile(`postHallOfFameRequestHeaders.txt`, []byte(fmt.Sprintf("%+v", header)))
	writeToFile(`postHallOfFameResponseHeaders.txt`, []byte(fmt.Sprintf("%+v", respHeaders)))

	return nil
}
