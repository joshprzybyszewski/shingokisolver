package compete

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"

	"github.com/joshprzybyszewski/shingokisolver/reader"
)

type websitePuzzle struct {
	id string
	pd reader.PuzzleDef

	secret string
}

func getPuzzle(
	size int,
) (websitePuzzle, error) {

	url, header := requestNewPuzzle(size)

	resp, err := get(url, header)
	if err != nil {
		return websitePuzzle{}, err
	}
	writeToFile(`./temp/getPuzzleResp.html`, resp)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return websitePuzzle{}, err
	}

	puzzID, taskString, secret, err := getPuzzleInfo(doc)

	pd, err := reader.FromWebsiteTask(
		size,
		puzzID,
		taskString,
	)
	if err != nil {
		return websitePuzzle{}, err
	}

	return websitePuzzle{
		id:     puzzID,
		pd:     pd,
		secret: secret,
	}, nil
}
