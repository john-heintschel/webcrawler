# webcrawler
A barebones web crawler in Go which finds web pages containing a provided search term. Supports concurrent searching.

## Installation
You'll need to first [install the Go runtime](https://golang.org/doc/install). After installing Go you can do:

> go get github.com/john-heintschel/webcrawler

## Usage
To begin using webcrawler you can run:

> go run github.com/john-heintschel/webcrawler/

You'll be asked to provide some configuration options before webcrawler begins searching. In particular you'll need to provide:
 - the maximum concurrency you'd like to use, as an integer
 - the maximum amount of times to query any one particular domain, as an integer
 - the maximum depth of webcrawler's search tree, as an integer
 - the search term you'd like to find
 - a starting website to begin the search

 webcrawler will then begin at the website specified and follow linked pages, printing out to stdout whether or not your search term was found at that page, and gathering additional links to search from each page traversed.
