package controllers

import (
	"context"

	"github.com/Ygnas/FoodLog/models"
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

func (s *Storage) Create(listing *models.Listing) error {
	if err := s.NewRef("listings/"+listing.ID.String()).Set(context.Background(), listing); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Delete(id string) error {
	return s.NewRef("listings/" + id).Delete(context.Background())
}

func (s *Storage) GetListing(id string) (*models.Listing, error) {
	var listing models.Listing
	if err := s.NewRef("listings/"+id).Get(context.Background(), &listing); err != nil {
		return nil, err
	}
	return &listing, nil
}