package auth

import (
	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router, authHandler *AuthHandler) {
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
}
