package main

import "fmt"

func main() {
	array := []int{1, 6, 23, 76, 11, 8, 9, 26, 19}
	fmt.Printf("sort result: %v", QuickSort(array))
}

// QuickSort 快速排序。
func QuickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	splitData := arr[0]
	low := make([]int, 0)
	high := make([]int, 0)
	mid := make([]int, 0)
	mid = append(mid, splitData)

	for i := 1; i < len(arr); i++ {
		if arr[i] < splitData {
			low = append(low, arr[i])
			fmt.Printf("low %d: %v\n", i, low)
		} else if arr[i] > splitData {
			high = append(high, arr[i])
			fmt.Printf("high %d: %v\n", i, high)
		} else {
			mid = append(mid, arr[i])
			fmt.Printf("mid %d: %v\n", i, mid)
		}
	}

	low = QuickSort(low)
	high = QuickSort(high)
	myArray := append(append(low, mid...), high...)
	return myArray
}
