package fetcher

import (
	"net/http"
	"fmt"
	"golang.org/x/net/html"
	"io"
)

type WebFetcher struct {

}

type FetchedPage struct {
	Links  []string
	Assets []string
}

func (WebFetcher) getPage(url string, base string) (error, *FetchedPage) {
	resp, err := http.Get(url)
	if err != nil {
		return err, nil
	}

	links, _ := getLinksAndAssets(resp.Body)
	fmt.Println(links)
	return nil, &FetchedPage{Links: links}
}

func getLinksAndAssets(body io.ReadCloser) (links []string, assets []string) {
	links = make([]string, 0, 10)

	defer body.Close()
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			//end
			return
		case tt == html.StartTagToken:
			t := z.Token()

			ok, url := getHref(t)
			if ok {
				links = append(links, url)
			}
		}
	}
}

func getHref(t html.Token) (ok bool, href string) {
	isAnchor := t.Data == "a"
	if !isAnchor {
		return
	}

	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

