package concurrency

import "testing"

func TestNewMarmot(t *testing.T) {
	m := NewMarmot(3, 1)
	m.AddProcessor(&MockProcessor{})
	go m.StartWork()
	m.WaitForClose()
}

type MockProcessor struct {
}

func (m *MockProcessor) PreProcess()   {}
func (m *MockProcessor) DoProcess()    {}
func (m *MockProcessor) AfterProcess() {}
