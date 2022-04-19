package concurrency

type Processor interface {
	PreProcess()
	DoProcess()
	AfterProcess()
}

type Worker interface {
	AddProcessor(processor Processor)
	StartWork()
	WaitForClose()
	WorkDone()
}
