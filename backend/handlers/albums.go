package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dhowden/tag"
)

type Album struct {
	Title  string   `json:"title"`
	Artist string   `json:"artist"`
	Tracks []string `json:"tracks"`
}

func GetAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := query.Get("filter")
	sortBy := query.Get("sort_by")

	albumsDir := "D:/Music/MP3"

	albums := make([]Album, 0)

	err := filepath.Walk(albumsDir, func(path string, info os.FileInfo, err error) error {

		if strings.HasSuffix(strings.ToLower(path), ".mp3") {

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			mp3File, err := tag.ReadFrom(file)
			if err != nil {
				if err != tag.ErrNoTagsFound {
					return err
				}

				title := filepath.Base(path)
				artist := "Unknown Artist"

				albums = append(albums, Album{
					Title:  title,
					Artist: artist,
					Tracks: []string{title},
				})

				return nil
			}

			title := mp3File.Title()
			artist := mp3File.Artist()
			album := mp3File.Album()

			var existingAlbum *Album
			for i := range albums {
				if albums[i].Title == album && albums[i].Artist == artist {
					existingAlbum = &albums[i]
					break
				}
			}

			if existingAlbum == nil {
				albums = append(albums, Album{
					Title:  album,
					Artist: artist,
					Tracks: []string{title},
				})
			} else {

				existingAlbum.Tracks = append(existingAlbum.Tracks, title)
			}
		}
		return nil
	})
	if err != nil {
		http.Error(w, "Error reading albums", http.StatusInternalServerError)
		return
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
