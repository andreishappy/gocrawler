package main

import (
	"fetcher"
	"net/http"
	"github.com/davecgh/go-spew/spew"
	"configuration"
	"os"
	"time"
	"log"
	"crawler"
	"flag"
	"fmt"
)

func main() {
	start := time.Now()

	url := os.Args[1]
	var howToRun = flag.String("type", "channel", "Choose to use channels or sync")
	var 
	flag.Parse()

	fmt.Println(*ip)
	isValid := configuration.HostUrlValidator(url)
	linkBuilder := configuration.AbsolutePathBuilder(url)

	f := fetcher.NewWebFetcher(http.Get, isValid, linkBuilder)
	p := crawler.NewChanler(f, isValid)
	channelResult := p.CrawlUsingChannels(url)
	spew.Dump(channelResult)
	elapsed := time.Since(start)

	log.Printf("Took %s to do %d nodes 1", elapsed, len(channelResult))

	time.Sleep(300 * time.Millisecond)
	start = time.Now()
	pWG := crawler.NewCrawler(f, isValid)
	syncResult := pWG.CrawlUsingSync(url)
	elapsed = time.Since(start)
	spew.Dump(syncResult)

	log.Printf("Took %s to do %d nodes 2", elapsed, len(syncResult))
	log.Printf("Same: %t", same(syncResult, channelResult))
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

