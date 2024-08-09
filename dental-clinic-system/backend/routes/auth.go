package routes

import (
    "github.com/gorilla/mux"
    "dental-clinic-system/handlers"
)

func RegisterAuthRoutes(router *mux.Router, authHandler *handlers.AuthHandler) {
    router.HandleFunc("/login", authHandler.Login).Methods("POST")
}
