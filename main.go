package main

import (
	"net/http"

	"github.com/Ygnas/FoodLog/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
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
	controllers.NewJwt()
	jwt := controllers.GetTokenAuth()
	r.Router.Use(middleware.Logger)

	r.Router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwt.TokenAuth))
		r.Use(jwtauth.Authenticator(jwt.TokenAuth))

		r.Get("/listings", controllers.GetAllListings)
		r.Get("/listings/{id}", controllers.GetListing)
		r.Post("/listings", controllers.CreateListing)
		r.Put("/listings/{id}", controllers.UpdateListing)
		r.Delete("/listings/{id}", controllers.DeleteListing)
	})

	r.Router.Group(func(r chi.Router) {
		r.Post("/users/register", controllers.Register)
		r.Post("/users/login", controllers.Login)
	})
}
