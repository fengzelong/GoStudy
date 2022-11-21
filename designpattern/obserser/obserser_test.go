package obserser

import (
	"fmt"
	"testing"
)

func TestObserser(t *testing.T) {

	sub := NewSubject()
	for i := 1; i <= 5; i++ {
		name := fmt.Sprintf("reader %d", i)
		reader := NewReader(name)
		sub.attach(reader)
	}

	sub.SendContext("everything change!!")
}
