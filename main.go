package main

import (
	"net/http"

	"github.com/Ygnas/FoodLog/controllers"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/listings", controllers.GetListing)
	r.Post("/listings", controllers.CreateListing)

	http.ListenAndServe(":3000", r)
}
