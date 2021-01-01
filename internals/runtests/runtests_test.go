package runtests

import (
	"testing"
)

func BenchmarkRunTests(b *testing.B) {
	Dir("../../bench/multiple", "../../bench/harness/")
}

func TestRunTests(t *testing.T) {
	Dir("../../../test262/test/", "../../bench/harness/")
	t.Logf("Success !")
}
