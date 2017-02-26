package chanler

import (
	"fmt"
	"sync"
	"fetcher"
)

type Chanler struct {
	fetcher     fetcher.Fetcher
	shouldCrawl func(string) bool
}

func NewChanler(fetcher fetcher.Fetcher, shouldCrawl func(string) bool) *Chanler {
	return &Chanler{fetcher: fetcher, shouldCrawl: shouldCrawl}
}

func (c *Chanler) crawlWithVisited(url string, ret chan fetcher.FetchedPage, visited map[string]bool, mutex *sync.Mutex) {
	defer close(ret)
	page, e := c.fetcher.GetPage(url)
	if e != nil {
		return
	}
	ret <- *page

	results := make([]chan fetcher.FetchedPage, 0, len(page.Links))
	for _, link := range page.Links {
		if !c.alreadyDispatched(link, visited, mutex) && c.shouldCrawl(link) {
			fmt.Println("Spawning go routine for " + link)
			ch := make(chan fetcher.FetchedPage)
			c.recordDispatch(link, visited, mutex)
			go c.crawlWithVisited(link, ch, visited, mutex)
			results = append(results, ch)
		} else {
			fmt.Println("Not starting go routine for " + link)
		}
	}

	for _, ch := range results {
		for p := range ch {
			ret <- p
		}
	}
	fmt.Println("Finished with all ", len(results), " spawns of " + url)
}

func (c *Chanler) alreadyDispatched(url string, visited map[string]bool, mutex *sync.Mutex) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, alreadyVisited := visited[url]
	return alreadyVisited
}

func (c *Chanler) recordDispatch(url string, visited map[string]bool, mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
	visited[url] = true
}

func (c *Chanler) Crawl(url string) map[string]fetcher.FetchedPage {
	ret := make(chan fetcher.FetchedPage)
	visited := map[string]bool{}
	mutex := &sync.Mutex{}
	c.alreadyDispatched(url, visited, mutex)
	go c.crawlWithVisited(url, ret, visited, mutex)

	result := map[string]fetcher.FetchedPage{}

	for {
		p, more := <- ret
		if more {
			fmt.Println("Adding page ", p)
			result[p.Url] = p
		} else {
			fmt.Println("Done with pages")
			return result
		}
	}
}
