package utils

import (
	"sync"
)

type Cache interface {
	GetBaseUrlCount(string) int
	AlreadySeen(string) bool
}

type UrlCache struct {
	BaseMem  map[string]int
	ExactMem map[string]bool
	mux      sync.Mutex
}

func (v *UrlCache) AlreadySeen(url string) bool {
	// this function also marks as seen while it has the lock
	v.mux.Lock()
	defer v.mux.Unlock()
	value := v.ExactMem[url]
	v.ExactMem[url] = true
	return value
}

func (v *UrlCache) GetBaseUrlCount(url string) int {
	// this function also increments the counter while it has the lock
	v.mux.Lock()
	defer v.mux.Unlock()
	if val, exists := v.BaseMem[url]; exists {
		v.BaseMem[url] = val + 1
		return val
	} else {
		v.BaseMem[url] = 1
		return 0
	}
}
