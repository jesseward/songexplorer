package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/campoy/apiai"

	"github.com/jesseward/songexplorer/app"
	"github.com/jesseward/songexplorer/config"
	"github.com/jesseward/songexplorer/internal/caches/redis"
	"github.com/jesseward/songexplorer/internal/sources/lastfm"
	"gopkg.in/natefinch/lumberjack.v2"
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

	// abort if we're unable to create a new source API client
	src, err := lastfm.New(cfg)
	if err != nil {
		log.Fatalf("unable to create metadata client, %v", err)
	}

	// create our cache client
	cache := redis.New(cfg)

	app := app.App{
		Source: src,
		Cache:  cache,
		Config: cfg,
	}

	h := apiai.NewHandler()
	h.Register("artist-recommendation", app.ArtistSimilar)
	h.Register("artist-bio", app.ArtistBio)
	h.Register("artist-top-tracks", app.ArtistTopTracks)
	h.Register("song-recomendations", app.TrackSimilar)

	log.Printf("INFO :: main :: starting http service on %s:%d", cfg.HTTPBindAddress, cfg.HTTPBindPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.HTTPBindAddress, cfg.HTTPBindPort), h))
}
