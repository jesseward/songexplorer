package config

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	HTTPBindPort        uint16 // port we listen and serve
	HTTPBindAddress     string // IP or hostname to bind to
	HTTPDebugBindAddres string // Bind address for internal metrics/admin server
	HTTPDebugPort       uint16 // TCP port for admin/stats service.
	SourceAPIKey        string // LAST.FM api key
	SourceSharedSecret  string // LAST.fm shared secret
	CacheHost           string // Cache Host
	CachePort           uint16 // Cache port
	CacheTTL            string // default cache expiration in human readable format.
	LogFileLocation     string // logfile location
	MaxTopSimArtists    int
	MaxTopSimTracks     int
	MaxTopArtistTracks  int
	CacheTTLDuration    time.Duration // used internally after sting->duration conversion
}

// New returns configuration data read from the target config file
func New(path string) (*Config, error) {
	var conf Config

	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return &conf, fmt.Errorf("unable to load config: %v", err)
	}

	// cache duration is 0 if not explicitly set.
	conf.CacheTTLDuration = 0
	if t, err := time.ParseDuration(conf.CacheTTL); err == nil {
		conf.CacheTTLDuration = t
	}
	return &conf, nil
}
