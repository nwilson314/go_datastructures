package arraylist

import "fmt"

type ArrayList[T any] struct {
	data []T
}

func New[T any]() *ArrayList[T] {
	return &ArrayList[T]{}
}

func (a *ArrayList[T]) Len() int {
	return len(a.data)
}

func (a *ArrayList[T]) Cap() int {
	return cap(a.data)
}

func (a *ArrayList[T]) Append(v T) {
	a.data = append(a.data, v)
}

func (a *ArrayList[T]) Get(i int) T {
	if i < 0 || i >= len(a.data) {
		panic(fmt.Sprintf("index out of range: %d (len=%d)", i, len(a.data)))
	}
	return a.data[i]
}

func (a *ArrayList[T]) Set(i int, v T) {
	if i < 0 || i >= len(a.data) {
		panic(fmt.Sprintf("index out of range: %d (len=%d)", i, len(a.data)))
	}
	a.data[i] = v
}

func (a *ArrayList[T]) InsertAt(i int, v T) {
	if i < 0 || i > len(a.data) {
		panic(fmt.Sprintf("index out of range: %d (len=%d)", i, len(a.data)))
	}

	a.ensureCapacity(len(a.data) + 1)
	a.data = a.data[:len(a.data)+1]
	copy(a.data[i+1:], a.data[i:])
	a.data[i] = v
}

func (a *ArrayList[T]) RemoveAt(i int) T {
	if i < 0 || i >= len(a.data) {
		panic(fmt.Sprintf("index out of range: %d (len=%d)", i, len(a.data)))
	}

	returned := a.data[i]
	last := len(a.data) - 1
	var zero T
	a.data[last] = zero
	a.data = append(a.data[:i], a.data[i+1:]...)
	return returned
}

func (a *ArrayList[T]) ensureCapacity(needed int) {
	if needed <= cap(a.data) {
		return
	}

	newCap := cap(a.data)
	if newCap == 0 {
		newCap = 1
	}

	for newCap < needed {
		newCap *= 2
	}

	newData := make([]T, len(a.data), newCap)
	copy(newData, a.data)
	a.data = newData
}
