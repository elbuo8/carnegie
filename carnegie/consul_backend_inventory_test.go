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
