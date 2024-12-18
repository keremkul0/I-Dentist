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
	"dental-clinic-system/config"
	"dental-clinic-system/helpers"
	"dental-clinic-system/init_func"
	"dental-clinic-system/middleware/authMiddleware"
	"dental-clinic-system/middleware/contextTimeoutMiddleware"
	"dental-clinic-system/repository/appointmentRepository"
	"dental-clinic-system/repository/clinicRepository"
	"dental-clinic-system/repository/loginRepository"
	"dental-clinic-system/repository/patientRepository"
	"dental-clinic-system/repository/procedureRepository"
	"dental-clinic-system/repository/redisRepository"
	"dental-clinic-system/repository/roleRepository"
	"dental-clinic-system/repository/tokenRepository"
	"dental-clinic-system/repository/userRepository"
	"dental-clinic-system/vault"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {

	configModel := init_func.SetConfig("resources")

	zerolog.SetGlobalLevel(configModel.Log.Level)

	clientVault, err := vault.ConnectVault(configModel.Vault)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to vault")
		panic("Error connecting to vault")
	}

	err = config.ReadConfigFromVault(clientVault, configModel)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading config from vault")
		panic("Error reading config from vault")
	}

	db := init_func.ConnectDatabase(configModel.Database)
	Rdb := init_func.ConnectRedis(configModel.Redis)
	mailDialer := init_func.SetupMailDialer(configModel.Email)

	init_func.MigrateDatabase(db)
	helpers.SetJWTKey(configModel.JWT.SecretKey)

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
	newRedisRepository := redisRepository.NewRedisRepository(Rdb)

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
	newSignUpUserService := signUpUserService.NewSignUpUserService(newUserRepository, newRedisRepository)
	newEmailService := emailService.NewEmailService(newUserRepository, newTokenRepository, mailDialer)

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

	//Create a new router
	router := mux.NewRouter()

	//Create subRouters
	securedRouter := router.PathPrefix("/api").Subrouter()

	//Middlewares
	newAuthMiddleware := authMiddleware.NewAuthMiddleware(newTokenService)

	//Middleware injection
	router.Use(contextTimeoutMiddleware.TimeoutMiddleware(5))
	securedRouter.Use(newAuthMiddleware.Authenticate)

	// Register Routes
	login.RegisterAuthRoutes(router, newLoginHandler)
	signUpClinic.RegisterSignupClinicRoutes(router, newSignUpClinicHandler)
	singUpUser.RegisterSignupUserRoutes(router, newSignUpUserHandler)
	verifyEmail.RegisterVerifyEmailRoutes(router, newVerifyEmailHandler)

	// Register Secured Routes
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

	log.Info().Msg(fmt.Sprintf("Server started on port %d", configModel.Server.Port))
	log.Fatal().Err(http.ListenAndServe(fmt.Sprintf(":%d", configModel.Server.Port), router)).
		Msg("Server failed to start")
}
