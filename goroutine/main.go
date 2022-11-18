package main

import (
	"sync"
)

var wg sync.WaitGroup

// 互斥锁
var lock sync.Mutex

// 读写锁
var wrLock sync.RWMutex

var x int64

type Job struct {
	// id
	Id int
	// 需要计算的随机数
	RandNum int
}

type Result struct {
	// 这里必须传对象实例
	job *Job
	// 求和
	sum int
}

func main() {
	//for i := 0; i < 5; i++ {
	//	wg.Add(1)
	//	go Hello(i)
	//}
	//wg.Wait()

	//res := PrintChan()
	//fmt.Printf("res = %t", res)

	// 无缓冲通道
	//ch := make(chan int)
	//go Rect(ch) // 启用goroutine从通道接收值
	//ch <- 10
	//ch <- 20
	//ch <- 30
	//ch <- 40
	//ch <- 50
	//fmt.Println("发送成功")

	// 有缓冲通道
	//ch := make(chan int, 1) // 创建一个容量为1的有缓冲区通道
	//ch <- 10
	//fmt.Println("发送成功")

	// 通道关闭
	//ChanClose()

	// 调用工作池
	//CallPool()

	// Select用法
	//SelectFunc()

	// waitGroup & 互斥锁
	//start := time.Now()
	//for i := 0; i < 10; i++ {
	//	wg.Add(1)
	//	go WriteFunc()
	//}
	//
	//for i := 0; i < 1000; i++ {
	//	wg.Add(1)
	//	go ReadFunc()
	//}
	//
	//wg.Wait()
	//end := time.Now()
	//fmt.Println(end.Sub(start))

	// sync.Map
	// CallMapFunc()

	// atomic
	CallAtomicAdd()
}
