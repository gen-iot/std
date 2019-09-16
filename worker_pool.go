package std

import (
	"container/list"
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

type WorkerFunc = func()

type WorkerPool struct {
	name      string
	cond      *sync.Cond
	workN     int
	workQ     *list.List
	startFlag int32
}

func NewWorkerPool(name string, num int) *WorkerPool {
	if num <= 0 {
		num = 1
	}
	return &WorkerPool{
		name:      name,
		cond:      sync.NewCond(new(sync.Mutex)),
		workN:     MinInt(num, runtime.NumCPU()),
		workQ:     list.New(),
		startFlag: 0,
	}
}

func (this *WorkerPool) Submit(f WorkerFunc) {
	Assert(f != nil, fmt.Sprintf("worker pool [%s],submit nil task", this.name))
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	this.workQ.PushBack(f)
	this.cond.Broadcast()
}

func (this *WorkerPool) workerTask(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		this.cond.L.Lock()
		if this.workQ.Len() == 0 {
			this.cond.Wait()
		}
		if this.workQ.Len() == 0 {
			this.cond.L.Unlock()
			continue
		}
		front := this.workQ.Front()
		this.workQ.Remove(front)
		this.cond.L.Unlock()
		fn := front.Value.(WorkerFunc)
		fn()
	}
}

func (this *WorkerPool) GetName() string {
	return this.name
}

func (this *WorkerPool) Execute(ctx context.Context) {
	Assert(atomic.CompareAndSwapInt32(&this.startFlag, 0, 1),
		fmt.Sprintf("worker pool [%s] start multiple times...", this.name))
	group := NewAsyncTaskGroup()
	defer group.Wait()
	for i := 0; i < this.workN; i++ {
		group.AddTask(func() {
			this.workerTask(ctx)
		})
	}
	group.AddTask(func() {
		// watcher to unlock,if ctx done
		<-ctx.Done()
		//
		this.cond.L.Lock()
		defer this.cond.L.Unlock()
		// drop all task
		for this.workQ.Len() != 0 {
			this.workQ.Remove(this.workQ.Front())
		}
		// send exit signal
		this.cond.Broadcast()
	})
}
