package models

import (
	"time"

	"github.com/google/uuid"
)

type MealType string

const (
	Breakfast MealType = "Breakfast"
	Lunch     MealType = "Lunch"
	Dinner    MealType = "Dinner"
	Snack     MealType = "Snack"
	Dessert   MealType = "Dessert"
)

type Comment struct {
	Email     string    `json:"email"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type Like struct {
	Email string `json:"email"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Listing struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Shared      bool      `json:"shared"`
	Image       string    `json:"image"`
	Type        MealType  `json:"type"`
	Likes       []Like    `json:"likes"`
	Location    Location  `json:"location"`
	Comments    []Comment `json:"comments"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
