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
	"sync"
	"time"

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

var (
	albumsCache     []Album
	cacheExpiration = 1 * time.Hour // Adjust cache expiration time as needed
	cacheMutex      sync.RWMutex
	lastCacheUpdate time.Time
)

func GetAlbumSongsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	albumTitle := params["title"]

	var songs []Song
	var albumsDir = "D:/Music/MP3"
	albumFound := false

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
					return nil // Continue walking the directory
				}
				return err
			}
			if strings.Contains(metadata.Album(), albumTitle) {
				albumFound = true
				songs = append(songs, Song{
					Title:    metadata.Title(),
					Artist:   metadata.Artist(),
					Album:    metadata.Album(),
					Genre:    metadata.Genre(),
					Year:     metadata.Year(),
					FileName: filepath.Base(path),
				})
			}
		}
		return nil
	})

	if err != nil {
		http.Error(w, "Error reading songs", http.StatusInternalServerError)
		return
	}

	if !albumFound {
		// If no songs are found for the album, return an empty array
		songs = make([]Song, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func GetAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := strings.ToLower(query.Get("filter"))
	sortBy := query.Get("sort_by")

	cacheMutex.RLock()
	if len(albumsCache) > 0 && time.Since(lastCacheUpdate) < cacheExpiration {
		log.Println("Returning albums from cache")
		sendFilteredAlbums(w, albumsCache, filter, sortBy)
		cacheMutex.RUnlock()
		return
	}
	cacheMutex.RUnlock()

	log.Println("Fetching albums from file system")
	albums, err := fetchAlbums("D:/Music/MP3")
	if err != nil {
		http.Error(w, "Error reading albums", http.StatusInternalServerError)
		log.Println("Error fetching albums:", err)
		return
	}

	cacheMutex.Lock()
	albumsCache = albums
	lastCacheUpdate = time.Now()
	cacheMutex.Unlock()

	sendFilteredAlbums(w, albums, filter, sortBy)
}

func sendFilteredAlbums(w http.ResponseWriter, albums []Album, filter string, sortBy string) {
	filteredAlbums := albums
	if filter != "" {
		filteredAlbums = filterAlbums(albums, filter)
	}

	if sortBy != "" {
		sortAlbums(filteredAlbums, sortBy)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredAlbums)
}

func fetchAlbums(dir string) ([]Album, error) {
	var albums []Album

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(strings.ToLower(path), ".mp3") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			metadata, err := tag.ReadFrom(file)
			if err != nil {
				if err == tag.ErrNoTagsFound {
					return nil // Continue walking the directory
				}
				return err
			}

			title := metadata.Album()
			artist := metadata.Artist()
			artwork := getArtworkData(metadata)

			var existingAlbum *Album
			for i := range albums {
				if strings.EqualFold(albums[i].Title, title) && strings.EqualFold(albums[i].Artist, artist) {
					existingAlbum = &albums[i]
					break
				}
			}

			if existingAlbum == nil {
				newAlbum := Album{
					Title:   title,
					Artist:  artist,
					Artwork: artwork,
					Tracks:  []string{metadata.Title()},
				}
				albums = append(albums, newAlbum)
			} else {
				existingAlbum.Tracks = append(existingAlbum.Tracks, metadata.Title())
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return albums, nil
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
	var filteredAlbums []Album

	for _, album := range albums {
		if strings.Contains(strings.ToLower(album.Title), filter) || strings.Contains(strings.ToLower(album.Artist), filter) {
			filteredAlbums = append(filteredAlbums, album)
		}
	}

	return filteredAlbums
}

func sortAlbums(albums []Album, sortBy string) {
	switch sortBy {
	case "title":
		sort.Slice(albums, func(i, j int) bool {
			return strings.ToLower(albums[i].Title) < strings.ToLower(albums[j].Title)
		})
	case "artist":
		sort.Slice(albums, func(i, j int) bool {
			return strings.ToLower(albums[i].Artist) < strings.ToLower(albums[j].Artist)
		})
	}
}
