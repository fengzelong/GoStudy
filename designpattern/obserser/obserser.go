package obserser

import "fmt"

type Subject struct {
	observers []Observer //观察者切片
	context   string     //通知内容
}

// NewSubject 生成subject对象
func NewSubject() *Subject {
	return &Subject{
		observers: make([]Observer, 0),
	}
}

// attach 追加观察者
func (s *Subject) attach(o Observer) {
	s.observers = append(s.observers, o)
}

// notify 遍历观察者切片发送通知消息
func (s *Subject) notify() {
	for _, o := range s.observers {
		o.Receive(s)
	}
}

// SendContext 发送消息
func (s *Subject) SendContext(context string) {
	s.context = context
	s.notify()
}

type Observer interface {
	Receive(*Subject)
}

type Reader struct {
	name string
}

// NewReader 创建观察者
func NewReader(name string) *Reader {
	return &Reader{
		name: name,
	}
}

func (r *Reader) Receive(s *Subject) {
	fmt.Printf("%s receive %s\n", r.name, s.context)
}
