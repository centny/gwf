package test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCond(t *testing.T) {
	lck_ := &sync.RWMutex{}
	lck := sync.NewCond(lck_)
	for i := 0; i < 5; i++ {
		go func(iv int) {
			fmt.Println("a->", iv, "xxx")
			lck.L.Lock()
			lck.Wait()
			fmt.Println("a->", iv, "...")
			lck.L.Unlock()
		}(i)
	}
	time.Sleep(3 * time.Second)
	for i := 0; i < 5; i++ {
		lck.L.Lock()
		lck.Broadcast()
		lck.L.Unlock()
	}
	time.Sleep(time.Second)
}
