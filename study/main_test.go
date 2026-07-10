package main

import "testing"

func TestYoung(t *testing.T) {
	tests := []struct {
		name string
		age  int
		want bool
	}{
		{name: "未成年边界", age: 17, want: false},
		{name: "成年边界", age: 18, want: true},
		{name: "成年以后", age: 30, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := young(tt.age); got != tt.want {
				t.Fatalf("young(%d) = %t，期望 %t", tt.age, got, tt.want)
			}
		})
	}
}

func TestChooseFood(t *testing.T) {
	if got := ChooseFood(1); got == "" {
		t.Fatal("ChooseFood(1) 不应该返回空字符串")
	}
	if got := ChooseFood(99); got != ChooseFood(0) {
		t.Fatalf("未知食物类型返回不一致: %q vs %q", got, ChooseFood(0))
	}
	if ChooseFood(1) == ChooseFood(2) {
		t.Fatal("不同食物类型应该返回不同名称")
	}
}
