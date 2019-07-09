package std

import (
	"runtime"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool("tester", runtime.NumCPU())
	pool.Submit(func() {
		time.Sleep(time.Second * 5)
		_ = pool.Close()
	})
	pool.Execute()
}
