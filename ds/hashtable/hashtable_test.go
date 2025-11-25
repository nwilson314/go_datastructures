package hashtable

import (
	"testing"
)

func TestHashTable_Basic(t *testing.T) {
	ht := New[string, int](10)

	if ht.Len() != 0 {
		t.Errorf("expected len 0, got %d", ht.Len())
	}

	ht.Put("foo", 42)
	if val, found := ht.Get("foo"); !found || val != 42 {
		t.Errorf("expected 42, got %d (found=%v)", val, found)
	}

	ht.Put("bar", 100)
	if val, found := ht.Get("bar"); !found || val != 100 {
		t.Errorf("expected 100, got %d (found=%v)", val, found)
	}

	if ht.Len() != 2 {
		t.Errorf("expected len 2, got %d", ht.Len())
	}

	// Test missing key
	if _, found := ht.Get("baz"); found {
		t.Errorf("expected 'baz' to be missing")
	}
}

func TestHashTable_Update(t *testing.T) {
	ht := New[string, int](10)

	ht.Put("foo", 1)
	ht.Put("foo", 2) // Update

	if val, found := ht.Get("foo"); !found || val != 2 {
		t.Errorf("expected 2, got %d", val)
	}

	if ht.Len() != 1 {
		t.Errorf("expected len 1, got %d", ht.Len())
	}
}

func TestHashTable_Collisions(t *testing.T) {
	// Capacity 5, insert 4 items.
	// With resize threshold 0.7, this WILL trigger a resize when inserting the 4th item.
	// 4/5 = 0.8.
	ht := New[string, int](5)

	keys := []string{"a", "b", "c", "d"}
	for i, k := range keys {
		ht.Put(k, i)
	}

	// Verify all keys are present
	for i, k := range keys {
		if val, found := ht.Get(k); !found || val != i {
			t.Errorf("key %s: expected %d, got %d (found=%v)", k, i, val, found)
		}
	}
}

func TestHashTable_Resize(t *testing.T) {
	// Start small
	ht := New[int, int](4)
	// Threshold is 0.7 * 4 = 2.8. So 3rd or 4th item triggers resize.

	// Insert 10 items
	for i := 0; i < 10; i++ {
		ht.Put(i, i*10)
	}

	// Verify size
	if ht.Len() != 10 {
		t.Errorf("expected size 10, got %d", ht.Len())
	}

	// Verify capacity grew
	if ht.cap <= 4 {
		t.Errorf("expected capacity to grow > 4, got %d", ht.cap)
	}

	// Verify all data is still there (and correctly rehashed)
	for i := 0; i < 10; i++ {
		if val, found := ht.Get(i); !found || val != i*10 {
			t.Errorf("key %d: expected %d, got %d", i, i*10, val)
		}
	}
}

func TestHashTable_IntKeys(t *testing.T) {
	ht := New[int, string](10)
	ht.Put(1, "one")
	ht.Put(2, "two")

	if v, ok := ht.Get(1); !ok || v != "one" {
		t.Errorf("expected 'one', got %v", v)
	}
	if v, ok := ht.Get(2); !ok || v != "two" {
		t.Errorf("expected 'two', got %v", v)
	}
}

func TestHashTable_Delete(t *testing.T) {
	ht := New[string, int](10)
	ht.Put("foo", 1)
	ht.Put("bar", 2)

	if !ht.Delete("foo") {
		t.Errorf("expected Delete('foo') to return true")
	}

	if ht.Len() != 1 {
		t.Errorf("expected len 1, got %d", ht.Len())
	}

	if _, found := ht.Get("foo"); found {
		t.Errorf("expected 'foo' to be gone")
	}

	if !ht.Delete("bar") {
		t.Errorf("expected Delete('bar') to return true")
	}

	if ht.Len() != 0 {
		t.Errorf("expected len 0, got %d", ht.Len())
	}
}

func TestHashTable_TombstoneReuse(t *testing.T) {
	// Cap 3.
	ht := New[int, int](3)

	// 1. Fill table
	ht.Put(1, 1)
	ht.Put(2, 2)
	ht.Put(3, 3)
	// Note: Put(3) might trigger resize depending on threshold logic.
	// If (3/3) > 0.7, it resizes.
	// Let's ensure we rely on Delete to create holes.

	ht.Delete(2)

	// 4. Insert new item. It should reuse slot or work anyway.
	ht.Put(4, 4)

	if ht.Len() != 3 {
		t.Errorf("expected len 3, got %d", ht.Len())
	}

	// Verify 4 is there
	if val, found := ht.Get(4); !found || val != 4 {
		t.Errorf("expected 4, got %d", val)
	}
}
