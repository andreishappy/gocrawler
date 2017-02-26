package main

import (
	"crawler"
	"fetcher"
	"net/http"
	"github.com/davecgh/go-spew/spew"
	"path_utils"
)

func main() {
	url := "http://tomblomfield.com/"
	validator := path_utils.HostUrlValidator(url)
	f := fetcher.NewWebFetcher(http.Get, validator)
	p := crawler.NewCrawl(f, validator)
	spew.Dump(p.Crawl(url))
}

