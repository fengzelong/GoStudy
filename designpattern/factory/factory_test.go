package factory

import "testing"

func compute(factory OperatorFactory, strA, strB int32) int32 {
	op := factory.NewInstance()
	op.SetStrA(strA)
	op.SetStrB(strB)
	return op.GetResult()
}

var factory OperatorFactory

// TestPlusOperator 加法操作测试
func TestPlusOperatorSuccess(t *testing.T) {
	factory = PlusOperatorFactory{}
	if res := compute(factory, 32, 10); res != 42 {
		t.Fatal("加法操作失败")
	}
}

func TestPlusOperatorWithDifferentInputs(t *testing.T) {
	factory = PlusOperatorFactory{}
	if res := compute(factory, 22, 10); res != 32 {
		t.Fatal("加法操作失败")
	}
}

// TestMinusOperator 减法操作测试
func TestMinusOperatorSuccess(t *testing.T) {
	factory = MinusOperatorFactory{}
	if res := compute(factory, 22, 10); res != 12 {
		t.Fatal("减法操作失败")
	}
}

func TestMinusOperatorWithDifferentExpectedResult(t *testing.T) {
	factory = MinusOperatorFactory{}
	if res := compute(factory, 22, 10); res != 12 {
		t.Fatal("减法操作失败")
	}
}
