package builder

import (
	"strings"
	"testing"
)

func TestBuilderAFailed(t *testing.T) {
	builder := &builderA{}
	director := NewDirector(builder)
	director.ConstructA()
	res := builder.ActionResult()
	if !strings.Contains(res, "nike") {
		t.Fatal("no wear")
	}
}

func TestBuilderASuccess(t *testing.T) {
	builder := &builderA{}
	director := NewDirector(builder)
	director.ConstructB()
	res := builder.ActionResult()
	if !strings.Contains(res, "nike") {
		t.Fatal("no wear")
	}
}
