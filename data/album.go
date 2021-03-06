package data

// Album represents an album known to wavepipe, and contains information
// extracted from song tags
type Album struct {
	ID       int    `json:"id"`
	Artist   string `json:"artist"`
	ArtistID int    `db:"artist_id" json:"artistId"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
}

// AlbumFromSong creates a new Album from a Song model, extracting its
// fields as needed to build the struct
func AlbumFromSong(song *Song) *Album {
	return &Album{
		Artist: song.Artist,
		Title:  song.Album,
		Year:   song.Year,
	}
}

// Delete removes an existing Album from the database
func (a *Album) Delete() error {
	return DB.DeleteAlbum(a)
}

// Load pulls an existing Album from the database
func (a *Album) Load() error {
	return DB.LoadAlbum(a)
}

// Save creates a new Album in the database
func (a *Album) Save() error {
	return DB.SaveAlbum(a)
}
