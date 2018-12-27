package caches

import (
	"time"
)

// Cache defines the minimum contract required to support a 3rd party data store
type Cache interface {
	// Get receives a key and returns a string or error if not found
	Get(k string) (string, error)
	// Set accepts the key, value and expiry times
	Set(k, v string, x time.Duration) (bool, error)
}
