package controllers

import (
	"encoding/json"
	"io"
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

func GetAllUserListings(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())

	storage := NewStorage()
	listings, err := storage.GetAllUserListings(util.Base64Encode(claims["email"].(string)))
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

	if listing.ID == uuid.Nil {
		listing.ID = uuid.New()
	}
	if listing.CreatedAt.IsZero() {
		listing.CreatedAt = time.Now()
	}

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

func GetAllListings(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	storage := NewStorage()
	listings, err := storage.GetAllListings()
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

func LikeListing(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")
	email := chi.URLParam(r, "email")

	_, claims, _ := jwtauth.FromContext(r.Context())

	storage := NewStorage()
	err = storage.LikeListing(id, util.Base64Encode(email), claims["email"].(string))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Like updated"))
}

func CommentListing(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")
	email := chi.URLParam(r, "email")

	var comment models.Comment

	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	storage := NewStorage()
	err = storage.CommentListing(id, util.Base64Encode(email), comment)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Comment added"))
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	storage := NewStorage()

	imageBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	image, err := storage.UploadImage(id, imageBytes)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(image)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(responseJSON))
}

func DeleteImage(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	storage := NewStorage()
	err = storage.DeleteImage(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Image deleted"))
}
