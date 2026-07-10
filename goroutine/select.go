package main

import (
	"fmt"
)

func SelectFunc() {
	intChan := make(chan int, 1)
	stringChan := make(chan string, 1)

	stringChan <- "hello"
	intChan <- 1

	for received := 0; received < 2; {
		select {
		case value := <-intChan:
			fmt.Println("int=", value)
			received++
		case value := <-stringChan:
			fmt.Println("string=", value)
			received++
		default:
			return
		}
	}

	fmt.Println("main结束")
}
