package crawler

//package management
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type TestFeeder struct {
}

var m = map[string]Page{
	"hello": {Links: []string{"yo"}, Assets: []string{}},
}

func (t TestFeeder) getLinks(url string) Page {
	return m[url]
}

func TestFirst(t *testing.T) {

	f := TestFeeder{}

	p := Crawl{f}
	assert.Equal(t, p.crawl("hello"), PageNode{Url: "hello"}, "they should be equal")
}
