package lastfm

import (
	"fmt"

	"github.com/jesseward/songexplorer/config"
	"github.com/jesseward/songexplorer/internal/sources"
	lfm "github.com/shkh/lastfm-go/lastfm"
)

// LastFM config container for a LastFM source.
type LastFM struct {
	c   *config.Config
	api *lfm.Api
}

// GetSimilarTracks retrieves similiar songs based on a seed artist + track name
func (c *LastFM) GetSimilarTracks(artist, track string) (*sources.CacheableSimilarTracks, error) {
	songs, err := c.api.Track.GetSimilar(lfm.P{"artist": artist, "track": track, "limit": c.c.MaxTopSimTracks})
	if err != nil {
		return nil, fmt.Errorf("unable to GetSimilarTracks for %s - %s, error: %v", artist, track, err)
	}

	var st []string
	for _, tracks := range songs.Tracks {
		st = append(st, fmt.Sprintf("%s by %s", tracks.Name, tracks.Artist.Name))
	}
	return &sources.CacheableSimilarTracks{SimilarTracks: st}, nil
}

// GetSimilarArtists retrieves the top N similar artists based on seeded/user supplied artist
func (c *LastFM) GetSimilarArtists(artist string) (*sources.CacheableSimilarArtists, error) {
	artists, err := c.api.Artist.GetSimilar(lfm.P{"artist": artist, "limit": c.c.MaxTopSimArtists})
	if err != nil {
		return nil, fmt.Errorf("unable to GetSimilar for %s, %v", artist, err)
	}

	var a []string
	for _, artist := range artists.Similars {
		a = append(a, artist.Name)
	}
	return &sources.CacheableSimilarArtists{SimilarArtist: a}, nil
}

// GetArtistTopTracks fetches top N songs for the requested artist.
func (c *LastFM) GetArtistTopTracks(artist string) (*sources.CacheableTopTracks, error) {
	tracks, err := c.api.Artist.GetTopTracks(lfm.P{"artist": artist, "limit": c.c.MaxTopArtistTracks})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch top tracks for %s, %v", artist, err)
	}
	var t []string
	for _, tr := range tracks.Tracks {
		t = append(t, tr.Name)
	}
	return &sources.CacheableTopTracks{Tracks: t}, nil
}

// GetArtistBio retrieves artist biography information from the LAST.fm API.
func (c *LastFM) GetArtistBio(artist string) (*sources.CachableBio, error) {
	info, err := c.api.Artist.GetInfo(lfm.P{"artist": artist})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch bio for %s, %v", artist, err)
	}
	return &sources.CachableBio{Bio: info.Bio.Content}, nil
}

// New returns a LASTFM API client
func New(config *config.Config) (*LastFM, error) {
	api := lfm.New(config.SourceAPIKey, config.SourceSharedSecret)
	return &LastFM{config, api}, nil
}
