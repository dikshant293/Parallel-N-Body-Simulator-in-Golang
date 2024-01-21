package barrier

import (
	"sync"
	"sync/atomic"
)

type Barrier struct {
	mutex  sync.Mutex
	cond   *sync.Cond
	count  int32
	target int32
}

func NewBarrier(n int32) *Barrier {
	b := &Barrier{count: 0, target: n}
	b.cond = sync.NewCond(&b.mutex)
	return b
}

func (b *Barrier) Wait() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.count++
	if b.count < b.target {
		b.cond.Wait()
	} else {
		b.cond.Broadcast()
	}
}

func (b *Barrier) Await() {
	for i := 0; i < int(b.target); i++ {
		b.Wait()
	}
	atomic.StoreInt32(&b.count, 0)
}
