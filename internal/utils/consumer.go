package utils

import (
	"fmt"
	"time"
	"log"
	"net/url"
)


type UrlConsumer struct {
	MaxDepth           int
	MaxRequestsPerHost int
	UrlCache           Cache
	HttpClient         Client
}

type UrlItem struct{
    Url string
	Depth int
}

func (v *UrlConsumer) Consume(urls chan UrlItem, indexer func(uri, input string)) {
	docParser := PageParser{}
    var uri string
	var depth int
	timeout := false
	for {
		select{
			case urlitem := <-urls:
				uri = urlitem.Url
				depth = urlitem.Depth
			case <- time.After(time.Second):
                fmt.Println("timeout!")
                timeout = true
        }
		time.Sleep(time.Second)
		fmt.Println("%v", uri)
		if timeout {
			break
		}
		if v.UrlCache.AlreadySeen(uri){
		    fmt.Println("already seen uri")
            continue
	    }
		u, err := url.Parse(uri)
		if err != nil {
		    fmt.Println("html parsing error")
			log.Fatal(err)
			continue
		}
		baseurl := fmt.Sprintf("%v://%v/", u.Scheme, u.Host)
		cachedBaseUrlCount := v.UrlCache.GetBaseUrlCount(baseurl)
		if cachedBaseUrlCount >= v.MaxRequestsPerHost{
			fmt.Println("more than max requests per host")
			continue
	    }
		htmldoc, err := v.HttpClient.GetDocument(uri)
		if err != nil {
		    fmt.Println("http error")
			log.Fatal(err)
			continue
		}
		indexer(uri, htmldoc)
		fmt.Println("ran indexer")
		if depth +1 >= v.MaxDepth{
            fmt.Println("max depth reached")
		    continue
		}
		links := docParser.GetLinks(htmldoc, baseurl)
		for _, link := range links {
			urls <- UrlItem{Url: link, Depth: depth+1}
		}
		fmt.Println("last line reached")
	}
}
