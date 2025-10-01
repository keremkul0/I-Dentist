package resetPassword

import "github.com/gorilla/mux"

func RegisterResetPasswordRoutes(router *mux.Router, handler *ResetPasswordHandler) {
	router.HandleFunc("/reset-password", handler.ResetPassword).Methods("POST")
}
