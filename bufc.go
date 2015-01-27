package bufb

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
	sb       *SeqBuffer
	ackchan  chan int
	readchan chan int
	opchan   chan int
}

func NewChanBuffer() *ChanBuffer {
	bp := new(ChanBuffer)
	bp.sb = NewSeqBuffer()
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
