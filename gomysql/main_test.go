package main

import (
	"os"
	"testing"
)

func TestOpenDBWithConfiguredMySQL(t *testing.T) {
	if os.Getenv("MYSQL_DSN") == "" {
		t.Skip("未设置 MYSQL_DSN，跳过 MySQL 真实连接测试")
	}

	db, err := openDB()
	if err != nil {
		t.Fatalf("打开 MySQL 失败: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("连接 MySQL 失败: %v", err)
	}
}
