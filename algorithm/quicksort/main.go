package main

import "fmt"

func main() {
	array := []int{1, 6, 23, 76, 11, 8, 9, 26, 19}
	fmt.Printf("sort result: %v", QuickSort(array))
}

func QuickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	splitData := arr[0]       //第一个数据
	low := make([]int, 0, 0)  //比我小的数据
	high := make([]int, 0, 0) //比我大的数据
	mid := make([]int, 0, 0)  //与我一样大的数据
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
