package logout

import (
	"github.com/gorilla/mux"
)

func RegisterLogoutRoutes(router *mux.Router, handler *LogoutController) {
	router.HandleFunc("/logout", handler.Logout).Methods("POST")
}
