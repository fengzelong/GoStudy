package main

import "fmt"

func main() {
	max := 5000
	arr := make([]int, max, max)
	for i := 0; i < max; i++ {
		arr[i] = i + 1
	}
	idx := BinSearch(arr, 5200)
	if idx == -1 {
		fmt.Println("not found!")
	} else {
		fmt.Printf("found data index: %d\n", idx)
		fmt.Printf("found data: %d\n", arr[idx])
	}
}

// BinSearch 二分查找
func BinSearch(arr []int, findData int) int {
	low := 0
	high := len(arr) - 1
	for low <= high {
		mid := (low + high) / 2
		if arr[mid] > findData {
			high = arr[mid-1]
		} else if arr[mid] < findData {
			low = arr[mid+1]
		} else {
			return mid
		}
	}
	return -1
}
