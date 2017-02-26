package fetcher

import (
	"testing"
	"net/http"
	"errors"
	"github.com/stretchr/testify/assert"
	"fmt"
	"bytes"
	"io"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

var returnError = func(url string) (resp *http.Response, err error) {
	return nil, errors.New("fail")
}

func returnString(bodyString string) func(url string) (resp *http.Response, err error) {
	return func(url string) (resp *http.Response, err error) {
		return body(bodyString), nil
	}
}

func body(bodyString string) *http.Response {
	b := nopCloser{bytes.NewBufferString(bodyString)}
	fmt.Println(b)
	return &http.Response{Body: b}
}

func TestWhenClientReturnsErrorFetcherReturnsError(t *testing.T) {
	f := NewWebFetcher(returnError)
	_, err := f.getPage("hello", "hello")
	assert.NotNil(t, err)
}

func TestWhenClientReturnsEmptyBodyReturnsEmptyPage(t *testing.T) {
	f := NewWebFetcher(returnString(""))
	page, err := f.getPage("hello", "hello")

	expectedLinks := []string{}
	expectedAssets := []string{}

	assert.Nil(t, err)
	assert.Equal(t, expectedLinks, page.Links)
	assert.Equal(t, expectedAssets, page.Assets)
}

func TestWhenClientReturnsOneHrefReturnedAsLink(t *testing.T) {
	f := NewWebFetcher(returnString("<html><a href=\"link1\" </html>"))
	page, err := f.getPage("hello", "hello")

	expectedLinks := []string{"link1"}

	assert.Nil(t, err)
	assert.Equal(t, expectedLinks, page.Links)
}

func TestWhenClientReturnsHrefInOtherElementThanANotReturnedAsLink(t *testing.T) {
	f := NewWebFetcher(returnString("<html><b href=\"link1\" </html>"))
	page, err := f.getPage("hello", "hello")

	expectedLinks := []string{}

	assert.Nil(t, err)
	assert.Equal(t, expectedLinks, page.Links)
}

func TestYo(t *testing.T) {
	f := NewWebFetcher(http.Get)
	page, _ := f.getPage("http://tomblomfield.com/", "hello")

	fmt.Println(page)
	assert.NotNil(t, nil)
}


