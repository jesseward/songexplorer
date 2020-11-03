package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jesseward/songexplorer/metrics"

	"github.com/campoy/apiai"
	"github.com/jesseward/songexplorer/app"
	"github.com/jesseward/songexplorer/caches/redis"
	"github.com/jesseward/songexplorer/config"
	"github.com/jesseward/songexplorer/sources/lastfm"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	httpBindPort := flag.Int("port", config.LookupEnvOrInt("port", 64738), "service bind port")
	httpBindAddress := flag.String("address", config.LookupEnvOrStr("address", ""), "sevice bind address")
	debugBindPort := flag.Int("debugport", config.LookupEnvOrInt("debugport", 49152), "debug bind port")
	debugBindAddress := flag.String("debugaddress", config.LookupEnvOrStr("debugaddress", ""), "debug bind address")
	lastFMAPIKey := flag.String("lastfmapikey", config.LookupEnvOrStr("lastfmapikey", ""), "Source API key")
	lastFMSharedSecret := flag.String("lastfmsharedsecret", config.LookupEnvOrStr("lastfmsharedsecret", ""), "Shared secret for Source")
	cacheHost := flag.String("cachehost", config.LookupEnvOrStr("cachehost", ""), "address of the cache (redis) server")
	cachePort := flag.Int("cacheport", config.LookupEnvOrInt("cacheport", 6380), "redis port")
	cacheTTL := flag.String("cachettl", config.LookupEnvOrStr("cachettl", "30d"), "cache expiration in human readable format")
	cacheSecret := flag.String("cachesecret", config.LookupEnvOrStr("cachesecret", ""), "secret key for cache.")
	loglocation := flag.String("logfile", config.LookupEnvOrStr("logfile", ""), "log file location. If omitted logs are written to STDOUT")
	maxArtist := flag.Int("maxartists", config.LookupEnvOrInt("maxartists", 10), "maximum number of artists")
	maxTracks := flag.Int("maxtracks", config.LookupEnvOrInt("maxtracks", 10), "maximum number of tracks")
	maxArtistTracks := flag.Int("maxartisttracks", config.LookupEnvOrInt("maxartisttracks", 10), "max artists top tracks")
	flag.Parse()
	log.SetOutput(os.Stdout)
	// note that preference is to log to STDOUT to be caught by container logging.
	if *loglocation != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   *loglocation,
			MaxSize:    5, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	}

	cfg := &config.Config{
		HTTPBindAddress:     *httpBindAddress,
		HTTPBindPort:        *httpBindPort,
		HTTPDebugBindAddres: *debugBindAddress,
		HTTPDebugPort:       *debugBindPort,
		SourceAPIKey:        *lastFMAPIKey,
		SourceSharedSecret:  *lastFMSharedSecret,
		CacheHost:           *cacheHost,
		CachePort:           *cachePort,
		CacheSecret:         *cacheSecret,
		LogFileLocation:     *loglocation,
		MaxTopArtistTracks:  *maxArtistTracks,
		MaxTopSimArtists:    *maxArtist,
		MaxTopSimTracks:     *maxTracks,
	}
	cfg.SetCacheDuration(*cacheTTL)
	log.Printf("ls: %s, lk: %s, ch: %s, cs: %s", cfg.SourceSharedSecret, cfg.SourceAPIKey, cfg.CacheHost, cfg.CacheSecret)

	// abort boot if we're unable to create a new source API client
	src, err := lastfm.New(cfg)
	if err != nil {
		log.Fatalf("unable to create metadata client, %v", err)
	}

	// create a Redis cache instance
	cache := redis.New(cfg)

	m, mux := metrics.NewAdminServletMetrics()

	a := app.App{
		Source:  src,
		Cache:   cache,
		Config:  cfg,
		Metrics: m,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	// launches metrics httpd service
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.HTTPDebugBindAddres, cfg.HTTPDebugPort)
		log.Printf("INFO :: main :: starting metrics service on %s", addr)
		srv := http.Server{
			Addr:         addr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			Handler:      mux,
		}
		log.Fatal(srv.ListenAndServe())
		wg.Done()
	}()

	wg.Add(1)
	// launches api.ai handler.
	go func() {
		h := apiai.NewHandler()
		h.Register("artist-recommendation", a.ArtistSimilar)
		h.Register("artist-bio", a.ArtistBio)
		h.Register("artist-top-tracks", a.ArtistTopTracks)
		h.Register("song-recomendations", a.TrackSimilar)

		addr := fmt.Sprintf("%s:%d", cfg.HTTPBindAddress, cfg.HTTPBindPort)
		log.Printf("INFO :: main :: starting songexplorer service on %s", addr)
		srv := http.Server{
			Addr:         addr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			Handler:      h,
		}
		log.Fatal(srv.ListenAndServe())
		wg.Done()
	}()
	wg.Wait()
}
