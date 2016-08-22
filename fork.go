package fork

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

// 加入闭包
func (fo *Fork) Puah(f func()) {
	// 如果缓冲已满就在这里阻塞
	fo.buf <- f

	// 如果未达到最大线程则开启新线程执行
	select {
	case fo.max <- none:
		go fo.fork()
	default:
	}
	return
}

// 开启一个新的线程执行闭包 没有闭包时自动结束
func (fo *Fork) fork() {
	for {
		select {
		case f := <-fo.buf:
			f()
		default:
			<-fo.max
			fo.forkExit()
			return
		}
	}
}

// 线程结束信号
func (fo *Fork) forkExit() {
	select {
	case fo.sub <- none:
	default:
	}
	return
}

// 等待所有线程结束在返回
func (fo *Fork) Join() {
	for {
		if len(fo.max) == 0 {
			return
		}
		select {
		case <-fo.sub:
		}
	}
	return
}
