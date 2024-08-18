package signup

// Your other code goes here

import (
	"github.com/gorilla/mux"
)

func RegisterSignupRoutes(router *mux.Router, handler *SignupHandler) {
	router.HandleFunc("/signup", handler.Signup).Methods("POST")
}
