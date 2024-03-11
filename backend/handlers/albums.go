package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dhowden/tag"
)

type Album struct {
	Title   string   `json:"title"`
	Artist  string   `json:"artist"`
	Tracks  []string `json:"tracks"`
	Artwork string   `json:"artwork,omitempty"`
}

func GetAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := query.Get("filter")
	sortBy := query.Get("sort_by")

	albumsDir := "D:/Music/MP3"

	albumsMap := make(map[string]*Album)

	err := filepath.Walk(albumsDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(strings.ToLower(path), ".mp3") {

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			metadata, err := tag.ReadFrom(file)
			if err != nil {
				if err != tag.ErrNoTagsFound {
					return err
				}

				title := filepath.Base(path)
				artist := "Unknown Artist"

				if _, exists := albumsMap[title]; !exists {
					albumsMap[title] = &Album{
						Title:  title,
						Artist: artist,
						Tracks: []string{title},
					}
				}

				return nil
			}

			title := metadata.Album()
			artist := metadata.Artist()

			if _, exists := albumsMap[title]; !exists {
				albumsMap[title] = &Album{
					Title:   title,
					Artist:  artist,
					Artwork: getArtworkData(metadata),
				}
			} else {
				albumsMap[title].Tracks = append(albumsMap[title].Tracks, metadata.Title())
			}
		}
		return nil
	})

	if err != nil {
		http.Error(w, "Error reading albums", http.StatusInternalServerError)
		return
	}

	var albums []Album
	for _, album := range albumsMap {
		albums = append(albums, *album)
	}

	if filter != "" {
		albums = filterAlbums(albums, filter)
	}

	if sortBy != "" {
		sortAlbums(albums, sortBy)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(albums)
}

func getArtworkData(metadata tag.Metadata) string {
	if artwork := metadata.Picture(); artwork != nil {
		return base64.StdEncoding.EncodeToString(artwork.Data)
	}
	return ""
}

func filterAlbums(albums []Album, filter string) []Album {
	filteredAlbums := make([]Album, 0)

	for _, album := range albums {
		if strings.Contains(strings.ToLower(album.Title), strings.ToLower(filter)) || strings.Contains(strings.ToLower(album.Artist), strings.ToLower(filter)) {
			filteredAlbums = append(filteredAlbums, album)
		}
	}

	return filteredAlbums
}

func sortAlbums(albums []Album, sortBy string) {
	switch sortBy {
	case "title":
		sort.Slice(albums, func(i, j int) bool {
			return albums[i].Title < albums[j].Title
		})
	case "artist":
		sort.Slice(albums, func(i, j int) bool {
			return albums[i].Artist < albums[j].Artist
		})
	}
}
