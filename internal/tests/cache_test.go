package tests

import (
	"testing"
	"time"

	"github.com/john-heintschel/webcrawler/internal/utils"
)

func TestCache(t *testing.T) {
	cache := utils.UrlCache{Mem: make(map[string]int, 0)}
	testVal := cache.Get("hello")
	if testVal != 0 {
		t.Errorf("cache Get failed, expected %v got %v", 0, testVal)
	}
	cache.Increment("hello")
	testVal = cache.Get("hello")
	if testVal != 1 {
		t.Errorf("cache Get failed, expected %v got %v", 1, testVal)
	}

	for i := 0; i < 10000; i++ {
		go cache.Increment("hello")
		go cache.Increment("there")
	}
	time.Sleep(time.Second)

	if cache.Get("hello") != 10001 || cache.Get("there") != 10000 {
		t.Errorf("cache Increment failed")
	}

}
