package proxy

import "fmt"

// Subject 业务抽象层
type Subject interface {
	SubjectFunc() string
}

type RealSubject struct{}

// SubjectFunc 被代理原始接口实现
func (RealSubject) SubjectFunc() string {
	return fmt.Sprintf("real subject function")
}

type PSubject struct {
	real RealSubject
}

// SubjectFunc 代理接口实现
func (p PSubject) SubjectFunc() string {
	// do something
	var res string

	res += "pre -> "

	res += p.real.SubjectFunc()

	res += " -> after"

	// do something

	return res
}
