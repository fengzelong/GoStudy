package decorator

import "testing"

func TestDecoratorSuccess(t *testing.T) {
	var c Component = &ConcreteComponent{}
	c = WarpAddDecorator(c, 23)

	res := c.CalFunc()
	if res == 0 {
		t.Fatal("test failed")
	}
}

func TestDecoratorFailed(t *testing.T) {
	var c Component = &ConcreteComponent{}
	res := c.CalFunc()
	if res == 0 {
		t.Fatal("test failed")
	}
}
