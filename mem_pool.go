package tail

import (
	"fmt"
	"sync"
	"time"
)

type Pool interface {
	Request(n int64)
	Return(n int64)
}

type MemoryPool struct {
	limit   int64
	current int64
	mutex   sync.Mutex
	cond    *sync.Cond
}

func NewMemoryPool(limitBytes int64) Pool {
	p := &MemoryPool{
		limit:   limitBytes,
		current: 0,
		mutex:   sync.Mutex{},
		cond:    sync.NewCond(&sync.Mutex{}),
	}
	go p.gracefulReturn()
	return p
}

func (p *MemoryPool) Request(n int64) {
	fmt.Printf("start request %d bytes, current total requested %d\n", n, p.current)

	for {
		p.mutex.Lock()
		if p.current+n > p.limit {
			p.mutex.Unlock()

			p.cond.L.Lock()
			p.cond.Wait()
			p.cond.L.Unlock()
		} else {
			p.current += n
			p.mutex.Unlock()
			fmt.Printf("request %d bytes, current total requested %d\n", n, p.current)
			break
		}
	}
}

func (p *MemoryPool) Return(n int64) {
	fmt.Printf("start return %d bytes, current total requested %d\n", n, p.current)
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.current -= n
	if p.current < 0 {
		p.current = 0
	}
	p.cond.Broadcast()
	fmt.Printf("return %d bytes, current total requested %d\n", n, p.current)
}

func (p *MemoryPool) gracefulReturn() {
	tick := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-tick.C:
			needReturn := false
			p.mutex.Lock()
			if float64(p.current)/float64(p.limit) > 0.8 {
				needReturn = true
			}
			p.mutex.Unlock()
			if needReturn {
				p.Return(p.limit / 5 / 60)
			}
		}
	}
}
