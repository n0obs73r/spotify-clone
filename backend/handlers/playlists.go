package handlers

import (
	"encoding/json"

	"net/http"
)

type Playlist struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Tracks    []string `json:"tracks"`
	Owner     string   `json:"owner"`
	Collabor  bool     `json:"collaborative"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

var playlists []Playlist

func GetPlaylists(w http.ResponseWriter, r *http.Request) {
	playlists := []Playlist{
		{ID: "1", Name: "Playlist 1", Tracks: []string{"song1.mp3", "song2.mp3"}, Owner: "user1", Collabor: false},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlists)
}

func CreatePlaylist(w http.ResponseWriter, r *http.Request) {

	var newPlaylist Playlist
	if err := json.NewDecoder(r.Body).Decode(&newPlaylist); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	newPlaylist.ID = GeneratePlaylistID()

	playlists = append(playlists, newPlaylist)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPlaylist)
}

func UpdatePlaylist(w http.ResponseWriter, r *http.Request) {

}

func DeletePlaylist(w http.ResponseWriter, r *http.Request) {

}

func AddTrackToPlaylist(w http.ResponseWriter, r *http.Request) {

}

func RemoveTrackFromPlaylist(w http.ResponseWriter, r *http.Request) {

}

func GeneratePlaylistID() string {

	return "some_unique_id"
}
