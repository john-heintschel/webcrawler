package tests

import (
	"testing"

	"github.com/john-heintschel/webcrawler/internal/utils"
)

func testEqual(A, B []string) bool {
	if (A == nil) != (B == nil) {
		return false
	}

	if len(A) != len(B) {
		return false
	}
	for idx, val := range A {
		if val != B[idx] {
			return false
		}
	}
	return true
}
func TestHTMLParser(t *testing.T) {
	exampleString := "<p>Links:</p><ul><li><a href=\"http://google.com/foo.html\">Foo</a><li><a href=\"foo.html\">Foo</a><li><a href=\"/bar/baz.html\">BarBaz</a></ul>"
	links := utils.PageParser{}.GetLinks(exampleString, "http://crawler.biz/")
	expectedLinks := []string{"http://google.com/foo.html", "http://crawler.biz/foo.html", "http://crawler.biz/bar/baz.html"}
	if !testEqual(links, expectedLinks) {
		t.Errorf("getLinks failed, expected %v got %v", expectedLinks, links)
	}
}
