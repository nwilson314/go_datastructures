package queue

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

func TestFIFO(t *testing.T) {
	q := New[int]()
	for i := 1; i <= 1000; i++ {
		q.Enqueue(i)
	}
	for i := 1; i <= 1000; i++ {
		v := q.Dequeue()
		if v != i {
			t.Fatalf("expected %d, got %d", i, v)
		}
	}
	if !q.Empty() || q.Len() != 0 {
		t.Fatalf("expected empty queue; len=%d", q.Len())
	}
}

func TestPeekDoesNotRemove(t *testing.T) {
	q := New[int]()
	q.Enqueue(42)
	if got := q.Peek(); got != 42 {
		t.Fatalf("expected 42 from Peek, got %d", got)
	}
	if q.Len() != 1 {
		t.Fatalf("expected Len 1 after Peek, got %d", q.Len())
	}
}

func TestUnderflowPanics(t *testing.T) {
	q := New[int]()
	mustPanic(t, func() { _ = q.Dequeue() })
	mustPanic(t, func() { _ = q.Peek() })
}

func TestWrapAroundMixed(t *testing.T) {
	q := New[int]()
	next := 0
	// Build up some initial contents
	for i := 0; i < 64; i++ {
		q.Enqueue(i)
	}
	// Interleave enqueues/dequeues to force internal index wrap if ring buffer is used.
	for i := 64; i < 5000; i++ {
		q.Enqueue(i)
		if q.Len() > 50 {
			v := q.Dequeue()
			if v != next {
				t.Fatalf("expected %d, got %d", next, v)
			}
			next++
		}
	}
	// Drain remaining and validate order.
	for !q.Empty() {
		v := q.Dequeue()
		if v != next {
			t.Fatalf("expected %d, got %d", next, v)
		}
		next++
	}
}
