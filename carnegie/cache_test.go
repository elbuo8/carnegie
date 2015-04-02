package carnegie

import (
	LRU "github.com/hashicorp/golang-lru"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func TestNewCache(t *testing.T) {
	config := viper.New()
	// Test for unsupported backend
	cache, err := NewCache(config)
	if err == nil {
		t.Fatalf("should report unsupported cache")
	}
	// Supported backend
	config.Set("backend", "consul")
	cache, err = NewCache(config)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if cache == nil {
		t.Fatalf("cache shouldnt be nil")
	}
}

func TestGetAddresses(t *testing.T) {
	localCache, _ := LRU.New(64)
	backend := &backendInventoryMock{
		Hosts: map[string][]string{
			"google.com": []string{"https://google.com"},
		},
	}
	cache := &Cache{
		LocalInventory: localCache,
		Backend:        backend,
	}
	// Should return error if backend threw and error
	_, err := cache.GetAddresses("error")
	if err == nil {
		t.Fatalf("err should be returned if backend returns error")
	}
	// Should return error on a non existent backend
	_, err = cache.GetAddresses("test")
	if err == nil {
		t.Fatalf("err should be returned if no backends exist")
	}
	// Should retrieve from Backend and store in Cache
	backendURLS, err := cache.GetAddresses("google.com")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if cache.LocalInventory.Len() != 1 {
		t.Fatalf("local cache shouldnt be empty")
	}
	cacheURLS, err := cache.GetAddresses("google.com")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(cacheURLS, backendURLS) {
		t.Fatalf("arrays should match")
	}
}

func TestRotateAddresses(t *testing.T) {
	localCache, _ := LRU.New(64)
	backend := &backendInventoryMock{
		Hosts: map[string][]string{
			"google.com": []string{"https://google.com"},
		},
	}
	cache := &Cache{
		LocalInventory: localCache,
		Backend:        backend,
	}
	// Should throw error if no host is found
	err := cache.RotateAddresses("google.com")
	if err == nil {
		t.Fatalf("Error should exist")
	}
	cache.GetAddresses("google.com") // Stores into localcache
	err = cache.RotateAddresses("google.com")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestInvalidate(t *testing.T) {
	localCache, _ := LRU.New(64)
	backend := &backendInventoryMock{}
	cache := &Cache{
		LocalInventory: localCache,
		Backend:        backend,
	}
	cache.Invalidate("google.com")
}
