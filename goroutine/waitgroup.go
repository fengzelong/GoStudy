package main

import "time"

// WriteFunc 读写锁
// x值累加
func WriteFunc() {
	for i := 0; i < 5000; i++ {
		wrLock.Lock()
		x = x + 1
		wrLock.Unlock()
	}
	wg.Done()
}

// ReadFunc 读锁
func ReadFunc() {
	wrLock.RLock()
	time.Sleep(time.Millisecond)
	wrLock.RUnlock()
	wg.Done()
}
