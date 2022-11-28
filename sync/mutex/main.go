package main

import (
	"fmt"
	"sync"
	"time"
)

var books = []string{"东邪", "西毒", "南帝", "北丐"}

var blackBoard string
var page int

var rwMutex = sync.RWMutex{}

type Student struct {
	Name string
	rTxt string
}

type Teacher struct {
	Name string
	wTxt string
}

func (s *Student) read() {
	fmt.Println(s.Name + " request read ...")
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	fmt.Println(s.Name + " start read")
	s.rTxt = blackBoard
	time.Sleep(time.Second)
	fmt.Println(s.Name + " end read " + s.rTxt)
}

func (w *Teacher) write(idx int) {
	fmt.Println(w.Name + " request write")
	rwMutex.Lock()

	if idx == page {
		fmt.Println(w.Name + " start write ...")
		w.wTxt = books[page]
		time.Sleep(time.Second * 3)
		blackBoard = w.wTxt
		page++
		fmt.Println(w.Name+" finish. blackboard: ", blackBoard)
		rwMutex.Unlock()
		time.Sleep(time.Second * 3)
	} else {
		fmt.Println(w.Name + " end ....")
		rwMutex.Unlock()
	}
}

func main() {
	stuA := &Student{Name: "stuA"}
	stuB := &Student{Name: "stuB"}
	stuC := &Student{Name: "stuC"}
	stuD := &Student{Name: "stuD"}
	teaA := &Teacher{Name: "teaA"}
	teaB := &Teacher{Name: "teaB"}

	go StartWrite(teaA)
	go StartWrite(teaB)
	time.Sleep(time.Second * 1)
	go StartRead(stuA)
	go StartRead(stuB)
	go StartRead(stuC)
	go StartRead(stuD)

}

func StartRead(s *Student) {
	for i := 0; i < 1000; i++ {
		s.read()
	}
}

func StartWrite(t *Teacher) {
	for i := 0; i < 4; i++ {
		t.write(i)
	}
}
