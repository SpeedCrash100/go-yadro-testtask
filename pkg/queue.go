package pkg

import "errors"

var (
	ErrQueueFull  = errors.New("queue is full")
	ErrQueueEmpty = errors.New("queue is empty")
)

type Queue[T any] struct {
	slice []T
	n     int
	start int
	end   int
}

func NewQueue[T any](n int) Queue[T] {

	return Queue[T]{
		slice: make([]T, n),
		n:     n,
		start: -1,
		end:   -1,
	}
}

func (q *Queue[T]) IsEmpty() bool {
	return q.start == -1 && q.end == -1
}

func (q *Queue[T]) IsFull() bool {
	return (q.end+1)%q.n == q.start
}

func (q *Queue[T]) Push(val T) error {
	place_id := (q.end + 1) % q.n
	if q.IsFull() {
		return ErrQueueFull
	}

	if q.start == -1 {
		q.start = 0
	}

	q.slice[place_id] = val
	q.end = place_id

	return nil
}

func (q *Queue[T]) Pop() (T, error) {
	if q.IsEmpty() {
		return q.slice[0], ErrQueueEmpty
	}

	val := q.slice[q.start]

	if q.start == q.end {
		q.start, q.end = -1, -1
	} else {
		q.start = (q.start + 1) % q.n
	}

	return val, nil
}
