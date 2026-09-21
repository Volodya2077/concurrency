package rwmutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const writer = 1 << 31

type RWMutex struct {
	state uint32
}

func (rw *RWMutex) RLock() {
	for {
		a := atomic.LoadUint32(&rw.state)
		if a&writer != 0 {
			futex.Wait(&rw.state, a)
			continue

		}
		if atomic.CompareAndSwapUint32(&rw.state, a, a+1) {
			return
		}
	}
}

func (rw *RWMutex) RUnlock() {
	for {
		a := atomic.LoadUint32(&rw.state)

		if a == 0 || a&writer != 0 {
			panic("чуть паники")
		}
		if atomic.CompareAndSwapUint32(&rw.state, a, a-1) {
			if a-1 == 0 {
				futex.WakeAll(&rw.state)
			}
			return
		}
	}
}

func (rw *RWMutex) Lock() {

	for {
		if atomic.CompareAndSwapUint32(&rw.state, 0, writer) {
			return
		}
		a := atomic.LoadUint32(&rw.state)
		if a != 0 {
			futex.Wait(&rw.state, a)
		}
	}
}

func (rw *RWMutex) Unlock() {
	if !atomic.CompareAndSwapUint32(&rw.state, writer, 0) {
		panic("тут тоже чуть паники")
	}
	futex.WakeAll(&rw.state)
}
