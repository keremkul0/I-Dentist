package routes

import (
    "dental-clinic-system/handlers"
    "github.com/gorilla/mux"
)

func RegisterSecretaryRoutes(router *mux.Router, handler *handlers.SecretaryHandler) {
    router.HandleFunc("/secretaries", handler.GetSecretaries).Methods("GET")
    router.HandleFunc("/secretaries/{id}", handler.GetSecretary).Methods("GET")
    router.HandleFunc("/secretaries", handler.CreateSecretary).Methods("POST")
    router.HandleFunc("/secretaries/{id}", handler.UpdateSecretary).Methods("PUT")
    router.HandleFunc("/secretaries/{id}", handler.DeleteSecretary).Methods("DELETE")
}
