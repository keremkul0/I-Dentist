package singUpUser

import (
	"github.com/gorilla/mux"
)

func RegisterSignupUserRoutes(router *mux.Router, handler *SignUpUserHandler) {

	router.HandleFunc("/signup-user", handler.SignUpUser).Methods("POST")
}
