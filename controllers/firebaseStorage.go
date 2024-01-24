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

func (s *Storage) Delete(listing *models.Listing) error {
	return s.NewRef("listings/" + listing.ID.String()).Delete(context.Background())
}
