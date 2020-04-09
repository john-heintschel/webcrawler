package utils

import (
	"bufio"
	"bytes"
	"log"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type Parser interface {
	GetLinks(string, string) []string
	UrlIsForbidden(string, []string) bool
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
	// base url will always end with '/', remove from relative url if need be
	if urlLink[0] == 0x2f {
		urlLink = urlLink[1:]
	}
	return base + urlLink
}

func UrlIsForbidden(url string, disallowed []string) bool {
	for _, element := range disallowed {
		regularExpression, _ := regexp.Compile(element)
		match := regularExpression.Find([]byte(url))
		if match != nil {
			return true
		}
	}
	return false
}

// taken from https://stackoverflow.com/questions/41433422/
// read-lines-from-a-file-with-variable-line-endings-in-go
// need this special scan function because robots.txt can use \r, \n, or \r\n
func ScanRobotsLines(
	data []byte,
	atEOF bool,
) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		if data[i] == '\n' {
			return i + 1, data[0:i], nil
		}
		advance = i + 1
		if len(data) > i+1 && data[i+1] == '\n' {
			advance += 1
		}
		return advance, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func (v PageParser) GetDisallowedPatterns(robotsDoc string) []string {
	r, _ := regexp.Compile(
		`(?U)User-agent: ?\*(\n|\r\n|\r)(?s)(.*)(User-agent|\z)`,
	)
	wildCardMatch := r.FindString(robotsDoc)
	s := bufio.NewScanner(strings.NewReader(wildCardMatch))
	s.Split(ScanRobotsLines)
	var returnPatterns []string
	for s.Scan() {
		line := s.Text()
		if !strings.HasPrefix(line, "Disallow:") {
			continue
		}
		suffix := strings.TrimSpace(line[9:])
		if suffix == "/" {
			return []string{".*"}
		}
		suffix = strings.ReplaceAll(suffix, `.`, `\.`)
		suffix = strings.ReplaceAll(suffix, `*`, `.*`)
		// TODO return list of compiled regex instead of raw []strings
		returnPatterns = append(returnPatterns, suffix)
	}
	return returnPatterns
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
