package concurrency

import (
	"github.com/allen-lj/common-tools/util"
	"sync"
	"time"
)

type BumbleBee struct {
	workQueue     chan Processor
	token         chan token
	concurrentNum uint
	isWorkDone    chan bool
	isTokenClose  chan bool
	sync.WaitGroup
}

func NewBumbleBee(queueLength uint, concurrentNum uint) *BumbleBee {
	util.Assert(concurrentNum <= 0, "the %s must be greater than 0", "concurrentNum")
	util.Assert(queueLength <= 0, "the %s must be greater than 0", "queueLength")

	return &BumbleBee{
		workQueue:     make(chan Processor, queueLength),
		token:         make(chan token, concurrentNum),
		concurrentNum: concurrentNum,
		isWorkDone:    make(chan bool, 1),
		isTokenClose:  make(chan bool, 1),
	}
}

type token byte

var _ Worker = &BumbleBee{}

func (b *BumbleBee) AddProcessor(processor Processor) {
	b.workQueue <- processor
}

func (b *BumbleBee) StartWork() {
	go b.generateToken()
	for p := range b.workQueue {
		process := p
		<-b.token
		b.Add(1)
		go b.doWork(process)
	}
}

func (b *BumbleBee) WorkDone() {
	b.isWorkDone <- true
}

func (b *BumbleBee) WaitForClose() {
	<-b.isWorkDone
	b.Wait()
	b.isTokenClose <- true
	close(b.workQueue)
}

func (b *BumbleBee) doWork(p Processor) {
	defer b.Done()
	p.PreProcess()
	p.DoProcess()
	p.AfterProcess()
}

func (b *BumbleBee) generateToken() {
	t := time.NewTicker(time.Second / time.Duration(b.concurrentNum))

	for {
		select {
		case <-t.C:
			b.token <- token(0)
		case <-b.isTokenClose:
			close(b.token)
			return
		}
	}
}
