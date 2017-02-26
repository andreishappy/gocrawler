package main

import (
	"fetcher"
	"net/http"
	"github.com/davecgh/go-spew/spew"
	"path_utils"
	"os"
	"time"
	"log"
	"chanler"
	"crawler"
)

func main() {
	start := time.Now()

	url := os.Args[1]
	isValid := path_utils.HostUrlValidator(url)
	linkBuilder := path_utils.HostUrlRelativiser(url)
	f := fetcher.NewWebFetcher(http.Get, isValid, linkBuilder)

	p := chanler.NewChanler(f, isValid)

	nodes := p.Crawl(url)

	spew.Dump(nodes)

	elapsed := time.Since(start)
	log.Printf("Took %s to do %d nodes 1", elapsed, len(nodes))
	time.Sleep(300 * time.Millisecond)
	start = time.Now()
	pWG := crawler.NewCrawl(f, isValid)
	graphWG := pWG.Crawl(url)
	elapsed = time.Since(start)
	log.Printf("Took %s to do %d nodes 2", elapsed, len(graphWG))

}

