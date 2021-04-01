package compete

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/joshprzybyszewski/shingokisolver/reader"
)

type websitePuzzle struct {
	id string
	pd reader.PuzzleDef
}

func getPuzzle(
	size int,
) (websitePuzzle, error) {

	// TODO:#
	log.Printf("TODO getPuzzle from website")
	resp, err := get(`https://www.puzzle-shingoki.com/`, nil)
	if err != nil {
		return websitePuzzle{}, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return websitePuzzle{}, err
	}

	puzzID := doc.Find(`#puzzleID`).First().Text()

	taskString, err := getTaskFromScriptText(
		doc.Find(`#rel`).Find(`script`).Text(),
	)
	if err != nil {
		return websitePuzzle{}, err
	}

	pd, err := reader.FromWebsiteTask(
		size,
		puzzID,
		taskString,
	)
	if err != nil {
		return websitePuzzle{}, err
	}

	return websitePuzzle{
		id: puzzID,
		pd: pd,
	}, nil
}

const (
	expGameScriptPrefix = ` var Game = {}; var Puzzle = {}; var task = '`
)

func getTaskFromScriptText(
	gameScript string,
) (string, error) {
	if gameScript[:len(expGameScriptPrefix)] != expGameScriptPrefix {
		return ``, fmt.Errorf(`unexpected prefix! %q`, gameScript)
	}
	taskString := gameScript[len(expGameScriptPrefix):]
	end := strings.Index(taskString, `'`)
	return taskString[:end], nil
}
