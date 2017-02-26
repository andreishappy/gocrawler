package fetcher

import (
	"net/http"
	"golang.org/x/net/html"
	"io"
)

type HttpGetter func(url string) (resp *http.Response, err error)

type Fetcher interface {
	GetPage(url string) (*FetchedPage, error)
}

type FetchedPage struct {
	Url    string
	Links  []string
	Assets []string
}

func NewWebFetcher(getter HttpGetter, linkValidator func(string) bool, buildLink func(string) string) *WebFetcher {
	return &WebFetcher{get: getter, linkIsValid: linkValidator, buildLink: buildLink}
}

type WebFetcher struct {
	get         HttpGetter
	linkIsValid func(string) bool
	buildLink   func(string) string
}

func (f *WebFetcher) GetPage(url string) (*FetchedPage, error) {
	resp, err := f.get(url)
	if err != nil {
		return nil, err
	}

	links, assets := f.getLinksAndAssets(resp.Body)
	return &FetchedPage{Url: url, Links: links, Assets: assets}, nil
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
			link = f.buildLink(link)
			if ok && f.linkIsValid(link) {
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

