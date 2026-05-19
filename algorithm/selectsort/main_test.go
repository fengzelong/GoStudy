package main

import (
	"reflect"
	"testing"
)

func TestSelectSort(t *testing.T) {
	arr := []int{5, 3, 1, 4, 2}
	got := SelectSort(arr)
	want := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SelectSort() = %v，期望 %v", got, want)
	}
}
