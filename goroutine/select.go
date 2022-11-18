package main

import (
	"fmt"
)

func SelectFunc() {
	intChan := make(chan int, 1)
	stringChan := make(chan string, 1)

	go func() {
		stringChan <- "hello"
	}()

	go func() {
		intChan <- 1
	}()

	for {
		select {
		case value := <-intChan:
			fmt.Println("int=", value)
		case value := <-stringChan:
			fmt.Println("string=", value)
		default:
			break
		}
	}

	fmt.Println("main结束")

}
