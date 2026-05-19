package main

import (
	"reflect"
	"testing"
)

func TestHeapSort(t *testing.T) {
	arr := []int{5, 1, 4, 2, 3}
	got := HeapSort(arr)
	want := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("HeapSort() = %v，期望 %v", got, want)
	}
}
