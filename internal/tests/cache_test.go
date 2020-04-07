package tests

import (
	"testing"

	"github.com/john-heintschel/webcrawler/internal/utils"
)

func TestCache(t *testing.T) {
	cache := utils.UrlCache{BaseMem: make(map[string]int), ExactMem: make(map[string]bool)}
	testVal := cache.GetBaseUrlCount("hello")
	if testVal != 0 {
		t.Errorf("cache GetBaseUrlCount failed, expected %v got %v", 0, testVal)
	}
	testVal = cache.GetBaseUrlCount("hello")
	if testVal != 1 {
		t.Errorf("cache GetBaseUrlCount failed, expected %v got %v", 1, testVal)
	}

	for i := 0; i < 100; i++ {
		cache.GetBaseUrlCount("hello")
		cache.GetBaseUrlCount("there")
	}

	if cache.GetBaseUrlCount("hello") != 102 || cache.GetBaseUrlCount("there") != 100 {
		t.Errorf("cache Increment failed")
	}

}
