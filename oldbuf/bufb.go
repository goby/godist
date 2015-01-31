package oldbuf

// DESCRIPTION
//  Implementation of FIFO buffer using linked list + tailed pointer.
//      With single-threaded operation SqeBuffer and thread-safe struct
//      Buffer(with traditional sync ).
// FEATURE
//  NewBuffer: Create a new buffer
//  Insert: Insert a node into buffer
//  Front: Get first node in buffer without changing buffer
//  Remove: Get and remove first node from buffer
//  Flush:  Clear buffer
// ROADMAP
//  Done: Multi-threaded support
import (
	"errors"
	"sync"
)

// BufferNode
type BufferNode struct {
	value interface{}
	next  *BufferNode
}

// Sequencial Buffer
type SeqBuffer struct {
	head *BufferNode
	tail *BufferNode
	size int
}

func NewSeqBuffer() *SeqBuffer {
	return &SeqBuffer{size: 0, head: nil, tail: nil}
}

func (bp *SeqBuffer) Insert(value interface{}) {
	node := &BufferNode{value: value}
	if bp.head == nil {
		bp.head = node
		bp.tail = node
	} else {
		bp.tail.next = node
		bp.tail = node
	}
	bp.size++
}

func (bp *SeqBuffer) Front() (interface{}, error) {
	if bp.head == nil {
		err := errors.New("Empty Buffer")
		return nil, err
	} else {
		return bp.head.value, nil
	}
}

func (bp *SeqBuffer) Remove() (interface{}, error) {
	if bp.head == nil {
		err := errors.New("Empty Buffer")
		return nil, err
	} else {
		var ret interface{} = bp.head.value
		if bp.tail == bp.head {
			bp.tail = nil
		}
		bp.size--
		bp.head = bp.head.next
		return ret, nil
	}
}

func (bp *SeqBuffer) Flush() {
	bp.head = nil
	bp.tail = nil
	bp.size = 0
}

func (bp *SeqBuffer) Empty() bool {
	return bp.size == 0
}

func (bp *SeqBuffer) Size() int {
	return bp.size
}

// Thread-safe Buffer
type Buffer struct {
	sb    *SeqBuffer
	mutex *sync.Mutex
	cvar  *sync.Cond
}

func NewBuffer() *Buffer {
	bp := &Buffer{}
	bp.sb = NewSeqBuffer()
	bp.mutex = &sync.Mutex{}
	bp.cvar = sync.NewCond(bp.mutex)
	return bp
}

func (bp *Buffer) Insert(value interface{}) {
	bp.mutex.Lock()
	bp.sb.Insert(value)
	bp.cvar.Signal()
	bp.mutex.Unlock()
}

// Block until some value will be return
func (bp *Buffer) Remove() (interface{}, error) {
	bp.mutex.Lock()
	for bp.sb.Empty() {
		bp.cvar.Wait()
	}
	x, err := bp.sb.Remove()
	bp.mutex.Unlock()
	return x, err
}

func (bp *Buffer) Flush() {
	bp.mutex.Lock()
	bp.sb.Flush()
	bp.mutex.Unlock()
}

func (bp *Buffer) Size() int {
	bp.mutex.Lock()
	size := bp.sb.Size()
	bp.mutex.Unlock()
	return size
}

func (bp *Buffer) Empty() bool {
	bp.mutex.Lock()
	isempty := bp.sb.Empty()
	bp.mutex.Unlock()
	return isempty
}
