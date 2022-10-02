package atomic_data

import "sync"

type WaitNotify struct {
	Lock *sync.Mutex
	Cond *sync.Cond
}

func NewWaitNotify() *WaitNotify {
	waitNotify := &WaitNotify{}
	waitNotify.Lock = &sync.Mutex{}
	waitNotify.Cond = sync.NewCond(waitNotify.Lock)
	return waitNotify
}
