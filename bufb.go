// DESCRIPTION
//  Implementation of FIFO buffer using linked list + tailed pointer. Single-threaded operation.
// FEATURE
//  NewBuffer: Create a new buffer
//  Insert: Insert a node into buffer
//  Front: Get first node in buffer without changing buffer
//  Remove: Get and remove first node from buffer
//  Flush:  Clear buffer
// ROADMAP
//  TODO: Multi-threaded support
package bufb

import (
	"errors"
)

type BufferNode struct {
	value []byte
	next  *BufferNode
}

type Buffer struct {
	head *BufferNode
	tail *BufferNode
	size int
}

func NewBuffer() *Buffer {
	return &Buffer{size: 0, head: nil, tail: nil}
}

func (bp *Buffer) Insert(value []byte) {
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

func (bp *Buffer) Front() ([]byte, error) {
	if bp.head == nil {
		err := errors.New("Empty Buffer")
		return nil, err
	} else {
		return bp.head.value, nil
	}
}

func (bp *Buffer) Remove() ([]byte, error) {
	if bp.head == nil {
		err := errors.New("Empty Buffer")
		return nil, err
	} else {
		var ret []byte = bp.head.value
		if bp.tail == bp.head {
			bp.tail = nil
		}
		bp.size--
		bp.head = bp.head.next
		return ret, nil
	}
}

func (bp *Buffer) Flush() {
	bp.head = nil
	bp.tail = nil
	bp.size = 0
}

func (bp *Buffer) Empty() bool {
	return bp.size == 0
}

func (bp *Buffer) Size() int {
	return bp.size
}
