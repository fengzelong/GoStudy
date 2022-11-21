package mediator

import (
	"fmt"
	"strings"
)

type CdDriver struct {
	Data string
}

func (c *CdDriver) ReadData() {
	c.Data = "music,image"
	fmt.Printf("CdDriver: reading data: %s\n", c.Data)
	// 驱动后续动作
	GetMediatorInstance().action(c)
}

type Cpu struct {
	Video string
	Sound string
}

func (c *Cpu) Process(data string) {
	strSlices := strings.Split(data, ",")
	c.Sound = strSlices[0]
	c.Video = strSlices[1]

	fmt.Printf("Cpu: process data with video %s, sound %s\n", c.Video, c.Sound)
	GetMediatorInstance().action(c)
}

type VideoCard struct {
	Data string
}

func (v *VideoCard) Display(data string) {
	v.Data = data
	fmt.Printf("VideoCard: display %s\n", v.Data)
	GetMediatorInstance().action(v)
}

type SoundCard struct {
	Data string
}

func (s *SoundCard) Play(data string) {
	s.Data = data
	fmt.Printf("SoundCard: play %s\n", s.Data)
	GetMediatorInstance().action(s)
}

// Mediator 中介者对象应用结构体组合
type Mediator struct {
	CD    *CdDriver
	CPU   *Cpu
	Video *VideoCard
	Sound *SoundCard
}

var mediator *Mediator

// GetMediatorInstance 单例，获取中介者实例
func GetMediatorInstance() *Mediator {
	if mediator == nil {
		mediator = &Mediator{}
	}
	return mediator
}

// action 指针方法联动相关动作
func (m *Mediator) action(i interface{}) {
	switch ins := i.(type) {
	case *CdDriver:
		m.CPU.Process(ins.Data)
	case *Cpu:
		if ins.Video != "" {
			m.Video.Display(ins.Video)
		}

		if ins.Sound != "" {
			m.Sound.Play(ins.Sound)
		}
	}
}
