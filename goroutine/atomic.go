package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

// 普通版加函数
func addFunc() {
	// x = x + 1
	x++ // 等价于上面的操作
	wg.Done()
}

// 互斥锁版加函数
func mutexAddFunc() {
	lock.Lock()
	x++
	lock.Unlock()
	wg.Done()
}

// 原子操作版加函数
func atomicAddFunc() {
	atomic.AddInt64(&x, 1)
	wg.Done()
}

func CallAtomicAdd() {
	start := time.Now()
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		//go addFunc() // 普通版add函数 不是并发安全的
		//go mutexAddFunc() // 加锁版add函数 是并发安全的，但是加锁性能开销大
		go atomicAddFunc() // 原子操作版add函数 是并发安全，性能优于加锁版
	}
	wg.Wait()
	end := time.Now()
	fmt.Println(x)
	fmt.Println(end.Sub(start))
}
