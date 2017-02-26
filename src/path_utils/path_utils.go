package path_utils

import (
	"net/url"
	"fmt"
)


//improve efficiency by parsing host only once
func HostUrlValidator(host string) func(url string) bool {
	hostUrl, err := url.Parse(host)
	if err != nil {
		panic(fmt.Sprintf("Invalid url %s", host))
	}

	return func(url string) bool {
		return hasHost(hostUrl.Host, url)
	}
}

func hasHost(host string, urlString string) bool {
	u, err := url.Parse(urlString)

	if err != nil {
		return false
	}

	return u.Host == host
}