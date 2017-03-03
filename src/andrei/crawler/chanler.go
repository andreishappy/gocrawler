package crawler

import (
	"fmt"
	"andrei/concurrentset"
	"andrei/fetcher"
)

type Chanler struct {
	fetcher     fetcher.Fetcher
	shouldCrawl func(string) bool
}

func NewChanler(fetcher fetcher.Fetcher, shouldCrawl func(string) bool) *Chanler {
	return &Chanler{fetcher: fetcher, shouldCrawl: shouldCrawl}
}

func (c *Chanler) crawlUsingChannels(url string, ret chan fetcher.Page, dispatched concurrentset.ConcurrentStringSet) {
	defer close(ret)
	page, e := c.fetcher.GetPage(url)
	if e != nil {
		fmt.Printf("Error fetching %s\n error: %s", url, e)
		return
	}

	ret <- page

	results := make([]chan fetcher.Page, 0, len(page.Links))
	for _, link := range page.Links {
		if dispatched.Contains(link) || !c.shouldCrawl(link) {
			continue
		}
		ch := make(chan fetcher.Page, 10)
		dispatched.Put(link)
		go c.crawlUsingChannels(link, ch, dispatched)
		results = append(results, ch)
	}

	for _, ch := range results {
		for p := range ch {
			ret <- p
		}
	}
}

func (c *Chanler) CrawlUsingChannels(url string) map[string]fetcher.Page {
	ret := make(chan fetcher.Page, 10)
	dispatched := concurrentset.NewConcurrentStringSet()
	dispatched.Put(url)
	go c.crawlUsingChannels(url, ret, dispatched)

	result := map[string]fetcher.Page{}

	for p := range ret {
		result[p.Url] = p
	}

	return result
}
