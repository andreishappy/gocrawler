package chanler

//package management
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"errors"
	"fmt"
	"fetcher"
)

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
	a: fetcher.NewPage(a).WithLinks(a, b),
	b: fetcher.NewPage(b).WithLinks(b, a),
}

var withAssets = map[string]fetcher.Page{
	a: fetcher.NewPage(a).WithLinks(b).WithAssets(assetA1, assetA2),
	b: pageWithLinks(b, c).WithAssets(assetB1),
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

func TestDoesNotHangWhenAPageReferencesItself(t *testing.T) {
	f := TestFeeder{m: selfReference}
	p := Chanler{f, returnTrue}
	assert.Equal(t, map[string]fetcher.Page{a: pageWithLinks(a, b)}, p.Crawl(a))
}

func TestDoesNotHangWhenA2PageCircularDependencyExists(t *testing.T) {
	f := TestFeeder{m: circularReferences}
	p := Chanler{f, returnTrue}
	expected := map[string]fetcher.Page{
		"hello": {Url: "hello", Links: []string{"hi"}, Assets: []string{}},
		"hi": {Url: "hi", Links: []string{"hello"}, Assets: []string{}},
	}
	assert.Equal(t, expected, p.Crawl("hello"))
}

func TestDoesNotFollowInvalidLinks(t *testing.T) {
	f := TestFeeder{m: circularReferences}
	p := Chanler{f, returnFalse}
	expected := map[string]fetcher.Page{
		"hello": {Url: "hello", Links: []string{"hi"}, Assets: []string{}},
	}
	assert.Equal(t, expected, p.Crawl("hello"))
}

func TestAddPageWhenFetcherReturnsError(t *testing.T) {
	f := TestFeeder{m: map[string]fetcher.Page{}}
	p := Chanler{f, returnTrue}
	expected := map[string]fetcher.Page{}
	assert.Equal(t, expected, p.Crawl("hello"))
}

func TestAddsAssetsToPages(t *testing.T) {
	f := TestFeeder{m: withAssets}
	p := Chanler{f, returnTrue}
	expected := map[string]fetcher.Page{
		"hello": {Links: []string{"hi"}, Assets: []string{"helloAsset"}},
		"hi": {Links: []string{"go"}, Assets: []string{"hiAsset1", "hiAsset2"}},
	}
	assert.Equal(t, expected, p.Crawl("hello"))
}
