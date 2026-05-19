package main

import "fmt"

func main() {
	max := 5000
	arr := make([]int, max)
	for i := 0; i < max; i++ {
		arr[i] = i + 1
	}
	idx := BinSearch(arr, 450)
	if idx == -1 {
		fmt.Println("not found!")
		return
	}
	fmt.Printf("found data index: %d\n", idx)
	fmt.Printf("found data: %d\n", arr[idx])
}

// BinSearch 二分查找，返回目标值在升序切片中的下标。
func BinSearch(arr []int, findData int) int {
	low := 0
	high := len(arr) - 1
	for low <= high {
		mid := (low + high) / 2
		if arr[mid] > findData {
			high = mid - 1
		} else if arr[mid] < findData {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}
