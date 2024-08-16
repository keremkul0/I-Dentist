package routes

import (
	"dental-clinic-system/handlers"

	"github.com/gorilla/mux"
)

func RegisterPatientsRoutes(router *mux.Router, patientHandler *handlers.PatientHandler) {

	router.HandleFunc("/patients", patientHandler.GetPatients).Methods("GET")
	router.HandleFunc("/patients/{id}", patientHandler.GetPatient).Methods("GET")
	router.HandleFunc("/patients", patientHandler.CreatePatient).Methods("POST")
	router.HandleFunc("/patients/{id}", patientHandler.UpdatePatient).Methods("PUT")
	router.HandleFunc("/patients/{id}", patientHandler.DeletePatient).Methods("DELETE")

}
