package main

import (
	"crawler"
	"fetcher"
	"net/http"
	"github.com/davecgh/go-spew/spew"
	"path_utils"
	"os"
	"time"
	"log"
)

func main() {
	start := time.Now()

	url := os.Args[1]
	isValid := path_utils.HostUrlValidator(url)
	linkBuilder := path_utils.HostUrlRelativiser(url)
	f := fetcher.NewWebFetcher(http.Get, isValid, linkBuilder)
	p := crawler.NewCrawl(f, isValid)
	graph := p.Crawl(url)
	spew.Dump(graph)

	elapsed := time.Since(start)
	log.Printf("Took %s to do %d nodes", elapsed, len(graph))
}

