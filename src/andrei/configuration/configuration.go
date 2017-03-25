package configuration

import (
	"net/url"
	"fmt"
	"strings"
)


func HostUrlValidator(host string) func(url string) bool {
	hostUrl, err := url.Parse(host)
	if err != nil {
		panic(fmt.Sprintf("Invalid url %s", host))
	}

	return func(url string) bool {
		return hasHost(hostUrl.Host, url)
	}
}

func AllValid() func(url string) bool {
	return func(url string) bool {
		return true
	}
}

func AbsolutePathBuilder(host string) func(url string) string {
	hostUrl, err := url.Parse(host)
	if err != nil {
		panic(fmt.Sprintf("Invalid url %s", host))
	}

	return func(url string) string {
		return relative(hostUrl, url)
	}
}

func hasHost(host string, urlString string) bool {
	u, err := url.Parse(urlString)

	if err != nil {
		return false
	}

	return u.Host == host
}

func relative(base *url.URL, pathOrUrl string) string {
	if (strings.HasPrefix(pathOrUrl, "/")) {
		return fmt.Sprintf("%s://%s%s", base.Scheme, base.Host, pathOrUrl)
	} else {
		return pathOrUrl
	}
}
