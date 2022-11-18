package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var keys [5]string

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "124.223.8.183:9905", redis.DialClientName(""), redis.DialPassword("fl666@2022"))
		},
	}
}

func main() {
	conn := pool.Get()
	//if err != nil {
	//	fmt.Println("conn redis failed,", err)
	//	return
	//}
	fmt.Println("redis conn success")

	//redis 设置值
	//res := SetFunc(conn, "def", 130)
	//fmt.Println("res = ", res)

	//redis 获取值
	//res1 := GetFunc(conn, "abc")
	//fmt.Println("res = ", res1)

	keys = [5]string{
		"abc",
		"def",
	}
	GetByKeysFunc(conn, keys)

	defer conn.Close()
	pool.Close() //连接池关闭
}

// SetFunc set
func SetFunc(conn redis.Conn, key string, value int) bool {
	_, err := conn.Do("Set", key, value)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// GetFunc get
func GetFunc(conn redis.Conn, key string) int {
	res, err := redis.Int(conn.Do("Get", key))
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return res
}

// GetByKeysFunc get
// conn redis连接
// keys redis缓存键数组
func GetByKeysFunc(conn redis.Conn, keys [5]string) {
	res, err := redis.Ints(conn.Do("MGet", keys[0], keys[1]))
	if err != nil {
		fmt.Println("get abc failed,", err)
		return
	}
	for _, v := range res {
		fmt.Println(v)
	}
}
