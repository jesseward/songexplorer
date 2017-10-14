package app

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/campoy/apiai"
	"github.com/jesseward/songexplorer/config"
	"github.com/jesseward/songexplorer/internal/caches"
	"github.com/jesseward/songexplorer/internal/sources"
)

type App struct {
	Source sources.Source // remote data source
	Cache  caches.Cache   // cache store client
	Config *config.Config // songexplorer configuration data
}

// TrackSimilar handles the similar tracks intent. Performs a cache lookup and then falls through to the
// external source on a cache miss.
func (a *App) TrackSimilar(ctx context.Context, req *apiai.Request) (*apiai.Response, error) {
	artist := req.Param("artist")
	if artist == "" {
		return nil, fmt.Errorf("missing artist parameter")
	}

	track := req.Param("track")
	if track == "" {
		return nil, fmt.Errorf("missing track parameter")
	}

	similarKey := "track_similar:" + artist + ":" + track
	var s string
	s, err := a.Cache.Get(similarKey)
	if err != nil || s == "" {
		log.Printf(":: INFO :: TrackSimilar :: fetching similar tracks from source for %s - %v", artist, track)
		similar, err := a.Source.GetSimilarTracks(artist, track)

		if err != nil || len(similar.SimilarTracks) < 1 {
			log.Printf(":: WARN :: TrackSimilar :: fetch failure from source for %s - %s, error=%v", artist, track, err)
			s = fmt.Sprintf("unable to fetch similar tracks for %s by %s", track, artist)
		} else {
			s = strings.Join(similar.SimilarTracks, ".")
			_, err := a.Cache.Set(similarKey, s, a.Config.CacheTTLDuration)
			if err != nil {
				log.Printf(":: ERROR :: TrackSimilar :: %v", err)
			}
		}

	}
	return &apiai.Response{
		Speech:      s,
		DisplayText: s,
	}, nil
}

// ArtistSimilar performs a similar artist lookup requested within the artist-recommendation intent.
func (a *App) ArtistSimilar(ctx context.Context, req *apiai.Request) (*apiai.Response, error) {

	artist := req.Param("artist")
	if artist == "" {
		return nil, fmt.Errorf("missing artist parameter")
	}

	artistKey := "artist_similar:" + strings.ToLower(artist)

	var s string
	s, err := a.Cache.Get(artistKey)
	if err != nil || s == "" {
		log.Printf(":: INFO :: ArtistSimilar :: fetching similar artist from source for artist=%s", artist)
		similar, err := a.Source.GetSimilarArtists(artist)

		if err != nil || len(similar.SimilarArtist) < 1 {
			log.Printf(":: WARN :: ArtistSimilar :: failure retrieving %s from source, %v", artist, err)
			s = fmt.Sprintf("unable to find similar artists for %s", artist)
		} else {
			s = strings.Join(similar.SimilarArtist, ",")
			_, err := a.Cache.Set(artistKey, s, a.Config.CacheTTLDuration)
			if err != nil {
				log.Printf(":: ERROR :: ArtistSimilar :: %v", err)
			}
		}
	}
	return &apiai.Response{
		Speech:      s,
		DisplayText: s,
	}, nil
}

func (a *App) ArtistTopTracks(ctx context.Context, req *apiai.Request) (*apiai.Response, error) {
	artist := req.Param("artist")

	if artist == "" {
		return nil, fmt.Errorf("missing artist parameter")
	}

	artistKey := "artist_top_tracks:" + strings.ToLower(artist)
	var s string

	s, err := a.Cache.Get(artistKey)
	if err != nil || s == "" {
		log.Printf(":: INFO :: ArtistTopTracks :: fetching top artist tracks from source for artist=%s", artist)
		tracks, err := a.Source.GetArtistTopTracks(artist)
		if err != nil || len(tracks.Tracks) < 1 {
			log.Printf(":: WARN :: ArtistTopTracks :: failure retrieving %s, %v", artist, err)
			s = fmt.Sprintf("unable to locate top tracks for artist %s", artist)
		} else {
			s = strings.Join(tracks.Tracks, ",")
			_, err := a.Cache.Set(artistKey, s, a.Config.CacheTTLDuration)
			if err != nil {
				log.Printf(":: ERROR :: ArtistTopTracks :: %v", err)
			}
		}
	}

	return &apiai.Response{
		Speech:      s,
		DisplayText: s,
	}, nil
}

func (a *App) ArtistBio(ctx context.Context, req *apiai.Request) (*apiai.Response, error) {
	artist := req.Param("artist")

	if artist == "" {
		return nil, fmt.Errorf("missing artist parameter")
	}

	artistKey := "artist_bio:" + strings.ToLower(artist)
	var s string

	s, err := a.Cache.Get(artistKey)
	if err != nil || s == "" {
		bio, err := a.Source.GetArtistBio(artist)
		log.Printf(":: INFO :: ArtistBio :: fetching bio from source for %s", artist)
		if err != nil || bio.Bio == "" {
			log.Printf(":: WARN :: ArtistBio :: unable to fetch artist %s, %v", artist, err)
			s = fmt.Sprintf("unable to locate biography for %s", artist)
		} else {
			s = bio.Bio
			_, err := a.Cache.Set(artistKey, s, a.Config.CacheTTLDuration)
			if err != nil {
				log.Printf(":: ERROR :: ArtistBio :: %v", err)
			}
		}
	}
	return &apiai.Response{
		Speech:      s,
		DisplayText: s,
	}, nil
}
