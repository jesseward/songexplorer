package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/jesseward/songexplorer/config"
)

type Redis struct {
	cli    *redis.Client
	expiry time.Duration
}

func (c *Redis) Get(k string) (string, error) {

	v, err := c.cli.Get(k).Result()

	if err != nil {
		return "", err
	}
	return v, nil
}

func (c *Redis) Set(k, v, x string) (bool, error) {
	expiry, err := time.ParseDuration(x)
	if err != nil {
		expiry = 0
	}
	err = c.cli.Set(k, v, expiry).Err()
	if err != nil {
		return false, fmt.Errorf("unable to set key: %s, value: %v, error: %v", k, v, err)
	}
	return true, nil
}

func New(cfg *config.Config) *Redis {
	r := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)})
	x, err := time.ParseDuration(cfg.RedisKeyExpiry)

	if err != nil {
		x = 0
	}
	return &Redis{r, x}
}
