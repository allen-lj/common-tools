package concurrency

import (
	"github.com/allen-lj/common-tools/util"
	"sync"
	"time"
)

// Marmot 采用同时并发执行，一般用于压测等同时并发场景
type Marmot struct {
	workQueue     chan Processor
	concurrencyCh chan int
	concurrentNum uint
	baselineTime  time.Duration
	isWorkDone    chan bool
	sync.WaitGroup
}

func NewMarmot(queueLength int, concurrentNum uint) *Marmot {
	util.Assert(concurrentNum <= 0, "the %s must be greater than 0", "concurrentNum")
	util.Assert(queueLength <= 0, "the %s must be greater than 0", "queueLength")

	return &Marmot{
		workQueue:     make(chan Processor, queueLength),
		concurrencyCh: make(chan int, concurrentNum),
		concurrentNum: concurrentNum,
		baselineTime:  calculateBaselineTime(concurrentNum),
		isWorkDone:    make(chan bool, 1),
	}
}

var _ Worker = &Marmot{}

func (m *Marmot) AddProcessor(processor Processor) {
	m.workQueue <- processor
}

func (m *Marmot) StartWork() {
	for p := range m.workQueue {
		process := p
		m.Add(1)
		m.concurrencyCh <- 1

		go m.doWork(process)
	}
}

func (m *Marmot) WorkDone() {
	m.isWorkDone <- true
}

func (m *Marmot) WaitForClose() {
	<-m.isWorkDone
	m.Wait()
	m.closeWorkQueue()
}

func (m *Marmot) doWork(p Processor) {
	start := time.Now()
	defer func(start time.Time) {
		cost := time.Since(start)
		if cost < m.baselineTime {
			time.Sleep(m.baselineTime - cost)
		}

		m.Done()
		<-m.concurrencyCh
	}(start)

	p.PreProcess()
	p.DoProcess()
	p.AfterProcess()
}

func (m *Marmot) closeWorkQueue() {
	close(m.workQueue)
}

func calculateBaselineTime(num uint) time.Duration {
	if num <= 0 {
		return time.Second
	}
	return time.Second / time.Duration(num)
}
