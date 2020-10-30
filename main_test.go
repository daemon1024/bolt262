package main

import (
	"testing"
)

func BenchmarkWalkDir(b *testing.B) {
	walkDir("../test262/test/annexB/")
}
