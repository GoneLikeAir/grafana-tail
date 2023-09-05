package tail

import (
	"sync"
	"testing"
	"time"
)

func TestMemoryPool(t *testing.T) {
	p := NewMemoryPool(100)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 150; i++ {
			p.Request(70)
			time.Sleep(100)
			p.Return(15)
			time.Sleep(100)
			p.Return(45)
			time.Sleep(100)
			p.Return(10)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			p.Request(80)
			time.Sleep(100)
			p.Return(21)
			time.Sleep(100)
			p.Return(21)
			time.Sleep(100)
			p.Return(38)
		}
	}()
	wg.Wait()
}
