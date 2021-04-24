package reader

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func CachedWebsitePuzzles() ([]model.Definition, error) {
	file, err := ioutil.ReadFile(websiteCachePuzzlesFilename)
	if err != nil {
		return nil, err
	}
	pds := []model.Definition{}
	includedPuzzleIDs := map[string]struct{}{}

	lines := strings.Split(string(file), "\n")
	for _, l := range lines {
		if l == `` {
			continue
		}
		parts := strings.Split(l, ":")
		numEdges, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		desc := parts[1]
		task := parts[2]

		descParts := strings.Split(desc, "_")
		puzzID := descParts[0]
		diff := model.NewDifficulty(descParts[1])

		if _, ok := includedPuzzleIDs[puzzID]; ok {
			continue
		}

		pd, err := fromWebsiteTask(numEdges, diff, puzzID, task)
		if err != nil {
			return nil, err
		}
		pds = append(pds, pd)
		includedPuzzleIDs[puzzID] = struct{}{}
	}

	return pds, nil
}
