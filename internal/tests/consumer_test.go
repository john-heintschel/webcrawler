package tests

import (
	"fmt"
	"testing"

	"github.com/john-heintschel/webcrawler/internal/utils"
    "github.com/stretchr/testify/assert"
)

const TEST_HTMLDOC_ONE string = "<p>Links:</p><ul><li><a href=\"http://www.google.com/foo.html\">Foo</a><li><a href=\"foo.html\">Foo</a><li><a href=\"/bar/baz.html\">BarBaz</a></ul>"
const TEST_HTMLDOC_TWO string = "<p>Links:</p><ul><li><a href=\"http://www.google.com/foo.html\">Foo</a><li><a href=\"http://www.linktest.com/foo.html\">Foo</a><li><a href=\"/bar/baz.html\">BarBaz</a></ul>"

type mockClient struct{}

func (v mockClient) GetDocument(input string) (string, error) {
	switch input {
	case "http://www.selflink.com":
		return "", nil
	case "http://www.linktest.com":
		return TEST_HTMLDOC_ONE, nil
	case "http://www.google.com/foo.html":
		return TEST_HTMLDOC_TWO, nil
	case "error":
		return "", fmt.Errorf("500 error returned from url %v", input)
	default:
		return "", nil
	}
}

type mockCache struct {
	Count int
}

func (v *mockCache) Get(input string) int {
	return v.Count
}
func (v *mockCache) Increment(string) {}

type indexArg struct {
	url     string
	htmldoc string
}

func TestConsumerDepth(t *testing.T) {

	// test that consumer doesn't go beyond the depth its supposed to
	httpclient := mockClient{}
	cache := &mockCache{Count: 0}
	urls := make(chan string, 10)
	urls <- "http://www.linktest.com"

	testConsumer := utils.UrlConsumer{
		MaxDepth:           1,
		MaxRequestsPerHost: 5,
		UrlCache:           cache,
		HttpClient:         httpclient,
	}
	var indexOutput []indexArg
	mockIndexer := func(url string, input string) { indexOutput = append(indexOutput, indexArg{url, input}) }

	testConsumer.Consume(urls, mockIndexer)
	assert.Equal(t, indexOutput[0], indexArg{url: "http://www.linktest.com", htmldoc: TEST_HTMLDOC_ONE})
	assert.Equal(t, <-urls, "http://www.google.com/foo.html")
	assert.Equal(t, <-urls, "http://www.linktest.com/foo.html")
	assert.Equal(t, <-urls, "http://www.linktest.com/bar/baz.html")
}

func TestConsumerMaxRequestsPerDomain(t *testing.T) {
	// test that the consumer doesn't go beyond the max requests per domain
	httpclient := mockClient{}
	cache := &utils.UrlCache{Mem: make(map[string]int, 0)}
	urls := make(chan string, 10)
	urls <- "http://www.linktest.com"

	testConsumer := utils.UrlConsumer{
		MaxDepth:           3,
		MaxRequestsPerHost: 1,
		UrlCache:           cache,
		HttpClient:         httpclient,
	}
	var indexOutput []indexArg
	mockIndexer := func(url string, input string) { indexOutput = append(indexOutput, indexArg{url, input}) }

	testConsumer.Consume(urls, mockIndexer)
	assert.Equal(t, indexArg{url: "http://www.linktest.com", htmldoc: TEST_HTMLDOC_ONE}, indexOutput[0])
	assert.Equal(t, indexArg{url: "http://www.google.com/foo.html", htmldoc: TEST_HTMLDOC_TWO}, indexOutput[1])
	assert.Equal(t, 2, len(indexOutput))
}

// test that the consumer continues if it gets an error fetching a page 500
