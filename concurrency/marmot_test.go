package concurrency

import (
	"testing"
	"time"
)

func TestNewMarmot(t *testing.T) {
	m := NewMarmot(3, 1)
	m.AddProcessor(&MockProcessor{name: "one"})
	m.AddProcessor(&MockProcessor{name: "two"})
	m.AddProcessor(&MockProcessor{name: "three"})
	go m.StartWork()

	time.Sleep(time.Second * 3)
	m.WorkDone()
	m.WaitForClose()
}
