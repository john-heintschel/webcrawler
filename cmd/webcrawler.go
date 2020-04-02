package main

import (
	"net/http"
	"time"

	"github.com/john-heintschel/webcrawler/internal/utils"
)

func main() {

	httpclient := utils.WebCrawlerClient{HTTPClient: &http.Client{Timeout: time.Duration(5) * time.Second}}
	cache := &utils.UrlCache{Mem: make(map[string]int, 0)}
	urls := make(chan string, 100)
	urls <- "http://www.wikipedia"

	utils.Consume(urls, cache, httpclient)
}
