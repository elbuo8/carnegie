package carnegie

import (
	"github.com/spf13/viper"
	"testing"
)

func TestNewConsulBackend(t *testing.T) {
	config := viper.New()
	backend, err := NewConsulBackend(config)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if backend == nil {
		t.Fatalf("backend is nil")
	}
}

/*
func TestGetBackends(t *testing.T) {
	config := viper.New()
	backend, err := NewConsulBackend(config)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// Figure out testing for failure in consul
	// Test for pass
	_, err = backend.GetBackends("test.com")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}
*/
