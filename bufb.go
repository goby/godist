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

type BufferNode struct {
    value [] byte
    next *BufferNode
}

func NewBuffer() *Buffer {
    return new(Buffer);
}

func (bp *Buffer) Insert(value []byte) {
    node := &BufferNode(value: value)
    if bp.head == nil {
        
    }
}
