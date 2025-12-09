package lru

import "testing"

func TestLRU_Basic(t *testing.T) {
	lru := New[string, int](2)

	lru.Put("a", 1)
	lru.Put("b", 2)

	if v, ok := lru.Get("a"); !ok || v != 1 {
		t.Errorf("expected 1, got %v", v)
	}
	if v, ok := lru.Get("b"); !ok || v != 2 {
		t.Errorf("expected 2, got %v", v)
	}
}

func TestLRU_Eviction(t *testing.T) {
	lru := New[string, int](2)
	lru.Put("a", 1)
	lru.Put("b", 2)

	// Access 'a' to make it MRU. 'b' becomes LRU.
	lru.Get("a")

	// Add 'c'. Should evict 'b'.
	lru.Put("c", 3)

	if _, ok := lru.Get("b"); ok {
		t.Errorf("expected 'b' to be evicted")
	}
	if _, ok := lru.Get("a"); !ok {
		t.Errorf("expected 'a' to stay")
	}
	if _, ok := lru.Get("c"); !ok {
		t.Errorf("expected 'c' to exist")
	}
}

func TestLRU_Update(t *testing.T) {
	lru := New[string, int](2)
	lru.Put("a", 1)
	lru.Put("b", 2)

	// Update 'a'. Should move it to front. 'b' becomes LRU.
	lru.Put("a", 10)

	// Add 'c'. Should evict 'b'.
	lru.Put("c", 3)

	if v, ok := lru.Get("a"); !ok || v != 10 {
		t.Errorf("expected 'a' to be 10, got %v", v)
	}
	if _, ok := lru.Get("b"); ok {
		t.Errorf("expected 'b' to be evicted")
	}
}

func TestLRU_CapacityOne(t *testing.T) {
	lru := New[int, int](1)
	lru.Put(1, 1)

	if _, ok := lru.Get(1); !ok {
		t.Error("expected 1")
	}

	lru.Put(2, 2) // Evicts 1

	if _, ok := lru.Get(1); ok {
		t.Error("expected 1 to be evicted")
	}
	if _, ok := lru.Get(2); !ok {
		t.Error("expected 2")
	}
}
