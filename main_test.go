package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ygnas/FoodLog/models"
	"github.com/Ygnas/FoodLog/util"
	"github.com/stretchr/testify/require"
)

func executeRequest(req *http.Request, r *Router) *httptest.ResponseRecorder {
	httptest := httptest.NewRecorder()
	r.Router.ServeHTTP(httptest, req)

	return httptest
}

var newListing = models.Listing{
	Title:       "Test",
	Description: "Test",
	Shared:      true,
	Image:       "Test",
	Type:        models.Snack,
	Likes: []models.Like{
		{Email: "test@test.com"},
	},
	Comments: []models.Comment{
		{Email: "test@test.com", Comment: "Test", CreatedAt: time.Now()},
	},
	Location:  models.Location{Latitude: 0, Longitude: 0},
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

var newUser = models.User{
	Email:    "gotest@gotest.com",
	Name:     "gotest",
	Password: "gotest",
}

var testToken string

func TestRegister(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	jsonInput, err := json.Marshal(newUser)
	require.NoError(t, err)

	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonInput))
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	require.NotEmpty(t, response.Body.String())
}

func TestLogin(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	jsonInput, err := json.Marshal(newUser)
	require.NoError(t, err)

	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonInput))
	response := executeRequest(req, r)

	testToken = response.Body.String()
	require.Equal(t, http.StatusOK, response.Code)
	require.NotEmpty(t, response.Body.String())
}

func TestGetListingEmpty(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	req, _ := http.NewRequest("GET", "/listings/00000", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	require.NotEmpty(t, response.Body.String())
}

func TestCreateListing(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	jsonInput, err := json.Marshal(newListing)
	require.NoError(t, err)

	req, _ := http.NewRequest("POST", "/listings", bytes.NewBuffer(jsonInput))
	req.Header.Set("Authorization", "Bearer "+testToken)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	require.NotEmpty(t, response.Body.String())

	err = json.Unmarshal(response.Body.Bytes(), &newListing)
	if err != nil {
		t.Error(err)
	}
}

func TestGetListing(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	var listing models.Listing

	req, _ := http.NewRequest("GET", "/listings/"+newListing.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	require.NotEmpty(t, response.Body.String())
	json.NewDecoder(response.Body).Decode(&listing)
	require.Equal(t, "Test", listing.Title)
}

func TestGetAllUserListings(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	req, _ := http.NewRequest("GET", "/listings", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	require.NotEmpty(t, response.Body.String())
}

func TestUpdateListing(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	var listing models.Listing

	newListing.Title = "Test-updated"
	jsonInput, err := json.Marshal(newListing)
	require.NoError(t, err)

	req, _ := http.NewRequest("PUT", "/listings/"+newListing.ID.String(), bytes.NewBuffer(jsonInput))
	req.Header.Set("Authorization", "Bearer "+testToken)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	json.NewDecoder(response.Body).Decode(&listing)
	require.Equal(t, "Test-updated", listing.Title)
}

func TestDeleteListing(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	req, _ := http.NewRequest("DELETE", "/listings/"+newListing.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "Listing deleted", response.Body.String())
}

func TestGetAllListings(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	req, _ := http.NewRequest("GET", "/all-listings", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	require.NotEmpty(t, response.Body.String())
}

func TestDeleteUserByID(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	req, _ := http.NewRequest("DELETE", "/users/delete/"+util.Base64Encode(newUser.Email), nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "User deleted", response.Body.String())
}
