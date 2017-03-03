package crawler

import (
	"fmt"
	"sync"
	"andrei/concurrentset"
	"andrei/fetcher"
)

type Crawl struct {
	fetcher     fetcher.Fetcher
	shouldCrawl func(string) bool
}

func NewCrawler(fetcher fetcher.Fetcher, shouldCrawl func(string) bool) *Crawl {
	return &Crawl{fetcher: fetcher, shouldCrawl: shouldCrawl}
}

func (c *Crawl) CrawlUsingSync(url string) map[string]fetcher.Page {
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	result := map[string]fetcher.Page{}
	dispatched := concurrentset.NewConcurrentStringSet()

	wg.Add(1)
	dispatched.Put(url)
	go c.crawlUsingSyc(url, wg, mutex, dispatched, result)
	wg.Wait()
	return result
}

func (c *Crawl) crawlUsingSyc(url string, wg *sync.WaitGroup, mutex *sync.Mutex, dispatched concurrentset.ConcurrentStringSet, result map[string]fetcher.Page) {
	defer wg.Done()
	//get the links and assets
	page, e := c.fetcher.GetPage(url)

	if (e != nil) {
		fmt.Println("Error when fetching " + url)
		return
	}


	//Spawn off go routines for the links
	for _, link := range page.Links {
		if (dispatched.Contains(link) || !c.shouldCrawl(link)) {
			continue
		}
		dispatched.Put(link)
		wg.Add(1)
		go c.crawlUsingSyc(link, wg, mutex, dispatched, result)
	}

	mutex.Lock()
	defer mutex.Unlock()
	result[url] = page
}
