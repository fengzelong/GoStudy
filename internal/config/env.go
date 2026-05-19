package config

import "os"

// Env 获取环境变量；未设置时返回默认值。
func Env(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
