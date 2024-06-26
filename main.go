package main

import (
	"net/http"
	"time"

	"github.com/Ygnas/FoodLog/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
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
	r.Router.Use(httprate.Limit(
		100,
		time.Minute*1,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	r.Router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwt.TokenAuth))
		r.Use(jwtauth.Authenticator(jwt.TokenAuth))

		r.Get("/listings", controllers.GetAllUserListings)
		r.Get("/all-listings", controllers.GetAllListings)
		r.Get("/listings/{id}", controllers.GetListing)
		r.Post("/listings", controllers.CreateListing)
		r.Put("/listings/{id}", controllers.UpdateListing)
		r.Delete("/listings/{id}", controllers.DeleteListing)
		r.Delete("/users/delete/{id}", controllers.DeleteUserByID)
		r.Post("/listings/{id}/{email}/like", controllers.LikeListing)
		r.Post("/listings/{id}/{email}/comment", controllers.CommentListing)

		r.Post("/upload/{id}", controllers.UploadImage)
		r.Delete("/images/{id}/delete", controllers.DeleteImage)
	})

	r.Router.Group(func(r chi.Router) {
		r.Post("/users/register", controllers.Register)
		r.Post("/users/login", controllers.Login)
	})
}
