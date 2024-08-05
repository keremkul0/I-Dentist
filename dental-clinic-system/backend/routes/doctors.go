package routes

import (
	"dental-clinic-system/handlers"

	"github.com/gorilla/mux"
)

func RegisterDoctorRoutes(router *mux.Router, handler *handlers.DoctorHandler) {
	router.HandleFunc("/doctors", handler.GetDoctors).Methods("GET")
	router.HandleFunc("/doctors/{id}", handler.GetDoctor).Methods("GET")
	router.HandleFunc("/doctors", handler.CreateDoctor).Methods("POST")
	router.HandleFunc("/doctors/{id}", handler.UpdateDoctor).Methods("PUT")
	router.HandleFunc("/doctors/{id}", handler.DeleteDoctor).Methods("DELETE")
}
