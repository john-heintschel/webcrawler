package utils

import (
	"fmt"
	"log"
	"net/url"
	"strings"
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

func (v *UrlConsumer) Consume(urls chan string, indexer func(uri, input string)) {
	depth := 0 // TODO this needs to be implemented as depth and in the cache
	docParser := PageParser{}

	for {
		if depth >= v.MaxDepth {
			break
		}
		uri := <-urls

		u, err := url.Parse(uri)
		if err != nil {
			log.Fatal(err)
			continue
		}
		baseurl := fmt.Sprintf("%v://%v/", u.Scheme, u.Host)
		if count := v.UrlCache.Get(baseurl); count >= v.MaxRequestsPerHost {
			break
		} else {
			v.UrlCache.Increment(baseurl)
		}
		htmldoc, err := v.HttpClient.GetDocument(uri)
		if err != nil {
			log.Fatal(err)
			continue
		}
		indexer(uri, htmldoc)
		links := docParser.GetLinks(htmldoc, baseurl)

		for _, link := range links {
			urls <- link
		}
		depth++
	}
}
