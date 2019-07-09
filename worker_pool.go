package std

import (
	"container/list"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

type WorkerFunc = func()

type WorkerPool struct {
	name      string
	cond      *sync.Cond
	exitSig   chan bool
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
		exitSig:   make(chan bool),
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

func (this *WorkerPool) workerTask(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-this.exitSig:
			return
		default:
		}
		this.cond.L.Lock()
		if this.workQ.Len() == 0{
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

func (this *WorkerPool) Execute() {
	Assert(atomic.CompareAndSwapInt32(&this.startFlag, 0, 1),
		fmt.Sprintf("worker pool [%s] start multiple times...", this.name))
	wg := new(sync.WaitGroup)
	wg.Add(this.workN)
	for i := 0; i < this.workN; i++ {
		go this.workerTask(wg)
	}
	wg.Wait()
}

func (this *WorkerPool) Close() error {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	closed := false
	select {
	case <-this.exitSig:
		closed = true
	default:
		closed = false
	}
	if !closed {
		close(this.exitSig)
		this.cond.Broadcast()
	}
	return nil
}
