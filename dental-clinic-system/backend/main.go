package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"

	"dental-clinic-system/handlers"
	"dental-clinic-system/models"
	"dental-clinic-system/routes"
)

func main() {
	dsn := "host=localhost user=clinicuser password=clinicpassword dbname=clinicdb port=5432 sslmode=disable TimeZone=Asia/Shanghai search_path=public"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		log.Fatal("Database connection test failed:", err)
	}
	fmt.Println("Database connection successful")
	db.Exec("CREATE SCHEMA IF NOT EXISTS public")

	// Transaction olmadan AutoMigrate işlemi
	err = db.AutoMigrate(
		&models.Appointment{},
		&models.Clinic{},
		&models.Group{},
		&models.Patient{},
		&models.Procedure{},
		&models.Role{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("AutoMigrate error:", err)
	}

	router := mux.NewRouter()

	// SingUp Handler
	singupHandler := &handlers.SignupHandler{DB: db}
	routes.RegisterSignupRoutes(router, singupHandler)

	// Auth Handler
	authHandler := &handlers.AuthHandler{DB: db}
	routes.RegisterAuthRoutes(router, authHandler)

	// Auth Middleware
	securedRouter := router.PathPrefix("/api").Subrouter()
	securedRouter.Use(authHandler.AuthMiddleware)

	// Group Handler
	groupHandler := &handlers.GroupHandler{DB: db}
	routes.RegisterGroupRoutes(securedRouter, groupHandler)

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

	// Procedure Handler
	procedureHandler := &handlers.ProcedureHandler{DB: db}
	routes.RegisterProcedureRoutes(securedRouter, procedureHandler)

	// Patient Handler
	patientHandler := &handlers.PatientHandler{DB: db}
	routes.RegisterPatientsRoutes(securedRouter, patientHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
