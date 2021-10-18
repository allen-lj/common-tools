package concurrency

type Processor interface {
	PreProcess()
	DoProcess()
	AfterProcess()
}
