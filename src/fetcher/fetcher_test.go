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

func returnTrue(string) bool {
	return true
}

func returnFalse(string) bool {
	return false
}

func identity(input string) string {
	return input
}

func body(bodyString string) *http.Response {
	b := nopCloser{bytes.NewBufferString(bodyString)}
	fmt.Println(b)
	return &http.Response{Body: b}
}

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

func TestWhenClientReturnsSameLinkTwiceReturnsOnlyOneLink(t *testing.T) {
	f := NewWebFetcher(returnString("<html> <a href=\"link1\"> <a href=\"link1\"> </html>"), returnTrue, identity)
	page, err := f.GetPage("hello")

	expectedLinks := []string{"link1"}

	assert.Nil(t, err)
	assert.Equal(t, expectedLinks, page.Links)
}

func TestWhenClientReturnsOneHrefAndValidatorReturnsFalseDoesNotAddLink(t *testing.T) {
	f := NewWebFetcher(returnString("<html> <a href=\"link1\"> <a href=\"link2\"> </html>"), returnFalse, identity)
	page, err := f.GetPage("hello")

	expectedLinks := []string{}

	assert.Nil(t, err)
	assert.Equal(t, expectedLinks, page.Links)
}

func TestWhenClientReturnsHrefInOtherElementThanANotReturnedAsLink(t *testing.T) {
	f := NewWebFetcher(returnString("<html> <b href=\"link1\"> </html>"), returnTrue, identity)
	page, err := f.GetPage("hello")

	expectedLinks := []string{}

	assert.Nil(t, err)
	assert.Equal(t, expectedLinks, page.Links)
}

func TestWhenTheLinkAddedIsTheRelativisedOne(t *testing.T) {
	relativiser := func(string) string { return "relative" }

	f := NewWebFetcher(returnString("<html> <a href=\"link1\"> </html>"), returnTrue, relativiser)
	page, err := f.GetPage("hello")

	expectedLinks := []string{"relative"}

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


