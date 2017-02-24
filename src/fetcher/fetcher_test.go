package fetcher

import (
	"testing"
)

func TestYo(t *testing.T) {
	w := WebFetcher{}
	w.getPage("http://tomblomfield.com/")
	t.Error("yo")
}
