package carnegie

import (
	"github.com/spf13/viper"
	"testing"
)

func TestNewBackend(t *testing.T) {
	config := viper.New()
	// Test unsupported backend
	backend, err := NewBackend("c", config)
	if err == nil {
		t.Fatalf("backend should not exist")
	}
	// Consul backend
	backend, err = NewBackend("consul", config)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if backend == nil {
		t.Fatalf("backend is nil")
	}
}
