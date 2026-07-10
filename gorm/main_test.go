package main

import (
	"os"
	"testing"
)

func TestOpenDBWithConfiguredGormDSN(t *testing.T) {
	if os.Getenv("GORM_DSN") == "" {
		t.Skip("未设置 GORM_DSN，跳过 Gorm 真实连接测试")
	}

	db, err := openDB()
	if err != nil {
		t.Fatalf("打开 Gorm 连接失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("获取底层数据库连接失败: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("连接 Gorm 数据库失败: %v", err)
	}
}
