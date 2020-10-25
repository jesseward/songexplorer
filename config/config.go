package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPBindPort        int    // port we listen and serve
	HTTPBindAddress     string // IP or hostname to bind to
	HTTPDebugBindAddres string // Bind address for internal metrics/admin server
	HTTPDebugPort       int    // TCP port for admin/stats service.
	SourceAPIKey        string // LAST.FM api key
	SourceSharedSecret  string // LAST.fm shared secret
	CacheHost           string // Cache Host
	CachePort           int    // Cache port
	CacheSecret         string // secret key for cache connectivity
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

// LookupEnvOrStr attempts to fetch configuration parameter from environmnet variable and assign result as a string
func LookupEnvOrStr(k string, d string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return d
}

// LookupEnvOrInt attempts to fetch a configuration parameter via environment variable and convert to an int
func LookupEnvOrInt(k string, d int) int {
	if val, ok := os.LookupEnv(k); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("failed to convert LookupEnvOrInt(%s), val=%s", k, val)
		}
		return v
	}
	return d
}
