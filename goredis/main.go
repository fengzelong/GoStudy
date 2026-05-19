package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

var keys [5]string

var pool *redis.Pool

func init() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
		fmt.Println("未设置 REDIS_ADDR，使用本地默认连接示例")
	}
	password := os.Getenv("REDIS_PASSWORD")

	pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			options := []redis.DialOption{redis.DialClientName("")}
			if password != "" {
				options = append(options, redis.DialPassword(password))
			}
			return redis.Dial("tcp", addr, options...)
		},
	}
}

func main() {
	conn := pool.Get()
	defer conn.Close()
	defer pool.Close()

	if err := conn.Err(); err != nil {
		fmt.Println("连接 Redis 失败:", err)
		return
	}
	fmt.Println("连接 Redis 成功")

	// 设置 Redis 值
	// res := SetFunc(conn, "def", 130)
	// fmt.Println("res = ", res)

	// 获取 Redis 值
	// res := GetFunc(conn, "abc")
	// fmt.Println("res = ", res)

	keys = [5]string{
		"abc",
		"def",
	}
	GetByKeysFunc(conn, keys)
}

// SetFunc 设置 Redis 值。
func SetFunc(conn redis.Conn, key string, value int) bool {
	_, err := conn.Do("Set", key, value)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// GetFunc 获取 Redis 值。
func GetFunc(conn redis.Conn, key string) int {
	res, err := redis.Int(conn.Do("Get", key))
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return res
}

// GetByKeysFunc 根据多个 key 获取 Redis 值。
func GetByKeysFunc(conn redis.Conn, keys [5]string) {
	res, err := redis.Ints(conn.Do("MGet", keys[0], keys[1]))
	if err != nil {
		fmt.Println("批量获取失败:", err)
		return
	}
	for _, v := range res {
		fmt.Println(v)
	}
}
