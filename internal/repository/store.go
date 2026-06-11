package repository

import (
	"fmt"
	"strings"
)

const (
	StorageMemory = "memory"
	StorageMySQL  = "mysql"
)

type Store interface {
	UserRepository
	TaskRepository
}

// NewStore 根据配置创建仓储实现，默认使用内存模式方便本地学习。
func NewStore(storage string, mysqlDSN string) (Store, error) {
	switch strings.ToLower(strings.TrimSpace(storage)) {
	case "", StorageMemory:
		return NewMemoryStore(), nil
	case StorageMySQL:
		return NewGormStore(mysqlDSN)
	default:
		return nil, fmt.Errorf("unsupported storage: %s", storage)
	}
}
