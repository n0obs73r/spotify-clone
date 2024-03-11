package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dhowden/tag"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
)

type Album struct {
	Title   string   `json:"title"`
	Artist  string   `json:"artist"`
	Tracks  []string `json:"tracks"`
	Artwork string   `json:"artwork,omitempty"`
}

func GetAlbumSongsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	albumTitle := params["title"]

	var songs []Song
	var albumsDir = "D:/Music/MP3"
	albumFound := false

	err := filepath.Walk(albumsDir, func(path string, info os.FileInfo, err error) error {
		// log.Println("Walking path:", path)
		if strings.HasSuffix(strings.ToLower(path), ".mp3") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			metadata, err := tag.ReadFrom(file)
			if err != nil {
				if err == tag.ErrNoTagsFound {
					// log.Printf("No tags found for file: %s\n", path)
					return nil // Continue walking the directory
				}
				return err
			}
			if strings.Contains(metadata.Album(), albumTitle) {
				println(metadata.Album())
				albumFound = true
				songs = append(songs, Song{
					Title:  metadata.Title(),
					Artist: metadata.Artist(),
				})
			}
		}
		return nil
	})

	if err != nil {
		http.Error(w, "Error reading songs", http.StatusInternalServerError)
		// log.Println("Error walking directory:", err)
		return
	}

	if !albumFound {
		songs = append(songs, Song{
			Title:  "Unknown",
			Artist: "Unknown Artist",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func GetAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := strings.ToLower(query.Get("filter"))
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
				if err == tag.ErrNoTagsFound {
					log.Printf("No tags found for file: %s\n", path)
					return nil // Continue walking the directory
				}
				return err
			}

			title := strings.ToLower(metadata.Album())

			if _, exists := albumsMap[title]; !exists {
				albumsMap[title] = &Album{
					Title:   metadata.Album(),
					Artist:  metadata.Artist(),
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
		log.Println("2) Error walking directory:", err)
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
		img, err := imaging.Decode(bytes.NewReader(artwork.Data))
		if err != nil {
			log.Println("Error decoding image:", err)
			return ""
		}

		resizedImg := imaging.Resize(img, 100, 100, imaging.Lanczos)

		var buf bytes.Buffer
		if err := imaging.Encode(&buf, resizedImg, imaging.PNG); err != nil {
			log.Println("Error encoding image:", err)
			return ""
		}

		return base64.StdEncoding.EncodeToString(buf.Bytes())
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
