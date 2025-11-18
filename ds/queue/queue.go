package queue

type Queue[T any] struct {
	data  []T
	front int
	back  int
	size  int
}

func New[T any]() *Queue[T] {
	data := make([]T, 10)
	return &Queue[T]{
		data:  data,
		front: 0,
		back:  0,
		size:  10,
	}
}

func (q *Queue[T]) Len() int {
	if q.front <= q.back {
		return q.back - q.front
	}
	return q.size - q.front + q.back
}

func (q *Queue[T]) Empty() bool {
	return q.Len() == 0
}

func (q *Queue[T]) Peek() T {
	if q.Len() == 0 {
		panic("Queue is empty!")
	}

	return q.data[q.front]
}

func (q *Queue[T]) Enqueue(v T) {
	if q.Len() == q.size-1 {
		resizeQueue(q)
	}

	q.data[q.back] = v
	q.back = (q.back + 1) % q.size
}

func (q *Queue[T]) Dequeue() T {
	if q.Len() == 0 {
		panic("Queue is empty!")
	}

	returnV := q.data[q.front]
	q.front = (q.front + 1) % q.size
	return returnV
}

func resizeQueue[T any](q *Queue[T]) {
	newData := make([]T, q.size*2)
	copy(newData, q.data[q.front:])
	copy(newData[q.size-q.front:], q.data[:q.front])
	q.data = newData
	q.front = 0
	q.back = q.size - 1
	q.size *= 2
}
