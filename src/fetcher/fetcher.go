package fetcher

import (
	"net/http"
	"golang.org/x/net/html"
	"io"
)

type HttpGetter func(url string) (resp *http.Response, err error)

type Fetcher interface {
	GetPage(url string) (Page, error)
}

type Page struct {
	Url    string
	Links  []string
	Assets []string
}

func NewPage(url string) Page {
	return Page{Url: url}
}

func (f Page)WithAssets(assets ...string) Page {
	if len(assets) == 0 {
		f.Assets = []string{}
	} else {
		f.Assets = assets
	}
	return f
}

func (f Page)WithLinks(links ...string) Page {
	if len(links) == 0 {
		f.Links = []string{}
	} else {
		f.Links = links
	}
	return f
}

func NewWebFetcher(getter HttpGetter, linkValidator func(string) bool, absoluteUrl func(string) string) *WebFetcher {
	return &WebFetcher{get: getter, shouldRecordLink: linkValidator, absoluteUrl: absoluteUrl}
}

type WebFetcher struct {
	get              HttpGetter
	shouldRecordLink func(string) bool
	absoluteUrl      func(string) string
}

func (f *WebFetcher) GetPage(url string) (p Page, e error) {
	resp, err := f.get(url)
	if err != nil {
		e = err
		return
	}

	links, assets := f.getLinksAndAssets(resp.Body)
	p = Page{Url: url, Links: links, Assets: assets}
	return
}

func (f *WebFetcher) getLinksAndAssets(body io.ReadCloser) (links []string, assets []string) {
	linkMap := map[string]bool{}
	assets = make([]string, 0, 10)

	defer body.Close()
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			//end
			links = sliceFromMap(linkMap)
			return
		case tt == html.StartTagToken:
			t := z.Token()

			ok, link := f.getHref(t)
			link = f.absoluteUrl(link)
			if ok && f.shouldRecordLink(link) {
				_, alreadyAdded := linkMap[link]
				if !alreadyAdded {
					linkMap[link] = true
				}
			}

			ok, asset := f.getAsset(t)
			if ok {
				assets = append(assets, asset)
			}
		}
	}
}

func sliceFromMap(m map[string]bool) []string {
	result := make([]string, len(m))

	i := 0
	for k := range m {
		result[i] = k
		i++
	}

	return result
}

func (f *WebFetcher) getAsset(t html.Token) (ok bool, asset string) {
	isImage := t.Data == "img"
	if isImage {
		for _, a := range t.Attr {
			if a.Key == "src" {
				return true, a.Val
			}
		}
	}

	//link -> href
	//script -> src
	//img -> src


	//figure out video and audio
	return
}

func (f *WebFetcher) getHref(t html.Token) (ok bool, link string) {
	isAnchor := t.Data == "a"
	if !isAnchor {
		return
	}

	for _, a := range t.Attr {
		if a.Key == "href" {
			link = a.Val
			ok = true
		}
	}
	return
}

