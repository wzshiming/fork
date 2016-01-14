package fork

import "runtime"

type Fork struct {
	max chan int
	buf chan func()
}

func NewForkBuf(max int, buf int) *Fork {
	fo := &Fork{
		max: make(chan int, max),
		buf: make(chan func(), buf),
	}
	for i := 0; i != max; i++ {
		fo.max <- 0
	}
	return fo
}

func NewFork(max int) *Fork {
	return NewForkBuf(max, max*10)
}

func (fo *Fork) Puah(f func()) {
	fo.buf <- f
	select {
	case <-fo.max:
		go fo.fork()
	default:
	}
}

func (fo *Fork) fork() {
	for {
		select {
		case f := <-fo.buf:
			f()
		default:
			fo.max <- 0
			return
		}
	}
}

func (fo *Fork) Join() {
	for {
		runtime.Gosched()
		select {
		case fo.max <- 0:
			<-fo.max
		default:
			return
		}
	}
}
