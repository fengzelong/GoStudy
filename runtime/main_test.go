package main

import (
	"runtime"
	"testing"
)

func TestGOMAXPROCSCanBeAdjusted(t *testing.T) {
	old := runtime.GOMAXPROCS(1)
	t.Cleanup(func() {
		runtime.GOMAXPROCS(old)
	})

	if got := runtime.GOMAXPROCS(2); got != 1 {
		t.Fatalf("上一次 GOMAXPROCS = %d，期望 1", got)
	}
	if got := runtime.GOMAXPROCS(old); got != 2 {
		t.Fatalf("上一次 GOMAXPROCS = %d，期望 2", got)
	}
}

func TestRuntimeBasicInfo(t *testing.T) {
	if runtime.NumCPU() <= 0 {
		t.Fatal("NumCPU 应该大于 0")
	}
	if runtime.GOROOT() == "" {
		t.Fatal("GOROOT 不应为空")
	}
}
