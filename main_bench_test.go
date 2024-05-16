package main

import "testing"

func Benchmark_checkContex(t *testing.B) {
	for range t.N {
		checkSwitch()
	}
}
