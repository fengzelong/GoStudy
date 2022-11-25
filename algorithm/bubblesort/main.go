package main

import "fmt"

func main() {
	arr := []int{1, 9, 10, 30, 2, 5, 45, 8, 63, 234, 12}
	fmt.Println(BubbleSort(arr))
}

// BubbleSort 冒泡排序并返回最大值
func BubbleSort(arr []int) ([]int, int) {
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
				fmt.Printf("arr :%v\n", arr)
			}
		}
	}
	return arr, arr[len(arr)-1]
}
