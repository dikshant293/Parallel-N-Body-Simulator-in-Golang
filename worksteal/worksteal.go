package worksteal

import (
	"math/rand"
	"proj3-redesigned/barrier"
	"proj3-redesigned/bdequeue"
	"proj3-redesigned/qTree"
	"sync"
	"sync/atomic"
)

type update_param func(start int, end int, pArray []qTree.Particle, root *qTree.QTNode, dt float64, theta float64, bar *barrier.Barrier)
type move_param func(start int, end int, pArray []qTree.Particle, dt float64, bar *barrier.Barrier)

type MyStruct struct {
	Idx    int
	PArray *[]qTree.Particle
	Root   *qTree.QTNode
	Dt     float64
	Theta  float64
	Bar    *barrier.Barrier
}

type WorkStealingThread struct {
	QueueArray []*bdequeue.Queue
	Threshold  int
	Lk         sync.Mutex
	Count      int64
}

func NewWorkStealingThread(q []*bdequeue.Queue, t int) *WorkStealingThread {
	return &WorkStealingThread{QueueArray: q, Threshold: t, Count: 0}
}

func (w *WorkStealingThread) Run(me int, str MyStruct, up update_param, move move_param, workArr *[]int64, checkArr *[]int64) {
	for {
		if w.Count >= int64(len(*str.PArray)) {
			str.Bar.Wait()
			return
		}
		var r *int64 = w.QueueArray[me].PopBottom()
		if r != nil && atomic.CompareAndSwapInt64(&(*checkArr)[*r], 0, 1) {
			if str.Root != nil {
				up(int(*r), int(*r+1), *str.PArray, str.Root, str.Dt, str.Theta, nil)
			} else {
				move(int(*r), int(*r+1), *str.PArray, str.Dt, nil)
			}
			atomic.AddInt64(&w.Count, 1)
		}
		size := w.QueueArray[me].Size()
		if size == 0 && str.Root != nil {
			victim := rand.Intn(len(w.QueueArray))
			w.Lk.Lock()
			w.Balance(w.QueueArray[me], w.QueueArray[victim], me, victim)
			w.Lk.Unlock()
		}
	}
}

func (w *WorkStealingThread) Balance(q0, q1 *bdequeue.Queue, me int, victim int) {
	var qMin, qMax *bdequeue.Queue
	if q0.Size() < q1.Size() {
		qMin, qMax = q0, q1
	} else {
		qMin, qMax = q1, q0
	}
	if qMax.Size()-qMin.Size() > w.Threshold {
		for qMax.Size() > qMin.Size() {
			qMin.PushBottom(qMax.PopTop())
		}
	}
}
