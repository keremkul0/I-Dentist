package routes

import (
	"dental-clinic-system/handlers"

	"github.com/gorilla/mux"
)

func RegisterGroupRoutes(router *mux.Router, handler *handlers.GroupHandler) {
	router.HandleFunc("/groups", handler.GetGroups).Methods("GET")
	router.HandleFunc("/groups/{id}", handler.GetGroup).Methods("GET")
	router.HandleFunc("/groups", handler.CreateGroup).Methods("POST")
	router.HandleFunc("/groups/{id}", handler.UpdateGroup).Methods("PUT")
	router.HandleFunc("/groups/{id}", handler.DeleteGroup).Methods("DELETE")
	router.HandleFunc("/groups/{id}/clinics", handler.GetClinicsByGroup).Methods("GET")
}
