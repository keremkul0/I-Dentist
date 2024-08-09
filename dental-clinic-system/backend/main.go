package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"dental-clinic-system/handlers"
	"dental-clinic-system/models"
	"dental-clinic-system/routes"
)

func main() {
	dsn := "host=localhost user=clinicuser password=clinicpassword dbname=clinicdb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Gerekli migrasyonları çalıştırıyoruz
	db.AutoMigrate(&models.Clinic{}, &models.Appointment{}, &models.User{}, &models.Role{}, &models.Doctor{}, &models.Assistant{}, &models.Secretary{}, &models.Procedure{})

	router := mux.NewRouter()

	// Auth Handler
	authHandler := &handlers.AuthHandler{DB: db}
	routes.RegisterAuthRoutes(router, authHandler)

	// Auth Middleware
	securedRouter := router.PathPrefix("/api").Subrouter()
	securedRouter.Use(authHandler.AuthMiddleware)

	// Clinic Handler
	clinicHandler := &handlers.ClinicHandler{DB: db}
	routes.RegisterClinicRoutes(securedRouter, clinicHandler)

	// Appointment Handler
	appointmentHandler := &handlers.AppointmentHandler{DB: db}
	routes.RegisterAppointmentRoutes(securedRouter, appointmentHandler)

	// User Handler
	userHandler := &handlers.UserHandler{DB: db}
	routes.RegisterUserRoutes(securedRouter, userHandler)

	// Role Handler
	roleHandler := &handlers.RoleHandler{DB: db}
	routes.RegisterRoleRoutes(securedRouter, roleHandler)

	// Doctor Handler
	doctorHandler := &handlers.DoctorHandler{DB: db}
	routes.RegisterDoctorRoutes(securedRouter, doctorHandler)

	// Assistant Handler
	assistantHandler := &handlers.AssistantHandler{DB: db}
	routes.RegisterAssistantRoutes(securedRouter, assistantHandler)

	// Secretary Handler
	secretaryHandler := &handlers.SecretaryHandler{DB: db}
	routes.RegisterSecretaryRoutes(securedRouter, secretaryHandler)

	// Procedure Handler
	procedureHandler := &handlers.ProcedureHandler{DB: db}
	routes.RegisterProcedureRoutes(securedRouter, procedureHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
