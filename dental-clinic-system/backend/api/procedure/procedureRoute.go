package procedure

import (
	"github.com/gorilla/mux"
)

func RegisterProcedureRoutes(router *mux.Router, handler *ProcedureHandler) {
	router.HandleFunc("/procedures", handler.GetProcedures).Methods("GET")
	router.HandleFunc("/procedures/{id}", handler.GetProcedure).Methods("GET")
	router.HandleFunc("/procedures", handler.CreateProcedure).Methods("POST")
	router.HandleFunc("/procedures/{id}", handler.UpdateProcedure).Methods("PUT")
	router.HandleFunc("/procedures/{id}", handler.DeleteProcedure).Methods("DELETE")
}
