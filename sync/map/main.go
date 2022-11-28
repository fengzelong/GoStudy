package main

import (
	"fmt"
	"reflect"
	"sync"
)

type ConcurrentMap struct {
	m         *sync.Map
	keyType   reflect.Type
	valueType reflect.Type
}

func main() {
	key := "test"
	cMap := ConcurrentMap{
		keyType:   reflect.TypeOf(""),
		valueType: reflect.TypeOf(""),
		m:         new(sync.Map),
	}
	cMap.store(key, "test value ...")
	//cMap.store(key, 123)

	res, ok := cMap.load(key)
	if ok {
		fmt.Println(res)
	} else {
		fmt.Println("get error")
	}
}

func (cMap *ConcurrentMap) load(key interface{}) (value interface{}, ok bool) {
	if reflect.TypeOf(key) != cMap.keyType {
		return
	}
	return cMap.m.Load(key)
}

// store 需要做类型检查
func (cMap *ConcurrentMap) store(key, value interface{}) (ok bool) {
	if reflect.TypeOf(key) != cMap.keyType {
		panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
	}
	if reflect.TypeOf(value) != cMap.valueType {
		panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
	}
	cMap.m.Store(key, value)
	return true
}

// delete 删除
func (cMap *ConcurrentMap) delete(key interface{}) bool {
	if reflect.TypeOf(key) != cMap.keyType {
		return false
	}
	cMap.m.Delete(key)
	return true
}
