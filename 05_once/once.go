package once

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	statChill uint32 = iota
	statRun
	statDone
)

type Once struct {
	state uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.state) == statDone {
		return
	}
	if atomic.CompareAndSwapUint32(&o.state, statChill, statRun) {
		defer func() {
			atomic.StoreUint32(&o.state, statDone)
			futex.WakeAll(&o.state)
		}()
		f()
		return
	}

	for {
		if atomic.LoadUint32(&o.state) == statDone {
			return
		}
		futex.Wait(&o.state, statRun)
	}
}

func (o *Once) Done() bool {
	return atomic.LoadUint32(&o.state) == statDone
}
