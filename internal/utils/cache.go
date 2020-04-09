package utils

import (
	"sync"
)

type Cache interface {
	GetBaseUrlCount(string) int
	AlreadySeen(string) bool
	GetRobotsPatterns(string) ([]string, bool)
	SetRobotsPatterns(string, []string)
}

type UrlCache struct {
	BaseMem   map[string]int
	ExactMem  map[string]bool
	RobotsMem map[string][]string
	mux       sync.Mutex
}

func NewUrlCache() *UrlCache {
	return &UrlCache{
		BaseMem:   make(map[string]int),
		ExactMem:  make(map[string]bool),
		RobotsMem: make(map[string][]string),
	}
}

func (v *UrlCache) GetRobotsPatterns(baseurl string) ([]string, bool) {
	v.mux.Lock()
	defer v.mux.Unlock()
	patterns, exists := v.RobotsMem[baseurl]
	return patterns, exists
}

func (v *UrlCache) SetRobotsPatterns(baseurl string, patterns []string) {
	v.mux.Lock()
	defer v.mux.Unlock()
	v.RobotsMem[baseurl] = patterns
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
