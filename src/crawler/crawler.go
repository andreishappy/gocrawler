package crawler

import (
	"fmt"
	"sync"
)

type PageNode struct {
	Links []string
	Assets []string
}

type Crawl struct {
	fetcher Fetcher
}

func NewCrawl(fetcher Fetcher) *Crawl {
	return &Crawl{fetcher: fetcher}
}

type Fetcher interface {
	GetPage(url string, baseUrl string) (*FetchedPage, error)
}

type FetchedPage struct {
	Links  []string
	Assets []string
}

func (c *Crawl) Crawl(url string) map[string]PageNode {
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	result := map[string]PageNode{}

	var crawlFunc func(string, *Crawl)

	crawlFunc = func(url string, c *Crawl) {
		defer wg.Done()
		defer fmt.Println("done with " + url)

		//get the links and assets
		page, e := c.fetcher.GetPage(url, "")

		if (e != nil) {
			fmt.Println("Error when fetching " + url)
			return
		}

		mutex.Lock()
		//Check if url has already been visited
		_, ok := result[url]
		if (ok) {
			fmt.Println("Already saw ", url)
			mutex.Unlock()
			return
		} else {
			fmt.Println("Saw ", url, " for the first time")

			//Spawn off go routines for the links
			for _, link := range page.Links {
				_, ok = result[link]
				if(!ok) {
					wg.Add(1)
					fmt.Println("spawning go routine for ", link)
					go crawlFunc(link, c)
				}
			}

			result[url] = PageNode{Links: page.Links, Assets: page.Assets}
		}
		mutex.Unlock()

	}

	wg.Add(1)
	crawlFunc(url, c)
	wg.Wait()
	return result
}
