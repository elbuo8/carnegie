package carnegie

import (
	"net/url"
)

type BackendInventoryMock struct{}

func (mock *BackendInventoryMock) GetBackends(host string) ([]*url.URL, error) {
	var parsedURL *url.URL
	var urls []*url.URL
	if host == "google.com" {
		parsedURL, _ = url.Parse("http://google.com")
		urls = append(urls, parsedURL)
	}
	return urls, nil
}
