package role

import (
	"github.com/gorilla/mux"
)

func RegisterRoleRoutes(router *mux.Router, handler *RoleHandler) {
	router.HandleFunc("/roles", handler.GetRoles).Methods("GET")
	router.HandleFunc("/roles/{id}", handler.GetRole).Methods("GET")
}
