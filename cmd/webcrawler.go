package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/john-heintschel/webcrawler/internal/utils"
)

func indexer(searchString string) func(string, string) {
	return func(uri string, htmldoc string) {
		// TODO replace this with more interesting indexing
		hasString := strings.Contains(htmldoc, searchString)
		fmt.Printf(
			"url: %v contains the string \"%v\": %v\n",
			uri,
			searchString,
			hasString,
		)
	}
}

func main() {
	var concurrency int
	fmt.Print("Specify max concurrency: ")
	_, err := fmt.Scanf("%d", &concurrency)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup

	var searchTxt string
	fmt.Print("Enter exact string to search: ")
	_, err = fmt.Scanln(&searchTxt)
	if err != nil {
		log.Fatal(err)
	}

	var website string
	fmt.Print("Enter website to start search: ")
	_, err = fmt.Scanln(&website)
	if err != nil {
		log.Fatal(err)
	}

	cache := utils.NewUrlCache()

	urls := make(chan utils.UrlItem, 100000)
	urls <- utils.UrlItem{Url: website, Depth: 0}

	indexFunction := indexer(searchTxt)

	for idx := 1; idx <= concurrency; idx++ {
		wg.Add(1)
		consumer := utils.NewUrlConsumer(2, 1, cache)
		go consumer.Consume(urls, indexFunction, &wg)
	}
	wg.Wait()
}
