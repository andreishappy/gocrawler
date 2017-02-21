package crawler

//package management
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"errors"
	"fmt"
)

type TestFeeder struct {
	m map[string]*FetchedPage
}

var m = map[string]*FetchedPage{
	"hello": {Links: []string{"yo"}, Assets: []string{}},
	"yo": {Links: []string{"my"}, Assets: []string{}},
}

var selfReference = map[string]*FetchedPage{
	"hello": {Links: []string{"hello"}, Assets: []string{}},
}

var circularReferences = map[string]*FetchedPage{
	"hello": {Links: []string{"hi"}, Assets: []string{}},
	"hi": {Links: []string{"hello"}, Assets: []string{}},
}

var withAssets = map[string]*FetchedPage{
	"hello": {Links: []string{"hi"}, Assets: []string{"helloAsset"}},
	"hi": {Links: []string{"go"}, Assets: []string{"hiAsset1", "hiAsset2"}},
}

func (t TestFeeder) getLinks(url string) (error, *FetchedPage) {
	elem, ok := t.m[url]
	if (ok) {
		return nil, elem
	} else {
		return errors.New(fmt.Sprint("Error when fetching ", url)), nil
	}
}

func TestDoesNotHangWhenAPageReferencesItself(t *testing.T) {
	f := TestFeeder{m: selfReference}
	p := Crawl{f}
	assert.Equal(t, map[string]PageNode{"hello": {Links: []string{"hello"}, Assets: []string{}}}, p.crawl("hello"))
}

func TestDoesNotHangWhenA2PageCircularDependencyExists(t *testing.T) {
	f := TestFeeder{m: circularReferences}
	p := Crawl{f}
	expected := map[string]PageNode{
		"hello": {Links: []string{"hi"}, Assets: []string{}},
		"hi": {Links: []string{"hello"}, Assets: []string{}},
	}
	assert.Equal(t, expected, p.crawl("hello"))
}

func TestAddPageWhenFetcherReturnsError(t *testing.T) {
	f := TestFeeder{m: map[string]*FetchedPage{}}
	p := Crawl{f}
	expected := map[string]PageNode{}
	assert.Equal(t, expected, p.crawl("hello"))
}

func TestAddsAssetsToPages(t *testing.T) {
	f := TestFeeder{m: withAssets}
	p := Crawl{f}
	expected := map[string]PageNode{
		"hello": {Links: []string{"hi"}, Assets: []string{"helloAsset"}},
		"hi": {Links: []string{"go"}, Assets: []string{"hiAsset1", "hiAsset2"}},
	}
	assert.Equal(t, expected, p.crawl("hello"))
}
