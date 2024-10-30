package logout

import (
	"github.com/gorilla/mux"
)

func RegisterLogoutRoutes(router *mux.Router, handler *LogoutHandler) {
	router.HandleFunc("/logout", handler.Logout).Methods("POST")
}
