package appointment

import (
	"github.com/gorilla/mux"
)

func RegisterAppointmentRoutes(router *mux.Router, handler *AppointmentHandler) {
	router.HandleFunc("/appointments", handler.CreateAppointment).Methods("POST")
	router.HandleFunc("/appointments/", handler.GetAppointments).Methods("GET")
	router.HandleFunc("/appointment/{id}", handler.GetAppointment).Methods("GET")
}
