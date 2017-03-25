package main

import (
	"fmt"
	"time"
	"andrei/fetcher"
	"net/http"
	"andrei/configuration"
)

type Crawler struct {
	fetcher     fetcher.Fetcher
	shouldCrawl func(string) bool
}

type Result struct {
	Page fetcher.Page
	err  error
}

func NewCrawler(fetcher fetcher.Fetcher, shouldCrawl func(string) bool) *Crawler {
	return &Crawler{fetcher: fetcher, shouldCrawl: shouldCrawl}
}

func main() {

	url := "http://tomblomfield.com"
	isValid := configuration.AllValid()

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	absolutePathBuilder := configuration.AbsolutePathBuilder(url)
	f := fetcher.NewWebFetcher(client.Get, isValid, absolutePathBuilder)

	c := NewCrawler(f, isValid)
	start := time.Now()
	result := c.Crawl(url, 10000)
	elapsed := time.Since(start)
	fmt.Printf("done %d in %s\n", len(result), elapsed)
}

func (c *Crawler) Crawl(url string, limit int) map[string]struct{} {
	workerQueue := make(chan string)
	filterQueue := make(chan Result)

	for i := 0; i < 20; i++ {
		go c.worker(workerQueue, filterQueue)
	}

	return c.filter(workerQueue, filterQueue, url, limit)
}

func (c *Crawler) worker(input chan string, output chan Result) {
	for url := range input {
		page, e := c.fetcher.GetPage(url)

		if e != nil {
			fmt.Printf("Error %s\n", e)
			output <- Result{err: e}
			continue
		}

		output <- Result{Page: page}
	}
}

func (c *Crawler) filter(workerQueue chan string, results chan Result, url string, limit int) map[string]struct{} {

	dispatched := 0
	dispatchedSet := map[string]struct{}{}
	dispatchedSet[url] = struct{}{}
	dispatched++
	workerQueue <- url

	for result := range results {
		dispatched--
		fmt.Println(dispatched)
		if result.err == nil {
			for _, p := range result.Page.Links {
				_, alreadySeen := dispatchedSet[p]

				if (!alreadySeen && c.shouldCrawl(p) && len(dispatchedSet) < limit) {
					dispatchedSet[p] = struct{}{}
					dispatched++
					go dispatch(workerQueue, p)
				}
			}
		}

		if dispatched == 0 {
			close(workerQueue)
			break
		}
	}

	return dispatchedSet
}

func dispatch(workerQueue chan string, p string) {
	workerQueue <- p
}


