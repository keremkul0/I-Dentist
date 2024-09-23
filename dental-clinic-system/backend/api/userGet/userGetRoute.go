package userGet

import (
	"github.com/gorilla/mux"
)

func RegisterUserGetRoutes(router *mux.Router, handler *UserGetHandler) {
	router.HandleFunc("/userGetModels", handler.GetUsers).Methods("GET")
	router.HandleFunc("/userGetModels/{id}", handler.GetUser).Methods("GET")
	router.HandleFunc("/userGetModels", handler.CreateUser).Methods("POST")
	router.HandleFunc("/userGetModels/{id}", handler.UpdateUser).Methods("PUT")
	router.HandleFunc("/userGetModels/{id}", handler.DeleteUser).Methods("DELETE")
}
