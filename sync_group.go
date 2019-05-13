package std

import (
	"sync"
	"time"
)

type SyncId string

type SyncGroup struct {
	hangs map[SyncId]CoSync
	lock  sync.Locker
}

func NewSyncGroup() *SyncGroup {
	ret := &SyncGroup{
		hangs: make(map[SyncId]CoSync),
		lock:  &sync.Mutex{},
	}
	return ret
}

func (this *SyncGroup) FinishCoSyncWithData(id SyncId, data interface{}) {
	defer this.lock.Unlock()
	this.lock.Lock()
	if c, ok := this.hangs[id]; ok {
		c.DoneWithData(data)
	}
	delete(this.hangs, id)
}

func (this *SyncGroup) AddCoSync(id SyncId, c CoSync, timeout time.Duration) {
	defer this.lock.Unlock()
	this.lock.Lock()
	this.hangs[id] = c;
	if timeout <= 0 {
		return
	}
	go func() {
		time.Sleep(timeout)
		this.FinishCoSyncWithData(id, nil)
	}()
}
