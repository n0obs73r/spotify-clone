// song.model.go

package handlers

// Song represents a song with its metadata.
type Song struct {
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	TrackNumber int    `json:"trackNumber"`
	Genre       string `json:"genre"`
	Year        int    `json:"year"`
	Duration    int    `json:"duration"` // Duration in seconds
}
