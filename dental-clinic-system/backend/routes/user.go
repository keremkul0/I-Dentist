package routes

import (
    "github.com/gorilla/mux"
    "dental-clinic-system/handlers"
)

func RegisterUserRoutes(router *mux.Router, handler *handlers.UserHandler) {
    router.HandleFunc("/register", handler.Register).Methods("POST")
    router.HandleFunc("/login", handler.Login).Methods("POST")
}
