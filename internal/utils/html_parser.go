package utils

import (
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Parser interface {
	GetLinks(string, string) []string
}

type PageParser struct{}

func makeUrlAbsolute(urlLink string, base string) string {
	parsedUrl, err := url.Parse(urlLink)
	if err != nil {
		log.Fatal(err)
	}
	if parsedUrl.IsAbs() {
		return urlLink
	}
	// base url will always start with '/', remove from relative url
	if urlLink[0] == 0x2f {
		urlLink = urlLink[1:]
	}
	return base + urlLink
}

func (v PageParser) GetLinks(htmldoc string, baseurl string) []string {
	doc, err := html.Parse(strings.NewReader(htmldoc))
	if err != nil {
		log.Fatal(err)
	}
	links := make([]string, 0)

	var f func(*html.Node, *[]string)
	f = func(node *html.Node, links *[]string) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					*links = append(*links, makeUrlAbsolute(attr.Val, baseurl))
					break
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c, links)
		}
	}
	f(doc, &links)
	return links
}
