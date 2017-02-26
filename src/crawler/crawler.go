package crawler

import (
	"fmt"
	"sync"
)

type Crawl struct {
	fetcher Fetcher
	linkValidator func(string) bool
}

func NewCrawl(fetcher Fetcher, linkValidator func(string) bool) *Crawl {
	return &Crawl{fetcher: fetcher, linkValidator: linkValidator}
}

type Fetcher interface {
	GetPage(url string) (*FetchedPage, error)
}

type FetchedPage struct {
	Links  []string
	Assets []string
}

type PageNode struct {
	Links []string
	Assets []string
}

func (c *Crawl) Crawl(url string) map[string]PageNode {
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	result := map[string]PageNode{}

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
		_, ok := result[url]
		if (ok) {
			return
		} else {
			//Spawn off go routines for the links
			for _, link := range page.Links {
				_, ok = result[link]
				if(!ok && c.linkValidator(link)) {
					wg.Add(1)
					fmt.Println("spawning go routine for ", link)
					go crawlFunc(link, c)
				}
			}

			result[url] = PageNode{Links: page.Links, Assets: page.Assets}
		}
	}

	wg.Add(1)
	crawlFunc(url, c)
	wg.Wait()
	return result
}
