package forgotPassword

import "github.com/gorilla/mux"

func RegisterForgotPasswordRoutes(router *mux.Router, handler *ForgotPasswordHandler) {
	router.HandleFunc("/forgot-password", handler.ForgotPassword).Methods("POST")
}
