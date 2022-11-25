package main

import "fmt"

func main() {
	arr := []int{1, 9, 10, 30, 2, 5, 45, 8, 63, 234, 12}
	selectSort := SelectSort(arr)
	fmt.Println(selectSort)
}

// SelectSort 选择排序
func SelectSort(arr []int) []int {
	length := len(arr)
	if length <= 1 {
		return arr
	}

	for i := 1; i < length; i++ {
		min := i
		for j := i + 1; j < length; j++ {
			if arr[min] > arr[j] {
				min = j
			}
		}
		if i != min {
			arr[i], arr[min] = arr[min], arr[i]
		}
	}
	return arr
}
