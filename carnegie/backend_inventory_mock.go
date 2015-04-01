package carnegie

import (
	"net/url"
)

type BackendInventoryMock struct {
	Hosts map[string][]string
}

func (mock *BackendInventoryMock) GetBackends(host string) ([]*url.URL, error) {
	var parsedURL *url.URL
	var urls []*url.URL
	if backends, ok := mock.Hosts[host]; ok {
		for i := 0; i < len(backends); i++ {
			parsedURL, _ = url.Parse(backends[i])
			urls = append(urls, parsedURL)
		}
	}
	return urls, nil
}
