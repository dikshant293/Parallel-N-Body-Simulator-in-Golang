package bdequeue

import (
	"fmt"
	"sync/atomic"
)

type Queue struct {
	Arr    []*int64
	Bottom int64
	Top    int64
	Shift  int64
}

func NewBDEQueue(capacity int) *Queue {
	return &Queue{Arr: make([]*int64, capacity), Bottom: 0, Top: 0, Shift: 32}
}

func (q *Queue) PushBottom(i *int64) {
	q.Arr[q.Bottom] = i
	atomic.AddInt64(&q.Bottom, 1)
}

func (q *Queue) IsEmpty() bool {
	return (q.Top >> q.Shift) >= q.Bottom
}

func (q *Queue) PopTop() *int64 {
	var oldTop, oldStamp, newTop, newStamp int64
	oldTop = q.Top >> q.Shift
	newTop = oldTop + 1
	oldStamp = ^(^int64(0) << q.Shift) & q.Top
	newStamp = oldStamp + 1
	if q.Bottom <= oldTop {
		return nil
	}
	r := q.Arr[oldTop]
	if atomic.CompareAndSwapInt64(&q.Top, (oldTop<<q.Shift)|oldStamp, (newTop<<q.Shift)|newStamp) {
		return r
	}
	return nil
}

func (q *Queue) PopBottom() *int64 {
	if q.Bottom == 0 {
		return nil
	}
	atomic.AddInt64(&q.Bottom, -1)
	r := q.Arr[q.Bottom]
	var oldTop, oldStamp, newTop, newStamp int64
	oldTop = q.Top >> q.Shift
	newTop = 0
	oldStamp = ^(^int64(0) << q.Shift) & q.Top
	newStamp = oldStamp + 1
	if q.Bottom > oldTop {
		return r
	}
	if q.Bottom == oldTop {
		q.Bottom = 0
		if atomic.CompareAndSwapInt64(&q.Top, (oldTop<<q.Shift)|oldStamp, (newTop<<q.Shift)|newStamp) {
			return r
		}
	}
	atomic.StoreInt64(&q.Top, (newTop<<q.Shift)|newStamp)
	return nil
}

func (q *Queue) PrintArr() {
	for a, b := range q.Arr {
		if b != nil {
			fmt.Print("{", a, *b, "} ")
		} else {
			fmt.Print("{", a, b, "} ")
		}
	}
	fmt.Println()
}

func (q *Queue) Size() int {
	return int(q.Bottom - (q.Top >> q.Shift))
}
