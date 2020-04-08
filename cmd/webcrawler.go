package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/john-heintschel/webcrawler/internal/utils"
)

func indexer(searchString string) func(string, string) {
	return func(uri string, htmldoc string) {
		// TODO replace this with more interesting indexing
		hasString := strings.Contains(htmldoc, searchString)
		fmt.Printf("url: %v contains the string \"%v\": %v\n", uri, searchString, hasString)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter exact string to search: ")
	search_txt, _ := reader.ReadString('\n')
	search_txt = strings.Replace(search_txt, "\n", "", -1)
	fmt.Print("Enter website to start search: ")
	website, _ := reader.ReadString('\n')
	website = strings.Replace(website, "\n", "", -1)

	cache := utils.NewUrlCache()

	urls := make(chan utils.UrlItem, 100000)
	urls <- utils.UrlItem{Url: website, Depth: 0}

	indexFunction := indexer(search_txt)

	consumer := utils.NewUrlConsumer(2, 1, cache)
	consumer.Consume(urls, indexFunction)
}
