package main

import (
	"testing"
)

func BenchmarkRunTests(b *testing.B) {
	runTests("./bench/multiple", "./bench/harness/")
}
