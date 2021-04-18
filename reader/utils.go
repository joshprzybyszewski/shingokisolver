package reader

import (
	"log"
)

func GetAllPuzzles() []PuzzleDef {
	allPuzzles, err := CachedWebsitePuzzles()
	if err != nil {
		log.Printf("CachedWebsitePuzzles err: %+v\n", err)
		return nil
	}

	return append(allPuzzles, DefaultPuzzles()...)
}

func GetPuzzleWithSize(
	all []PuzzleDef,
	numEdges int,
) []PuzzleDef {
	filtered := make([]PuzzleDef, 0, len(all))
	for _, pd := range all {
		if pd.NumEdges == numEdges {
			filtered = append(filtered, pd)
		}
	}
	return filtered
}
