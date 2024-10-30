package signupClinic

import (
	"github.com/gorilla/mux"
)

func RegisterSignupClinicRoutes(router *mux.Router, handler *SignUpClinicHandler) {

	router.HandleFunc("/singup", handler.SignUpClinic).Methods("Post")
}
