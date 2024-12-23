package signUpClinic

import (
	"github.com/gorilla/mux"
)

func RegisterSignupClinicRoutes(router *mux.Router, handler *SignUpClinicController) {

	router.HandleFunc("/singup-clinic", handler.SignUpClinic).Methods("Post")
}
