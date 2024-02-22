package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ygnas/FoodLog/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	var user models.User

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user.ID = uuid.New()
	user.CreatedAt = time.Now()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	storage := NewStorage()
	err = storage.RegisterUser(&user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(responseJSON))
}

func Login(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}

	var user models.User

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	storage := NewStorage()
	storedUser, err := storage.LoginUser(&user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	jwt := GetTokenAuth()
	tokenString := jwt.GetToken(map[string]interface{}{
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"name":  storedUser.Name,
		"email": storedUser.Email,
	})

	_, err = json.Marshal(storedUser)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(tokenString))
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	err := GetFirebaseDatabase().FirebaseConnect()
	if err != nil {
		http.Error(w, "Could not connect to Firebase", http.StatusInternalServerError)
		return
	}
	id := chi.URLParam(r, "id")

	storage := NewStorage()
	err = storage.DeleteUser(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User deleted"))
}
