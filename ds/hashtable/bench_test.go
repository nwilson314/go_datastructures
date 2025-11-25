package hashtable

import (
	"fmt"
	"testing"
)

// BenchmarkInsert measures the speed of inserting new keys.
// We compare our HashTable vs Go's native map.
func BenchmarkHashTable_Insert(b *testing.B) {
	// Setup
	ht := New[string, int](b.N)

	// Pre-generate keys to avoid measuring Sprintf overhead too much
	keys := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ht.Put(keys[i], i)
	}
}

func BenchmarkNativeMap_Insert(b *testing.B) {
	m := make(map[string]int, b.N)
	keys := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[keys[i]] = i
	}
}

// BenchmarkGet measures successful lookups.
func BenchmarkHashTable_Get(b *testing.B) {
	numItems := 10000
	ht := New[string, int](numItems)
	keys := make([]string, numItems)
	for i := 0; i < numItems; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
		ht.Put(keys[i], i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Cycle through the known keys
		key := keys[i%numItems]
		ht.Get(key)
	}
}

func BenchmarkNativeMap_Get(b *testing.B) {
	numItems := 10000
	m := make(map[string]int, numItems)
	keys := make([]string, numItems)
	for i := 0; i < numItems; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
		m[keys[i]] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := keys[i%numItems]
		_ = m[key]
	}
}

