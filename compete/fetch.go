package compete

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"

	"github.com/joshprzybyszewski/shingokisolver/reader"
)

type difficulty int

const (
	easy   difficulty = 0
	medium difficulty = 1
	hard   difficulty = 2
)

func (d difficulty) String() string {
	switch d {
	case easy:
		return `easy`
	case medium:
		return `medium`
	case hard:
		return `hard`
	default:
		return `unknown difficulty`
	}
}

type websitePuzzle struct {
	secret map[string]string
	id     string
	pd     reader.PuzzleDef
}

func (wp websitePuzzle) String() string {
	return fmt.Sprintf("websitePuzzle{id: %s, pd: %s}", wp.id, wp.pd)
}

func getPuzzle(
	size int,
	diff difficulty,
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
		puzzID, diff.String(),
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
