package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type UrlConsumer struct {
	MaxDepth           int
	MaxRequestsPerHost int
	UrlCache           Cache
	HttpClient         Client
	DocumentParser     PageParser
}

type UrlItem struct {
	Url   string
	Depth int
}

func NewUrlConsumer(
	maxDepth int,
	maxRequestsPerHost int,
	urlCache Cache,
) *UrlConsumer {
	httpclient := WebCrawlerClient{
		HTTPClient: &http.Client{Timeout: time.Duration(5) * time.Second},
	}
	docParser := PageParser{}
	return &UrlConsumer{
		MaxDepth:           maxDepth,
		MaxRequestsPerHost: maxRequestsPerHost,
		UrlCache:           urlCache,
		HttpClient:         httpclient,
		DocumentParser:     docParser,
	}
}

func (v *UrlConsumer) Consume(
	urls chan UrlItem,
	indexer func(uri, input string),
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	var uri string
	var depth int
	timeout := false
	for {
		select {
		case urlitem := <-urls:
			uri = urlitem.Url
			depth = urlitem.Depth
		case <-time.After(time.Second / 2):
			timeout = true
		}
		if timeout {
			break
		}
		if v.UrlCache.AlreadySeen(uri) {
			continue
		}
		u, err := url.Parse(uri)
		if err != nil {
			continue
		}
		baseurl := fmt.Sprintf("%v://%v/", u.Scheme, u.Host)
		cachedBaseUrlCount := v.UrlCache.GetBaseUrlCount(baseurl)
		if cachedBaseUrlCount >= v.MaxRequestsPerHost {
			continue
		}
		disallowed, exists := v.UrlCache.GetRobotsPatterns(baseurl)
		if !exists {
			robots, err := v.HttpClient.GetDocument(baseurl + "robots.txt")
			if err != nil {
				robots = ""
			}
			disallowed = v.DocumentParser.GetDisallowedPatterns(robots)
			v.UrlCache.SetRobotsPatterns(baseurl, disallowed)
		}
		if UrlIsForbidden(uri, disallowed) {
			continue
		}
		htmldoc, err := v.HttpClient.GetDocument(uri)
		if err != nil {
			continue
		}
		indexer(uri, htmldoc)
		if depth+1 >= v.MaxDepth {
			continue
		}
		links := v.DocumentParser.GetLinks(htmldoc, baseurl)
		for _, link := range links {
			urls <- UrlItem{Url: link, Depth: depth + 1}
		}
	}
}
