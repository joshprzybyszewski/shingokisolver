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
		writeToFile(`./temp/answer.txt`, []byte(res.Puzzle.Solution()))
	}()

	postURL, header, data := getPostSolutionData(wp, res)

	resp, respHeaders, err := post(postURL, header, data)
	if err != nil {
		return err
	}
	writeToFile(`./temp/postAnswer.html`, resp)
	writeToFile(`./temp/postRequestHeaders.txt`, []byte(fmt.Sprintf("%+v", header)))
	writeToFile(`./temp/postAnswerHeaders.txt`, []byte(fmt.Sprintf("%+v", respHeaders)))

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
	writeToFile(`./temp/postHallOfFame.html`, resp)
	writeToFile(`./temp/postHallOfFameRequestHeaders.txt`, []byte(fmt.Sprintf("%+v", header)))
	writeToFile(`./temp/postHallOfFameResponseHeaders.txt`, []byte(fmt.Sprintf("%+v", respHeaders)))

	return nil
}
