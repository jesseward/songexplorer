package config

import (
	"time"
)

type Config struct {
	HTTPBindPort        uint   // port we listen and serve
	HTTPBindAddress     string // IP or hostname to bind to
	HTTPDebugBindAddres string // Bind address for internal metrics/admin server
	HTTPDebugPort       uint   // TCP port for admin/stats service.
	SourceAPIKey        string // LAST.FM api key
	SourceSharedSecret  string // LAST.fm shared secret
	CacheHost           string // Cache Host
	CachePort           uint   // Cache port
	LogFileLocation     string // logfile location
	MaxTopSimArtists    int
	MaxTopSimTracks     int
	MaxTopArtistTracks  int
	CacheTTLDuration    time.Duration // used internally after sting->duration conversion
}

// SetCacheDuration returns configuration data read from the target config file
func (c *Config) SetCacheDuration(ttl string) {

	c.CacheTTLDuration = 0
	if t, err := time.ParseDuration(ttl); err == nil {
		c.CacheTTLDuration = t
	}
}
