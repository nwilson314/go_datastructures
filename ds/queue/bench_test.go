package queue

import "testing"

func BenchmarkEnqueue(b *testing.B) {
	q := New[int]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}

func BenchmarkEnqueueDequeueBalanced(b *testing.B) {
	q := New[int]()
	// Keep size roughly constant while exercising both ops.
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
		_ = q.Dequeue()
	}
}

func BenchmarkBatchEnqueueThenDequeue(b *testing.B) {
	const batch = 1024
	q := New[int]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < batch; j++ {
			q.Enqueue(j)
		}
		for j := 0; j < batch; j++ {
			_ = q.Dequeue()
		}
	}
}
