package tests

import (
	"fmt"
	"sync"
	"testing"

	"github.com/john-heintschel/webcrawler/internal/utils"
	"github.com/stretchr/testify/assert"
)

const TEST_HTMLDOC_ONE string = `
    <p>Links:</p>
        <ul>
            <li><a href="http://www.google.com/foo.html">Foo</a>
        	<li><a href="foo.html">Foo</a>
        	<li><a href="/bar/baz.html">BarBaz</a>
    	</ul>
`
const TEST_HTMLDOC_TWO string = `
    <p>Links:</p>
	<ul>
        <li><a href="http://www.google.com/foo.html">Foo</a>
    	<li><a href="http://www.linktest.com/foo.html">Foo</a>
    	<li><a href="/bar/baz.html">BarBaz</a>
	</ul>
`
const TEST_HTMLDOC_THREE string = `
    <p>Links:</p>
	<ul>
        <li><a href="http://www.robotstest.com/foo.html">Foo</a>
    	<li><a href="http://www.linktest.com/foo.html">Foo</a>
    	<li><a href="/bar/baz.html">BarBaz</a>
	</ul>
`

type mockClient struct{}

func (v mockClient) GetDocument(input string) (string, error) {
	switch input {
	case "http://www.selflink.com":
		return "", nil
	case "http://www.linktest.com":
		return TEST_HTMLDOC_ONE, nil
	case "http://www.google.com/foo.html":
		return TEST_HTMLDOC_TWO, nil
	case "http://www.robotstest.com":
		return TEST_HTMLDOC_THREE, nil
	case "http://www.robotstest.com/robots.txt":
		return "User-agent: *\nDisallow: /", nil
	case "error":
		return "", fmt.Errorf("500 error returned from url %v", input)
	default:
		return "", nil
	}
}

type indexArg struct {
	url     string
	htmldoc string
}

func TestConsumerDepth(t *testing.T) {
	// test that consumer doesn't go beyond the depth its supposed to
	cache := utils.NewUrlCache()
	urls := make(chan utils.UrlItem, 10)
	urls <- utils.UrlItem{Url: "http://www.linktest.com", Depth: 0}
	var wg sync.WaitGroup
	wg.Add(1)

	testConsumer := utils.NewUrlConsumer(1, 5, cache)
	testConsumer.HttpClient = mockClient{}
	var indexOutput []indexArg
	mockIndexer := func(url string, input string) {
		indexOutput = append(indexOutput, indexArg{url, input})
	}

	testConsumer.Consume(urls, mockIndexer, &wg)
	assert.Equal(
		t,
		indexOutput[0],
		indexArg{url: "http://www.linktest.com", htmldoc: TEST_HTMLDOC_ONE},
	)
}

func TestConsumerMaxRequestsPerDomain(t *testing.T) {
	// test that the consumer doesn't go beyond the max requests per domain
	cache := utils.NewUrlCache()
	urls := make(chan utils.UrlItem, 10)
	urls <- utils.UrlItem{Url: "http://www.linktest.com", Depth: 0}
	var wg sync.WaitGroup
	wg.Add(1)

	testConsumer := utils.NewUrlConsumer(3, 1, cache)
	testConsumer.HttpClient = mockClient{}
	var indexOutput []indexArg
	mockIndexer := func(url string, input string) {
		indexOutput = append(indexOutput, indexArg{url, input})
	}

	testConsumer.Consume(urls, mockIndexer, &wg)
	assert.Equal(
		t,
		indexArg{
			url:     "http://www.linktest.com",
			htmldoc: TEST_HTMLDOC_ONE,
		},
		indexOutput[0],
	)
	assert.Equal(
		t,
		indexArg{
			url:     "http://www.google.com/foo.html",
			htmldoc: TEST_HTMLDOC_TWO,
		},
		indexOutput[1],
	)
	assert.Equal(t, 2, len(indexOutput))
}

func TestConsumerObeysRobotsTxt(t *testing.T) {
	cache := utils.NewUrlCache()
	urls := make(chan utils.UrlItem, 10)
	urls <- utils.UrlItem{Url: "http://www.robotstest.com", Depth: 0}
	var wg sync.WaitGroup
	wg.Add(1)

	testConsumer := utils.NewUrlConsumer(2, 3, cache)
	testConsumer.HttpClient = mockClient{}
	var indexOutput []indexArg
	mockIndexer := func(url string, input string) {
		indexOutput = append(indexOutput, indexArg{url, input})
	}

	testConsumer.Consume(urls, mockIndexer, &wg)
	assert.Equal(t, 0, len(indexOutput))
}
