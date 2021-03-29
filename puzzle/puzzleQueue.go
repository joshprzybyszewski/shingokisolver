package puzzle

type puzzleQueue struct {
	items []*puzzle
}

func newPuzzleQueue() *puzzleQueue {
	return &puzzleQueue{
		items: make([]*puzzle, 0),
	}
}

func (q *puzzleQueue) isEmpty() bool {
	return len(q.items) == 0
}

func (q *puzzleQueue) pop() *puzzle {
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *puzzleQueue) push(items ...*puzzle) {
	q.items = append(q.items, items...)
}
