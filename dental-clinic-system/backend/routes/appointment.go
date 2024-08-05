package routes

import (
    "github.com/gorilla/mux"
    "dental-clinic-system/handlers"
)

func RegisterAppointmentRoutes(router *mux.Router, handler *handlers.AppointmentHandler) {
    router.HandleFunc("/appointments", handler.GetAppointments).Methods("GET")
    router.HandleFunc("/appointments/{id}", handler.GetAppointment).Methods("GET")
    router.HandleFunc("/appointments", handler.CreateAppointment).Methods("POST")
    router.HandleFunc("/appointments/{id}", handler.UpdateAppointment).Methods("PUT")
    router.HandleFunc("/appointments/{id}", handler.DeleteAppointment).Methods("DELETE")
}
