package main

import (
	"fetcher"
	"net/http"
	"github.com/davecgh/go-spew/spew"
	"urlhelper"
	"os"
	"time"
	"log"
	"chanler"
	"crawler"
)

func main() {
	start := time.Now()

	url := os.Args[1]
	isValid := urlhelper.HostUrlValidator(url)
	linkBuilder := urlhelper.AbsolutePathBuilder(url)
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
	spew.Dump(graphWG)
	log.Printf("Took %s to do %d nodes 2", elapsed, len(graphWG))

	log.Printf("Same: %t", same(graphWG, nodes))

}

func same(left map[string]fetcher.Page, right map[string]fetcher.Page) bool {
	if len(left) != len(right) {
		return false
	}

	for k, pLeft := range left {
		pRight := right[k]
		if len(pLeft.Links) != len(pRight.Links) {
			return false
		}
	}
	return true
}

