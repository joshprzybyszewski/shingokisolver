package reader

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func CachedWebsitePuzzles() ([]PuzzleDef, error) {
	file, err := ioutil.ReadFile(websiteCachePuzzlesFilename)
	if err != nil {
		return nil, err
	}
	pds := []PuzzleDef{}
	includedPuzzleIDs := map[string]struct{}{}

	lines := strings.Split(string(file), "\n")
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		parts := strings.Split(l, ":")
		numEdges, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		puzzID := parts[1]
		task := parts[2]

		if _, ok := includedPuzzleIDs[puzzID]; ok {
			continue
		}

		pd, err := fromWebsiteTask(numEdges, puzzID, task)
		if err != nil {
			return nil, err
		}
		pds = append(pds, pd)
		includedPuzzleIDs[puzzID] = struct{}{}
	}
	log.Printf("got puzzles: %+v\n", pds)

	return pds, nil
}
