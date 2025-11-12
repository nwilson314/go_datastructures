package arraylist

import "testing"

func BenchmarkAppend(b *testing.B) {
	a := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// placeholder: will append later
		_ = a
	}
}
