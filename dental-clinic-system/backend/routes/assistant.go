package routes

import (
	"dental-clinic-system/handlers"

	"github.com/gorilla/mux"
)

func RegisterAssistantRoutes(router *mux.Router, handler *handlers.AssistantHandler) {
	router.HandleFunc("/assistants", handler.GetAssistants).Methods("GET")
	router.HandleFunc("/assistants/{id}", handler.GetAssistant).Methods("GET")
	router.HandleFunc("/assistants", handler.CreateAssistant).Methods("POST")
	router.HandleFunc("/assistants/{id}", handler.UpdateAssistant).Methods("PUT")
	router.HandleFunc("/assistants/{id}", handler.DeleteAssistant).Methods("DELETE")
}
