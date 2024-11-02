package main

import (
	"dental-clinic-system/api/appointment"
	"dental-clinic-system/api/clinic"
	"dental-clinic-system/api/login"
	"dental-clinic-system/api/logout"
	"dental-clinic-system/api/patient"
	"dental-clinic-system/api/procedure"
	"dental-clinic-system/api/role"
	"dental-clinic-system/api/signUpClinic"
	"dental-clinic-system/api/singUpUser"
	"dental-clinic-system/api/user"
	"dental-clinic-system/application/appointmentService"
	"dental-clinic-system/application/clinicService"
	"dental-clinic-system/application/loginService"
	"dental-clinic-system/application/patientService"
	"dental-clinic-system/application/procedureService"
	"dental-clinic-system/application/roleService"
	"dental-clinic-system/application/signUpClinicService"
	"dental-clinic-system/application/singUpUserService"
	"dental-clinic-system/application/tokenService"
	"dental-clinic-system/application/userService"
	background_jobs "dental-clinic-system/background-jobs"
	"dental-clinic-system/helpers"
	"dental-clinic-system/middleware/authMiddleware"
	"dental-clinic-system/models"
	"dental-clinic-system/redisService"
	"dental-clinic-system/repository/appointmentRepository"
	"dental-clinic-system/repository/clinicRepository"
	"dental-clinic-system/repository/loginRepository"
	"dental-clinic-system/repository/patientRepository"
	"dental-clinic-system/repository/procedureRepository"
	"dental-clinic-system/repository/roleRepository"
	"dental-clinic-system/repository/tokenRepository"
	"dental-clinic-system/repository/userRepository"
	"dental-clinic-system/vault"
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
		&models.ExpiredTokens{},
	)
	if err != nil {
		log.Fatal("AutoMigrate error:", err)
	}

	jwtKey, err := vault.GetJWTKeyFromVault()
	if err != nil {
		log.Fatalf("Error retrieving JWT key from Vault: %v", err)
	}
	log.Printf("Retrieved JWT key: %s", jwtKey)
	helpers.NewJWTKeyHelper(jwtKey)

	router := mux.NewRouter()

	//Repositories
	newClinicRepository := clinicRepository.NewClinicRepository(db)
	newAppointmentRepository := appointmentRepository.NewAppointmentRepository(db)
	newPatientRepository := patientRepository.NewPatientRepository(db)
	newProcedureRepository := procedureRepository.NewProcedureRepository(db)
	newRoleRepository := roleRepository.NewRoleRepository(db)
	newUserRepository := userRepository.NewUserRepository(db)
	newLoginRepository := loginRepository.NewLoginRepository(db)
	newTokenRepository := tokenRepository.NewTokenRepository(db)

	//Services
	newClinicService := clinicService.NewClinicService(newClinicRepository)
	newAppointmentService := appointmentService.NewAppointmentService(newAppointmentRepository)
	newPatientService := patientService.NewPatientService(newPatientRepository)
	newProcedureService := procedureService.NewProcedureService(newProcedureRepository)
	newRoleService := roleService.NewRoleService(newRoleRepository)
	newUserService := userService.NewUserService(newUserRepository)
	newLoginService := loginService.NewLoginService(newLoginRepository)
	newSignUpClinicService := signUpClinicService.NewSignUpClinicService(newClinicRepository, newUserRepository)
	newTokenService := tokenService.NewTokenService(newTokenRepository)
	newSignUpUserService := singUpUserService.NewSignUpUserService(newUserRepository)

	//Handlers
	newClinicHandler := clinic.NewClinicHandlerController(newClinicService, newUserService)
	newAppointmentHandler := appointment.NewAppointmentHandlerController(newAppointmentService, newUserService)
	newPatientHandler := patient.NewPatientController(newPatientService, newUserService)
	newProcedureHandler := procedure.NewProcedureController(newProcedureService, newUserService)
	newRoleHandler := role.NewRoleController(newRoleService)
	newUserHandler := user.NewUserController(newUserService)
	newLoginHandler := login.NewLoginController(newLoginService)
	newSignUpClinicHandler := signUpClinic.NewSignUpClinicController(newSignUpClinicService)
	newSignUpUserHandler := singUpUser.NewSignUpUserHandler(newSignUpUserService)
	newLogoutHandler := logout.NewLogoutController(newTokenService)

	//Middleware
	securedRouter := router.PathPrefix("/api").Subrouter()
	securedRouter.Use(authMiddleware.AuthMiddleware)

	//Routes
	login.RegisterAuthRoutes(router, newLoginHandler)
	signUpClinic.RegisterSignupClinicRoutes(router, newSignUpClinicHandler)
	singUpUser.RegisterSignupUserRoutes(router, newSignUpUserHandler)

	//Secured Routes
	clinic.RegisterClinicRoutes(securedRouter, newClinicHandler)
	appointment.RegisterAppointmentRoutes(securedRouter, newAppointmentHandler)
	patient.RegisterPatientsRoutes(securedRouter, newPatientHandler)
	procedure.RegisterProcedureRoutes(securedRouter, newProcedureHandler)
	role.RegisterRoleRoutes(securedRouter, newRoleHandler)
	user.RegisterUserRoutes(securedRouter, newUserHandler)
	logout.RegisterLogoutRoutes(securedRouter, newLogoutHandler)

	//background services
	background_jobs.StartCleanExpiredJwtTokens(newTokenService)

	//Initialize Redis
	redisService.InitializeRedis()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
