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

func (w *Marmot) AddProcessor(processor Processor) {
	w.workQueue <- processor
}

func (w *Marmot) StartWork() {
	qps := 0
	for p := range w.workQueue {
		process := p
		if qps >= w.concurrentNum {
			time.Sleep(time.Second)
			qps = 0
		}

		w.Add(1)
		qps++
		w.concurrency <- 1

		go w.doWork(process)
	}
}

func (w *Marmot) doWork(p Processor) {
	start := time.Now()
	defer func(start time.Time) {
		cost := time.Since(start)
		if cost < w.baselineTime {
			time.Sleep(w.baselineTime - cost)
		}
		w.Done()
		<-w.concurrency
	}(start)

	p.PreProcess()
	p.DoProcess()
	p.AfterProcess()
}

func (w *Marmot) CloseWorkQueue() {
	close(w.workQueue)
}

func (w *Marmot) WaitForClose() {
	w.Wait()
}

func calculateBaselineTime(num int) time.Duration {
	if num <= 0 {
		return time.Second
	}
	return time.Second / time.Duration(num)
}
