package carnegie

import (
	consulApi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"net/url"
	"strconv"
)

// ConsulBackend is a BackendInventory that interfaces with Consul.
type ConsulBackend struct {
	Client *consulApi.Client
}

// NewConsulBackend returns a ConsulBackend with the provided configuration.
func NewConsulBackend(config *viper.Viper) (*ConsulBackend, error) {
	consulConfig := &consulApi.Config{
		Address: config.GetString("address"),
		Scheme:  config.GetString("scheme"),
		Token:   config.GetString("token"),
	}

	consulClient, err := consulApi.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	return &ConsulBackend{
		Client: consulClient,
	}, nil
}

// GetBackends returns accessible backends for a specified service in Consul
func (cb *ConsulBackend) GetBackends(host string) ([]*url.URL, error) {
	health := cb.Client.Health()
	serviceEntries, _, err := health.Service(host, "", true, nil)
	if err != nil {
		return nil, err
	}
	totalEntries := len(serviceEntries)
	urls := make([]*url.URL, totalEntries)
	for i := 0; i < totalEntries; i++ {
		address := "http://"
		address += serviceEntries[i].Node.Address
		if port := serviceEntries[i].Service.Port; port != 0 {
			address += ":" + strconv.Itoa(port)
		}
		parsedURL, err := url.Parse(address)
		if err != nil {
			return nil, err
		}
		urls[i] = parsedURL
	}
	return urls, nil
}
