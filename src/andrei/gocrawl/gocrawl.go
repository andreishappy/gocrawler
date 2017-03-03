package main

import (
	"andrei/fetcher"
	"net/http"
	"andrei/configuration"
	"os"
	"time"
	"log"
	"andrei/crawler"
	"andrei/fileutils"
)

func main() {
	if len(os.Args) < 3 {
		panic("Usage: command ( sync | channel | both ) url")
	}
	mode := os.Args[1]
	url := os.Args[2]

	isValid := configuration.HostUrlValidator(url)
	absolutePathBuilder := configuration.AbsolutePathBuilder(url)

	switch mode {
	case "sync":
		runSync(url, isValid, absolutePathBuilder)
	case "channel":
		runChannel(url, isValid, absolutePathBuilder)
	case "both":
		runSync(url, isValid, absolutePathBuilder)
		runChannel(url, isValid, absolutePathBuilder)
	default:
		panic("Need to choose ( sync | channel | both ) as the second argument")
	}
}

func runChannel(url string, isValid func(url string) bool, absolutePathBuilder func(string) string) {
	log.Printf("Crawling %s using channels", url)
	start := time.Now()

	f := fetcher.NewWebFetcher(http.Get, isValid, absolutePathBuilder)
	p := crawler.NewChanler(f, isValid)
	channelResult := p.CrawlUsingChannels(url)

	fileutils.WriteToFileInJson("channel.json", channelResult)

	elapsed := time.Since(start)
	log.Printf("Took %s to do %d nodes using channels", elapsed, len(channelResult))
}

func runSync(url string, isValid func(url string) bool, absolutePathBuilder func(string) string) {
	log.Printf("Crawling %s using sync", url)
	start := time.Now()

	f := fetcher.NewWebFetcher(http.Get, isValid, absolutePathBuilder)
	syncCrawler := crawler.NewCrawler(f, isValid)
	syncResult := syncCrawler.CrawlUsingSync(url)
	fileutils.WriteToFileInJson("sync.json", syncResult)

	elapsed := time.Since(start)
	log.Printf("Took %s to do %d nodes using sync", elapsed, len(syncResult))
}
