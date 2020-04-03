package utils

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

var SEARCH_STRING = "wham" // TODO this should be taken in as argument

type UrlConsumer struct {
	MaxDepth           int
	MaxRequestsPerHost int
	UrlCache           Cache
	HttpClient         Client
}

func index(uri string, htmldoc string) {
	// TODO replace this with more interesting indexing
	hasString := strings.Contains(htmldoc, SEARCH_STRING)
	fmt.Printf("url: %v contains the string \"%v\": %v\n", uri, SEARCH_STRING, hasString)
}

type UrlItem struct {
	Url   string
	Depth int
}

func (v *UrlConsumer) Consume(urls chan UrlItem, indexer func(uri, input string)) {
	docParser := PageParser{}
	var uri string
	var depth int
	timeout := false
	for {
		select {
		case urlitem := <-urls:
			uri = urlitem.Url
			depth = urlitem.Depth
		case <-time.After(time.Second):
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
			log.Fatal(err)
			continue
		}
		baseurl := fmt.Sprintf("%v://%v/", u.Scheme, u.Host)
		cachedBaseUrlCount := v.UrlCache.GetBaseUrlCount(baseurl)
		if cachedBaseUrlCount >= v.MaxRequestsPerHost {
			continue
		}
		htmldoc, err := v.HttpClient.GetDocument(uri)
		if err != nil {
			log.Fatal(err)
			continue
		}
		indexer(uri, htmldoc)
		if depth+1 >= v.MaxDepth {
			continue
		}
		links := docParser.GetLinks(htmldoc, baseurl)
		for _, link := range links {
			urls <- UrlItem{Url: link, Depth: depth + 1}
		}
	}
}
