package main

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestWriteFuncWithRWMutex(t *testing.T) {
	x = 0
	wg = sync.WaitGroup{}
	wg.Add(1)
	go WriteFunc()
	wg.Wait()

	if x != 5000 {
		t.Fatalf("x = %d，期望 5000", x)
	}
}

func TestAtomicAddFunc(t *testing.T) {
	x = 0
	wg = sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go atomicAddFunc()
	}
	wg.Wait()

	if x != 100 {
		t.Fatalf("x = %d，期望 100", x)
	}
}

func TestCreatePool(t *testing.T) {
	jobChan := make(chan *Job, 2)
	resultChan := make(chan *Result, 2)
	CreatePool(2, jobChan, resultChan)

	jobChan <- &Job{Id: 1, RandNum: 123}
	jobChan <- &Job{Id: 2, RandNum: 405}
	close(jobChan)

	got := map[int]int{}
	deadline := time.After(time.Second)
	for len(got) < 2 {
		select {
		case result := <-resultChan:
			got[result.job.Id] = result.sum
		case <-deadline:
			t.Fatal("等待工作池结果超时")
		}
	}

	if got[1] != 6 || got[2] != 9 {
		t.Fatalf("工作池结果 = %v，期望 map[1:6 2:9]", got)
	}
}

func TestSelectFuncDoesNotBlock(t *testing.T) {
	done := make(chan struct{})
	go func() {
		SelectFunc()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("SelectFunc 执行超时")
	}
}

func TestCallMapFuncStoresValues(t *testing.T) {
	m = sync.Map{}
	CallMapFunc()

	for i := 0; i < 20; i++ {
		if _, ok := m.Load(strconv.Itoa(i)); !ok {
			t.Fatalf("缺少 key %d", i)
		}
	}
}
