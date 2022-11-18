package main

import (
	"fmt"
	"time"
)

func Hello(i int) {
	defer wg.Done()
	fmt.Println("hello goroutine", i)
}

// PrintChan 打印通道
func PrintChan() bool {
	var ch chan int
	fmt.Printf("ch = %v\n", ch)
	return true
}

// Rect 接受通道值
func Rect(c chan int) {
	for {
		ret := <-c
		fmt.Println("接收成功", ret)
	}
}

// ChanClose 关闭通道
func ChanClose() bool {
	c := make(chan int, 2)
	go func() {
		for i := 0; i < 5; i++ {
			c <- i
		}
		close(c)
	}()

	for {
		if data, ok := <-c; ok {
			fmt.Printf("data = %d\n", data)
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	//for i := range c {
	//	fmt.Println(i)
	//}

	fmt.Println("chanClose end")
	return true
}
