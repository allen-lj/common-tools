package concurrency

type MockProcessor struct {
	name string
}

func (m *MockProcessor) PreProcess()   {}
func (m *MockProcessor) DoProcess()    {
	//println(m.name,"run")
	//println(time.Now().UnixMilli())
}
func (m *MockProcessor) AfterProcess() {}
