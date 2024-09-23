package main

import (
	"dental-clinic-system/api/appointment"
	"dental-clinic-system/api/auth"
	"dental-clinic-system/api/clinic"
	"dental-clinic-system/api/patient"
	"dental-clinic-system/api/procedure"
	"dental-clinic-system/api/role"
	"dental-clinic-system/api/signupClinic"
	"dental-clinic-system/api/user"
	"dental-clinic-system/api/userGet"
	"dental-clinic-system/application/appointmentService"
	"dental-clinic-system/application/authService"
	"dental-clinic-system/application/clinicService"
	"dental-clinic-system/application/patientService"
	"dental-clinic-system/application/procedureService"
	"dental-clinic-system/application/roleService"
	"dental-clinic-system/application/signupClinicService"
	"dental-clinic-system/application/userGetModelService"
	"dental-clinic-system/application/userService"
	"dental-clinic-system/models"
	"dental-clinic-system/repository/appointmentRepository"
	"dental-clinic-system/repository/authRepository"
	"dental-clinic-system/repository/clinicRepository"
	"dental-clinic-system/repository/patientRepository"
	"dental-clinic-system/repository/procedureRepository"
	"dental-clinic-system/repository/roleRepository"
	"dental-clinic-system/repository/userGetModelRepository"
	"dental-clinic-system/repository/userRepository"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		//&gorm.Config{
		//	Logger: logger.Default.LogMode(logger.Info),
		//},
	)

	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		log.Fatal("Database connection test failed:", err)
	}
	fmt.Println("Database connection successful")
	db.Exec("CREATE SCHEMA IF NOT EXISTS public")

	err = db.AutoMigrate(
		&models.Appointment{},
		&models.Clinic{},
		&models.Patient{},
		&models.Procedure{},
		&models.Role{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("AutoMigrate error:", err)
	}

	router := mux.NewRouter()

	// newSignUpClinicService
	newClinicRepository := clinicRepository.NewClinicRepository(db)
	newUserRepository := userRepository.NewUserRepository(db)
	newSignUpClinicService := signupClinicService.NewSignUpClinicService(newClinicRepository, newUserRepository)
	newSignUpClinic := signupClinic.NewSignUpClinicHandler(newSignUpClinicService)
	signupClinic.RegisterSignupClinicRoutes(router, newSignUpClinic)

	// Auth Handler
	newAuthRepository := authRepository.NewAuthRepository(db)
	newAuthService := authService.NewAuthService(newAuthRepository)
	authHandler := auth.NewAuthHandlerController(newAuthService)
	auth.RegisterAuthRoutes(router, authHandler)

	// Auth Middleware
	securedRouter := router.PathPrefix("/api").Subrouter()
	securedRouter.Use(authHandler.AuthMiddleware)

	// Clinic Handler
	newClinicRepository2 := clinicRepository.NewClinicRepository(db)
	newClinicService := clinicService.NewClinicService(newClinicRepository2)
	clinicHandler := clinic.NewClinicHandlerController(newClinicService)
	clinic.RegisterClinicRoutes(securedRouter, clinicHandler)

	// Appointment Handler
	newAppointmentRepository := appointmentRepository.NewAppointmentRepository(db)
	newAppointmentService := appointmentService.NewAppointmentService(newAppointmentRepository)
	appointmentHandler := appointment.NewAppointmentHandlerController(newAppointmentService)
	appointment.RegisterAppointmentRoutes(securedRouter, appointmentHandler)

	// User Handler
	newUserRepository2 := userRepository.NewUserRepository(db)
	newUserService := userService.NewUserService(newUserRepository2)
	userHandler := user.NewUserController(newUserService)
	user.RegisterUserRoutes(securedRouter, userHandler)

	// User Get Handler
	newUserGetRepository := userGetModelRepository.NewUserGetModelRepository(db)
	newUserGetService := userGetModelService.NewUserGetModelService(newUserGetRepository)
	userGetHandler := userGet.NewUserGetController(newUserGetService)
	userGet.RegisterUserGetRoutes(securedRouter, userGetHandler)

	// Role Handler
	newRoleRepository := roleRepository.NewRoleRepository(db)
	newRoleService := roleService.NewRoleService(newRoleRepository)
	roleHandler := role.NewRoleController(newRoleService)
	role.RegisterRoleRoutes(securedRouter, roleHandler)

	// Procedure Handler
	newProcedureRepository := procedureRepository.NewProcedureRepository(db)
	newProcedureService := procedureService.NewProcedureService(newProcedureRepository)
	procedureHandler := procedure.NewProcedureController(newProcedureService)
	procedure.RegisterProcedureRoutes(securedRouter, procedureHandler)

	// Patient Handler
	newPatientRepository := patientRepository.NewPatientRepository(db)
	newPatientService := patientService.NewPatientService(newPatientRepository)
	patientHandler := patient.NewPatientController(newPatientService)
	patient.RegisterPatientsRoutes(securedRouter, patientHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
