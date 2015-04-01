package carnegie

import (
	consulApi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"net/url"
	"strconv"
)

// BILL: Why not embed the pointer since this is the only
//       type associated with this struct.
type ConsulBackend struct {
	*consulApi.Client
}

func NewConsulBackend(config *viper.Viper) (*ConsulBackend, error) {
	consulConfig := consulApi.Config{
		Address: config.GetString("address"),
		Scheme:  config.GetString("scheme"),
		Token:   config.GetString("token"),
	}

	consulClient, err := consulApi.NewClient(&consulConfig)
	if err != nil {
		return nil, err
	}

	return &ConsulBackend{
		Client: consulClient,
	}, nil
}

func (cb *ConsulBackend) GetBackends(host string) ([]*url.URL, error) {
	serviceEntries, _, err := cb.Health().Service(host, "", true, nil)
	if err != nil {
		return nil, err
	}
	
	// BILL: Why not use a for/range on serviceEntries?
	urls := make([]*url.URL, len(serviceEntries))
	for i, se := range serviceEntries {
		address := fmt.Sprintf("http://%s", se.Node.Address)
		if se.Service.Port != 0 {
			address += ":" + strconv.Itoa(se.Service.Port)
		}
		if urls[i], err = url.Parse(address); err != nil {
			return nil, err
		}
	}

	//totalEntries := len(serviceEntries)
	//urls := make([]*url.URL, totalEntries)
	//for i := 0; i < totalEntries; i++ {
	//	address := fmt.Sprintf("http://%s", serviceEntries[i].Node.Address)
	//	if port := serviceEntries[i].Service.Port; port != 0 {
	//		address += ":" + strconv.Itoa(port)
	//	}
	//	if urls[i], err = url.Parse(address); err != nil {
	//		return nil, err
	//	}
	//}
	
	return urls, nil
}
