package main

import (
	"fmt"
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
	dsn := "host=localhost user=clinicuser password=clinicpassword dbname=clinicdb port=5432 sslmode=disable TimeZone=Asia/Shanghai search_path=public"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		log.Fatal("Database connection test failed:", err)
	}
	fmt.Println("Database connection successful")
	db.Exec("CREATE SCHEMA IF NOT EXISTS public")
	// Gerekli migrasyonları çalıştırıyoruz
	err = db.Transaction(func(tx *gorm.DB) error {
		return tx.AutoMigrate(
			&models.Appointment{},
			&models.Assistant{},
			&models.Clinic{},
			&models.Doctor{},
			&models.Group{},
			&models.Patient{},
			&models.Procedure{},
			&models.Role{},
			&models.Secretary{},
			&models.UserRole{},
			&models.User{},
		)
	})

	if err != nil {
		log.Fatal("AutoMigrate error:", err)
	}
	router := mux.NewRouter()

	// Singup Handeler
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
