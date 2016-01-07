package fork

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Fork struct {
	size int
	max  int
	cb   chan func()
	lock sync.RWMutex
}

func NewFork(max int) *Fork {
	return &Fork{
		max: max,
		cb:  make(chan func(), max*10),
	}
}
func (fo *Fork) Puah(f func()) {
	fo.cb <- f
	fo.lock.RLock()
	s := fo.size
	fo.lock.RUnlock()
	if s < fo.max {
		go fo.fork()
	}
}

func (fo *Fork) fork() {
	fo.lock.Lock()
	fo.size++
	fo.lock.Unlock()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		fo.lock.Lock()
		fo.size--
		fo.lock.Unlock()
	}()
	for {
		select {
		case f := <-fo.cb:
			f()
		case <-time.After(time.Second / 10):

			return
		}
	}
	return
}

func (fo *Fork) Join() {
	for {
		runtime.Gosched()
		fo.lock.RLock()
		s := fo.size
		fo.lock.RUnlock()
		if s == 0 {
			return
		}
	}
}
