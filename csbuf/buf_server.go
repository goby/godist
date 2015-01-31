package csbuf

import (
	"github.com/goby/godist/oldbuf"
)

const (
	doinsert = iota
	doremove
	doflush
	doempty
	dofront
)

// Implement a set in go
var deferOnEmpty = map[int]bool{doremove: true}

type request struct {
	op     int
	val    interface{}
	replyc chan interface{}
}

type Buffer struct {
	requestc chan *request
}

func NewBuffer() *Buffer {
	bp := &Buffer{make(chan *request)}
	go bp.runServer()
	return bp
}

func (bp *Buffer) runServer() {
	// Buffer to hold data
	sb := oldbuf.NewSeqBuffer()
	// Buffer to hold deferred request
	db := oldbuf.NewSeqBuffer()
	for {
		var r *request
		// No need to select, We do our own scheduling
		if !sb.Empty() && !db.Empty() {
			// Revisit deferred operation
			b, _ := db.Remove()
			r = b.(*request)
		} else {
			r = <-bp.requestc
			if sb.Empty() && deferOnEmpty[r.op] {
				// Must defer this operation
				db.Insert(r)
				continue
			}
		}
		switch r.op {
		case doinsert:
			sb.Insert(r.val)
			r.replyc <- nil
		case doremove:
			v, _ := sb.Remove()
			r.replyc <- v
		case doempty:
			b := sb.Empty()
			r.replyc <- b
		case doflush:
			sb.Flush()
			r.replyc <- nil
		case dofront:
			v, _ := sb.Front()
			r.replyc <- v
		}
	}
}

func (bp *Buffer) dorequest(op int, val interface{}) interface{} {
	r := &request{op, val, make(chan interface{})}
	bp.requestc <- r
	v := <-r.replyc
	return v
}

func (bp *Buffer) Insert(val interface{}) {
	bp.dorequest(doinsert, val)
}

func (bp *Buffer) Empty() bool {
	e := bp.dorequest(doempty, nil)
	return e.(bool)
}

func (bp *Buffer) Remove() interface{} {
	v := bp.dorequest(doremove, nil)
	return v
}

func (bp *Buffer) Flush() {
	bp.dorequest(doflush, nil)
}

func (bp *Buffer) Front() interface{} {
	v := bp.dorequest(dofront, nil)
	return v
}
