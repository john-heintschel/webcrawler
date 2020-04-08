package tests

import (
	"testing"

	"github.com/john-heintschel/webcrawler/internal/utils"
)

func TestCache(t *testing.T) {
	cache := utils.NewUrlCache()
	testVal := cache.GetBaseUrlCount("hello")
	if testVal != 0 {
		t.Errorf(
			"cache GetBaseUrlCount failed, expected %v got %v",
			0,
			testVal,
		)
	}
	testVal = cache.GetBaseUrlCount("hello")
	if testVal != 1 {
		t.Errorf(
			"cache GetBaseUrlCount failed, expected %v got %v",
			1,
			testVal,
		)
	}

	for i := 0; i < 10; i++ {
		cache.GetBaseUrlCount("hello")
		cache.GetBaseUrlCount("there")
	}

	if cache.GetBaseUrlCount("hello") != 12 ||
		cache.GetBaseUrlCount("there") != 10 {
		t.Errorf("cache Increment failed")
	}

}
