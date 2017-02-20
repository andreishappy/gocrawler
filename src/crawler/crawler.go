package crawler

import "fmt"

type Page struct {
	Links  []string
	Assets []string
}

type PageNode struct {
	Url string
	Links []PageNode
	Assets []string
}

type Feeder interface {
  getLinks(url string) Page
}

type Crawl struct {
	feeder Feeder
}

func (c *Crawl) crawl(url string) PageNode {

	feederResult := c.feeder.getLinks(url)
	fmt.Println(feederResult)
	return PageNode{Url: url}

	// create a set of visited urls
	// instantiate a graph

	// issue go routines that check if a url has already been visited
	//   - if it has stop
	//   -
}
