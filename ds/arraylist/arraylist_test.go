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
