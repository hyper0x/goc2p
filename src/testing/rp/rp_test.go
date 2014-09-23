package rp

import (
	"fmt"
	"testing"
)

func BenchmarkParallel(b *testing.B) {
	b.SetParallelism(3)
	fmt.Printf("\nN: %d\n", b.N)
	var index int
	b.RunParallel(func(pb *testing.PB) {
		i := index
		index++
		var count int
		for pb.Next() {
			count++
		}
		fmt.Printf("count[%d]: %d \n", i, count)
	})
}
