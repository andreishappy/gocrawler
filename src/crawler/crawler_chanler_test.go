package crawler

//package management
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"errors"
	"fmt"
	"fetcher"
)

//
// Setup
//

type TestFeeder struct {
	m map[string]fetcher.Page
}

var a = "a"
var b = "b"
var c = "c"
var assetA1 = "assetA1"
var assetA2 = "assetA2"
var assetB1 = "assetB1"

var selfReference = map[string]fetcher.Page{
	a: fetcher.NewPage(a).WithLinks(a),
}

var circularReferences = map[string]fetcher.Page{
	a: fetcher.NewPage(a).WithLinks(b),
	b: fetcher.NewPage(b).WithLinks(a),
}

var withAssets = map[string]fetcher.Page{
	a: fetcher.NewPage(a).WithLinks(b).WithAssets(assetA1, assetA2),
	b: fetcher.NewPage(b).WithLinks(c).WithAssets(assetB1),
}

func (t TestFeeder) GetPage(url string) (page fetcher.Page, err error) {
	elem, ok := t.m[url]
	if (ok) {
		page = elem
	} else {
		err = errors.New(fmt.Sprint("Error when fetching ", url))
	}
	return
}

func returnTrue(string) bool {
	return true
}

func returnFalse(string) bool {
	return false
}

//
// TestCrawl_CrawlUsingSync
//

func TestCrawl_CrawlUsingSync_DoesNotHangWhenAPageReferencesItself(t *testing.T) {
	f := TestFeeder{m: selfReference}
	p := Crawl{f, returnTrue}
	expected := selfReference
	assert.Equal(t, expected, p.CrawlUsingSync(a))
}

func TestCrawl_CrawlUsingSync_DoesNotHangWhenA2PageCircularDependencyExists(t *testing.T) {
	f := TestFeeder{m: circularReferences}
	p := Crawl{f, returnTrue}
	expected := circularReferences
	assert.Equal(t, expected, p.CrawlUsingSync(a))
}

func TestCrawl_CrawlUsingSync_DoesNotFollowInvalidLinks(t *testing.T) {
	f := TestFeeder{m: circularReferences}
	p := Crawl{f, returnFalse}
	expected := map[string]fetcher.Page{
		a: fetcher.NewPage(a).WithLinks(b),
	}
	assert.Equal(t, expected, p.CrawlUsingSync(a))
}

func TestCrawl_CrawlUsingSync_AddPageWhenFetcherReturnsError(t *testing.T) {
	f := TestFeeder{m: map[string]fetcher.Page{}}
	p := Crawl{f, returnTrue}
	expected := map[string]fetcher.Page{}
	assert.Equal(t, expected, p.CrawlUsingSync(a))
}

func TestCrawl_CrawlUsingSync_AddsAssetsToPages(t *testing.T) {
	f := TestFeeder{m: withAssets}
	p := Crawl{f, returnTrue}
	expected := withAssets
	assert.Equal(t, expected, p.CrawlUsingSync(a))
}

//
// TestChanler_CrawlUsingChannels
//

func TestChanler_CrawlUsingChannels_DoesNotHangWhenAPageReferencesItself(t *testing.T) {
	f := TestFeeder{m: selfReference}
	p := Chanler{f, returnTrue}
	expected :=  selfReference
	assert.Equal(t, expected, p.CrawlUsingChannels(a))
}

func TestChanler_CrawlUsingChannels_DoesNotHangWhenA2PageCircularDependencyExists(t *testing.T) {
	f := TestFeeder{m: circularReferences}
	p := Chanler{f, returnTrue}
	expected := circularReferences
	assert.Equal(t, expected, p.CrawlUsingChannels(a))
}

func TestChanler_CrawlUsingChannels_DoesNotFollowInvalidLinks(t *testing.T) {
	f := TestFeeder{m: circularReferences}
	p := Chanler{f, returnFalse}
	expected := map[string]fetcher.Page{
		a: fetcher.NewPage(a).WithLinks(b),
	}
	assert.Equal(t, expected, p.CrawlUsingChannels(a))
}

func TestChanler_CrawlUsingChannels_AddPageWhenFetcherReturnsError(t *testing.T) {
	f := TestFeeder{m: map[string]fetcher.Page{}}
	p := Chanler{f, returnTrue}
	expected := map[string]fetcher.Page{}
	assert.Equal(t, expected, p.CrawlUsingChannels(a))
}

func TestChanler_CrawlUsingChannels_AddPageWhenFetcherReturnsErrorOnSecondLink(t *testing.T) {
	f := TestFeeder{m: map[string]fetcher.Page{}}
	p := Chanler{f, returnTrue}
	expected := map[string]fetcher.Page{}
	assert.Equal(t, expected, p.CrawlUsingChannels(a))
}

func TestChanler_CrawlUsingChannels_AddsAssetsToPages(t *testing.T) {
	f := TestFeeder{m: withAssets}
	p := Chanler{f, returnTrue}
	expected := withAssets
	assert.Equal(t, expected, p.CrawlUsingChannels(a))
}
