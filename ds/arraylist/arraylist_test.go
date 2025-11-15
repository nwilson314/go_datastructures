package arraylist

import "testing"

func TestSmoke_NewAndLen(t *testing.T) {
	a := New[int]()
	if a == nil {
		t.Fatal("New returned nil")
	}
	if a.Len() != 0 {
		t.Fatalf("expected Len 0, got %d", a.Len())
	}
}

func TestAppendSetGet(t *testing.T) {
	a := New[int]()

	for i := 0; i < 3; i++ {
		a.Append(i)
		if a.Len() != i+1 {
			t.Fatalf("expected Len %d, got %d", i+1, a.Len())
		}

		a.Set(i, i-1)
		if a.Get(i) != i-1 {
			t.Fatalf("expected %d, got %d", i-1, a.Get(i))
		}
	}
}

func TestInsertAt(t *testing.T) {
	a := New[int]()

	for i := 0; i < 3; i++ {
		a.Append(i)
	}

	a.InsertAt(0, 5)

	if a.Len() != 4 {
		t.Fatalf("expected Len 4, got %d", a.Len())
	}

	if a.Get(0) != 5 {
		t.Fatalf("expected %d, got %d", 5, a.Get(0))
	}

	a.InsertAt(3, 5)

	if a.Len() != 5 {
		t.Fatalf("expected Len 5, got %d", a.Len())
	}

	if a.Get(3) != 5 {
		t.Fatalf("expected %d, got %d", 5, a.Get(3))
	}

	a.InsertAt(5, 5)

	if a.Len() != 6 {
		t.Fatalf("expected Len 6, got %d", a.Len())
	}

	if a.Get(5) != 5 {
		t.Fatalf("expected %d, got %d", 5, a.Get(5))
	}
}

func TestRemoveAt(t *testing.T) {
	a := New[int]()

	for i := 0; i < 3; i++ {
		a.Append(i)
	}

	a.RemoveAt(0)

	if a.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", a.Len())
	}

	if a.Get(0) != 1 {
		t.Fatalf("expected %d, got %d", 1, a.Get(0))
	}

	a.RemoveAt(1)

	if a.Len() != 1 {
		t.Fatalf("expected Len 1, got %d", a.Len())
	}
}

func TestClear(t *testing.T) {
	a := New[int]()
	for i := 0; i < 3; i++ {
		a.Append(i)
	}
	oldCap := a.Cap()

	a.Clear()

	if a.Len() != 0 {
		t.Fatalf("expected Len 0, got %d", a.Len())
	}
	if a.Cap() != oldCap {
		t.Fatalf("expected Cap %d, got %d", oldCap, a.Cap())
	}
	a.Append(42)
	if a.Len() != 1 || a.Get(0) != 42 {
		t.Fatalf("append after Clear failed")
	}
}

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

func TestBoundsPanics(t *testing.T) {
	a := New[int]()
	// Get/Set on empty and negative
	mustPanic(t, func() { _ = a.Get(0) })
	mustPanic(t, func() { _ = a.Get(-1) })
	mustPanic(t, func() { a.Set(0, 1) })
	mustPanic(t, func() { a.Set(-1, 1) })

	// InsertAt out of range (len = 0 allows only i==0)
	mustPanic(t, func() { a.InsertAt(-1, 1) })
	mustPanic(t, func() { a.InsertAt(1, 1) })

	// RemoveAt out of range
	mustPanic(t, func() { _ = a.RemoveAt(0) })
	mustPanic(t, func() { _ = a.RemoveAt(-1) })

	// After some appends, verify upper-bound panics
	for i := 0; i < 5; i++ {
		a.Append(i)
	}
	mustPanic(t, func() { _ = a.Get(a.Len()) })
	mustPanic(t, func() { a.Set(a.Len(), 0) })
	mustPanic(t, func() { a.InsertAt(a.Len()+1, 123) })
	mustPanic(t, func() { _ = a.RemoveAt(a.Len()) })
}

func TestRemoveAt_Tail(t *testing.T) {
	a := New[int]()
	a.Append(10)
	a.Append(20)
	a.Append(30)

	removed := a.RemoveAt(2)
	if removed != 30 {
		t.Fatalf("expected removed 30, got %d", removed)
	}
	if a.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", a.Len())
	}
	if a.Get(1) != 20 {
		t.Fatalf("expected tail 20, got %d", a.Get(1))
	}
}

func TestInsertRemove_RoundTrip(t *testing.T) {
	a := New[int]()
	a.Append(0)
	a.Append(1)
	a.Append(2)

	a.InsertAt(1, 99) // [0,99,1,2]
	if got := a.Get(1); got != 99 {
		t.Fatalf("expected 99 at index 1, got %d", got)
	}
	removed := a.RemoveAt(1) // back to [0,1,2]
	if removed != 99 {
		t.Fatalf("expected removed 99, got %d", removed)
	}
	if a.Len() != 3 || a.Get(0) != 0 || a.Get(1) != 1 || a.Get(2) != 2 {
		t.Fatalf("list not restored to [0,1,2]; len=%d, vals=%d,%d,%d", a.Len(), a.Get(0), a.Get(1), a.Get(2))
	}
}

func TestCapacityInvariants(t *testing.T) {
	a := New[int]()
	prevCap := a.Cap()
	for i := 0; i < 4096; i++ {
		a.Append(i)
		if a.Cap() < prevCap {
			t.Fatalf("capacity decreased from %d to %d", prevCap, a.Cap())
		}
		if a.Cap() < a.Len() {
			t.Fatalf("capacity %d less than length %d", a.Cap(), a.Len())
		}
		prevCap = a.Cap()
	}
}
