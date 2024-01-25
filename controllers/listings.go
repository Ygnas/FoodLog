package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Ygnas/FoodLog/models"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func GetListing(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	storage := NewStorage()
	listing, err := storage.GetListing(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(listing)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(responseJSON))
}

func CreateListing(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	var listing models.Listing

	err = json.NewDecoder(r.Body).Decode(&listing)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	listing.ID = uuid.New()

	storage := NewStorage()
	err = storage.Create(&listing)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(listing)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(responseJSON))
}

func DeleteListing(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	storage := NewStorage()
	err = storage.Delete(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Listing deleted"))
}
