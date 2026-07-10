package main

import "testing"

func TestTeacherWriteSkipsWhenIndexNotCurrentPage(t *testing.T) {
	page = 1
	blackBoard = ""
	teacher := &Teacher{Name: "teacher"}

	teacher.write(0)

	if blackBoard != "" {
		t.Fatalf("blackBoard = %q，期望保持为空", blackBoard)
	}
	if page != 1 {
		t.Fatalf("page = %d，期望保持 1", page)
	}
}
