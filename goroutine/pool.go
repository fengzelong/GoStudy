package main

import (
	"fmt"
	"math/rand"
)

// CreatePool 创建工作池
// 参数1：开几个协程
func CreatePool(num int, jobChan chan *Job, resultChan chan *Result) {
	// 根据开协程个数，去跑运行
	for i := 0; i < num; i++ {
		go func(jobChan chan *Job, resultChan chan *Result) {
			// 执行运算
			// 遍历job管道所有数据，进行相加
			for job := range jobChan {
				// 随机数接过来
				rNum := job.RandNum
				// 随机数每一位相加
				// 定义返回值
				var sum int
				for rNum != 0 {
					tmp := rNum % 10
					sum += tmp
					rNum /= 10
				}
				// 想要的结果是Result
				r := &Result{
					job: job,
					sum: sum,
				}
				//运算结果扔到管道
				resultChan <- r
			}
		}(jobChan, resultChan)
	}
}

// CallPool 调用工作池
func CallPool() {
	// 1.job管道
	jobChan := make(chan *Job, 128)
	// 2.结果管道
	resultChan := make(chan *Result, 128)
	// 3.创建工作池
	CreatePool(64, jobChan, resultChan)
	// 4.开个打印的协程
	go func(resultChan chan *Result) {
		// 遍历结果管道打印
		for result := range resultChan {
			fmt.Printf("job id:%v randnum:%v result:%d\n", result.job.Id,
				result.job.RandNum, result.sum)
			//time.Sleep(100 * time.Nanosecond)
		}
	}(resultChan)
	var id int
	// 循环创建job，输入到管道
	for {
		id++
		// 生成随机数
		rNum := rand.Int()
		job := &Job{
			Id:      id,
			RandNum: rNum,
		}
		jobChan <- job

		if id > 10000 {
			fmt.Println("sum end")
			break
		}
	}
}
