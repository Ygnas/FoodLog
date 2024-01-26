package main

import (
	"net/http"

	"github.com/Ygnas/FoodLog/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := CreateNewRouter()
	r.MountRoutes()

	http.ListenAndServe(":3000", r.Router)
}

type Router struct {
	Router *chi.Mux
}

func CreateNewRouter() *Router {
	r := &Router{}
	r.Router = chi.NewRouter()
	return r
}

func (r *Router) MountRoutes() {
	r.Router.Use(middleware.Logger)

	r.Router.Get("/listings", controllers.GetAllListings)
	r.Router.Get("/listings/{id}", controllers.GetListing)
	r.Router.Post("/listings", controllers.CreateListing)
	r.Router.Delete("/listings/{id}", controllers.DeleteListing)
}
