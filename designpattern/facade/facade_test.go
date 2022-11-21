package facade

import (
	"fmt"
	"strings"
	"testing"
)

func TestFacadeAPI(t *testing.T) {
	api := NewFaceApi()
	res := api.FinalCall()
	fmt.Println(res)
	if !strings.Contains(res, "B") {
		t.Fatal("failed")
	}
}
