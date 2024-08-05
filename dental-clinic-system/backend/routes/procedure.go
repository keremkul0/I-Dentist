package routes

import (
    "github.com/gorilla/mux"
    "dental-clinic-system/handlers"
)

func RegisterProcedureRoutes(router *mux.Router, handler *handlers.ProcedureHandler) {
    router.HandleFunc("/procedures", handler.GetProcedures).Methods("GET")
    router.HandleFunc("/procedures/{id}", handler.GetProcedure).Methods("GET")
    router.HandleFunc("/procedures", handler.CreateProcedure).Methods("POST")
    router.HandleFunc("/procedures/{id}", handler.UpdateProcedure).Methods("PUT")
    router.HandleFunc("/procedures/{id}", handler.DeleteProcedure).Methods("DELETE")
}
