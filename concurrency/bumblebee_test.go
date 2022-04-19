package concurrency

import (
	"testing"
	"time"
)

func TestNewBumbleBee(t *testing.T) {
	b := NewBumbleBee(3, 1)
	b.AddProcessor(&MockProcessor{name: "one"})
	b.AddProcessor(&MockProcessor{name: "two"})
	b.AddProcessor(&MockProcessor{name: "three"})
	go b.StartWork()
	go func() {
		time.Sleep(time.Second * 3)
		b.WorkDone()
	}()
	b.WaitForClose()
}
