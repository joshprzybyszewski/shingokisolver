package puzzle

type puzzleCache struct {
	attempted map[string]struct{}
}

func newPuzzleCache() puzzleCache {
	return puzzleCache{
		attempted: map[string]struct{}{},
	}
}

func (pc puzzleCache) contains(p *puzzle) bool {
	// TODO don't use a String on the puzzle to determine this...
	_, ok := pc.attempted[p.String()]
	return ok
}

func (pc puzzleCache) add(p *puzzle) {
	pc.attempted[p.String()] = struct{}{}
}
