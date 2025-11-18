package stack

type Stack[T any] struct {
	data []T
}

func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}

func (s *Stack[T]) Empty() bool {
	return s.Len() == 0
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) Peek() T {
	if s.Len() == 0 {
		panic("Stack is empty!")
	}

	return s.data[s.Len()-1]
}

func (s *Stack[T]) Pop() T {
	if s.Len() == 0 {
		panic("Stack is empty!")
	}

	removed := s.data[s.Len()-1]
	s.data = s.data[:s.Len()-1]
	return removed
}
