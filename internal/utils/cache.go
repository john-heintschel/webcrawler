package utils

import (
	"sync"
)

type Cache interface {
	Get(string) int // TODO change this to include depth in search tree as well as frequency
	Increment(string)
}

type UrlCache struct {
	Mem map[string]int
	mux sync.Mutex
}

func (v *UrlCache) Get(url string) int {
	v.mux.Lock()
	defer v.mux.Unlock()
	if val, exists := v.Mem[url]; exists {
		return val
	} else {
		return 0
	}
}

func (v *UrlCache) Increment(url string) {
	v.mux.Lock()
	if _, exists := v.Mem[url]; exists {
		v.Mem[url]++
	} else {
		v.Mem[url] = 1
	}
	v.mux.Unlock()
}
