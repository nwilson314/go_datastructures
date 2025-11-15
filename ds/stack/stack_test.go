package stack

import "testing"

// mustPanic asserts that the provided function panics.
func mustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic, got none")
		}
	}()
	f()
}

func TestPushPopLIFO(t *testing.T) {
	s := New[int]()
	for i := 1; i <= 1000; i++ {
		s.Push(i)
	}
	for i := 1000; i >= 1; i-- {
		v := s.Pop()
		if v != i {
			t.Fatalf("expected %d, got %d", i, v)
		}
	}
	if !s.Empty() || s.Len() != 0 {
		t.Fatalf("expected empty stack after pops; len=%d", s.Len())
	}
}

func TestPeekDoesNotRemove(t *testing.T) {
	s := New[int]()
	s.Push(42)
	if got := s.Peek(); got != 42 {
		t.Fatalf("expected 42 from Peek, got %d", got)
	}
	if s.Len() != 1 {
		t.Fatalf("expected Len 1 after Peek, got %d", s.Len())
	}
}

func TestUnderflowPanics(t *testing.T) {
	s := New[int]()
	mustPanic(t, func() { _ = s.Pop() })
	mustPanic(t, func() { _ = s.Peek() })
}

func TestMixedOps(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if v := s.Pop(); v != 3 {
		t.Fatalf("expected 3, got %d", v)
	}
	s.Push(4)
	if v := s.Peek(); v != 4 {
		t.Fatalf("expected top 4, got %d", v)
	}
	if s.Len() != 3 {
		t.Fatalf("expected Len 3, got %d", s.Len())
	}
	if v := s.Pop(); v != 4 {
		t.Fatalf("expected 4, got %d", v)
	}
	if v := s.Pop(); v != 2 {
		t.Fatalf("expected 2, got %d", v)
	}
	if v := s.Pop(); v != 1 {
		t.Fatalf("expected 1, got %d", v)
	}
	if !s.Empty() {
		t.Fatalf("expected empty at end")
	}
}
