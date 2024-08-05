package routes

import (
    "dental-clinic-system/handlers"
    "github.com/gorilla/mux"
)

func RegisterRoleRoutes(router *mux.Router, handler *handlers.RoleHandler) {
    router.HandleFunc("/roles", handler.GetRoles).Methods("GET")
    router.HandleFunc("/roles/{id}", handler.GetRole).Methods("GET")
    router.HandleFunc("/roles", handler.CreateRole).Methods("POST")
    router.HandleFunc("/roles/{id}", handler.UpdateRole).Methods("PUT")
    router.HandleFunc("/roles/{id}", handler.DeleteRole).Methods("DELETE")
}
