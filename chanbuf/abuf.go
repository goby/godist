package chanbuf

// Have goroutine for buffer that acts as traffic director
// * Receives requests on incoming channels
// * Selects on that may proceed
// * Calling function does operation
// * Tells director that it is done
// Find that need two request channels:
// 1. Operations that can proceed in any case
// 2. Operations that block if buffer is empty
//    Read Ops    -->|
//    Other Ops   -->|                 Director
//                   |<-- Ack channel

import (
	"github.com/goby/godist/tranbuf"
)

type Buffer struct {
	sb       *tranbuf.Buffer
	ackchan  chan int
	readchan chan int
	opchan   chan int
}

func NewBuffer() *Buffer {
	bp := new(Buffer)
	bp.sb = tranbuf.NewBuffer()
	bp.ackchan = make(chan int)
	bp.readchan = make(chan int)
	bp.opchan = make(chan int)
	go bp.director()
}

// Go routine to respond to requests
func (bp *Buffer) director() {
	for {
		if bp.sb.Empty() {
			bp.opchan <- 1
		} else {
			select {
			case bp.readchan <- 1:
			case bp.opchan <- 1:
			}
		}
		<-bp.ackchan
	}
}
