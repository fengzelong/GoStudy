package chain

import "testing"

func TestRequestChain(t *testing.T) {
	c1 := NewProjectManagerChain()
	c2 := NewDepManagerChain()
	c3 := NewGeneralManagerChain()

	c1.SetSuccessor(c2)
	c2.SetSuccessor(c3)

	var c Manager = c1

	c.HandleFeeRequest("lily", 350)
	c.HandleFeeRequest("tom", 500)

}
