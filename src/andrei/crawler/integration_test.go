package crawler

import (
	"testing"
	"net/http"
	"io"
	"errors"
	"bytes"
	"github.com/stretchr/testify/assert"
	"andrei/configuration"
	"andrei/collectionutils"
	"andrei/fetcher"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error {
	return nil
}

var rootUrl = "http://root.com"
var absoluteUrl = "http://root.com/absolute"
var relativeUrl = "http://root.com/relative"
var imageUrl = "http://image.png"
var errorNotUnderDomain = "http://error"
var urlUnderDomainButError = "http://root.com/error"
var urlUnderOtherDomainNoError = "http://otherroot.com"

var returnedByClient = map[string] string {
	rootUrl: roothtml,
	absoluteUrl: absolutehtml,
	relativeUrl: relativehtml,
	urlUnderOtherDomainNoError: emptyhtml,
	urlUnderDomainButError: emptyhtml,
}

func graph(url string) (*http.Response, error) {
	html, ok := returnedByClient[url]
	if ok {
		return body(html), nil
	}
	return nil, errors.New("error")
}

var expected = map[string]fetcher.Page {
	rootUrl: fetcher.NewPage(rootUrl).WithLinks(absoluteUrl, urlUnderDomainButError).WithAssets(imageUrl),
	absoluteUrl: fetcher.NewPage(absoluteUrl).WithLinks(relativeUrl).WithAssets(),
	relativeUrl: fetcher.NewPage(relativeUrl).WithLinks().WithAssets(),
}

func linkWithUrl(url string) string {
	return "<a href=\"" + url + "\"></a>"
}

var roothtml = "<html>" +
	linkWithUrl(absoluteUrl) +
	linkWithUrl(errorNotUnderDomain) +
	linkWithUrl(urlUnderDomainButError) +
	linkWithUrl(urlUnderOtherDomainNoError) +
	"<img src=\"http://image.png\">" +
	"</html>"

var absolutehtml = "<html>" +
	"<a href=/relative></a>" +
	"</html>"

var relativehtml = "<html>" +
	linkWithUrl(urlUnderOtherDomainNoError) +
	"</html>"

var emptyhtml = "<html></html>";

func body(bodyString string) *http.Response {
	b := nopCloser{bytes.NewBufferString(bodyString)}
	return &http.Response{Body: b}
}

func TestIntegrationCrawler(t *testing.T) {
	isValid := configuration.HostUrlValidator(rootUrl)
	linkBuilder := configuration.AbsolutePathBuilder(rootUrl)
	f := fetcher.NewWebFetcher(graph, isValid, linkBuilder)
	c := NewCrawler(f, isValid)
	result := c.CrawlUsingSync(rootUrl)

	assertSameGraph(t, expected, result)
}

func TestIntegrationChanler(t *testing.T) {
	isValid := configuration.HostUrlValidator(rootUrl)
	linkBuilder := configuration.AbsolutePathBuilder(rootUrl)
	f := fetcher.NewWebFetcher(graph, isValid, linkBuilder)
	c := NewChanler(f, isValid)
	result := c.CrawlUsingChannels(rootUrl)
	assertSameGraph(t, expected, result)
}

func assertSameGraph(t *testing.T, expected map[string]fetcher.Page, actual map[string]fetcher.Page) {
	for k, expectedPage := range expected {
		actualPage := actual[k]
		assert.Equal(t, expectedPage.Url, actualPage.Url)
		assertSameElements(t, expectedPage.Links, actualPage.Links, k)
		assertSameElements(t, expectedPage.Assets, actualPage.Assets, k)
	}
}

func assertSameElements(t *testing.T, expected []string, actual []string, identifier string) {
	expectedMap := collectionutils.MapFromArray(expected)
	actualMap := collectionutils.MapFromArray(actual)
	assert.Equal(t, expectedMap, actualMap, "Arrays for " + identifier + " do not have the same elements")
}
