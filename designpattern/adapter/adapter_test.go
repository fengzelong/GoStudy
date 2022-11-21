package adapter

import (
	"strings"
	"testing"
)

func TestAdapter(t *testing.T) {
	source := &sourceImpl{}
	target := NewAdapter(source)
	res, ok := target.CommonFunc(2)
	if !ok {
		t.Fatal("test failed")
	}
	if !strings.Contains(res, "source method") {
		t.Fatal("test failed")
	}
}
