package sources

type CachableBio struct {
	Bio string // Artist biography
}

type CacheableTopTracks struct {
	Tracks []string // Artist (top) tracks
}

type CacheableSimilarArtists struct {
	SimilarArtist []string // Artists that are similar
}

type CacheableSimilarTracks struct {
	SimilarTracks []string
}

// Source defines the contract required to implement another 3rd party source for look-up data.
type Source interface {
	GetSimilarTracks(artist, track string) (*CacheableSimilarTracks, error)
	GetArtistBio(artist string) (*CachableBio, error)
	GetArtistTopTracks(artist string) (*CacheableTopTracks, error)
	GetSimilarArtists(artist string) (*CacheableSimilarArtists, error)
}
