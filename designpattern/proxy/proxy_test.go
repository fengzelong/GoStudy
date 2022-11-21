package proxy

import (
	"fmt"
	"strings"
	"testing"
)

func TestSubject(t *testing.T) {
	var subject Subject
	subject = &PSubject{}

	res := subject.SubjectFunc()
	fmt.Println(res)

	if !strings.HasPrefix(res, "pre") {
		t.Fatal("test failed")
	}
}
