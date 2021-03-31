package compete

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/joshprzybyszewski/shingokisolver/model"
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

	wp := websitePuzzle{
		pd: reader.PuzzleDef{
			Description: `WebsitePuzzle`,
			NumEdges:    size,
		},
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return websitePuzzle{}, err
	}

	// Find the review items
	nodes, err := convertTaskToNodeLocations(
		wp.pd.NumEdges,
		doc.Find(`#rel`).Find(`script`).Text(),
	)
	if err != nil {
		return websitePuzzle{}, err
	}
	wp.pd.Nodes = nodes

	wp.id = doc.Find(`#puzzleID`).First().Text()
	wp.pd.Description = `PuzzleID: ` + wp.id

	return wp, nil
}

const (
	expGameScriptPrefix = ` var Game = {}; var Puzzle = {}; var task = '`
)

func convertTaskToNodeLocations(
	numEdges int,
	gameScript string,
) ([]model.NodeLocation, error) {
	if gameScript[:len(expGameScriptPrefix)] != expGameScriptPrefix {
		return nil, fmt.Errorf(`unexpected prefix! %q`, gameScript)
	}
	taskString := gameScript[len(expGameScriptPrefix):]
	end := strings.Index(taskString, `'`)
	taskString = taskString[:end]

	return reader.FromWebsiteTask(numEdges, taskString)
}
