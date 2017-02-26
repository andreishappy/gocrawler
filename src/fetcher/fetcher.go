package fetcher

import (
	"net/http"
	"golang.org/x/net/html"
	"io"
	"crawler"
	"strings"
)

type HttpGetter func(url string) (resp *http.Response, err error)

func NewWebFetcher(getter HttpGetter) *WebFetcher {
	return &WebFetcher{get: getter}
}

type WebFetcher struct {
	get HttpGetter
}

func (f *WebFetcher) GetPage(url string, base string) (*crawler.FetchedPage, error) {
	resp, err := f.get(url)
	if err != nil {
		return nil, err
	}

	links, _ := getLinksAndAssets(resp.Body)
	return &crawler.FetchedPage{Links: links, Assets: []string{}}, nil
}

func getLinksAndAssets(body io.ReadCloser) (links []string, assets []string) {
	links = make([]string, 0, 10)
	assets = make([]string, 0, 10)

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

			ok, link := getHref(t)
			if ok {
				links = append(links, link)
			}

			ok, asset := getAsset(t)
			if ok {
				assets = append(assets, asset)
			}
		}
	}
}

func getAsset(t html.Token) (ok bool, asset string) {
	return

	//link -> href
	//script -> src
	//img -> src


	//figure out video and audio
}

func getHref(t html.Token) (ok bool, link string) {
	isAnchor := t.Data == "a"
	if !isAnchor {
		return
	}

	for _, a := range t.Attr {
		if a.Key == "href" && strings.HasPrefix(a.Val, "http://tomblomfield.com/") {
			link = a.Val
			ok = true
		}
	}
	return
}

