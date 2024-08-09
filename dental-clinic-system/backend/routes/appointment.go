package routes

import (
	"dental-clinic-system/handlers"

	"github.com/gorilla/mux"
)

func RegisterAppointmentRoutes(router *mux.Router, handler *handlers.AppointmentHandler) {
	router.HandleFunc("/appointments", handler.CreateAppointment).Methods("POST")
	router.HandleFunc("/appointments", handler.GetAppointments).Methods("GET")
}
