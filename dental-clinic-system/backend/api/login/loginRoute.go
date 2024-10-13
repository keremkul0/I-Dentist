package login

import (
	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router, handler *LoginHandler) {
	router.HandleFunc("/login", handler.Login).Methods("POST")
}
