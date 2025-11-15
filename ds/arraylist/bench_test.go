package arraylist

import "testing"

func BenchmarkAppend(b *testing.B) {
	a := New[int]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Append(i)
	}
}

func prefill(a *ArrayList[int], n int) {
	for i := 0; i < n; i++ {
		a.Append(i)
	}
}

func BenchmarkInsertAtHead_1024(b *testing.B) {
	const size = 1024
	a := New[int]()
	prefill(a, size)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.InsertAt(0, i)
		_ = a.RemoveAt(0)
	}
}

func BenchmarkInsertAtMiddle_1024(b *testing.B) {
	const size = 1024
	a := New[int]()
	prefill(a, size)
	mid := a.Len() / 2
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.InsertAt(mid, i)
		_ = a.RemoveAt(mid)
	}
}

func BenchmarkInsertAtTail_1024(b *testing.B) {
	const size = 1024
	a := New[int]()
	prefill(a, size)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.InsertAt(a.Len(), i)
		_ = a.RemoveAt(a.Len() - 1)
	}
}

func BenchmarkRemoveAtHead_1024(b *testing.B) {
	const size = 1024
	a := New[int]()
	prefill(a, size)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.RemoveAt(0)
		a.InsertAt(0, i)
	}
}

func BenchmarkRemoveAtTail_1024(b *testing.B) {
	const size = 1024
	a := New[int]()
	prefill(a, size)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.RemoveAt(a.Len() - 1)
		a.InsertAt(a.Len(), i)
	}
}

func BenchmarkClearReuse_4096(b *testing.B) {
	const size = 4096
	a := New[int]()
	prefill(a, size)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Clear()
		// minimal reuse to ensure capacity is preserved
		a.Append(i)
	}
}
