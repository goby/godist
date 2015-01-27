package chanbuf

import (
	"github.com/goby/godist/tranbuf"
)

const (
	doinsert = iota
	doremove
	doflush
	doempty
)

type request struct {
	op     int
	val    interface{}
	replyc chan interface{}
}

type Buffer struct {
	opc   chan *request
	readc chan *request
}

func NewBuffer() *Buffer {
	bp := &Buffer{make(chan *request), make(chan *request)}
	go bp.runServer()
	return bp
}

func (bp *Buffer) runServer() {
	sb := tranbuf.NewBuffer()
	for {
		var r *request
		if sb.Empty() {
			r = <-bp.opc
		} else {
			select {
			case r1 := <-bp.opc:
				r = r1
			case r2 := <-bp.readc:
				r = r2
			}
		}
	}
}
