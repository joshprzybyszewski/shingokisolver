package solvers

import "github.com/joshprzybyszewski/shingokisolver/puzzle"

type puzzleCache struct {
	attempted map[string]struct{}
}

func newPuzzleCache() puzzleCache {
	return puzzleCache{
		attempted: map[string]struct{}{},
	}
}

func (pc puzzleCache) contains(p *puzzle.Puzzle) bool {
	// TODO figure out if we really want a cache or not:#
	return false
	// TODO don't use a String on the puzzle to determine this...
	_, ok := pc.attempted[p.String()]
	return ok
}

func (pc puzzleCache) add(p *puzzle.Puzzle) {
	// TODO figure out if we really want a cache or not:#
	return
	pc.attempted[p.String()] = struct{}{}
}
