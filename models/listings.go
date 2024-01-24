package models

import "github.com/google/uuid"

type MealType string

const (
	Breakfast MealType = "Breakfast"
	Lunch     MealType = "Lunch"
	Dinner    MealType = "Dinner"
	Snack     MealType = "Snack"
	Dessert   MealType = "Dessert"
)

type Listing struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Shared      bool      `json:"shared"`
	Image       string    `json:"image"`
	Type        MealType  `json:"type"`
	Rating      int       `json:"rating"`
	Location    string    `json:"location"`
}
