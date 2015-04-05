package carnegie

import (
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"net/url"
	"strconv"
)

// ConsulBackend is a BackendInventory that interfaces with Consul.
type ConsulBackend struct {
	*consulApi.Client
}

// NewConsulBackend returns a ConsulBackend with the provided configuration.
func NewConsulBackend(config *viper.Viper) (*ConsulBackend, error) {
	consulConfig := consulApi.Config{
		Address: config.GetString("address"),
		Scheme:  config.GetString("scheme"),
		Token:   config.GetString("token"),
	}
	// NewClient never returns an error
	consulClient, _ := consulApi.NewClient(&consulConfig)

	return &ConsulBackend{
		Client: consulClient,
	}, nil
}

// GetBackends returns accessible backends for a specified service in Consul
func (cb *ConsulBackend) GetBackends(host string) ([]*url.URL, error) {
	serviceEntries, _, err := cb.Health().Service(host, "", true, nil)
	if err != nil {
		return nil, err
	}

	urls := make([]*url.URL, len(serviceEntries))
	for i, se := range serviceEntries {
		address := fmt.Sprintf("http://%s", se.Node.Address)
		if se.Service.Port != 0 {
			address += ":" + strconv.Itoa(se.Service.Port)
		}
		// Addresses from Consul are always valid
		urls[i], _ = url.Parse(address)
	}

	return urls, nil
}
