package main

import (
	"crawler"
	"fetcher"
	"net/http"
	"github.com/davecgh/go-spew/spew"
	"path_utils"
	"os"
)

func main() {
	url := os.Args[1]
	validator := path_utils.HostUrlValidator(url)
	relativiser := path_utils.HostUrlRelativiser(url)
	f := fetcher.NewWebFetcher(http.Get, validator, relativiser)
	p := crawler.NewCrawl(f, validator)
	spew.Dump(p.Crawl(url))
}

