package simplefactory

import (
	"strings"
	"testing"
)

func TestXiaoMing(t *testing.T) {
	api := NewTalkApi(1)
	msg := api.Say("hello")
	if !strings.Contains(msg, "xiaoming") {
		t.Fatal("xiaoming test fail")
	}
}

func TestXiaoHua(t *testing.T) {
	api := NewTalkApi(2)
	msg := api.Say("hello")
	if !strings.Contains(msg, "xiaohua") {
		t.Fatal("xiaohua test fail")
	}
}
