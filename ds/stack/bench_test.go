package stack

import "testing"

func BenchmarkPush(b *testing.B) {
	s := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
}

func BenchmarkPushThenPop(b *testing.B) {
	s := New[int]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
	for i := 0; i < b.N; i++ {
		_ = s.Pop()
	}
}
