package algorithms

import (
	"strconv"
	"testing"
)

func BenchmarkIssubset(b *testing.B) {
	sliceLen := 10_000
	b.StopTimer()
	var (
		first, second []string
	)
	for i := 0; i < sliceLen; i++ {
		first = append(first, strconv.Itoa(i))
		second = append(second, strconv.Itoa(i))
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Issubset(first, second)
	}
}
