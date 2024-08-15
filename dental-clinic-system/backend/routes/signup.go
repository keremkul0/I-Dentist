package routes

import (
    "dental-clinic-system/handlers"
    "github.com/gorilla/mux"
)

func RegisterSignupRoutes(router *mux.Router, handler *handlers.SignupHandler) {
    router.HandleFunc("/signup", handler.Signup).Methods("POST")
}
