package routes

import (
	"dental-clinic-system/handlers"
	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router, authHandler *handlers.AuthHandler) {
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
}
