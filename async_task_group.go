package std

import "sync"

// best practice for wait_group
type AsyncTaskGroup struct {
	wg *sync.WaitGroup
}

func NewAsyncTaskGroup() *AsyncTaskGroup {
	return NewAsyncTaskGroup2(new(sync.WaitGroup))
}

func NewAsyncTaskGroup2(wg *sync.WaitGroup) *AsyncTaskGroup {
	Assert(wg != nil, "waitGroup is nil")
	return &AsyncTaskGroup{wg: wg}
}

func (this *AsyncTaskGroup) AddTask(exe func()) {
	Assert(exe != nil, "exe cant be nil")
	this.wg.Add(1)
	go func() {
		defer this.wg.Done()
		exe()
	}()
}

func (this *AsyncTaskGroup) Wait() {
	this.wg.Wait()
}
