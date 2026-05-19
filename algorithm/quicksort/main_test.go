package main

import (
	"reflect"
	"testing"
)

func TestQuickSort(t *testing.T) {
	arr := []int{3, 1, 2, 3}
	got := QuickSort(arr)
	want := []int{1, 2, 3, 3}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("QuickSort() = %v，期望 %v", got, want)
	}
}
