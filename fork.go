package fork // import "gopkg.in/wzshiming/fork.v2"

import (
	"sync"
)

var none = struct{}{}

type Fork struct {
	buf chan func()    // 缓冲闭包
	max chan struct{}  // 最大线程
	wg  sync.WaitGroup // 等待组
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

// 执行任务的 长度 最大 只能是线程数加缓冲数 其他的阻塞看不到
func (fo *Fork) Len() int {
	return len(fo.buf) + len(fo.max)
}

// 加入闭包 可以的话立即执行
func (fo *Fork) Push(f func()) {
	// 如果未达到最大线程则开启新线程执行 缓冲已满就在这里阻塞
	fo.wg.Add(1)
	select {
	case fo.max <- none:
		go fo.fork(f)
	default:
		fo.buf <- f
	}
	return
}

// 把当前线程加入 线程执行队列
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

// 等待所有线程结束在返回
func (fo *Fork) Join() {
	fo.wg.Wait()
	return
}
