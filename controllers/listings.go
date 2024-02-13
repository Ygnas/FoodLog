package controllers

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/Ygnas/FoodLog/models"
	"github.com/Ygnas/FoodLog/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

func GetListing(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")
	_, claims, _ := jwtauth.FromContext(r.Context())

	storage := NewStorage()
	listing, err := storage.GetListing(util.Base64Encode(claims["email"].(string)), id)
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

func GetAllListings(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())

	storage := NewStorage()
	listings, err := storage.GetAllListings(util.Base64Encode(claims["email"].(string)))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sort.Slice(listings, func(i, j int) bool {
		return listings[i].CreatedAt.After(listings[j].CreatedAt)
	})

	responseJSON, err := json.Marshal(listings)
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
	listing.CreatedAt = time.Now()
	_, claims, _ := jwtauth.FromContext(r.Context())

	storage := NewStorage()
	err = storage.Create(util.Base64Encode(claims["email"].(string)), &listing)
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
	_, claims, _ := jwtauth.FromContext(r.Context())

	storage := NewStorage()
	err = storage.Delete(util.Base64Encode(claims["email"].(string)), id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Listing deleted"))
}

func UpdateListing(w http.ResponseWriter, r *http.Request) {
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

	listing.UpdatedAt = time.Now()
	_, claims, _ := jwtauth.FromContext(r.Context())

	storage := NewStorage()
	err = storage.UpdateListing(util.Base64Encode(claims["email"].(string)), &listing)
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
