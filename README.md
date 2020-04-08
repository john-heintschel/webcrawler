# Webcrawler
A barebones web crawler in Go which finds web pages containing a provided search term. Supports concurrent searching.

## Installation

> go get github.com/john-heintschel/webcrawler

## Usage

> go run webcrawler.go

You'll be asked to provide some configuration options before webcrawler begins searching. In particular you'll need to provide:
 - the maximum concurrency you'd like to use, as an integer
 - the search term you'd like to find
 - a starting website to begin the search
 - the maximum amount of times to query any one particular domain (TODO)
 - the maximum depth of webcrawler's search tree (TODO)
