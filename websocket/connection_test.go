package main

import "testing"

func TestDel(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		user string
		want []string
	}{
		{name: "空列表", in: []string{}, user: "a", want: []string{}},
		{name: "删除唯一用户", in: []string{"a"}, user: "a", want: []string{}},
		{name: "删除中间用户", in: []string{"a", "b", "c"}, user: "b", want: []string{"a", "c"}},
		{name: "用户不存在", in: []string{"a", "b"}, user: "c", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := del(tt.in, tt.user)
			if len(got) != len(tt.want) {
				t.Fatalf("长度 = %d，期望 %d，结果 %#v", len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("got[%d] = %q，期望 %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
