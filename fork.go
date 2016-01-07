package fork

import (
	"fmt"
	"runtime"
	"time"
)

type Fork struct {
	size int
	max  int
	cb   chan func()
}

func NewFork(max int) *Fork {
	return &Fork{
		size: 0,
		max:  max,
		cb:   make(chan func(), max*10),
	}
}
func (fo *Fork) Puah(f func()) {
	if fo.size < fo.max {
		go fo.fork()
	}
	fo.cb <- f
}

func (fo *Fork) fork() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		fo.size--
	}()
	fo.size++
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
		if fo.size == 0 {
			return
		}
	}
}
