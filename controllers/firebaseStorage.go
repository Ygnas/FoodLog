package controllers

import (
	"context"
	"log"

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
