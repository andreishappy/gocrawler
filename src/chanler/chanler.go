package chanler

import (
	"fmt"
	"fetcher"
	"concurrentset"
)

type Chanler struct {
	fetcher     fetcher.Fetcher
	shouldCrawl func(string) bool
}

func NewChanler(fetcher fetcher.Fetcher, shouldCrawl func(string) bool) *Chanler {
	return &Chanler{fetcher: fetcher, shouldCrawl: shouldCrawl}
}

func (c *Chanler) crawlWithVisited(url string, ret chan fetcher.FetchedPage, dispatched concurrentset.ConcurrentStringSet) {
	defer close(ret)
	page, e := c.fetcher.GetPage(url)
	if e != nil {
		return
	}
	ret <- *page

	results := make([]chan fetcher.FetchedPage, 0, len(page.Links))
	for _, link := range page.Links {
		if dispatched.HasAlreadySeen(link) || !c.shouldCrawl(link) {
			continue
		}
		fmt.Println("Spawning go routine for " + link)
		ch := make(chan fetcher.FetchedPage)
		dispatched.RecordSeen(link)
		go c.crawlWithVisited(link, ch, dispatched)
		results = append(results, ch)
	}

	for _, ch := range results {
		for p := range ch {
			ret <- p
		}
	}
}

func (c *Chanler) Crawl(url string) map[string]fetcher.FetchedPage {
	ret := make(chan fetcher.FetchedPage)
	dispatched := concurrentset.NewConcurrentStringSet()
	dispatched.RecordSeen(url)
	go c.crawlWithVisited(url, ret, dispatched)

	result := map[string]fetcher.FetchedPage{}

	for p := range ret {
		result[p.Url] = p
	}

	return result
}
