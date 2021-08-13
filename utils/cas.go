package utils

import "sync"

type AtomicReference struct {
	lock      *sync.Mutex
	reference interface{}
}

func NewAtomicReference(obj interface{}) *AtomicReference {
	return &AtomicReference{
		lock:      &sync.Mutex{},
		reference: obj,
	}
}

func (a *AtomicReference) CompareAndSet(old, new interface{}) bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	if old == a.reference {
		a.reference = new
		return true
	}
	return false
}

func (a AtomicReference) Get() interface{} {
	return a.reference
}

func (a *AtomicReference) Set(new interface{}) {
	a.reference = new
}
