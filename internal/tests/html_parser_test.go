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
	exampleString := `
	<p>Links:</p>
	<ul>
		<li><a href="http://google.com/foo.html">Foo</a>
		<li><a href="foo.html">Foo</a>
		<li><a href="/bar/baz.html">BarBaz</a>
	</ul>
	`
	links := utils.PageParser{}.GetLinks(exampleString, "http://crawler.biz/")
	expectedLinks := []string{
		"http://google.com/foo.html",
		"http://crawler.biz/foo.html",
		"http://crawler.biz/bar/baz.html",
	}
	if !testEqual(links, expectedLinks) {
		t.Errorf("getLinks failed, expected %v got %v", expectedLinks, links)
	}
}

func TestRobotsParser(t *testing.T) {
	// robots.txt could use \n or \r\n for newline
	exampleString := "User-agent: test\n" +
		"Disallow: /.png\n\n" +
		"User-agent: *\n" +
		"Disallow: /\n" +
		"Disallow: /*.json\n" +
		"Disallow: /*.xml\n\n" +
		"User-agent: test2\n" +
		"Disallow: /.mp3"

	disallowed := utils.PageParser{}.GetDisallowedSuffixes(exampleString)
	expectedDisallowed := []string{"", ".json", ".xml"}
	if !testEqual(disallowed, expectedDisallowed) {
		t.Errorf(
			"getDisallowedSuffixes failed, expected %v got %v",
			expectedDisallowed,
			disallowed,
		)
	}

	exampleString2 := "User-agent: test\r\n" +
		"Disallow: /.png\r\n\r\n" +
		"User-agent: *\r\n" +
		"Disallow: /\r\n" +
		"Disallow: /*.json\r\n" +
		"Disallow: /*.xml\r\n\r\n" +
		"User-agent: test2\r\n" +
		"Disallow: /.mp3"

	disallowed = utils.PageParser{}.GetDisallowedSuffixes(exampleString2)
	expectedDisallowed = []string{"", ".json", ".xml"}
	if !testEqual(disallowed, expectedDisallowed) {
		t.Errorf(
			"getDisallowedSuffixes failed, expected %v got %v",
			expectedDisallowed,
			disallowed,
		)
	}
}
