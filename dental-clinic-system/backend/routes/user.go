package routes

import (
	"dental-clinic-system/handlers"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, handler *handlers.UserHandler) {
	router.HandleFunc("/users", handler.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	router.HandleFunc("/users", handler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")
}
