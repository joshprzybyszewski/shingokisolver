package compete

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"

	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/reader"
)

type websitePuzzle struct {
	secret map[string]string
	id     string
	pd     model.Definition
}

func (wp websitePuzzle) String() string {
	return fmt.Sprintf("websitePuzzle{id: %s, pd: %s}", wp.id, wp.pd)
}

func getPuzzle(
	size int,
	diff model.Difficulty,
) (websitePuzzle, error) {

	secret := map[string]string{}
	url, header := requestNewPuzzle(size, diff, secret)

	resp, respHeaders, err := get(url, header)
	if err != nil {
		return websitePuzzle{}, err
	}
	writeToFile(`getPuzzleResp.html`, resp)
	writeToFile(`getPuzzleRespHeaders.txt`, []byte(fmt.Sprintf("%+v", respHeaders)))

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return websitePuzzle{}, err
	}

	puzzID, taskString, err := getPuzzleInfo(doc, secret)
	if err != nil {
		return websitePuzzle{}, err
	}

	pd, err := reader.FromWebsiteTask(
		size,
		diff,
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
