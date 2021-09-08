package router

import (
	"go-postgres-api/handlers"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
// define all the api endpoints.
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/", handlers.HomeLink)
	router.HandleFunc("/user/{id}", handlers.GetUser).Methods("GET")
	router.HandleFunc("/user", handlers.GetAllUser).Methods("GET")
	router.HandleFunc("/newuser", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", handlers.UpdateUser).Methods("PUT")
	router.HandleFunc("/deleteuser/{id}", handlers.DeleteUser).Methods("DELETE")

	return router
}
