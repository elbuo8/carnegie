package carnegie

import (
	"errors"
	LRU "github.com/hashicorp/golang-lru"
	"github.com/spf13/viper"
	"net/url"
)

// Cache uses an internal LRU cache along with a BackendInventory
// for quick access to different VHOSTs.
type Cache struct {
	Backend        BackendInventory
	LocalInventory *LRU.Cache
}

// NewCache returns a new Cache with the provided configuration.
func NewCache(config *viper.Viper) (*Cache, error) {
	cache, _ := LRU.New(128)
	backend, err := NewBackend(config.GetString("backend"), config)
	if err != nil {
		return nil, err
	}
	return &Cache{
		Backend:        backend,
		LocalInventory: cache,
	}, nil
}

// GetAddresses returns a list of URLs pointing to the VHOST provided.
// If no VHOST was found an error is thrown. It will attempt to look in the
// LRU cache then move towards the BackedInventory.
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

// RotateAddresses moves the most recent used backend for a VHOST towards the
// end of the queue.
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

// Invalidate removes a VHOST from the LRU.
func (c *Cache) Invalidate(host string) {
	c.LocalInventory.Remove(host)
}
