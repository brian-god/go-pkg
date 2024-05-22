package goroutine

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync/atomic"
)

// Panic 子协程 panic 会被重新包装，添加调用栈信息
type Panic struct {
	R     interface{} // recover() 返回值
	Stack []byte      // 当时的调用栈
}

func (p Panic) String() string {
	return fmt.Sprintf("%v\n%s", p.R, p.Stack)
}

type PanicGroup struct {
	panics chan Panic // 协程 panic 通知信道
	dones  chan int   // 协程完成通知信道
	jobN   int32      // 协程并发数量
}

func NewPanicGroup() *PanicGroup {
	return &PanicGroup{
		panics: make(chan Panic, 8),
		dones:  make(chan int, 8),
	}
}

func (g *PanicGroup) Go(f func()) *PanicGroup {
	atomic.AddInt32(&g.jobN, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.panics <- Panic{R: r, Stack: debug.Stack()}
				return
			}
			g.dones <- 1
		}()
		f()
	}()

	return g // 方便链式调用
}

func (g *PanicGroup) Wait(ctx context.Context) error {
	for {
		select {
		case <-g.dones:
			if atomic.AddInt32(&g.jobN, -1) == 0 {
				return nil
			}
		case p := <-g.panics:
			panic(p)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
