package arith

import "testing"

func TestMultiply(t *testing.T) {
	service := &Arith{}
	req := Request{A: 9, B: 2}
	var res Response

	if err := service.Multiply(req, &res); err != nil {
		t.Fatalf("Multiply() 返回错误: %v", err)
	}
	if res.Pro != 18 {
		t.Fatalf("乘积 = %d，期望 18", res.Pro)
	}
}

func TestDivide(t *testing.T) {
	service := &Arith{}
	req := Request{A: 9, B: 2}
	var res Response

	if err := service.Divide(req, &res); err != nil {
		t.Fatalf("Divide() 返回错误: %v", err)
	}
	if res.Quo != 4 || res.Rem != 1 {
		t.Fatalf("商和余数 = %d, %d，期望 4, 1", res.Quo, res.Rem)
	}
}

func TestDivideByZero(t *testing.T) {
	service := &Arith{}
	req := Request{A: 9, B: 0}
	var res Response

	if err := service.Divide(req, &res); err == nil {
		t.Fatal("除数为零时应该返回错误")
	}
}
