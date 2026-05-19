package main

import "testing"

func TestFibonaci(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{name: "第 0 项", n: 0, want: 0},
		{name: "第 1 项", n: 1, want: 1},
		{name: "第 2 项", n: 2, want: 1},
		{name: "第 10 项", n: 10, want: 55},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Fibonaci(tt.n); got != tt.want {
				t.Fatalf("Fibonaci(%d) = %d，期望 %d", tt.n, got, tt.want)
			}
		})
	}
}

func TestAnonymousFunc(t *testing.T) {
	if !anonymousFunc(0) {
		t.Fatal("合法下标应该返回 true")
	}
	if anonymousFunc(2) {
		t.Fatal("越界下标应该返回 false")
	}
}

func TestNewMap(t *testing.T) {
	if !newMap() {
		t.Fatal("newMap 应该返回 true")
	}
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
