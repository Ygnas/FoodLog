package main

import (
	"bytes"
	"encoding/json"
	"io"
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
	UserEmail: "gotest@gotest.com",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

var newUser = models.User{
	Email:    "gotest@gotest.com",
	Name:     "gotest",
	Password: "gotest",
}

var comment = models.Comment{
	Email:     "gotest@gotest.com",
	Comment:   "Test comment",
	CreatedAt: time.Now(),
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

func TestUploadImage(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	req, _ := http.NewRequest("POST", "/upload/"+newListing.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	req.Body = io.NopCloser(bytes.NewBuffer([]byte("notimage")))

	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
}

func TestDeleteImage(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	req, _ := http.NewRequest("DELETE", "/images/"+newListing.ID.String()+"/delete", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)
}

func TestLikeListing(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	var listing models.Listing

	req, _ := http.NewRequest("POST", "/listings/"+newListing.ID.String()+"/"+newListing.UserEmail+"/like", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/listings/"+newListing.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	response = executeRequest(req, r)

	json.NewDecoder(response.Body).Decode(&listing)

	require.Equal(t, 2, len(listing.Likes))
}

func TestCommentListing(t *testing.T) {
	r := CreateNewRouter()

	r.MountRoutes()

	var listing models.Listing

	jsonInput, err := json.Marshal(comment)
	require.NoError(t, err)

	req, _ := http.NewRequest("POST", "/listings/"+newListing.ID.String()+"/"+newListing.UserEmail+"/comment", bytes.NewBuffer(jsonInput))
	req.Header.Set("Authorization", "Bearer "+testToken)
	response := executeRequest(req, r)

	require.Equal(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/listings/"+newListing.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	response = executeRequest(req, r)

	json.NewDecoder(response.Body).Decode(&listing)

	require.Equal(t, 2, len(listing.Comments))
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

// Random sample data generation

// var testTokens map[string]string

// func TestCreateSampleListings(t *testing.T) {
// 	r := CreateNewRouter()
// 	r.MountRoutes()

// 	testTokens = make(map[string]string)

// 	titles := []string{"Cornflakes", "Chicken Salad", "Lamb Chop", "Tiramisu", "Sandwitch", "Pasta and pork", "Chowder in Bread Bowl", "Ice Cream", "Pancakes", "Burger", "Pizza", "Sushi", "Ramen", "Tacos", "Burrito", "Tofu Poke Bowl", "Pasta", "Wild Salmon", "Steak and asparagus", "Creme Brulle"}
// 	descriptions := []string{"", "", "Very delicious", "A little bit too sweet", "", "Not bad", "Need to try this again", "Yummy!", "My favorite breakfast", "Very tasty", "Delicious", "I love sushi", "I love ramen", "", "", "First time, but will try again!", "", "", "Not the fan of asparagus", "Crunchy and delicious"}
// 	locations := []models.Location{
// 		{Latitude: 51.897426, Longitude: -8.465763},
// 		{Latitude: 51.898514, Longitude: -8.475604},
// 		{Latitude: 51.899204, Longitude: -8.470565},
// 		{Latitude: 51.897928, Longitude: -8.473089},
// 		{Latitude: 51.903614, Longitude: -8.468399},
// 		{Latitude: 51.896891, Longitude: -8.486316},
// 		{Latitude: 51.897349, Longitude: -8.491550},
// 		{Latitude: 51.898202, Longitude: -8.495064},
// 		{Latitude: 51.900405, Longitude: -8.472643},
// 		{Latitude: 51.901567, Longitude: -8.459267},
// 	}

// 	// All the following images are CC0 (aka CC Zero) is a public dedication tool,
// 	// which enables creators to give up their copyright and put their works into
// 	// the worldwide public domain. CC0 enables reusers to distribute, remix, adapt,
// 	// and build upon the material in any medium or format, with no conditions. https://openverse.org/
// 	images := []string{"https://images.rawpixel.com/editor_1024/czNmcy1wcml2YXRlL3Jhd3BpeGVsX2ltYWdlcy93ZWJzaXRlX2NvbnRlbnQvbHIvcHUyMzM0MjE2LWltYWdlLWt3dndqbzd0LmpwZw.jpg",
// 		"https://live.staticflickr.com/5458/18065696241_c902746455_b.jpg",
// 		"https://live.staticflickr.com/5308/5674634678_d852770273_b.jpg",
// 		"https://images.rawpixel.com/editor_1024/czNmcy1wcml2YXRlL3Jhd3BpeGVsX2ltYWdlcy93ZWJzaXRlX2NvbnRlbnQvbHIvcHg2ODkzMDItaW1hZ2Uta3d2eGozbWYuanBn.jpg",
// 		"https://live.staticflickr.com/2293/32833088496_c4b6570ce1_b.jpg",
// 		"https://live.staticflickr.com/46/133683033_16b63c8c1f_b.jpg",
// 		"https://live.staticflickr.com/266/31991157410_353c81f2e7_b.jpg",
// 		"https://live.staticflickr.com/1090/1229438458_ba9423b26b_b.jpg",
// 		"https://cdn.stocksnap.io/img-thumbs/960w/RLOLDF9IQ4.jpg",
// 		"https://live.staticflickr.com/1481/25039204229_f5d0b6c857_b.jpg",
// 		"https://live.staticflickr.com/4339/36457904755_e14bae248a_b.jpg",
// 		"https://live.staticflickr.com/7146/13464656265_33607612fd_b.jpg",
// 		"https://live.staticflickr.com/1/178901_88efb0c812.jpg",
// 		"https://live.staticflickr.com/4550/24926715648_5f1bcbd5d6_b.jpg",
// 		"https://upload.wikimedia.org/wikipedia/commons/6/60/Burrito.JPG",
// 		"https://upload.wikimedia.org/wikipedia/commons/7/7a/Fried_Tofu_Poke_Bowl_%28M%29_with_Wasabi_Soy_sauce_-_Kitokito_2023-06-09.jpg",
// 		"https://live.staticflickr.com/1338/1151299870_b713be2ffa_b.jpg",
// 		"https://live.staticflickr.com/13/16612423_949b98c78c.jpg",
// 		"https://upload.wikimedia.org/wikipedia/commons/0/08/Steak_and_asparagus.jpg",
// 		"https://live.staticflickr.com/7285/9687444563_93fc3469f2_b.jpg"}
// 	likes := []models.Like{
// 		{Email: "lana@test.com"},
// 		{Email: "lukas@test.com"},
// 		{Email: "ignas@test.com"},
// 	}
// 	comments := []models.Comment{
// 		{Email: "lana@test.com", Comment: "Looks very nice!", CreatedAt: time.Now()},
// 		{Email: "ignas@test.com", Comment: "I would like to try this!", CreatedAt: time.Now()},
// 		{Email: "lukas@test.com", Comment: "I have tried this before, it was delicious!", CreatedAt: time.Now()},
// 		{Email: "ignas@test.com", Comment: "Nice.", CreatedAt: time.Now()},
// 		{Email: "lana@test.com", Comment: "Looks delicious!", CreatedAt: time.Now()},
// 		{Email: "lukas@test.com", Comment: "This is my favorite comfort food.", CreatedAt: time.Now()},
// 		{Email: "lukas@test.com", Comment: "This dish is absolutely delicious!", CreatedAt: time.Now()},
// 		{Email: "lana@test.com", Comment: "I love the unique blend of flavors in this!", CreatedAt: time.Now()},
// 		{Email: "ignas@test.com", Comment: "The presentation is stunning!", CreatedAt: time.Now()},
// 	}
// 	mealTypes := []models.MealType{models.Breakfast, models.Lunch, models.Dinner, models.Snack}

// 	now := time.Now()
// 	daysSinceMonday := (int(now.Weekday()) - int(time.Monday) + 7) % 7
// 	monday := now.AddDate(0, 0, -daysSinceMonday)

// 	sharedCounter := 0

// 	for _, userEmail := range []string{"ignas@test.com", "lukas@test.com", "lana@test.com"} {
// 		registerUser(t, userEmail)
// 		loginUser(t, userEmail)

// 		listingCount := 0

// 		skipProbabilitySnack := 0.3
// 		skipProbabilityOther := 0.2

// 		for i := 0; i < 20; i++ {
// 			date := monday.AddDate(0, 0, listingCount/4)
// 			mealType := mealTypes[i%len(mealTypes)]

// 			newListing := models.Listing{
// 				Title:       titles[i%len(titles)],
// 				Description: descriptions[i%len(descriptions)],
// 				Shared:      false,
// 				Image:       images[i%len(images)],
// 				Type:        mealType,
// 				Likes:       []models.Like{likes[i%len(likes)]},
// 				Comments:    []models.Comment{comments[i%len(comments)]},
// 				Location:    locations[i%len(locations)],
// 				UserEmail:   userEmail,
// 				CreatedAt:   date.Add(time.Duration(i) * time.Minute),
// 				UpdatedAt:   time.Now(),
// 			}

// 			if sharedCounter%6 == 0 {
// 				newListing.Shared = true
// 			}

// 			sharedCounter++

// 			if rand.Float64() < skipProbabilitySnack && mealType == models.Snack {
// 				continue
// 			}

// 			if rand.Float64() < skipProbabilityOther && mealType != models.Snack {
// 				continue
// 			}

// 			jsonListing, err := json.Marshal(newListing)
// 			require.NoError(t, err)
// 			createListingReq, _ := http.NewRequest("POST", "/listings", bytes.NewBuffer(jsonListing))
// 			createListingReq.Header.Set("Authorization", "Bearer "+testTokens[userEmail])
// 			response := executeRequest(createListingReq, r)
// 			require.Equal(t, http.StatusOK, response.Code)
// 			listingCount++
// 		}
// 	}
// }

// func registerUser(t *testing.T, email string) {
// 	r := CreateNewRouter()

// 	r.MountRoutes()

// 	index := strings.Index(email, "@")
// 	name := email[:index]

// 	newUser := models.User{
// 		Email:    email,
// 		Name:     name,
// 		Password: "test",
// 	}

// 	jsonUser, err := json.Marshal(newUser)
// 	require.NoError(t, err)
// 	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonUser))
// 	response := executeRequest(req, r)
// 	require.Equal(t, http.StatusOK, response.Code)
// }

// func loginUser(t *testing.T, email string) {
// 	r := CreateNewRouter()

// 	r.MountRoutes()
// 	newUser := models.User{
// 		Email:    email,
// 		Password: "test",
// 	}

// 	jsonUser, err := json.Marshal(newUser)
// 	require.NoError(t, err)
// 	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonUser))
// 	response := executeRequest(req, r)
// 	require.Equal(t, http.StatusOK, response.Code)

// 	testTokens[email] = response.Body.String()
// }
