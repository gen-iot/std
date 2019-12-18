package std

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	cancel, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	pool := NewWorkerPool("tester", runtime.NumCPU())
	pool.Submit(func() {
		time.Sleep(time.Second * 5)
		cancelFunc()
	})
	pool.Execute(cancel)
}
