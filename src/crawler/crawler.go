package crawler

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	getPage(url string) (error, *FetchedPage)
}

type PageNode struct {
	Links []string
	Assets []string
}

type Crawl struct {
	fetcher Fetcher
}

func (c *Crawl) crawl(url string) map[string]PageNode {
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	visited := map[string]bool{}
	result := map[string]PageNode{}

	var crawlFunc func(string, map[string]bool, *Crawl)
	crawlFunc = func(url string, visited map[string]bool, c *Crawl) {
		defer wg.Done()

		//get the links and assets
		e, page := c.fetcher.getPage(url)

		if (e != nil) {
			fmt.Println("Error when fetching " + url)
			return
		}

		mutex.Lock()

		//Check if url has already been visited
		_, ok := visited[url]
		if (ok) {
			fmt.Println("Already saw ", url)
			return
		} else {
			fmt.Println("Saw ", url, " for the first time")

			//Spawn off go routines for the links
			for _, link := range page.Links {
				wg.Add(1)
				fmt.Println("spawning go routine for ", link)
				go crawlFunc(link, visited, c)
				fmt.Println("after spawn for ", link)
			}

			visited[url] = true
			result[url] = PageNode{Links: page.Links, Assets: page.Assets}
		}

		mutex.Unlock()

	}

	wg.Add(1)
	crawlFunc(url, visited, c)
	wg.Wait()
	fmt.Println(result)
	return result
}
