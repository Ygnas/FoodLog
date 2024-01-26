package main

import (
	"net/http"

	"github.com/Ygnas/FoodLog/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/listings", controllers.GetAllListings)
	r.Get("/listings/{id}", controllers.GetListing)
	r.Post("/listings", controllers.CreateListing)
	r.Delete("/listings/{id}", controllers.DeleteListing)

	http.ListenAndServe(":3000", r)
}
