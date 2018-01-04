package fork

import (
	"runtime"
	"time"
)

var none = struct{}{}

type Fork struct {
	buf chan func()   // 缓冲闭包
	max chan struct{} // 最大线程
	sub chan struct{} // 线程结束 信号
}

func NewForkBuf(max int, buf int) *Fork {
	return &Fork{
		buf: make(chan func(), buf),
		max: make(chan struct{}, max),
		sub: make(chan struct{}, 1),
	}
}

func NewFork(max int) *Fork {
	return NewForkBuf(max, 1)
}

// 清空缓冲
func (fo *Fork) CleanBuf() {
	for {
		select {
		case <-fo.buf:
		default:
			runtime.Gosched()
			return
		}
	}
	return
}

// 执行任务的 长度 最大 只能是线程数加缓冲数 其他的阻塞看不到
func (fo *Fork) Len() int {
	return len(fo.buf) + len(fo.max)
}

// 加入闭包 可以的话立即执行
func (fo *Fork) Push(f func()) {
	// 如果未达到最大线程则开启新线程执行 缓冲已满就在这里阻塞
	select {
	case fo.max <- none:
		go fo.fork(f)
	default:
		fo.buf <- f
	}
	return
}

// 加入闭包
func (fo *Fork) PushMerge(f func()) {
	// 如果缓冲已满就在这里阻塞
	fo.buf <- f
	return
}

// 把当前线程加入 线程执行队列
func (fo *Fork) forkMerge() {
	select {
	case fo.max <- none:
		fo.fork(nil)
	}
}

// 把当前线程加入 线程执行队列
func (fo *Fork) fork(f0 func()) {
	if f0 != nil {
		f0()
	}

loop:
	for {
		select {
		case f := <-fo.buf:
			f()
		default:
			if len(fo.buf) == 0 {
				break loop
			}
		}
	}

	fo.forkExit()
}

// 线程结束信号
func (fo *Fork) forkExit() {
	<-fo.max
	select {
	case fo.sub <- none:
	default:
	}
	return
}

// 等待所有线程结束在返回
func (fo *Fork) Join() {
	fo.join(false)
	return
}

// 等待所有线程结束在返回 把当前线程加入线程执行队列
func (fo *Fork) JoinMerge() {
	fo.join(true)
	return
}

func (fo *Fork) join(merge bool) {
	for {
		if len(fo.max) == 0 {
			if len(fo.buf) == 0 {
				return
			}
			if !merge {
				fo.max <- none
				go fo.fork(nil)
			}
		}
		if merge {
			fo.forkMerge()
		}
		select {
		case <-fo.sub:
		case <-time.After(time.Second):
		}
	}
	return
}
