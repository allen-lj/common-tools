package concurrency

import (
	"fmt"
	"sync"
	"time"
)

var freshTime = time.Millisecond * 100

// Marmot 采用同时并发执行，一般用于压测等同时并发场景
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

func (m *Marmot) WaitForClose() {
	m.isEmptyQueue()
	m.Wait()
	m.closeWorkQueue()
}

// TODO: to be optimized
func (m *Marmot) isEmptyQueue() {
	t := time.NewTicker(freshTime)
	for {
		select {
		case <-t.C:
			fmt.Println(len(m.workQueue))
			if len(m.workQueue) == 0 {
				return
			}
		}
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

func (m *Marmot) closeWorkQueue() {
	close(m.workQueue)
}

func calculateBaselineTime(num int) time.Duration {
	if num <= 0 {
		return time.Second
	}
	return time.Second / time.Duration(num)
}
