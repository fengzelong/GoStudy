package main

import (
	"reflect"
	"sync"
	"testing"
)

func newStringMap() ConcurrentMap {
	return ConcurrentMap{
		keyType:   reflect.TypeOf(""),
		valueType: reflect.TypeOf(""),
		m:         new(sync.Map),
	}
}

func TestConcurrentMapStoreLoadDelete(t *testing.T) {
	cMap := newStringMap()
	if ok := cMap.store("name", "gostudy"); !ok {
		t.Fatal("store 应返回 true")
	}

	value, ok := cMap.load("name")
	if !ok {
		t.Fatal("load 应返回 true")
	}
	if value != "gostudy" {
		t.Fatalf("value = %v，期望 gostudy", value)
	}

	if ok := cMap.delete("name"); !ok {
		t.Fatal("delete 应返回 true")
	}
	if _, ok := cMap.load("name"); ok {
		t.Fatal("删除后不应该再读取到值")
	}
}

func TestConcurrentMapRejectsWrongType(t *testing.T) {
	cMap := newStringMap()
	defer func() {
		if recover() == nil {
			t.Fatal("错误类型应该触发 panic")
		}
	}()
	cMap.store("age", 18)
}
