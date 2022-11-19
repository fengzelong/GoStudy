package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

type TestJob struct {
}

func (t TestJob) Run() {
	fmt.Println("testJob1 running ...")
}

type Test2Job struct {
}

func (t Test2Job) Run() {
	fmt.Println("testJob2 running ...")
}

func main() {
	i := 0
	c := NewWithSeconds()
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		i++
		fmt.Println("cron running:", i)
	})

	c.AddJob(spec, TestJob{})
	c.AddJob(spec, Test2Job{})

	c.Start()

	defer c.Stop()

	select {}
}

// NewWithSeconds 生成cron实例
func NewWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}
