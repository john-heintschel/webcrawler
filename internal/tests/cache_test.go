package tests

import (
	"testing"
	"time"

	"github.com/john-heintschel/webcrawler/internal/utils"
)

func TestCache(t *testing.T) {
	cache := utils.UrlCache{BaseMem: make(map[string]int, 0), ExactMem: make(map[string]bool, 0)}
	testVal := cache.GetBaseUrlCount("hello")
	if testVal != 0 {
		t.Errorf("cache GetBaseUrlCount failed, expected %v got %v", 0, testVal)
	}
	testVal = cache.GetBaseUrlCount("hello")
	if testVal != 1 {
		t.Errorf("cache GetBaseUrlCount failed, expected %v got %v", 1, testVal)
	}

	for i := 0; i < 10000; i++ {
		go cache.GetBaseUrlCount("hello")
		go cache.GetBaseUrlCount("there")
	}
	time.Sleep(time.Second)

	if cache.GetBaseUrlCount("hello") != 10002 || cache.GetBaseUrlCount("there") != 10000 {
		t.Errorf("cache Increment failed")
	}

}
