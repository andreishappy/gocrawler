package crawler

import (
	"fmt"
	"sync"
	"fetcher"
	"concurrentset"
)

type Crawl struct {
	fetcher     fetcher.Fetcher
	shouldCrawl func(string) bool
}

func NewCrawl(fetcher fetcher.Fetcher, shouldCrawl func(string) bool) *Crawl {
	return &Crawl{fetcher: fetcher, shouldCrawl: shouldCrawl}
}

type PageNode struct {
	Links  []string
	Assets []string
}

func (c *Crawl) Crawl(url string) map[string]PageNode {
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	result := map[string]PageNode{}
	dispatched := concurrentset.NewConcurrentStringSet()

	wg.Add(1)
	dispatched.RecordSeen(url)
	go c.crawlConcurrent(url, wg, mutex, &dispatched, result)
	wg.Wait()
	return result
}

func (c *Crawl) crawlConcurrent(url string, wg *sync.WaitGroup, mutex *sync.Mutex, dispatched *concurrentset.ConcurrentStringSet, result map[string]PageNode) {
	defer wg.Done()
	//get the links and assets
	page, e := c.fetcher.GetPage(url)

	if (e != nil) {
		fmt.Println("Error when fetching " + url)
		return
	}


	//Spawn off go routines for the links
	for _, link := range page.Links {
		if (dispatched.HasAlreadySeen(link) || !c.shouldCrawl(link)) {
			continue
		}
		dispatched.RecordSeen(link)
		wg.Add(1)
		fmt.Println("spawning go routine for ", link)
		go c.crawlConcurrent(link, wg, mutex, dispatched, result)
	}

	mutex.Lock()
	result[url] = PageNode{Links: page.Links, Assets: page.Assets}
	mutex.Unlock()
}
