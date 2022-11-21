package strategy

import "testing"

func TestPayment(t *testing.T) {
	ctx := &PaymentContext{
		Name:   "shopping",
		CardId: "424124124",
		Money:  200,
	}
	strategy := &Cash{}
	payment := NewPayment(ctx, strategy)
	payment.Pay()
}
