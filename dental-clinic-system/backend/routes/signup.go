package routes

// Your other code goes here

import (
	"dental-clinic-system/handlers"

	"github.com/gorilla/mux"
)

func RegisterSignupRoutes(router *mux.Router, handler *handlers.SignupHandler) {
	router.HandleFunc("/signup", handler.Signup).Methods("POST")
}