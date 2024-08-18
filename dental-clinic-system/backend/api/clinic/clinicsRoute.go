package clinic

import (
	"github.com/gorilla/mux"
)

func RegisterClinicRoutes(router *mux.Router, handler *ClinicHandler) {
	router.HandleFunc("/clinics", handler.GetClinics).Methods("GET")
	router.HandleFunc("/clinics/{id}", handler.GetClinic).Methods("GET")
	router.HandleFunc("/clinics", handler.CreateClinic).Methods("POST")
	router.HandleFunc("/clinics/{id}", handler.UpdateClinic).Methods("PUT")
	router.HandleFunc("/clinics/{id}", handler.DeleteClinic).Methods("DELETE")
}
