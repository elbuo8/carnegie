package carnegie

import (
	"net/url"
)

// BackendInventory interface mplements interactions with third party backends.
type BackendInventory interface {
	GetBackends(string) ([]*url.URL, error)
	//RemoveBackend(string) error
}
