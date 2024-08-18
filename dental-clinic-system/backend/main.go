package main

import (
	"dental-clinic-system/api/appointment"
	"dental-clinic-system/api/auth"
	"dental-clinic-system/api/clinic"
	"dental-clinic-system/api/group"
	"dental-clinic-system/api/patient"
	"dental-clinic-system/api/procedure"
	"dental-clinic-system/api/role"
	"dental-clinic-system/api/signup"
	"dental-clinic-system/api/user"
	appointmentService "dental-clinic-system/application/appointment"
	"dental-clinic-system/models"
	appointmentRepository "dental-clinic-system/repository/appointment"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("DNS")
	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)

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
	singupHandler := &signup.SignupHandler{DB: db}
	signup.RegisterSignupRoutes(router, singupHandler)

	// Auth Handler
	authHandler := &auth.AuthHandler{DB: db}
	auth.RegisterAuthRoutes(router, authHandler)

	// Auth Middleware
	securedRouter := router.PathPrefix("/api").Subrouter()
	securedRouter.Use(authHandler.AuthMiddleware)

	// Group Handler
	groupHandler := group.NewGroupHandlerService(db)
	group.RegisterGroupRoutes(securedRouter, groupHandler)

	// Clinic Handler
	clinicHandler := clinic.NewClinicHandlerService(db)
	clinic.RegisterClinicRoutes(securedRouter, clinicHandler)

	// Appointment Handler
	appointmentRepository := appointmentRepository.NewAppointmentRepository(db)
	appointmentService := appointmentService.NewAppointmentService(appointmentRepository)
	appointmentHandler := appointment.NewAppointmentHandlerService(appointmentService)
	appointment.RegisterAppointmentRoutes(securedRouter, appointmentHandler)

	// User Handler
	userHandler := &user.UserHandler{DB: db}
	user.RegisterUserRoutes(securedRouter, userHandler)

	// Role Handler
	roleHandler := &role.RoleHandler{DB: db}
	role.RegisterRoleRoutes(securedRouter, roleHandler)

	// Procedure Handler
	procedureHandler := &procedure.ProcedureHandler{DB: db}
	procedure.RegisterProcedureRoutes(securedRouter, procedureHandler)

	// Patient Handler
	patientHandler := &patient.PatientHandler{DB: db}
	patient.RegisterPatientsRoutes(securedRouter, patientHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
