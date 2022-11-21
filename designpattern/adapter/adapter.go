package adapter

import (
	"fmt"
)

type Target interface {
	CommonFunc(int32) (string, bool)
}

type Source interface {
	SourceFunc(int32) string
}

// NewAdapter 适配器工厂
func NewAdapter(source Source) Target {
	return &adapter{
		Source: source,
	}
}

// sourceImpl 适配源标记结构体
type sourceImpl struct{}

// SourceFunc 被适配接口实现
func (*sourceImpl) SourceFunc(t int32) string {
	switch t {
	case 1:
		return fmt.Sprintf("type %v source method", t)
	case 2:
		return fmt.Sprintf("type %v source method", t)
	default:
		return fmt.Sprintln("default")
	}
}

type adapter struct {
	Source
}

// CommonFunc 适配器接口实现
func (a *adapter) CommonFunc(t int32) (string, bool) {
	res := a.SourceFunc(t)
	if res == "" {
		return "", false
	}
	return res, true
}
