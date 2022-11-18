package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	//fmt.Println("--------------字符串---------------")
	//str := "总炮轰金门123123123"
	//fmt.Printf("长度%d\n", len(str))
	//fmt.Println("----------------------------------")
	//var str1 string = "你好吗"
	//fmt.Printf("长度1%d", len(str1))
	//fmt.Println("----------------------------------")
	//res := strings.Compare("aaa", "aaa")
	//fmt.Printf("比对结果：%d", res)
	//fmt.Println("----------------------------------")

	//res1 := strings.Contains("aaaa", "a")
	//fmt.Printf("是否包含：%t\n", res1)
	//
	//a, b := Name("sdf", "hgh")
	//fmt.Printf("name:%s，age：%s\n", a, b)
	//
	//_, c := Name("sdf", "hgh")
	//fmt.Printf("age：%s\n", c)

	//var name string = "xiaoming"
	//var age string = "18"
	//str := fmt.Sprintf("name：%s，age：%s", name, age)
	//fmt.Printf(str)

	//fmt.Printf("--------------------数组---------------\n")
	//a1 := [3]string{"aa", "bb", "cc"}
	//fmt.Printf("a1：%v，len：%d，cap：%d\n", a1, len(a1), cap(a1))
	//
	//a2 := [3]int{1, 2, 3}
	//fmt.Printf("a2：%v，len：%d，cap：%d\n", a2, len(a2), cap(a2))
	//fmt.Printf("a2：%d\n", a2[0])
	//
	//fmt.Printf("--------------------切片---------------\n")
	//a3 := []int{1, 2, 3}
	//fmt.Printf("a3：%v，len：%d，cap：%d\n", a3, len(a3), cap(a3))
	//
	//a4 := make([]int, 3, 5)
	//a4[0] = 15
	//fmt.Printf("a4：%v，len：%d，cap：%d\n", a4, len(a4), cap(a4))
	//fmt.Printf("a2：%d\n", a4[0])
	//a4 = append(a4, 7)
	//a4 = append(a4, 5)
	//a4 = append(a4, 6) //超过切片cao，触发扩容
	//fmt.Printf("a4：%v，len：%d，cap：%d\n", a4, len(a4), cap(a4))
	//
	//a5 := make([]int, 5)
	//a5[0] = 15
	//a5[1] = 16
	//a5[2] = 17
	//fmt.Printf("a5：%v，len：%d，cap：%d\n", a5, len(a5), cap(a5))
	//
	//fmt.Printf("--------------------子切片(左开右闭/只读)---------------\n")
	//a6 := a5[2:4]
	//fmt.Printf("a6：%v，len：%d，cap：%d\n", a6, len(a6), cap(a6))

	//fmt.Printf("--------------------For---------------\n")
	//arr := []int{1, 2, 3, 4, 5}
	//index := 0
	//for {
	//	if index == 3 {
	//		break
	//	}
	//	fmt.Printf("val：%d\n", arr[index])
	//	index++
	//}

	//fmt.Printf("--------------------ForI---------------\n")
	//arr1 := []int{1, 2, 3, 4, 5}
	//for i := 0; i < len(arr1); i++ {
	//	fmt.Printf("val：%d\n", arr1[i])
	//}

	//fmt.Printf("--------------------ForRange---------------\n")
	//arr2 := []int{1, 2, 3, 4, 5}
	//for key, value := range arr2 {
	//	fmt.Printf("key：%d，value：%d\n", key, value)
	//}

	//fmt.Printf("--------------------ifelse---------------\n")
	////young(18)
	////HowFar(15, 16)
	//foodName := ChooseFood(1)
	//fmt.Printf("食物名称：%s\n", foodName)

	for routine := 0; routine < 2; routine++ {
		Wait.Add(1)
		go Routine(routine)
	}

	Wait.Wait()
	fmt.Printf("Final Counter：%d\n", Counter)
}

var Wait sync.WaitGroup
var Counter int = 0

func Routine(id int) {
	for count := 0; count < 2; count++ {
		value := Counter
		time.Sleep(1 * time.Nanosecond)
		value++
		Counter = value
	}
	Wait.Done()
}

//func Name(name string, age string) (string, string) {
//	return name, age
//}

//young
func young(age int) bool {
	if age >= 18 {
		fmt.Printf("你已经成年了\n")
		return true
	} else {
		fmt.Printf("你还未成年\n")
		return false
	}
}

//HowFar if-else用法
func HowFar(start int, end int) {
	if res := start - end; res > 0 {
		fmt.Printf("你出门了\n")
	} else {
		fmt.Printf("你还在家\n")
	}
}

//ChooseFood switch用法
func ChooseFood(foodType int) string {
	switch foodType {
	case 1:
		return "苹果"
	case 2:
		return "香蕉"
	case 3:
		return "橘子"
	case 4:
		return "葡萄"
	case 5:
		return "柚子"
	default:
		return "水果"
	}
}
