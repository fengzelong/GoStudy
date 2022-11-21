package factory

import (
	"fmt"
)

// OperatorBase 基类 封装公共成员
type OperatorBase struct {
	strA, strB int32
}

// Operator 封装实际接口
type Operator interface {
	SetStrA(val int32) bool
	SetStrB(val int32) bool
	GetResult() int32
}

// OperatorFactory 工厂接口
type OperatorFactory interface {
	NewInstance() Operator
}

// PlusOperator 结构体组合（基础结构体）
type PlusOperator struct {
	*OperatorBase
}

// PlusOperatorFactory 加法工厂 标记结构体
type PlusOperatorFactory struct{}

// NewInstance 加法服务实例
func (PlusOperatorFactory) NewInstance() Operator {
	return &PlusOperator{
		OperatorBase: &OperatorBase{},
	}
}

// MinusOperator 结构体组合（基础结构体）
type MinusOperator struct {
	*OperatorBase
}

// MinusOperatorFactory 减法工厂 标记结构体
type MinusOperatorFactory struct{}

// NewInstance 减法服务实例
func (MinusOperatorFactory) NewInstance() Operator {
	return &MinusOperator{
		OperatorBase: &OperatorBase{},
	}
}

// SetStrA 设置公共结构体成员StrA，需要指针接收器，需要修改成员
func (o *OperatorBase) SetStrA(val int32) bool {
	if val <= 0 {
		fmt.Println("param invalid!")
		return false
	}
	o.strA = val
	return true
}

// SetStrB 设置公共结构体成员StrB
func (o *OperatorBase) SetStrB(val int32) bool {
	if val <= 0 {
		fmt.Println("param invalid!")
		return false
	}
	o.strB = val
	return true
}

// GetResult 加法接口实现
func (p PlusOperator) GetResult() int32 {
	return p.strA + p.strB
}

// GetResult 减法接口实现
func (m MinusOperator) GetResult() int32 {
	return m.strA - m.strB
}
