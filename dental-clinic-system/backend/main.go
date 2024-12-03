package main

import (
	"dental-clinic-system/api/appointment"
	"dental-clinic-system/api/clinic"
	"dental-clinic-system/api/login"
	"dental-clinic-system/api/logout"
	"dental-clinic-system/api/patient"
	"dental-clinic-system/api/procedure"
	"dental-clinic-system/api/role"
	"dental-clinic-system/api/sendEmail"
	"dental-clinic-system/api/signUpClinic"
	"dental-clinic-system/api/singUpUser"
	"dental-clinic-system/api/user"
	"dental-clinic-system/api/verifyEmail"
	"dental-clinic-system/application/appointmentService"
	"dental-clinic-system/application/clinicService"
	"dental-clinic-system/application/emailService"
	"dental-clinic-system/application/loginService"
	"dental-clinic-system/application/patientService"
	"dental-clinic-system/application/procedureService"
	"dental-clinic-system/application/roleService"
	"dental-clinic-system/application/signUpClinicService"
	"dental-clinic-system/application/singUpUserService"
	"dental-clinic-system/application/tokenService"
	"dental-clinic-system/application/userService"
	"dental-clinic-system/background-jobs"
	"dental-clinic-system/helpers"
	"dental-clinic-system/middleware/authMiddleware"
	"dental-clinic-system/models"
	"dental-clinic-system/redis/redisRepository"
	"dental-clinic-system/repository/appointmentRepository"
	"dental-clinic-system/repository/clinicRepository"
	"dental-clinic-system/repository/loginRepository"
	"dental-clinic-system/repository/patientRepository"
	"dental-clinic-system/repository/procedureRepository"
	"dental-clinic-system/repository/roleRepository"
	"dental-clinic-system/repository/tokenRepository"
	"dental-clinic-system/repository/userRepository"
	"dental-clinic-system/vault"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
)

func main() {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Int("port", 8080).Msg("Starting server on port")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
		panic(err)
	}

	dsn := os.Getenv("DNS")
	db, err := gorm.Open(
		postgres.Open(dsn),
	)

	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
		panic(err)
	}
	log.Info().Msg("Database connection successful")
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
		log.Fatal().Err(err).Msg("Error migrating models")
		panic(err)
	}

	var Rdb *redis.Client
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err = Rdb.Ping(Rdb.Context()).Result()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to Redis")
		panic(err)
	}

	emailPort, err := strconv.Atoi(os.Getenv("Email_port"))
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing email port")
		panic(err)
	}

	Mail := gomail.NewDialer(os.Getenv("Email_host"), emailPort, os.Getenv("Email_user"), os.Getenv("Email_password"))

	jwtKey, err := vault.GetJWTKeyFromVault()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting JWT key from vault")
		panic(err)
	}
	log.Info().Str("JWT Key", string(jwtKey)).Msg("Retrieved JWT Key from Vault")
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

	//Redis Repository
	newRedisRepository := redisRepository.NewRedisUserRepository(Rdb)

	//Services
	newClinicService := clinicService.NewClinicService(newClinicRepository)
	newAppointmentService := appointmentService.NewAppointmentService(newAppointmentRepository)
	newPatientService := patientService.NewPatientService(newPatientRepository)
	newProcedureService := procedureService.NewProcedureService(newProcedureRepository)
	newRoleService := roleService.NewRoleService(newRoleRepository)
	newUserService := userService.NewUserService(newUserRepository)
	newLoginService := loginService.NewLoginService(newLoginRepository)
	newSignUpClinicService := signUpClinicService.NewSignUpClinicService(newClinicRepository, newUserRepository, newRedisRepository)
	newTokenService := tokenService.NewTokenService(newTokenRepository)
	newSignUpUserService := singUpUserService.NewSignUpUserService(newUserRepository, newRedisRepository)
	newEmailService := emailService.NewEmailService(newUserRepository, newTokenRepository, *Mail)

	//Middleware
	newAuthMiddleware := authMiddleware.NewAuthMiddleware(newTokenService)

	//Handlers
	newClinicHandler := clinic.NewClinicHandlerController(newClinicService, newUserService)
	newAppointmentHandler := appointment.NewAppointmentHandlerController(newAppointmentService, newUserService, newPatientService)
	newPatientHandler := patient.NewPatientController(newPatientService, newUserService)
	newProcedureHandler := procedure.NewProcedureController(newProcedureService, newUserService)
	newRoleHandler := role.NewRoleController(newRoleService)
	newUserHandler := user.NewUserController(newUserService)
	newLoginHandler := login.NewLoginController(newLoginService)
	newSignUpClinicHandler := signUpClinic.NewSignUpClinicController(newSignUpClinicService)
	newSignUpUserHandler := singUpUser.NewSignUpUserHandler(newSignUpUserService)
	newLogoutHandler := logout.NewLogoutController(newTokenService)
	newVerifyEmailHandler := verifyEmail.NewVerifyEmailController(newEmailService)
	newSendEmailHandler := sendEmail.NewSendEmailController(newEmailService)

	//Middleware
	securedRouter := router.PathPrefix("/api").Subrouter()
	securedRouter.Use(newAuthMiddleware.Authenticate)

	//Routes
	login.RegisterAuthRoutes(router, newLoginHandler)
	signUpClinic.RegisterSignupClinicRoutes(router, newSignUpClinicHandler)
	singUpUser.RegisterSignupUserRoutes(router, newSignUpUserHandler)
	verifyEmail.RegisterVerifyEmailRoutes(router, newVerifyEmailHandler)

	//Secured Routes
	clinic.RegisterClinicRoutes(securedRouter, newClinicHandler)
	appointment.RegisterAppointmentRoutes(securedRouter, newAppointmentHandler)
	patient.RegisterPatientsRoutes(securedRouter, newPatientHandler)
	procedure.RegisterProcedureRoutes(securedRouter, newProcedureHandler)
	role.RegisterRoleRoutes(securedRouter, newRoleHandler)
	user.RegisterUserRoutes(securedRouter, newUserHandler)
	logout.RegisterLogoutRoutes(securedRouter, newLogoutHandler)
	sendEmail.RegisterSendEmailRoutes(securedRouter, newSendEmailHandler)

	//background services
	background_jobs.StartCleanExpiredJwtTokens(newTokenService)

	log.Info().Msg("Server started on port 8080")
	log.Fatal().Err(http.ListenAndServe(":8080", router))
}
