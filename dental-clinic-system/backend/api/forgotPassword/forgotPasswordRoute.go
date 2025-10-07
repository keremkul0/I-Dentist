package forgotPassword

import "github.com/gorilla/mux"

func RegisterForgotPasswordRoutes(router fiber.Router, handler *ForgotPasswordHandler) {
	router.Post("/forgot-password", handler.ForgotPassword)
}
