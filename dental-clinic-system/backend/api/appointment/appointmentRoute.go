package appointment

import (
	"github.com/gorilla/mux"
)

func RegisterAppointmentRoutes(router *mux.Router, handler *AppointmentHandler) {
	router.HandleFunc("/appointments", handler.GetAppointments).Methods("GET")
	router.HandleFunc("/appointment/{id}", handler.GetAppointment).Methods("GET")
	router.HandleFunc("/appointments", handler.CreateAppointment).Methods("POST")
	router.HandleFunc("/appointment/{id}", handler.UpdateAppointment).Methods("PUT")
	router.HandleFunc("/appointment/{id}", handler.DeleteAppointment).Methods("DELETE")
}
