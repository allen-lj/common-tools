package concurrency

import (
	"sync"
	"time"
)

type Processor interface {
	PreProcess()
	DoProcess()
	AfterProcess()
}

type Marmot struct {
	workQueue     chan Processor
	concurrency   chan int
	concurrentNum int
	baselineTime  time.Duration
	sync.WaitGroup
}

func NewMarmot(queueLength int, concurrentNum int) *Marmot {
	return &Marmot{
		workQueue:     make(chan Processor, queueLength),
		concurrency:   make(chan int, concurrentNum),
		concurrentNum: concurrentNum,
		baselineTime:  calculateBaselineTime(concurrentNum),
	}
}

func (m *Marmot) AddProcessor(processor Processor) {
	m.workQueue <- processor
}

func (m *Marmot) StartWork() {
	qps := 0
	for p := range m.workQueue {
		process := p
		if qps >= m.concurrentNum {
			time.Sleep(time.Second)
			qps = 0
		}

		m.Add(1)
		qps++
		m.concurrency <- 1

		go m.doWork(process)
	}
}

func (m *Marmot) doWork(p Processor) {
	start := time.Now()
	defer func(start time.Time) {
		cost := time.Since(start)
		if cost < m.baselineTime {
			time.Sleep(m.baselineTime - cost)
		}
		m.Done()
		<-m.concurrency
	}(start)

	p.PreProcess()
	p.DoProcess()
	p.AfterProcess()
}

func (m *Marmot) CloseWorkQueue() {
	close(m.workQueue)
}

func (m *Marmot) WaitForClose() {
	m.Wait()
}

func calculateBaselineTime(num int) time.Duration {
	if num <= 0 {
		return time.Second
	}
	return time.Second / time.Duration(num)
}
