package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ygnas/FoodLog/models"
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
	Rating:      5,
	Location:    "Test",
}

func TestGetListing(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	req, _ := http.NewRequest("GET", "/listings/01ed812a-d465-4b4c-b3e7-15e46a005910", nil)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "{\"id\":\"01ed812a-d465-4b4c-b3e7-15e46a005910\",\"title\":\"Example Listing\",\"description\":\"This is a sample listing\",\"shared\":true,\"image\":\"example.jpg\",\"type\":\"Snack\",\"rating\":5,\"location\":\"Sample Location\"}", response.Body.String())
}

func TestGetAllListings(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	req, _ := http.NewRequest("GET", "/listings", nil)
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
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	require.NotEmpty(t, response.Body.String())

	err = json.Unmarshal(response.Body.Bytes(), &newListing)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateListing(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	var listing models.Listing

	newListing.Title = "Test-updated"
	jsonInput, err := json.Marshal(newListing)
	require.NoError(t, err)

	req, _ := http.NewRequest("PUT", "/listings/"+newListing.ID.String(), bytes.NewBuffer(jsonInput))
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	json.NewDecoder(response.Body).Decode(&listing)
	require.Equal(t, "Test-updated", listing.Title)
}

func TestDeleteListing(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	req, _ := http.NewRequest("DELETE", "/listings/"+newListing.ID.String(), nil)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "Listing deleted", response.Body.String())
}
