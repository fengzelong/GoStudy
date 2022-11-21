package decorator

type Component interface {
	CalFunc() int32
}
type ConcreteComponent struct{}

func (*ConcreteComponent) CalFunc() int32 {
	return 0
}

// AddDecorator Go语言借助于匿名组合和非入侵式接口可以很方便实现装饰模式
type AddDecorator struct {
	Component
	num int32
}

func WarpAddDecorator(c Component, num int32) Component {
	return &AddDecorator{
		Component: c,
		num:       num,
	}
}

func (d *AddDecorator) CalFunc() int32 {
	return d.Component.CalFunc() + d.num
}
