package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("*****************test start************")
	retCode := m.Run()
	fmt.Println("*****************test end**************")
	os.Exit(retCode)
}

func benchmarkFibonaci(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		Fibonaci(n)
	}
}

func BenchmarkFib10(b *testing.B) {
	benchmarkFibonaci(b, 10)
}
func BenchmarkFib20(b *testing.B) {
	benchmarkFibonaci(b, 20)
}
func BenchmarkFib40(b *testing.B) {
	benchmarkFibonaci(b, 40)
}
