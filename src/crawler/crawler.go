package crawler

import (
	"fmt"
	"sync"
	"fetcher"
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
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	result := map[string]PageNode{}
	dispatched := map[string]bool{}

	var crawlFunc func(string, *Crawl)

	crawlFunc = func(url string, c *Crawl) {
		defer wg.Done()
		//get the links and assets
		page, e := c.fetcher.GetPage(url)

		if (e != nil) {
			fmt.Println("Error when fetching " + url)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()
		//Check if url has already been visited

		//Spawn off go routines for the links
		for _, link := range page.Links {
			_, ok := dispatched[link]
			if (!ok && c.shouldCrawl(link)) {
				dispatched[link] = true
				wg.Add(1)
				fmt.Println("spawning go routine for ", link)
				go crawlFunc(link, c)
			}
		}

		result[url] = PageNode{Links: page.Links, Assets: page.Assets}
	}

	wg.Add(1)
	dispatched[url] = true
	crawlFunc(url, c)
	wg.Wait()
	return result
}
