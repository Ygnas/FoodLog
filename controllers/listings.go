package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Ygnas/FoodLog/models"
	"github.com/google/uuid"
)

func GetListing(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	listing := models.Listing{
		ID:          uuid.New(),
		Title:       "Example Listing",
		Description: "This is a sample listing",
		Shared:      true,
		Image:       "example.jpg",
		Type:        models.Snack,
		Rating:      5,
		Location:    "Sample Location",
	}

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

func CreateListing(w http.ResponseWriter, r *http.Request) {
	var listing models.Listing

	err := json.NewDecoder(r.Body).Decode(&listing)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

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
