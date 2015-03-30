package carnegie

import (
	"errors"
	LRU "github.com/hashicorp/golang-lru"
	"github.com/spf13/viper"
	"net/url"
)

type BackendInventory interface {
	GetBackends(string) ([]*url.URL, error)
	//RemoveBackend(string) error
}

type Cache struct {
	Backend        BackendInventory
	LocalInventory *LRU.Cache
}

func NewCache(config *viper.Viper) (*Cache, error) {
	cache, err := LRU.New(128)
	if err != nil {
		return nil, err
	}
	backend, err := NewBackend(config.GetString("backend"), config)
	if err != nil {
		return nil, err
	}
	return &Cache{
		Backend:        backend,
		LocalInventory: cache,
	}, nil
}

func (c *Cache) GetAddresses(host string) ([]*url.URL, error) {
	if raw, ok := c.LocalInventory.Get(host); ok {
		cached := raw.([]*url.URL)
		return cached, nil
	}
	addresses, err := c.Backend.GetBackends(host)
	if err != nil {
		return nil, err
	}
	if len(addresses) == 0 {
		return nil, errors.New("no backend found")
	}
	c.LocalInventory.Add(host, addresses)
	return addresses, nil
}

func (c *Cache) RotateAddresses(host string) error {
	if raw, ok := c.LocalInventory.Get(host); ok {
		cached := raw.([]*url.URL)
		cached = append(cached, cached[0])
		cached = cached[1:]
		c.LocalInventory.Add(host, cached)
		return nil
	}
	return errors.New("No addresses associated with host")
}

func (c *Cache) Invalidate(host string) {
	c.LocalInventory.Remove(host)
}
