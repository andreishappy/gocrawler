package fetcher

import (
	"net/http"
	"fmt"
	"golang.org/x/net/html"
	"errors"
)

type WebFetcher struct {

}

type FetchedPage struct {

}

func (WebFetcher) getPage(url string) (error, *FetchedPage) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error")
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	for {
		fmt.Println("yo")
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			fmt.Println("Error token")
			// End of the document, we're done
			return errors.New("hello"), nil
		case tt == html.StartTagToken:
			t := z.Token()

			ok, url := getHref(t)
			if ok {
				fmt.Println(url)
			}

			isAnchor := t.Data == "a"
			if isAnchor {
				fmt.Println("We found a link!\n")
			}
		}
	}

	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(err)
	//fmt.Println(body)
	return nil, &FetchedPage{}
}

func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

