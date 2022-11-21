package strategy

import "fmt"

// PaymentContext 策略模式上下文对象
type PaymentContext struct {
	Name, CardId string
	Money        int32
}

// PaymentStrategy 策略模式接口签名
type PaymentStrategy interface {
	Pay(*PaymentContext)
}

type Cash struct{}

// Pay cash具体实现
func (*Cash) Pay(ctx *PaymentContext) {
	fmt.Printf("pay ￥%d to %s by cash", ctx.Money, ctx.Name)
}

type AliPay struct{}

func (*AliPay) Pay(ctx *PaymentContext) {
	fmt.Printf("pay ￥%d to %s by cash", ctx.Money, ctx.Name)
}

type WechatPay struct{}

func (*WechatPay) Pay(ctx *PaymentContext) {
	fmt.Printf("pay ￥%d to %s by cash", ctx.Money, ctx.Name)
}

type Payment struct {
	context  *PaymentContext
	strategy PaymentStrategy
}

// Pay 匹配策略类型，接口具体实现
func (p *Payment) Pay() {
	p.strategy.Pay(p.context)
}

// NewPayment 生成payment实例
func NewPayment(ctx *PaymentContext, strategy PaymentStrategy) *Payment {
	return &Payment{
		context: &PaymentContext{
			Name:   ctx.Name,
			CardId: ctx.CardId,
			Money:  ctx.Money,
		},
		strategy: strategy,
	}
}
