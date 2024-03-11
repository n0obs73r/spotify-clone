package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func GetAlbumArt(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	albumArtPath := filepath.Join("./mp3", title+".jpg")

	http.ServeFile(w, r, albumArtPath)
}
