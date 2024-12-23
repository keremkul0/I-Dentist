package verifyEmail

import (
	"github.com/gorilla/mux"
)

func RegisterVerifyEmailRoutes(router *mux.Router, handler *VerifyEmailHandler) {
	router.HandleFunc("/verify-email", handler.VerifyUserEmail).Methods("GET")
}
