package builder

import (
	"strings"
	"testing"
)

func TestBuilderAWithoutClothes(t *testing.T) {
	builder := &builderA{}
	director := NewDirector(builder)
	director.ConstructA()
	res := builder.ActionResult()
	if strings.Contains(res, "nike") {
		t.Fatal("ConstructA 不应该包含衣服")
	}
}

func TestBuilderASuccess(t *testing.T) {
	builder := &builderA{}
	director := NewDirector(builder)
	director.ConstructB()
	res := builder.ActionResult()
	if !strings.Contains(res, "nike") {
		t.Fatal("没有穿衣服")
	}
}
