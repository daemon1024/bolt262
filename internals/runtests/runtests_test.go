package runtests

import (
	"testing"
)

func BenchmarkRunTests(b *testing.B) {
	Dir("../../bench/multiple", "../../bench/harness/")
}
