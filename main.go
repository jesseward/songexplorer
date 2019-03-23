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

	cfgFlg := flag.String("config", "/tmp/config.toml", "songexplorer configuration file location")
	flag.Parse()
	log.SetOutput(os.Stdout)
	// abort if we're unable to parse configuration data.
	cfg, err := config.New(*cfgFlg)
	if err != nil {
		log.Fatalf("unable to load config file %s, %v", *cfgFlg, err)
	}

	log.SetOutput(&lumberjack.Logger{
		Filename:   cfg.LogFileLocation,
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

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
