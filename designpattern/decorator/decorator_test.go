package decorator

import "testing"

func TestDecoratorSuccess(t *testing.T) {
	var c Component = &ConcreteComponent{}
	c = WarpAddDecorator(c, 23)

	res := c.CalFunc()
	if res == 0 {
		t.Fatal("装饰器计算失败")
	}
}

func TestConcreteComponentDefault(t *testing.T) {
	var c Component = &ConcreteComponent{}
	res := c.CalFunc()
	if res != 0 {
		t.Fatal("默认组件应该返回零值")
	}
}
