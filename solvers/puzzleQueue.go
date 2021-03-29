package solvers

import "github.com/joshprzybyszewski/shingokisolver/puzzle"

type puzzleQueue struct {
	items []*puzzle.Puzzle
}

func newPuzzleQueue() *puzzleQueue {
	return &puzzleQueue{
		items: make([]*puzzle.Puzzle, 0),
	}
}

func (q *puzzleQueue) isEmpty() bool {
	return len(q.items) == 0
}

func (q *puzzleQueue) pop() *puzzle.Puzzle {
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *puzzleQueue) push(items ...*puzzle.Puzzle) {
	q.items = append(q.items, items...)
}
