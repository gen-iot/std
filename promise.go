package std

import (
	"github.com/pkg/errors"
	"sync/atomic"
	"time"
)

var ErrFutureTimeout = errors.New("future wait timeout")

type Future interface {
	Wait(timeout time.Duration) error
	WaitData(timeout time.Duration) (interface{}, error)
}

type Promise interface {
	Future
	Done(err error)
	DoneData(err error, data interface{})
	GetFuture() Future
}

type promise struct {
	msk  int32       // atomic op
	c    chan int    // 同步chan
	data interface{} // 存放交换的数据
	err  error
}

func NewPromise() Promise {
	return &promise{
		msk: 0,
		c:   make(chan int),
	}
}

// get future
func (this *promise) GetFuture() Future {
	return this
}

//等待
func (this *promise) Wait(timeout time.Duration) error {
	_, err := this.WaitData(timeout)
	return err
}

//等待,同时返回同步的数据
func (this *promise) WaitData(timeout time.Duration) (interface{}, error) {
	tm := time.NewTimer(timeout)
	defer tm.Stop()
	select {
	case <-tm.C:
		return nil, ErrFutureTimeout
	case <-this.c:
	}
	return this.data, this.err
}

//完成等待
func (this *promise) Done(err error) {
	this.DoneData(err, nil)
}

//完成等待,同时传递数据.只允许调用一次,多次调用,只生效一次
func (this *promise) DoneData(err error, data interface{}) {
	if atomic.CompareAndSwapInt32(&this.msk, 0, 1) {
		this.err = err
		this.data = data
		close(this.c)
	}
}
