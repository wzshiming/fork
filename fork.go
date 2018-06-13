package fork

import (
	"sync"
)

var none = struct{}{}

type Fork struct {
	buf chan func()    // pending
	max chan struct{}  // maximum fork
	wg  sync.WaitGroup // wait group
}

func NewForkBuf(max int, buf int) *Fork {
	return &Fork{
		buf: make(chan func(), buf),
		max: make(chan struct{}, max),
	}
}

func NewFork(max int) *Fork {
	return NewForkBuf(max, 1)
}

// Len returns the size of the wait and the current fork
func (fo *Fork) Len() int {
	return len(fo.buf) + len(fo.max)
}

// Push trying to perform
func (fo *Fork) Push(f func()) {
	// If the maximum fork is not reached, the new fork execution buffer is full and blocked here
	fo.wg.Add(1)
	select {
	case fo.max <- none:
		go fo.fork(f)
	default:
		fo.buf <- f
	}
	return
}

// fork a new fork
func (fo *Fork) fork(f0 func()) bool {
	if f0 != nil {
		f0()
		fo.wg.Done()
	}

	select {
	case f := <-fo.buf:
		return fo.fork(f)
	default:
	}

	<-fo.max
	return true
}

// Join wait blocks until the WaitGroup counter is zero.
func (fo *Fork) Join() {
	fo.wg.Wait()
	return
}
