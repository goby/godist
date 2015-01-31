package chanbuf

import (
	"github.com/goby/godist/oldbuf"
)

// Using go channels

type Mutex struct {
	mc chan int
}

func NewMutex() *Mutex {
	m := &Mutex{make(chan int, 1)}
	m.Unlock()
	return m
}

func (m *Mutex) Lock() {
	<-m.mc
}

func (m *Mutex) Unlock() {
	m.mc <- 1
}

// ChanBuffer
type ChanBuffer struct {
	sb       *oldbuf.SeqBuffer
	ackchan  chan int
	readchan chan int
	opchan   chan int
}

func NewChanBuffer() *ChanBuffer {
	bp := new(ChanBuffer)
	bp.sb = oldbuf.NewSeqBuffer()
	bp.ackchan = make(chan int)
	bp.readchan = make(chan int)
	bp.opchan = make(chan int)
	go bp.director()
	return bp
}

// Go routine to respond to requests
func (bp *ChanBuffer) director() {
	for {
		if bp.sb.Empty() {
			// Enable only nonblocking operations
			bp.opchan <- 1
		} else {
			// Enable reads and other operations
			// Will alow only on communication
			select {
			case bp.readchan <- 1:
			case bp.opchan <- 1:
			}
		}
		<-bp.ackchan // Wait until operations completed
	}
}

func (bp *ChanBuffer) Insert(value interface{}) {
	<-bp.opchan
	bp.sb.Insert(value)
	bp.ackchan <- 1
}

func (bp *ChanBuffer) Remove() interface{} {
	<-bp.readchan
	x, _ := bp.sb.Remove()
	bp.ackchan <- 1
	return x
}

func (bp *ChanBuffer) Flush() {
	<-bp.opchan
	bp.sb.Flush()
	bp.ackchan <- 1
}

func (bp *ChanBuffer) Front() interface{} {
	<-bp.readchan
	x, _ := bp.sb.Front()
	bp.ackchan <- 1
	return x
}

func (bp *ChanBuffer) Empty() bool {
	x := bp.sb.Empty()
	return x
}
