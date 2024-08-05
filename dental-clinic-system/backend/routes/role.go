package routes

import (
    "github.com/gorilla/mux"
    "dental-clinic-system/handlers"
)

func RegisterRoleRoutes(router *mux.Router, handler *handlers.RoleHandler) {
    router.HandleFunc("/roles", handler.GetRoles).Methods("GET")
    router.HandleFunc("/roles", handler.CreateRole).Methods("POST")
}
