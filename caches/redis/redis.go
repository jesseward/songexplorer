package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/jesseward/songexplorer/config"
)

var ctx = context.Background()

type Redis struct {
	cli    *redis.Client
	expiry time.Duration
}

// Get fetches and returns a (string) value from the store.
func (c *Redis) Get(k string) (string, error) {
	v, err := c.cli.Get(ctx, k).Result()
	if err != nil {
		return "", err
	}
	return v, nil
}

// Set writes a key->value to the key store.
func (c *Redis) Set(k, v string, x time.Duration) (bool, error) {
	err := c.cli.Set(ctx, k, v, x).Err()
	if err != nil {
		return false, fmt.Errorf("unable to set key: %s, error: %v", k, err)
	}
	return true, nil
}

// New returns a Redis (Cache interface) object.
func New(cfg *config.Config) *Redis {
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
	r := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", cfg.CacheHost, cfg.CachePort),
		Password: cfg.CacheSecret, TLSConfig: tlsConfig})
	return &Redis{r, cfg.CacheTTLDuration}
}
