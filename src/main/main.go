package main

import (
	"fmt"
	"crawler"
	"fetcher"
	"net/http"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	fmt.Println("hello")
	f := fetcher.NewWebFetcher(http.Get)
	p := crawler.NewCrawl(f)
	spew.Dump(p.Crawl("http://tomblomfield.com/"))
}

