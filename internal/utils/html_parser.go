package utils

import (
	"bufio"
	"log"
	"net/url"
	"regexp"
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

func (v PageParser) GetDisallowedSuffixes(robotsDoc string) []string {
	r, _ := regexp.Compile(`User-agent: \*(\n|\r\n)(?s)(.*)(\n\n|\r\n\r\n)`)
	afterStarRegexp, _ := regexp.Compile(`\*(?s)(.*)$`)
	wildCardMatch := r.FindString(robotsDoc)
	s := bufio.NewScanner(strings.NewReader(wildCardMatch))

	var returnSuffixes []string
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "User-agent: *") {
			continue
		}
		if line == "Disallow: /" {
			returnSuffixes = append(returnSuffixes, "")
			continue
		}
		if len(line) <= 10 {
			continue
		}
		suffix := afterStarRegexp.FindStringSubmatch(line)[1]
		returnSuffixes = append(returnSuffixes, suffix)
	}
	return returnSuffixes
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
