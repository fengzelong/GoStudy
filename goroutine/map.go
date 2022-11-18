package main

import (
	"fmt"
	"strconv"
	"sync"
)

var m = sync.Map{}

//func get(key string) int {
//	return m[key]
//}
//
//func set(key string, value int) {
//	m[key] = value
//}

func CallMapFunc() {
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			key := strconv.Itoa(n)
			//set(key, n)
			//fmt.Printf("k=:%v,v:=%v\n", key, get(key))
			m.Store(key, n)
			value, _ := m.Load(key)
			fmt.Printf("k=:%v,v:=%v\n", key, value)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
