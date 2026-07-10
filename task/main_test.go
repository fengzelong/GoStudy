package main

import (
	"testing"
	"time"
)

func TestNewWithSecondsAcceptsSecondFieldSpec(t *testing.T) {
	c := NewWithSeconds()
	id, err := c.AddFunc("*/5 * * * * ?", func() {})
	if err != nil {
		t.Fatalf("带秒字段的 cron 表达式应该可用: %v", err)
	}
	if id == 0 {
		t.Fatal("任务 ID 不应为 0")
	}
}

func TestNewWithSecondsRejectsInvalidSpec(t *testing.T) {
	c := NewWithSeconds()
	if _, err := c.AddFunc("bad spec", func() {}); err == nil {
		t.Fatal("非法 cron 表达式应该返回错误")
	}
}

func TestJobRunDoesNotPanic(t *testing.T) {
	done := make(chan struct{})
	go func() {
		TestJob{}.Run()
		Test2Job{}.Run()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Job Run 执行超时")
	}
}
