package main

import "testing"

func TestBinSearch(t *testing.T) {
	arr := []int{1, 3, 5, 7, 9}

	tests := []struct {
		name string
		find int
		want int
	}{
		{name: "查找第一个元素", find: 1, want: 0},
		{name: "查找中间元素", find: 5, want: 2},
		{name: "查找最后一个元素", find: 9, want: 4},
		{name: "目标不存在", find: 4, want: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BinSearch(arr, tt.find); got != tt.want {
				t.Fatalf("BinSearch() = %d，期望 %d", got, tt.want)
			}
		})
	}
}
