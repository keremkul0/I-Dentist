package signupClinic

import (
	"github.com/gorilla/mux"
)

func RegisterSignupClinicRoutes(router *mux.Router, handler *SignUpClinicHandler) {

	router.HandleFunc("/Singup", handler.SignUpClinic).Methods("Post")
}
