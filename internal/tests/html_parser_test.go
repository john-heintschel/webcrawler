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

func TestUrlIsForbidden(t *testing.T) {
	forbiddenUrl := "https://en.wikipedia.org/hello.json"
	forbiddenUrl2 := "https://en.wikipedia.org/web/hello.html"
	forbiddenUrl3 := "https://en.wikipedia.org/hello.json"
	allowedUrl := "https://en.wikipedia.org"
	disallowed := []string{`/.*\.json`, `/web/`}
	wildcardDisallowed := []string{`.*`}
	if !utils.UrlIsForbidden(forbiddenUrl, disallowed) {
		t.Errorf("UrlIsForbidden failed")
	}
	if !utils.UrlIsForbidden(forbiddenUrl2, disallowed) {
		t.Errorf("UrlIsForbidden failed")
	}
	if !utils.UrlIsForbidden(forbiddenUrl3, wildcardDisallowed) {
		t.Errorf("UrlIsForbidden failed")
	}
	if utils.UrlIsForbidden(allowedUrl, disallowed) {
		t.Errorf("UrlIsForbidden failed, allowed url forbidden")
	}
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
	// robots.txt could use \n, \r, or \r\n for newline
	exampleString := "User-agent: test\n" +
		"Disallow: /.png\n\n" +
		"User-agent: *\n" +
		"Disallow: /\n" +
		"Allow: /\n" +
		"Disallow: /*.json\n" +
		"Disallow: /*.xml\n\n" +
		"User-agent: test2\n" +
		"Disallow: /.mp3"

	disallowed := utils.PageParser{}.GetDisallowedPatterns(exampleString)
	expectedDisallowed := []string{".*"}
	if !testEqual(disallowed, expectedDisallowed) {
		t.Errorf(
			"getDisallowedPatterns failed, expected %v got %v",
			expectedDisallowed,
			disallowed,
		)
	}

	exampleString2 := "User-agent: test\r\n" +
		"Disallow: /.png\r\n\r\n" +
		"User-agent: *\r\n" +
		"Allow: /\r\n" +
		"Disallow: /*.json\r\n" +
		"Disallow: /*.xml\r\n\r\n"

	disallowed = utils.PageParser{}.GetDisallowedPatterns(exampleString2)
	expectedDisallowed = []string{`/.*\.json`, `/.*\.xml`}
	if !testEqual(disallowed, expectedDisallowed) {
		t.Errorf(
			"getDisallowedPatterns failed, expected %v got %v",
			expectedDisallowed,
			disallowed,
		)
	}
	exampleString3 := "User-agent: test\r" +
		"Disallow: /.png\r\r" +
		"User-agent: *\r" +
		"Allow: /\r" +
		"Disallow: /*.json\r" +
		"Disallow: /*.xml\r\r" +
		"User-agent: test2\r" +
		"Disallow: /.mp3"

	disallowed = utils.PageParser{}.GetDisallowedPatterns(exampleString3)
	if !testEqual(disallowed, expectedDisallowed) {
		t.Errorf(
			"getDisallowedPatterns failed, expected %v got %v",
			expectedDisallowed,
			disallowed,
		)
	}
}
