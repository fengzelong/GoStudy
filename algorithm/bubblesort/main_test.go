package main

import (
	"reflect"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	arr := []int{3, 1, 2}
	got, max := BubbleSort(arr)
	want := []int{1, 2, 3}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BubbleSort() = %v，期望 %v", got, want)
	}
	if max != 3 {
		t.Fatalf("最大值 = %d，期望 3", max)
	}
}

func TestBubbleSortEmpty(t *testing.T) {
	got, max := BubbleSort(nil)
	if len(got) != 0 {
		t.Fatalf("空切片排序结果 = %v，期望空切片", got)
	}
	if max != 0 {
		t.Fatalf("空切片最大值 = %d，期望 0", max)
	}
}
