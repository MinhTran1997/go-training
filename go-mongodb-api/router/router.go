package router

import (
	"go-mongodb-api/handlers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", handlers.HomeLink)
	router.HandleFunc("/employees", handlers.GetAllEmployees).Methods("GET")
	router.HandleFunc("/employee/{id}", handlers.GetEmployeeByID).Methods("GET")
	router.HandleFunc("/employee", handlers.AddEmployee).Methods("POST")
	router.HandleFunc("/employee/{id}", handlers.UpdateEmployeeById).Methods("PUT")
	router.HandleFunc("/employee/{id}", handlers.DeleteEmployeeByID).Methods("DELETE")

	return router
}
