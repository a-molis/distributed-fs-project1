package atomic_data

import "sync/atomic"

type AtomicInt struct {
	num *int32
}

func NewAtomicInt(num *int32) *AtomicInt {
	return &AtomicInt{num: num}
}

func (atomicInt *AtomicInt) Increment() int32 {
	return atomic.AddInt32(atomicInt.num, 1)
}

func (atomicInt *AtomicInt) Get() int32 {
	return atomic.LoadInt32(atomicInt.num)
}

func (atomicInt *AtomicInt) Decrement() int32 {
	return atomic.AddInt32(atomicInt.num, -1)
}
