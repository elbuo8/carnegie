package carnegie

import (
	"errors"
	"github.com/spf13/viper"
)

// NewBackend returns specified BackendInventory if supported.
func NewBackend(typ string, config *viper.Viper) (BackendInventory, error) {
	var backend BackendInventory
	var err error
	switch typ {
	case "consul":
		backend, err = NewConsulBackend(config)
	default:
		err = errors.New("Backend not supported")
	}
	return backend, err
}
