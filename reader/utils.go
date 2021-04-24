package reader

import (
	"log"

	"github.com/joshprzybyszewski/shingokisolver/model"
)

func GetAllPuzzles() []model.Definition {
	allPuzzles, err := CachedWebsitePuzzles()
	if err != nil {
		log.Printf("CachedWebsitePuzzles err: %+v\n", err)
		return nil
	}

	return append(allPuzzles, DefaultPuzzles()...)
}

func GetPuzzleWithSize(
	all []model.Definition,
	numEdges int,
) []model.Definition {
	filtered := make([]model.Definition, 0, len(all))
	for _, pd := range all {
		if pd.NumEdges == numEdges {
			filtered = append(filtered, pd)
		}
	}
	return filtered
}
