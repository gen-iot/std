package std

import "sync/atomic"

type cwChanType chan int

type CoSync interface {
	Done()
	DoneWithData(interface{})
	Wait()
	WaitWithData() interface{}
}

type CloseWaiterSync struct {
	msk  int32       // atomic op
	c    cwChanType  // 同步chan
	data interface{} // 存放交换的数据
}

func NewCloseWaiterSync() CoSync {
	return &CloseWaiterSync{
		msk: 0,
		c:   make(cwChanType),
	}
}

func (this *CloseWaiterSync) GetWaitChan() <-chan int {
	return this.c
}

//等待
func (this *CloseWaiterSync) Wait() {
	_ = this.WaitWithData()
}

//等待,同时返回同步的数据
func (this *CloseWaiterSync) WaitWithData() interface{} {
	_ = <-this.c
	return this.data
}

//完成等待
func (this *CloseWaiterSync) Done() {
	this.DoneWithData(nil)
}

//完成等待,同时传递数据.只允许调用一次,多次调用,只生效一次
func (this *CloseWaiterSync) DoneWithData(data interface{}) {
	if atomic.CompareAndSwapInt32(&this.msk, 0, 1) {
		this.data = data
		close(this.c)
	}
}
