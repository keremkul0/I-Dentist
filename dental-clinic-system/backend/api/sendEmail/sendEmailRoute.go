package sendEmail

import (
	"github.com/gorilla/mux"
)

func RegisterSendEmailRoutes(router *mux.Router, handler *SendEmailHandler) {
	router.HandleFunc("/send-verification-email", handler.SendVerificationEmail).Methods("POST")
}
