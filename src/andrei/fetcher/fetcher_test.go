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

//
// Setup
//

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func returnString(bodyString string) func(url string) (resp *http.Response, err error) {
	return func(url string) (resp *http.Response, err error) {
		return body(bodyString), nil
	}
}

var returnError = func(url string) (resp *http.Response, err error) {
	return nil, errors.New("fail")
}

func body(bodyString string) *http.Response {
	b := nopCloser{bytes.NewBufferString(bodyString)}
	return &http.Response{Body: b}
}

func returnTrue(string) bool {
	return true
}

func returnTrueFor(one string) func(string) bool {
	return func(other string) bool {
		return one == other
	}
}

func returnFalse(string) bool {
	return false
}

func identity(input string) string {
	return input
}

//
// WebFetcherTests
//

func TestWhenClientReturnsErrorFetcherReturnsError(t *testing.T) {
	f := NewWebFetcher(returnError, returnTrue, identity)
	_, err := f.GetPage("hello")
	assert.NotNil(t, err)
}

func TestWhenClientReturnsEmptyBodyReturnsEmptyPage(t *testing.T) {
	f := NewWebFetcher(returnString(""), returnTrue, identity)
	page, err := f.GetPage("hello")

	expectedLinks := []string{}
	expectedAssets := []string{}

	assert.Nil(t, err)
	assert.Equal(t, "hello", page.Url)
	assert.Equal(t, expectedLinks, page.Links)
	assert.Equal(t, expectedAssets, page.Assets)
}

func TestWhenClientReturnsHrefsReturnedAsLinks(t *testing.T) {
	f := NewWebFetcher(returnString("<html> <a href=\"link1\"> <a href=\"link2\"> </html>"), returnTrue, identity)
	page, err := f.GetPage("hello")

	assert.Nil(t, err)
	assert.Contains(t, page.Links, "link1")
	assert.Contains(t, page.Links, "link2")
	assert.Len(t, page.Links, 2)
}

func TestDoesNotAddDuplicateLinks(t *testing.T) {
	f := NewWebFetcher(returnString("<html> <a href=\"link1\"></a> <a href=\"link1\"></a> </html>"), returnTrue, identity)
	page, err := f.GetPage("hello")

	expectedLinks := []string{"link1"}

	assert.Nil(t, err)
	assert.Equal(t, expectedLinks, page.Links)
}

func TestWhenClientReturnsOneHrefAndValidatorReturnsFalseDoesNotAddLink(t *testing.T) {
	f := NewWebFetcher(returnString("<html> <a href=\"link1\"></a> <a href=\"link2\"></a> </html>"), returnFalse, identity)
	page, err := f.GetPage("hello")

	expectedLinks := []string{}

	assert.Nil(t, err)
	assert.Equal(t, expectedLinks, page.Links)
}

func TestWhenClientReturnsHrefInOtherElementThanANotReturnedAsLink(t *testing.T) {
	f := NewWebFetcher(returnString("<html> <b href=\"link1\"> </b> </html>"), returnTrue, identity)
	page, err := f.GetPage("hello")

	expectedLinks := []string{}

	assert.Nil(t, err)
	assert.Equal(t, expectedLinks, page.Links)
}

func TestLinksAddedAsAbsolute(t *testing.T) {
	absolute := func(string) string { return "absolute" }

	f := NewWebFetcher(returnString("<html> <a href=\"link1\"></a> </html>"), returnTrueFor("absolute"), absolute)
	page, err := f.GetPage("hello")

	expectedLinks := []string{"absolute"}

	assert.Nil(t, err)
	assert.Equal(t, expectedLinks, page.Links)
}

func TestWhenClientReturnsImageElementItIsReturnedAsAsset(t *testing.T) {
	f := NewWebFetcher(returnString("<html> <img src=\"asset1\"> </html>"), returnTrue, identity)
	page, err := f.GetPage("hello")

	expectedAssets := []string{"asset1"}

	assert.Nil(t, err)
	assert.Equal(t, expectedAssets, page.Assets)
}


