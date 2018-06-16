package bmt

import (
	"testing"
	"time"
)

func Benchmark(b *testing.B) {
	customTimerTag := false
	if customTimerTag {
		b.StopTimer()
	}
	b.SetBytes(12345678)
	time.Sleep(time.Second)
	if customTimerTag {
		b.StartTimer()
	}
}
