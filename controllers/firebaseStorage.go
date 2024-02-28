package controllers

import (
	"context"
	"log"
	"strings"

	"github.com/Ygnas/FoodLog/models"
	"github.com/Ygnas/FoodLog/util"
)

type Storage struct {
	*FirebaseDatabase
}

func NewStorage() *Storage {
	db := GetFirebaseDatabase()
	return &Storage{
		FirebaseDatabase: db,
	}
}

func (s *Storage) Create(emailHash string, listing *models.Listing) error {
	listing.UserEmail = util.Base64Decode(emailHash)
	if err := s.NewRef("listings/").Child(emailHash).Child(listing.ID.String()).Set(context.Background(), listing); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Delete(emailHash string, id string) error {
	return s.NewRef("listings/").Child(emailHash).Child(id).Delete(context.Background())
}

func (s *Storage) GetListing(emailHash string, id string) (*models.Listing, error) {
	var listing models.Listing
	if err := s.NewRef("listings/").Child(emailHash).Child(id).Get(context.Background(), &listing); err != nil {
		return nil, err
	}
	return &listing, nil
}

func (s *Storage) GetAllUserListings(emailHash string) ([]*models.Listing, error) {
	var listingsMap map[string]*models.Listing

	if err := s.NewRef("listings").Child(emailHash).Get(context.Background(), &listingsMap); err != nil {
		log.Println(err)
		return nil, err
	}

	var listings []*models.Listing
	for _, listing := range listingsMap {
		listings = append(listings, listing)
	}

	return listings, nil
}

func (s *Storage) UpdateListing(emailHash string, listing *models.Listing) error {
	if err := s.NewRef("listings/").Child(emailHash).Child(listing.ID.String()).Set(context.Background(), listing); err != nil {
		return err
	}
	return nil
}

func (s *Storage) RegisterUser(user *models.User) error {
	if err := s.NewRef("users/"+util.Base64Encode(user.Email)).Set(context.Background(), user); err != nil {
		return err
	}
	return nil

}

func (s *Storage) LoginUser(user *models.User) (*models.User, error) {
	var returnedUser models.User
	if err := s.NewRef("users/"+util.Base64Encode(user.Email)).Get(context.Background(), &returnedUser); err != nil {
		return nil, err
	}
	return &returnedUser, nil
}

func (s *Storage) GetAllListings() ([]*models.Listing, error) {
	var listingsMap map[string]map[string]*models.Listing

	if err := s.NewRef("listings").Get(context.Background(), &listingsMap); err != nil {
		log.Println(err)
		return nil, err
	}

	var listings []*models.Listing
	for _, userListing := range listingsMap {
		for _, listing := range userListing {
			listings = append(listings, listing)
		}
	}

	return listings, nil
}

func (s *Storage) DeleteUser(emailHash string) error {
	s.DeleteAllUserListings(emailHash)
	return s.NewRef("users").Child(emailHash).Delete(context.Background())
}

func (s *Storage) DeleteAllUserListings(emailHash string) error {
	return s.NewRef("listings").Child(emailHash).Delete(context.Background())
}

func (s *Storage) LikeListing(listingID string, listingEmail string, email string) error {
	var listing models.Listing

	if err := s.NewRef("listings").Child(listingEmail).Child(listingID).Get(context.Background(), &listing); err != nil {
		return err
	}

	existingIndex := -1

	for index, like := range listing.Likes {
		if like.Email == email {
			existingIndex = index
			break
		}
	}

	if existingIndex != -1 {
		listing.Likes = append(listing.Likes[:existingIndex], listing.Likes[existingIndex+1:]...)
	} else {
		listing.Likes = append(listing.Likes, models.Like{Email: email})
	}

	return s.NewRef("listings").Child(listingEmail).Child(listingID).Set(context.Background(), listing)

}

func (s *Storage) CommentListing(listingID string, listingEmail string, comment models.Comment) error {
	var listing models.Listing

	if err := s.NewRef("listings").Child(listingEmail).Child(listingID).Get(context.Background(), &listing); err != nil {
		return err
	}

	listing.Comments = append(listing.Comments, models.Comment{Email: comment.Email, Comment: comment.Comment, CreatedAt: comment.CreatedAt})

	return s.NewRef("listings").Child(listingEmail).Child(listingID).Set(context.Background(), listing)

}

func (s *Storage) UploadImage(listingID string, image []byte) (string, error) {
	imagePath := "listings/" + listingID + ".jpg"
	bucket, err := s.Storage.DefaultBucket()
	if err != nil {
		return "", err
	}
	imageRef := bucket.Object(imagePath)
	writer := imageRef.NewWriter(context.Background())
	writer.ContentType = "image/jpeg"
	if _, err := writer.Write(image); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	lowerImagePath := strings.ToLower(imagePath)
	lowerImagePath = strings.ReplaceAll(lowerImagePath, "/", "%2f")

	return "https://firebasestorage.googleapis.com/v0/b/foodlog-9c3fd.appspot.com/o/" + lowerImagePath + "?alt=media", nil
}

func (s *Storage) DeleteImage(listingID string) error {
	imagePath := "listings/" + listingID + ".jpg"
	bucket, err := s.Storage.DefaultBucket()
	if err != nil {
		return err
	}
	imageRef := bucket.Object(imagePath)
	return imageRef.Delete(context.Background())
}
