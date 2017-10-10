package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	HTTPBindPort       uint16 // port we listen and serve
	HTTPBindAddress    string // IP or hostname to bind to
	LastFMAPIKey       string // LAST.FM api key
	LastFMSharedSecret string // LAST.fm shared secret
	RedisHost          string // Redis Host
	RedisPort          uint16 // Redis port
	RedisKeyExpiry     string // default Redis key expiration
	LogFileLocation    string // logfile location
	MaxTopSimArtists   int
	MaxTopSimTracks    int
	MaxTopArtistTracks int
}

// New returns configuration data read from the target config file
func New(path string) (*Config, error) {
	var conf Config

	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return &conf, fmt.Errorf("unable to load config: %v", err)
	}
	return &conf, nil
}
