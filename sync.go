package po

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	done   uint32
	doneM  sync.Mutex
	doneC  chan struct{}
	doneCM sync.Mutex
	err    error
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

func (o *Once) ErrorDo(f func() error) error {
	o.Do(func() {
		o.err = f()
	})
	return o.err
}

func (o *Once) Wait() <-chan struct{} {
	o.doneCM.Lock()
	defer o.doneCM.Unlock()
	if o.doneC == nil {
		o.doneC = make(chan struct{})
	}
	if o.Done() {
		select {
		case <-o.doneC:
		default:
			close(o.doneC)
		}
	}
	return o.doneC
}

func (o *Once) Done() bool {
	return atomic.LoadUint32(&o.done) == 1
}

func (o *Once) doSlow(f func()) {
	o.doneM.Lock()
	defer o.doneM.Unlock()
	if o.done == 1 {
		return
	}
	defer func() {
		atomic.StoreUint32(&o.done, 1)
		o.Wait()
	}()
	f()
}
