package main

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

type person struct {
	name string
	city string
	age  int8
}

// Dream name:名称
func (p *person) Dream(name string) {
	p.name = name
	//fmt.Printf("%s想学好Go语言", p.name)
}

//NewPerson 结构体构造函数
func NewPerson(name string, city string, age int8) *person {
	return &person{
		name: name,
		city: city,
		age:  age,
	}
}

// add 可变参数
func add(a int, b int, args ...int) {
	fmt.Printf("args = %v;args[0] = %d", args, args[0])
}

// itfFunc interface应用
func itfFunc(a int, b int, args ...interface{}) {
	fmt.Printf("args = %v;type = %T;args[0] = %s", args, args, args[0])
}

//aliasFunc 返回参数别名
func aliasFunc(a int, b int) (err error) {
	fmt.Printf("err = %v\n", err)
	err = errors.New("happen ex")
	fmt.Printf("err = %v\n", err)
	defer func() {
		if err != nil {
			time.Sleep(3 * time.Second)
			err = errors.New("happen ex change")
		} else {
			fmt.Println("no error")
		}
	}()
	return
}

// anonymousFunc 匿名函数
func anonymousFunc(idx int) bool {
	fns := []func(x int) int{
		func(x int) int { return x + 1 },
		func(x int) int { return x + 2 },
	}
	if idx+1 > len(fns) {
		println("下标越界")
		return false
	}
	println(fns[idx](100))
	return true
}

// fibonaci 斐波那契数列
func Fibonaci(i int) int {
	if i == 0 {
		return 0
	}
	if i == 1 {
		return 1
	}
	return Fibonaci(i-1) + Fibonaci(i-2)
}

func main() {
	//array()

	//add(1, 2, 3, 4, 5)

	//itfFunc(1, 2, "3", "4", "5")

	//err := aliasFunc(1, 2)
	//fmt.Printf("err = %v\n", err)

	//anonymousFunc(1)

	//for i := 0; i < 10; i++ {
	//	fmt.Printf("%d\n", fibonaci(i))
	//}

	//var ch = make(chan int, 10)
	//ch <- 1
	//fmt.Printf("ch len %d\n", len(ch))
	//var i1 int
	//select {
	//case i1 = <-ch:
	//	fmt.Printf("i1 = %d\n", i1)
	//	fmt.Printf("ch len %d\n", len(ch))
	//}

	//newP := NewPerson("玥玥", "武汉", 3)
	//newP.Dream("yueyue")
	//fmt.Printf("%s想学好Go语言", newP.name)

	//var p1 person
	//p1.name = "张三"
	//p1.city = "武汉"
	//p1.age = 18
	//fmt.Printf("p1 = %v\n", p1)
	//fmt.Printf("%s住在%s，今年%d岁了\n", p1.name, p1.city, p1.age)

	// 匿名结构体
	//var p2 struct {
	//	Name string
	//	Age  int
	//}
	//p2.Name = "王伟"
	//p2.Age = 19
	//fmt.Printf("%s今年%d岁了\n", p2.Name, p2.Age)

	//var p3 = new(person)
	//p3.name = "赵六"
	////fmt.Printf("p3 = %v\n", p3)
	//fmt.Printf("%s是哥哥啊\n", p3.name)
	//
	//p4 := &person{
	//	name: "111",
	//	city: "222",
	//	age:  30,
	//}
	////p4.name = "liqi"
	////p4.city = "shang hai"
	////p4.age = 25
	//fmt.Printf("%s住在%s，今年%d岁了\n", p4.name, p4.city, p4.age)

}

//array 切片练习
func array() {
	arraya := [3]int{1, 2, 3}
	arrayb := arraya[1:]

	fmt.Printf("arraya:%p，%v\n", &arraya, arraya)
	fmt.Printf("arrayb:%p，%v\n", &arrayb, arrayb)
}

//newMap map练习
func newMap() bool {
	map1 := make(map[int]string, 5)
	map1[0] = "a"
	map1[1] = "b"
	map1[2] = "c"
	map1[3] = "d"
	map1[4] = "e"
	fmt.Printf("map1 = %v,len = %d\n", map1, len(map1))
	fmt.Printf("first item = %s\n", map1[3])
	map1[5] = "f"
	fmt.Printf("map1 = %v,len = %d\n", map1, len(map1))
	map1[3] = "dddd"
	fmt.Printf("first item = %s\n", map1[3])

	val, _ := map1[3]
	fmt.Printf("contains item = %s\n", val)

	var keys = make([]int, 0, 5)
	for index := range map1 {
		keys = append(keys, index)
	}
	sort.Ints(keys)
	for index, key := range keys {
		fmt.Printf("item%d = %s\n", index, map1[key])
	}
	return true
}
