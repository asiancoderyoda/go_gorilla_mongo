package routes

import (
	"go-gorilla-mongo/cmd/api/controllers"
	"go-gorilla-mongo/cmd/api/middlewares"

	"github.com/gorilla/mux"
)

func ConfigureRoutes() *mux.Router {
	router := mux.NewRouter()
	userRoutes := router.PathPrefix("/users").Subrouter()
	userRoutes.HandleFunc("/login", controllers.Login).Methods("POST")
	userRoutes.HandleFunc("/register", controllers.Register).Methods("POST")
	router.Use(middlewares.EnableCors)
	return router
}
