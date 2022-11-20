package simplefactory

import "fmt"

type TalkApi interface {
	Say(word string) string
}

func NewTalkApi(t int) TalkApi {
	switch t {
	case 1:
		return &xiaomingSayApi{}
	case 2:
		return &xiaohuaSayApi{}
	default:
		return &xiaomingSayApi{}
	}
}

type xiaomingSayApi struct{}

func (*xiaomingSayApi) Say(word string) string {
	return fmt.Sprintf("xiaoming say: %v", word)
}

type xiaohuaSayApi struct{}

func (*xiaohuaSayApi) Say(word string) string {
	return fmt.Sprintf("xiaohua say: %v", word)
}
